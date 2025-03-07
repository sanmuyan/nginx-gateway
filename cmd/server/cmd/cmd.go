package cmd

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"nginx-gateway/pkg/configpost"
)

var rootCtx context.Context

var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "Nginx Gateway Server Admin",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := initConfig(cmd)
		if err != nil {
			logrus.Fatalf("init config error: %v", err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		configpost.PostInit(rootCtx)
	},
}

var configFile string

const (
	logLevel          = 4
	nginxTemplateFile = "config/nginx.conf.tmpl"
	openrestyPath     = "openresty"
	balancerConfigApi = "http://127.0.0.1:9001/configs"
	redisAddr         = "127.0.0.1:6379"
	serverBind        = "127.0.0.1:9000"
)

func init() {
	// 初始化命令行参数
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file")
	rootCmd.PersistentFlags().IntP("log-level", "l", logLevel, "log level")
	rootCmd.PersistentFlags().BoolP("pprof-server", "", false, "enable pprof server")
	rootCmd.Flags().String("nginx-template-file", nginxTemplateFile, "nginx template file")
	rootCmd.Flags().String("openresty-path", openrestyPath, "openresty path")
	rootCmd.Flags().String("balancer-config-api", balancerConfigApi, "balancer config api")
	rootCmd.Flags().String("redis-addr", redisAddr, "redis addr")
	rootCmd.Flags().String("server-bind", serverBind, "server bind addr")
}

func Execute(ctx context.Context) {
	rootCtx = ctx
	if err := rootCmd.Execute(); err != nil {
		logrus.Tracef("cmd execute error: %v", err)
	}
}
