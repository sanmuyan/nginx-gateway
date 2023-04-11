package config

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path"
	"runtime"
)

const (
	GatewaysKey = "nginx_gateway:gateways"
	RoutesKey   = "nginx_gateway:routes"
	BackendsKey = "nginx_gateway:backends"
	ReloadKey   = "nginx_gateway:reload"
)

const defaultConfig = `
log_level: 4
nginx_template_file: config/nginx.conf.tmpl
openresty_path: openresty
balancer_config_api: http://127.0.0.1:9001/configs
redis_addr: 127.0.0.1:6379
band_addr: 127.0.0.1:9000
`

type Config struct {
	LogLevel          int    `yaml:"log_level"`
	NginxTemplateFile string `yaml:"nginx_template_file"`
	OpenrestyPath     string `yaml:"openresty_path"`
	BalancerConfigAPI string `yaml:"balancer_config_api"`
	RedisAddr         string `yaml:"redis_addr"`
	BandAddr          string `yaml:"band_addr"`
}

var Conf *Config

func NewConfig(configFile string) *Config {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File)
			return frame.Function, fileName
		},
	})

	configByte, err := ioutil.ReadFile(configFile)
	if err != nil {
		logrus.Fatalln("config", err)
	}
	config := &Config{}
	err = yaml.Unmarshal([]byte(defaultConfig), config)
	if err != nil {
		logrus.Fatalln("config", err)
	}

	if err := yaml.Unmarshal(configByte, config); err != nil {
		logrus.Fatalln("config", err)
	}

	logrus.SetLevel(logrus.Level(config.LogLevel))
	if logrus.Level(config.LogLevel) >= logrus.DebugLevel {
		logrus.SetReportCaller(true)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	logrus.Debugf("config %+v", config)
	Conf = config
	return config
}
