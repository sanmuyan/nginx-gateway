package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"nginx-gateway/pkg/config"
	"nginx-gateway/pkg/db"
	"nginx-gateway/server/controller"
	"path"
	"runtime"
)

var cmdReady bool

var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "Nginx Gateway Server Admin",
	Run: func(cmd *cobra.Command, args []string) {
		cmdReady = true
	},
	Example: "admin -c config.yaml",
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
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "config file")
	rootCmd.Flags().IntP("log-level", "l", logLevel, "log level")
	rootCmd.Flags().String("nginx-template-file", nginxTemplateFile, "nginx template file")
	rootCmd.Flags().String("openresty-path", openrestyPath, "openresty path")
	rootCmd.Flags().String("balancer-config-api", balancerConfigApi, "balancer config api")
	rootCmd.Flags().String("redis-addr", redisAddr, "redis addr")
	rootCmd.Flags().String("server-bind", serverBind, "server bind addr")
}

func initConfig() error {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File)
			return frame.Function, fileName
		},
	})

	if len(configFile) > 0 {
		viper.SetConfigFile(configFile)
		err := viper.ReadInConfig()
		if err != nil {
			return err
		}
	}
	_ = viper.BindPFlag("log_level", rootCmd.Flags().Lookup("log-level"))
	_ = viper.BindPFlag("nginx_template_file", rootCmd.Flags().Lookup("nginx-template-file"))
	_ = viper.BindPFlag("openresty_path", rootCmd.Flags().Lookup("openresty-path"))
	_ = viper.BindPFlag("balancer_config_api", rootCmd.Flags().Lookup("balancer-config-api"))
	_ = viper.BindPFlag("redis_addr", rootCmd.Flags().Lookup("redis-addr"))
	_ = viper.BindPFlag("server_bind", rootCmd.Flags().Lookup("server-bind"))

	err := viper.Unmarshal(&config.Conf)
	if err != nil {
		return err
	}
	logrus.SetLevel(logrus.Level(config.Conf.LogLevel))
	gin.SetMode(gin.ReleaseMode)
	if logrus.Level(config.Conf.LogLevel) >= logrus.DebugLevel {
		gin.SetMode(gin.DebugMode)
		logrus.SetReportCaller(true)
	}
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
	if cmdReady {
		err := initConfig()
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Debugf("config %+v", config.Conf)

		db.InitRedis()
		controller.RunServer(config.Conf.ServerBind)
	}
}
