// Copyright 2025 莫维龙 <kalandramo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/srxstack/gintpl. The professional
// version of this repository is https://github.com/srxstack/srxstack.

package frostmourne

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/srxstack/gintpl/internal/pkg/log"
	"github.com/srxstack/gintpl/pkg/configz"

	"github.com/swxctx/ghttp"

	apiv1 "github.com/srxstack/gintpl/pkg/api/apiserver/v1"
)

const (
	apiCreateDep  = "/api/monitor-api/department/create?_appId=frostmourne&"
	apiCreateTeam = "/api/monitor-api/team/create?_appId=frostmourne&"
	apiCreateUser = "/api/monitor-api/userinfo/create?_appId=frostmourne&"
	apiLogin      = "/api/monitor-api/user/login"

	contentType = "application/json; charset=utf-8"
)

var _ FrostmourneBiz = (*frostmourneBiz)(nil)

type FrostmourneBiz interface {
	CreateDep(ctx context.Context, rq *apiv1.CreateFrostmourneDepRequest) (*apiv1.CommonFrostmourneResponse, error)
	CreateTeam(ctx context.Context, rq *apiv1.CreateFrostmourneTeamRequest) (*apiv1.CommonFrostmourneResponse, error)
	CreateUser(ctx context.Context, rq *apiv1.CreateFrostmourneUserRequest) (*apiv1.CommonFrostmourneResponse, error)
}

type frostmourneBiz struct{}

func New() *frostmourneBiz {
	return &frostmourneBiz{}
}

// 创建部门.
func (b *frostmourneBiz) CreateDep(ctx context.Context, rq *apiv1.CreateFrostmourneDepRequest) (*apiv1.CommonFrostmourneResponse, error) {
	return b.sendRequest(ctx, apiCreateDep, rq)
}

// 创建团队.
func (b *frostmourneBiz) CreateTeam(ctx context.Context, rq *apiv1.CreateFrostmourneTeamRequest) (*apiv1.CommonFrostmourneResponse, error) {
	return b.sendRequest(ctx, apiCreateTeam, rq)
}

// 创建用户.
func (b *frostmourneBiz) CreateUser(ctx context.Context, rq *apiv1.CreateFrostmourneUserRequest) (*apiv1.CommonFrostmourneResponse, error) {
	return b.sendRequest(ctx, apiCreateUser, rq)
}

// 发送请求的公共方法.
func (b *frostmourneBiz) sendRequest(ctx context.Context, apiPath string, request interface{}) (*apiv1.CommonFrostmourneResponse, error) {
	var result apiv1.CommonFrostmourneResponse

	// 获取token
	token, err := b.getToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取token失败: %w", err)
	}

	// 构建请求URL
	furl := configz.C.FrostmourneOptions.URL + apiPath

	// 序列化请求数据
	data, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("序列化请求数据失败: %w", err)
	}

	// 构建并发送请求
	req := ghttp.Request{
		Method:    "POST",
		Url:       furl,
		ShowDebug: true,
		Body:      data,
	}
	req.AddHeader("frostmourne-token", token)
	req.AddHeader("Content-Type", contentType)

	log.Debugw("发送请求", "url", furl)
	res, err := req.Do()
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}

	// 处理响应
	err = res.Body.FromToJson(&result)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// 获取token的公共方法.
func (b *frostmourneBiz) getToken(ctx context.Context) (string, error) {
	var result apiv1.FrostmourneLoginResponse

	furl := configz.C.FrostmourneOptions.URL + apiLogin

	// 获取管理员账号密码
	basicAuth := apiv1.FrostmourneLoginRequest{
		Username: configz.C.FrostmourneOptions.Username,
		Password: configz.C.FrostmourneOptions.Password,
	}

	data, err := json.Marshal(&basicAuth)
	if err != nil {
		return "", fmt.Errorf("序列化登录请求失败: %w", err)
	}

	req := ghttp.Request{
		Method:    "POST",
		Url:       furl,
		ShowDebug: true,
		Body:      data,
	}
	req.AddHeader("Content-Type", contentType)

	res, err := req.Do()
	if err != nil {
		return "", fmt.Errorf("登录请求失败: %w", err)
	}

	err = res.Body.FromToJson(&result)
	if err != nil {
		return "", fmt.Errorf("解析登录响应失败: %w", err)
	}

	if result.GetResult() == "" {
		return "", errors.New("获取token失败，返回结果为空")
	}

	return result.GetResult(), nil
}
