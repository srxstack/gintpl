// Copyright 2025 莫维龙 <kalandramo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/srxstack/gintpl. The professional
// version of this repository is https://github.com/srxstack/srxstack.

package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/srxstack/gintpl/internal/pkg/contextx"
	"github.com/srxstack/gintpl/internal/pkg/known"
)

// RequestIDMiddleware 是一个 Gin 中间件，用于在每个 HTTP 请求的上下文和
// 响应中注入 `x-request-id` 键值对.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 `x-request-id`，如果不存在则生成新的 UUID
		requestID := c.Request.Header.Get(known.XRequestID)

		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 将 RequestID 保存到 context.Context 中，以便后续程序使用
		ctx := contextx.WithRequestID(c.Request.Context(), requestID)
		c.Request = c.Request.WithContext(ctx)

		// 将 RequestID 保存到 HTTP 返回头中，Header 的键为 `x-request-id`
		c.Writer.Header().Set(known.XRequestID, requestID)

		// 继续处理请求
		c.Next()
	}
}
