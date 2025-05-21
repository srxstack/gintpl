// Copyright 2025 莫维龙 <kalandramo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/srxstack/gintpl. The professional
// version of this repository is https://github.com/srxstack/srxstack.

package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/srxstack/srxstack/pkg/core"
)

func (h *Handler) CreateFrostmourneUser(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.FrostmourneV1().CreateUser, h.val.ValidateCreateFrostmourneUserRequest)
}

func (h *Handler) CreateFrostmourneDep(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.FrostmourneV1().CreateDep, h.val.ValidateCreateFrostmourneDepRequest)
}

func (h *Handler) CreateFrostmourneTeam(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.FrostmourneV1().CreateTeam, h.val.ValidateCreateFrostmourneTeamRequest)
}
