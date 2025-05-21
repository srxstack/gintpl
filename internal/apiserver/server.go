// Copyright 2025 莫维龙 <kalandramo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/srxstack/gintpl. The professional
// version of this repository is https://github.com/srxstack/srxstack.

package apiserver

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/srxstack/srxstack/pkg/authz"
	genericoptions "github.com/srxstack/srxstack/pkg/options"
	"github.com/srxstack/srxstack/pkg/server"
	"github.com/srxstack/srxstack/pkg/store/where"
	"github.com/srxstack/srxstack/pkg/token"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"k8s.io/utils/ptr"

	"github.com/srxstack/gintpl/internal/apiserver/biz"
	"github.com/srxstack/gintpl/internal/apiserver/model"
	"github.com/srxstack/gintpl/internal/apiserver/pkg/validation"
	"github.com/srxstack/gintpl/internal/apiserver/store"
	"github.com/srxstack/gintpl/internal/pkg/contextx"
	"github.com/srxstack/gintpl/internal/pkg/known"
	"github.com/srxstack/gintpl/internal/pkg/log"
	mw "github.com/srxstack/gintpl/internal/pkg/middleware"
)

type Config struct {
	ServerMode        string
	JWTKey            string
	Expiration        time.Duration
	EnableMemoryStore bool
	HTTPOptions       *genericoptions.HTTPOptions
	TLSOptions        *genericoptions.TLSOptions
	MySQLOptions      *genericoptions.MySQLOptions
}

type UnionServer struct {
	srv server.Server
}

type ServerConfig struct {
	cfg       *Config
	biz       biz.IBiz
	val       *validation.Validator
	retriever mw.UserRetriever
	authz     *authz.Authz
}

func (cfg *Config) NewUnionServer() (*UnionServer, error) {
	// 注册租户解析函数，通过上下文获取用户 ID
	//nolint: gocritic
	where.RegisterTenant("userID", func(ctx context.Context) string {
		return contextx.UserID(ctx)
	})

	// 初始化 token 包的签名密钥、认证 Key 及 Token 默认过期时间
	token.Init(cfg.JWTKey, known.XUserID, cfg.Expiration)

	log.Infow("Initializing federation server", "server-mode", cfg.ServerMode, "enable-memory-store", cfg.EnableMemoryStore)

	// 创建服务配置，这些配置可用来创建服务器
	srv, err := InitializeWebServer(cfg)
	if err != nil {
		return nil, err
	}

	return &UnionServer{srv: srv}, nil
}

func (s *UnionServer) Run() error {
	go s.srv.RunOrDie()

	// 创建一个 os.Signal 类型的 channel，用于接收系统信号
	quit := make(chan os.Signal, 1)
	// 当执行 kill 命令时（不带参数），默认会发送 syscall.SIGTERM 信号
	// 使用 kill -2 命令会发送 syscall.SIGINT 信号（例如按 CTRL+C 触发）
	// 使用 kill -9 命令会发送 syscall.SIGKILL 信号，但 SIGKILL 信号无法被捕获，因此无需监听和处理
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 阻塞程序，等待从 quit channel 中接收到信号
	<-quit

	log.Infow("Shutting down server ...")

	// 优雅关闭服务
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 先关闭依赖的服务，再关闭被依赖的服务
	s.srv.GracefulStop(ctx)

	log.Infow("Server exited")
	return nil
}

// NewDB 创建一个 *gorm.DB 实例.
func (cfg *Config) NewDB() (*gorm.DB, error) {
	if !cfg.EnableMemoryStore {
		log.Infow("Initializing database connection", "type", "mysql", "addr", cfg.MySQLOptions.Addr)
		return cfg.MySQLOptions.NewDB()
	}

	log.Infow("Initializing database connection", "type", "memory", "engine", "SQLite")
	// 使用SQLite内存模式配置数据库
	// ?cache=shared 用于设置 SQLite 的缓存模式为 共享缓存模式 (shared)。
	// 默认情况下，SQLite 的每个数据库连接拥有自己的独立缓存，这种模式称为 专用缓存 (private)。
	// 使用 共享缓存模式 (shared) 后，不同连接可以共享同一个内存中的数据库和缓存。
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Errorw("Failed to create database connection", "err", err)
		return nil, err
	}

	// 自动迁移数据库结构
	if err := db.AutoMigrate(&model.UserM{}, &model.CasbinRuleM{}); err != nil {
		log.Errorw("Failed to migrate database schema", "err", err)
		return nil, err
	}

	// 注意：这里仅仅为了实现快速部署，降低学习难度。
	// 在真实企业开发中，不能再代码中硬编码这些初始化配置，
	// 尤其是硬编码密码、密钥之类的信息.
	// 插入 casbin_rule 表记录
	adminR, userR := "role::admin", "role::user"
	casbinRules := []model.CasbinRuleM{
		{PType: ptr.To("g"), V0: ptr.To("user-000000"), V1: &adminR},
		{PType: ptr.To("p"), V0: &adminR, V1: ptr.To("*"), V2: ptr.To("*"), V3: ptr.To("allow")},
		{PType: ptr.To("p"), V0: &userR, V1: ptr.To("/v1.MiniBlog/DeleteUser"), V2: ptr.To("CALL"), V3: ptr.To("deny")},
		{PType: ptr.To("p"), V0: &userR, V1: ptr.To("/v1.MiniBlog/ListUser"), V2: ptr.To("CALL"), V3: ptr.To("deny")},
		{PType: ptr.To("p"), V0: &userR, V1: ptr.To("/v1/users"), V2: ptr.To("GET"), V3: ptr.To("deny")},
		{PType: ptr.To("p"), V0: &userR, V1: ptr.To("/v1/users/*"), V2: ptr.To("DELETE"), V3: ptr.To("deny")},
	}

	if err := db.Create(&casbinRules).Error; err != nil {
		log.Fatalw("Failed to insert casbin_rule records", "err", err)
		return nil, err
	}

	// 插入默认用户（root用户）
	user := model.UserM{
		UserID:    "user-000000",
		Username:  "root",
		Password:  "miniblog1234",
		Nickname:  "administrator",
		Email:     "colin404@foxmail.com",
		Phone:     "18110000000",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&user).Error; err != nil {
		log.Fatalw("Failed to insert default root user", "err", err)
		return nil, err
	}

	return db, nil
}

// UserRetriever 定义一个用户数据获取器. 用来获取用户信息.
type UserRetriever struct {
	store store.IStore
}

// GetUser 根据用户 ID 获取用户信息.
func (r *UserRetriever) GetUser(ctx context.Context, userID string) (*model.UserM, error) {
	return r.store.User().Get(ctx, where.F("userID", userID))
}

// ProvideDB 根据配置提供一个数据库实例。
func ProvideDB(cfg *Config) (*gorm.DB, error) {
	return cfg.NewDB()
}

func NewWebServer(serverMode string, serverConfig *ServerConfig) (server.Server, error) {
	switch serverMode {
	case gin.DebugMode:
		return serverConfig.NewGinServer(), nil
	case gin.ReleaseMode:
		return serverConfig.NewGinServer(), nil
	case gin.TestMode:
		return serverConfig.NewGinServer(), nil
	default:
		return serverConfig.NewGinServer(), nil
	}
}
