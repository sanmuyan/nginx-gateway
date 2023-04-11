package config

type LimitTraffic struct {
	ReqRate   int `json:"req_rate"`
	ReqBurst  int `json:"req_burst"`
	ConnMax   int `json:"conn_max"`
	ConnBurst int `json:"conn_burst"`
}

type Auth struct {
	Enable   bool   `json:"enable"`
	AuthType string `json:"auth_type"`
	Token    string `json:"token"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Annotations struct {
	Websocket      string            `json:"websocket"`
	ProxyHeader    map[string]string `json:"proxy_header"`
	CustomLocation []string          `json:"custom_location"`
	Whitelist      []string          `json:"whitelist"`
	Blacklist      []string          `json:"blacklist"`
	LimitTraffic   LimitTraffic      `json:"limit_traffic"`
	Auth           Auth              `json:"auth"`
}

type Route struct {
	RouteName   string      `json:"route_name" binding:"required"`
	Description string      `json:"description"`
	RoutePath   string      `json:"route_path" binding:"required"`
	BackendName string      `json:"backend_name" binding:"required"`
	Annotations Annotations `json:"annotations"`
}

type AuthRequest struct {
	Enable    bool   `json:"enable"`
	AuthPath  string `json:"auth_path"  binding:"required"`
	ProxyPass string `json:"proxy_pass"  binding:"required"`
}

type Gateway struct {
	GatewayName string      `json:"gateway_name" binding:"required"`
	Description string      `json:"description"`
	ListenPort  int         `json:"listen_port" binding:"required"`
	Host        string      `json:"host" binding:"required"`
	RouteNames  []string    `json:"route_names"`
	Routes      []Route     `json:"routes,omitempty"`
	AuthRequest AuthRequest `json:"auth_request"`
}

type HealthCheck struct {
	Enable   bool   `json:"enable"`
	Type     string `json:"type"`
	Timeout  int    `json:"timeout"`
	Interval int    `json:"interval"`
	Success  int    `json:"success"`
	Fail     int    `json:"fail"`
	Uri      string `json:"uri"`
}

type Server struct {
	Addr   string `json:"addr"`
	Port   int    `json:"port"`
	Weight int    `json:"weight"`
}

type Backend struct {
	BackendName string      `json:"backend_name" binding:"required"`
	Servers     []Server    `json:"servers" binding:"required"`
	HealthCheck HealthCheck `json:"health_check"`
}

type GatewayConfig struct {
	Backends []Backend `json:"backends,omitempty"`
	Gateways []Gateway `json:"gateways,omitempty"`
}
