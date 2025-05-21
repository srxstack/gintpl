// Copyright 2025 莫维龙 <kalandramo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/srxstack/gintpl. The professional
// version of this repository is https://github.com/srxstack/srxstack.

package handler

import (
	"github.com/srxstack/gintpl/internal/apiserver/biz"
	"github.com/srxstack/gintpl/internal/apiserver/pkg/validation"
)

type Handler struct {
	biz biz.IBiz
	val *validation.Validator
}

func NewHandler(biz biz.IBiz, val *validation.Validator) *Handler {
	return &Handler{
		biz: biz,
		val: val,
	}
}
