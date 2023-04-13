package config

const (
	GatewaysKey = "nginx_gateway:gateways"
	RoutesKey   = "nginx_gateway:routes"
	BackendsKey = "nginx_gateway:backends"
	ReloadKey   = "nginx_gateway:reload"
)

type Config struct {
	LogLevel          int    `mapstructure:"log_level"`
	NginxTemplateFile string `mapstructure:"nginx_template_file"`
	OpenrestyPath     string `mapstructure:"openresty_path"`
	BalancerConfigAPI string `mapstructure:"balancer_config_api"`
	RedisAddr         string `mapstructure:"redis_addr"`
	ServerBind        string `mapstructure:"server_bind"`
}

var Conf Config
