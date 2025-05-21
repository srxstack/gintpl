// Copyright 2025 莫维龙 <kalandramo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/srxstack/gintpl. The professional
// version of this repository is https://github.com/srxstack/srxstack.

package configz

import (
	"sync"

	gintploptions "github.com/srxstack/gintpl/pkg/options"
)

var (
	once sync.Once
	// 全局变量，方便其它包直接调用已初始化好的 Configz 实例.
	C *configz
)

// ConfigzOptions 包含全局配置选项。不再需要层层传递配置。
type configz struct {
	FrostmourneOptions gintploptions.FrostmourneOptions `json:"frostmourne" mapstructure:"frostmourne"`
}

// NewStore 创建一个 IStore 类型的实例.
func NewConfigz() *configz {
	// 确保 S 只被初始化一次
	once.Do(func() {
		C = &configz{}
	})

	return C
}
