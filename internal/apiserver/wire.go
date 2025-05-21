// Copyright 2025 莫维龙 <kalandramo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/srxstack/gintpl. The professional
// version of this repository is https://github.com/srxstack/srxstack.

//go:build wireinject
// +build wireinject

package apiserver

import (
	"github.com/google/wire"
	"github.com/srxstack/srxstack/pkg/authz"

	"github.com/srxstack/srxstack/pkg/server"

	"github.com/srxstack/gintpl/internal/apiserver/biz"
	"github.com/srxstack/gintpl/internal/apiserver/pkg/validation"
	"github.com/srxstack/gintpl/internal/apiserver/store"
	ginmw "github.com/srxstack/gintpl/internal/pkg/middleware"
)

func InitializeWebServer(*Config) (server.Server, error) {
	wire.Build(
		wire.NewSet(NewWebServer, wire.FieldsOf(new(*Config), "ServerMode")),
		wire.Struct(new(ServerConfig), "*"), // * 表示注入全部字段
		wire.NewSet(store.ProviderSet, biz.ProviderSet),
		ProvideDB, // 提供数据库实例
		validation.ProviderSet,
		wire.NewSet(
			wire.Struct(new(UserRetriever), "*"),
			wire.Bind(new(ginmw.UserRetriever), new(*UserRetriever)),
		),
		authz.ProviderSet,
	)
	return nil, nil
}
