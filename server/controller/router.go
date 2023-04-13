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
	apiGroup := r.Group("/api")
	{
		apiGroup.POST("/backend", Backend)
		apiGroup.POST("/route", Route)
		apiGroup.POST("/gateway", Gateway)
		apiGroup.GET("/reload", Reload)
	}
}
