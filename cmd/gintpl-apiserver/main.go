// Copyright 2025 莫维龙 <kalandramo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/srxstack/gintpl. The professional
// version of this repository is https://github.com/srxstack/srxstack.

package main

import (
	"os"

	"github.com/srxstack/gintpl/cmd/gintpl-apiserver/app"

	_ "go.uber.org/automaxprocs"
)

func main() {
	command := app.NewGinTplCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
