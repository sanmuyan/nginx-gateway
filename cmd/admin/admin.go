package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"nginx-gateway/pkg/config"
	"nginx-gateway/pkg/db"
	"nginx-gateway/pkg/server/controller"
)

func main() {
	// 初始化配置
	configFile := flag.String("c", "config.yaml", "config file")
	flag.Parse()
	conf := config.NewConfig(*configFile)
	db.InitRedis()

	// 启动gin服务器
	r := gin.Default()
	r.POST("/api/backend", controller.Backend)
	r.POST("/api/route", controller.Route)
	r.POST("/api/gateway", controller.Gateway)
	r.GET("/api/reload", controller.Reload)
	err := r.Run(conf.BandAddr)
	if err != nil {
		log.Fatal(err)
	}
}
