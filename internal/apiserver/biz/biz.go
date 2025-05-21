// Copyright 2025 莫维龙 <kalandramo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/srxstack/gintpl. The professional
// version of this repository is https://github.com/srxstack/srxstack.

package biz

//go:generate mockgen -destination mock_biz.go -package biz github.com/onexstack/miniblog/internal/apiserver/biz IBiz

import (
	"github.com/google/wire"
	"github.com/srxstack/srxstack/pkg/authz"

	frostmournev1 "github.com/srxstack/gintpl/internal/apiserver/biz/v1/frostmourne"
	userv1 "github.com/srxstack/gintpl/internal/apiserver/biz/v1/user"
	"github.com/srxstack/gintpl/internal/apiserver/store"
)

var ProviderSet = wire.NewSet(NewBiz, wire.Bind(new(IBiz), new(*biz)))

// IBiz 定义了业务层需要实现的方法.
type IBiz interface {
	// 获取用户业务接口.
	UserV1() userv1.UserBiz
	FrostmourneV1() frostmournev1.FrostmourneBiz
}

// biz 是 IBiz 的一个具体实现.
type biz struct {
	store store.IStore
	authz *authz.Authz
}

// 确保 biz 实现了 IBiz 接口.
var _ IBiz = (*biz)(nil)

// NewBiz 创建一个 IBiz 类型的实例.
func NewBiz(store store.IStore, authz *authz.Authz) *biz {
	return &biz{store: store, authz: authz}
}

// UserV1 返回一个实现了 UserBiz 接口的实例.
func (b *biz) UserV1() userv1.UserBiz {
	return userv1.New(b.store, b.authz)
}

func (b *biz) FrostmourneV1() frostmournev1.FrostmourneBiz {
	return frostmournev1.New()
}
