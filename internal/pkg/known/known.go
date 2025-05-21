// Copyright 2025 莫维龙 <kalandramo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/srxstack/gintpl. The professional
// version of this repository is https://github.com/srxstack/srxstack.

package known

const (
	// XRequestID 用来定义上下文中的键，代表请求 ID.
	XRequestID = "x-request-id"

	// XUserID 用来定义上下文的键，代表请求用户 ID. UserID 整个用户生命周期唯一.
	XUserID = "x-user-id"
)

// 定义其他常量.
const (
	// Admin 用户名.
	AdminUsername = "root"

	MaxErrGroupConcurrency = 1000
)
