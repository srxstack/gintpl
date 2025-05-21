// Copyright 2025 莫维龙 <kalandramo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/srxstack/gintpl. The professional
// version of this repository is https://github.com/srxstack/srxstack.

package validation

import (
	"context"

	genericvalidation "github.com/srxstack/srxstack/pkg/validation"

	"github.com/srxstack/gintpl/internal/pkg/errno"
	apiv1 "github.com/srxstack/gintpl/pkg/api/apiserver/v1"
)

// Validate 校验字段的有效性.
func (v *Validator) ValidateFrostmourneRules() genericvalidation.Rules {
	// 定义各字段的校验逻辑，通过一个 map 实现模块化和简化
	return genericvalidation.Rules{
		"Role": func(value any) error {
			if value.(string) == "" {
				return errno.ErrInvalidArgument.WithMessage("role cannot be empty")
			}
			return nil
		},
		"TeamId": func(value any) error {
			if value.(int32) == 0 {
				return errno.ErrInvalidArgument.WithMessage("teamId cannot be empty")
			}
			return nil
		},
		"Account": func(value any) error {
			if value.(string) == "" {
				return errno.ErrInvalidArgument.WithMessage("account cannot be empty")
			}
			return nil
		},
		"Password": func(value any) error {
			if value.(string) == "" {
				return errno.ErrInvalidArgument.WithMessage("password cannot be empty")
			}
			return nil
		},
		"FullName": func(value any) error {
			if value.(string) == "" {
				return errno.ErrInvalidArgument.WithMessage("fullName cannot be empty")
			}
			return nil
		},
		"DepartmentName": func(value any) error {
			if value.(string) == "" {
				return errno.ErrInvalidArgument.WithMessage("departmentName cannot be empty")
			}
			return nil
		},
		"DepartmentId": func(value any) error {
			if value.(int32) == 0 {
				return errno.ErrInvalidArgument.WithMessage("departmentId cannot be empty")
			}
			return nil
		},
		"TeamName": func(value any) error {
			if value.(string) == "" {
				return errno.ErrInvalidArgument.WithMessage("teamName cannot be empty")
			}
			return nil
		},
	}
}

// ValidateCreateFrostmourneUserRequest 校验 CreateFrostmourneUserRequest 结构体的有效性.
func (v *Validator) ValidateCreateFrostmourneUserRequest(ctx context.Context, rq *apiv1.CreateFrostmourneUserRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidateFrostmourneRules())
}

func (v *Validator) ValidateCreateFrostmourneDepRequest(ctx context.Context, rq *apiv1.CreateFrostmourneDepRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidateFrostmourneRules())
}

func (v *Validator) ValidateCreateFrostmourneTeamRequest(ctx context.Context, rq *apiv1.CreateFrostmourneTeamRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidateFrostmourneRules())
}
