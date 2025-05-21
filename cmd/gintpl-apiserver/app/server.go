// Copyright 2025 莫维龙 <kalandramo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/srxstack/gintpl. The professional
// version of this repository is https://github.com/srxstack/srxstack.

package app

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/srxstack/srxstack/pkg/core"
	"github.com/srxstack/srxstack/pkg/version"

	"github.com/srxstack/gintpl/cmd/gintpl-apiserver/app/options"
	"github.com/srxstack/gintpl/internal/pkg/log"
)

var configFile string // 配置文件路径

func NewGinTplCommand() *cobra.Command {
	opts := options.NewServerOptions()

	cmd := &cobra.Command{
		Use:   " gintpl-apiserver",
		Short: "A gin template",
		Long: `A gin template.

The project features include:
• Utilization of a clean architecture;
...
• High-quality code.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(opts)
		},
		Args: cobra.NoArgs,
	}

	cobra.OnInitialize(core.OnInitialize(&configFile, "GINTPL", searchDirs(), defaultConfigName))

	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", filePath(), "Path to the gintpl configuration file.")

	opts.AddFlags(cmd.PersistentFlags())

	version.AddFlags(cmd.PersistentFlags())

	return cmd
}

func run(opts *options.ServerOptions) error {
	version.PrintAndExitIfRequested()

	log.Init(logOptions())
	defer log.Sync()

	if err := viper.Unmarshal(opts); err != nil {
		return err
	}

	if err := opts.Validate(); err != nil {
		return err
	}

	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	opts.Configz()

	server, err := cfg.NewUnionServer()
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	return server.Run()
}

// logOptions 从 viper 中读取日志配置，构建 *log.Options 并返回.
// 注意：viper.Get<Type>() 中 key 的名字需要使用 . 分割，以跟 YAML 中保持相同的缩进.
func logOptions() *log.Options {
	opts := log.NewOptions()
	if viper.IsSet("log.disable-caller") {
		opts.DisableCaller = viper.GetBool("log.disable-caller")
	}
	if viper.IsSet("log.disable-stacktrace") {
		opts.DisableStacktrace = viper.GetBool("log.disable-stacktrace")
	}
	if viper.IsSet("log.level") {
		opts.Level = viper.GetString("log.level")
	}
	if viper.IsSet("log.format") {
		opts.Format = viper.GetString("log.format")
	}
	if viper.IsSet("log.output-paths") {
		opts.OutputPaths = viper.GetStringSlice("log.output-paths")
	}
	return opts
}
