package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RunServer(addr string) {
	r := gin.Default()
	router(r)
	err := r.Run(addr)
	if err != nil {
		logrus.Fatal(err)
	}
}

func router(r *gin.Engine) {
	r.POST("/api/backend", Backend)
	r.POST("/api/route", Route)
	r.POST("/api/gateway", Gateway)
	r.GET("/api/reload", Reload)
}
