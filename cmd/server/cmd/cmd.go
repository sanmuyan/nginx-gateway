package cmd

import (
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
	serverBandAddr    = "127.0.0.1:9000"
)

func init() {
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "config file")
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

	viper.SetConfigName("config")
	viper.SetDefault("log_level", logLevel)
	viper.SetDefault("nginx_template_file", nginxTemplateFile)
	viper.SetDefault("openresty_path", openrestyPath)
	viper.SetDefault("balancer_config_api", balancerConfigApi)
	viper.SetDefault("redis_addr ", redisAddr)
	viper.SetDefault("server_band_addr", serverBandAddr)

	if len(configFile) > 0 {
		viper.SetConfigFile(configFile)
		err := viper.ReadInConfig()
		if err != nil {
			return err
		}
	} else {
		_ = viper.BindPFlag("log_level", rootCmd.Flags().Lookup("log-level"))
	}

	err := viper.Unmarshal(&config.Conf)
	if err != nil {
		return err
	}
	logrus.SetLevel(logrus.Level(config.Conf.LogLevel))
	if logrus.Level(config.Conf.LogLevel) >= logrus.DebugLevel {
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
		controller.RunServer(config.Conf.ServerBandAddr)
	}
}
