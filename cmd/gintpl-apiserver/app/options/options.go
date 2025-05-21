// Copyright 2025 莫维龙 <kalandramo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/srxstack/gintpl. The professional
// version of this repository is https://github.com/srxstack/srxstack.

package options

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/pflag"
	genericoptions "github.com/srxstack/srxstack/pkg/options"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/srxstack/gintpl/internal/apiserver"
	"github.com/srxstack/gintpl/pkg/configz"

	gintploptions "github.com/srxstack/gintpl/pkg/options"
)

// 定义支持的服务器模式集合。
var availableServerModes = sets.New(
	"release",
	"debug",
	"test",
)

// ServerOptions 包含服务器配置选项。
type ServerOptions struct {
	// ServerMode 定义 gin 服务器模式：Release、Debug、Test。
	ServerMode         string                            `json:"server-mode"         mapstructure:"server-mode"`
	JWTKey             string                            `json:"jwt-key"             mapstructure:"jwt-key"`
	Expiration         time.Duration                     `json:"expiration"          mapstructure:"expiration"`
	EnableMemoryStore  bool                              `json:"enable-memory-store" mapstructure:"enable-memory-store"`
	HTTPOptions        *genericoptions.HTTPOptions       `json:"http"                mapstructure:"http"`
	TLSOptions         *genericoptions.TLSOptions        `json:"tls"                 mapstructure:"tls"`
	MySQLOptions       *genericoptions.MySQLOptions      `json:"mysql"               mapstructure:"mysql"`
	FrostmourneOptions *gintploptions.FrostmourneOptions `json:"frostmourne"         mapstructure:"frostmourne"`
}

func NewServerOptions() *ServerOptions {
	opts := &ServerOptions{
		ServerMode:         "debug",
		JWTKey:             "Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5",
		Expiration:         2 * time.Hour,
		EnableMemoryStore:  true,
		HTTPOptions:        genericoptions.NewHTTPOptions(),
		TLSOptions:         genericoptions.NewTLSOptions(),
		MySQLOptions:       genericoptions.NewMySQLOptions(),
		FrostmourneOptions: gintploptions.NewFrostmourneOptions(),
	}
	opts.HTTPOptions.Addr = ":5555"
	return opts
}

func (o *ServerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ServerMode, "server-mode", o.ServerMode, fmt.Sprintf("Server mode, available options: %v", availableServerModes.UnsortedList()))
	fs.StringVar(&o.JWTKey, "jwt-key", o.JWTKey, "JWT signing key. Must be at least 6 characters long.")
	fs.DurationVar(&o.Expiration, "expiration", o.Expiration, "The expiration duration of JWT tokens.")
	fs.BoolVar(&o.EnableMemoryStore, "enable-memory-store", o.EnableMemoryStore, "Enable in-memory database (useful for testing or development).")
	o.HTTPOptions.AddFlags(fs)
	o.TLSOptions.AddFlags(fs)
	o.MySQLOptions.AddFlags(fs)
	o.FrostmourneOptions.AddFlags(fs)
}

func (o *ServerOptions) Validate() error {
	errs := []error{}

	if !availableServerModes.Has(o.ServerMode) {
		errs = append(errs, fmt.Errorf("invalid server mode: must be one of %v", availableServerModes.UnsortedList()))
	}

	if len(o.JWTKey) < 6 {
		errs = append(errs, errors.New("JWTKey must be at least 6 characters long"))
	}

	errs = append(errs, o.HTTPOptions.Validate()...)
	errs = append(errs, o.TLSOptions.Validate()...)
	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.FrostmourneOptions.Validate()...)

	return utilerrors.NewAggregate(errs)
}

// 应用配置.
func (o *ServerOptions) Config() (*apiserver.Config, error) {
	return &apiserver.Config{
		ServerMode:        o.ServerMode,
		JWTKey:            o.JWTKey,
		Expiration:        o.Expiration,
		EnableMemoryStore: o.EnableMemoryStore,
		HTTPOptions:       o.HTTPOptions,
		TLSOptions:        o.TLSOptions,
		MySQLOptions:      o.MySQLOptions,
	}, nil
}

// 全局配置.
func (o *ServerOptions) Configz() {
	z := configz.NewConfigz()
	z.FrostmourneOptions = *o.FrostmourneOptions
}
