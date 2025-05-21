// Copyright 2025 莫维龙 <kalandramo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/srxstack/gintpl. The professional
// version of this repository is https://github.com/srxstack/srxstack.

package frostmourne

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/srxstack/gintpl/internal/pkg/log"
	apiv1 "github.com/srxstack/gintpl/pkg/api/apiserver/v1"
	"github.com/srxstack/gintpl/pkg/configz"
	"github.com/srxstack/gintpl/pkg/httpx"
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

func (b *frostmourneBiz) CreateDep(ctx context.Context, rq *apiv1.CreateFrostmourneDepRequest) (*apiv1.CommonFrostmourneResponse, error) {
	// 获取 token
	loginResult, err := b.login(ctx)
	if err != nil {
		return nil, err
	}
	toke := loginResult.GetResult()

	client := httpx.NewHTTPX()

	header := make(map[string]string)
	header["frostmourne-token"] = toke
	header["Content-Type"] = contentType
	client.SetHeader(header)

	data, err := json.Marshal(rq)
	if err != nil {
		return nil, err
	}

	url := configz.C.FrostmourneOptions.URL + apiCreateDep
	resp, err := client.POST(ctx, url, data)
	log.Debugw("Send create dept req", "url", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return b.resp(resp)
}

func (b *frostmourneBiz) CreateTeam(ctx context.Context, rq *apiv1.CreateFrostmourneTeamRequest) (*apiv1.CommonFrostmourneResponse, error) {
	// 获取 token
	loginResult, err := b.login(ctx)
	if err != nil {
		return nil, err
	}
	toke := loginResult.GetResult()

	client := httpx.NewHTTPX()

	header := make(map[string]string)
	header["frostmourne-token"] = toke
	header["Content-Type"] = "application/json; charset=utf-8"
	client.SetHeader(header)

	data, err := json.Marshal(rq)
	if err != nil {
		return nil, err
	}

	url := configz.C.FrostmourneOptions.URL
	resp, err := client.POST(ctx, url, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return b.resp(resp)
}

func (b *frostmourneBiz) CreateUser(ctx context.Context, rq *apiv1.CreateFrostmourneUserRequest) (*apiv1.CommonFrostmourneResponse, error) {
	// 获取 token
	loginResult, err := b.login(ctx)
	if err != nil {
		return nil, err
	}
	toke := loginResult.GetResult()

	client := httpx.NewHTTPX()

	header := make(map[string]string)
	header["frostmourne-token"] = toke
	header["Content-Type"] = "application/json; charset=utf-8"
	client.SetHeader(header)

	data, err := json.Marshal(rq)
	if err != nil {
		return nil, err
	}

	url := configz.C.FrostmourneOptions.URL
	resp, err := client.POST(ctx, url, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return b.resp(resp)
}

func (b *frostmourneBiz) resp(resp *http.Response) (*apiv1.CommonFrostmourneResponse, error) {
	var result *apiv1.CommonFrostmourneResponse

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}
	log.Debugw("Send create dept resp", "body", string(body))
	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (b *frostmourneBiz) login(ctx context.Context) (*apiv1.FrostmourneLoginResponse, error) {
	var result *apiv1.FrostmourneLoginResponse

	client := httpx.NewHTTPX()

	// url 组合
	url := configz.C.FrostmourneOptions.URL + apiLogin

	// 获取管理员账号密码
	var basicAuth apiv1.FrostmourneLoginRequest
	basicAuth.Username = configz.C.FrostmourneOptions.Username
	basicAuth.Password = configz.C.FrostmourneOptions.Password

	data, err := json.Marshal(&basicAuth)
	if err != nil {
		return result, err
	}

	resp, err := client.POST(ctx, url, data)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}
