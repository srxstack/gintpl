// Copyright 2025 莫维龙 <kalandramo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/srxstack/gintpl. The professional
// version of this repository is https://github.com/srxstack/srxstack.

package options

import (
	"fmt"

	"github.com/spf13/pflag"
	genericoptions "github.com/srxstack/srxstack/pkg/options"
)

var _ genericoptions.IOptions = (*FrostmourneOptions)(nil)

type FrostmourneOptions struct {
	URL      string `json:"url,omitempty" mapstructure:"url"`
	Username string `json:"username,omitempty" mapstructure:"username"`
	Password string `json:"password,omitempty" mapstructure:"password"`
}

func NewFrostmourneOptions() *FrostmourneOptions {
	return &FrostmourneOptions{}
}

func (o *FrostmourneOptions) Validate() []error {
	errs := []error{}

	if o.URL == "" {
		errs = append(errs, fmt.Errorf("url not be empty"))
	}
	if o.Username == "" {
		errs = append(errs, fmt.Errorf("username not be empty"))
	}
	if o.Password == "" {
		errs = append(errs, fmt.Errorf("password not be empty"))
	}

	return errs
}

func (o *FrostmourneOptions) AddFlags(fs *pflag.FlagSet, prefixes ...string) {
	fs.StringVar(&o.URL, "frostmourne.url", o.URL, ""+"URL is the http address of the frostmourne server.")
	fs.StringVar(&o.Username, "frostmourne.username", o.Username, ""+"Username for access to frostmourne server.")
	fs.StringVar(&o.Password, "frostmourne.password", o.Password, ""+"Password for access to frostmourne server.")
}
