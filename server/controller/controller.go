package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"nginx-gateway/pkg/config"
	"nginx-gateway/pkg/response"
	"nginx-gateway/server/service"
)

// API列表

var svc = service.NewService()

var respf = func() *response.Response {
	return response.NewResponse()
}

func Backend(c *gin.Context) {
	backend := &config.Backend{}
	err := c.ShouldBindJSON(backend)
	if err != nil {
		logrus.Println(err)
		respf().Fail(response.HttpBadRequest).SetGin(c)
		return
	}

	if err = svc.UpdateBackend(backend); err != nil {
		logrus.Println(err)
		respf().Fail(response.HttpInternalServerError).SetGin(c)
		return
	}

	respf().Ok().SetGin(c)
}

func Gateway(c *gin.Context) {
	gateway := &config.Gateway{}
	err := c.ShouldBindJSON(gateway)
	if err != nil {
		logrus.Println(err)
		respf().Fail(response.HttpBadRequest).SetGin(c)
		return
	}
	if err = svc.UpdateGateway(gateway); err != nil {
		logrus.Println(err)
		respf().Fail(response.HttpInternalServerError).SetGin(c)
		return
	}

	respf().Ok().SetGin(c)
}

func Route(c *gin.Context) {
	route := &config.Route{}
	err := c.ShouldBindJSON(route)
	if err != nil {
		logrus.Println(err)
		respf().Fail(response.HttpBadRequest).SetGin(c)
		return
	}
	if err = svc.UpdateRoute(route); err != nil {
		logrus.Println(err)
		respf().Fail(response.HttpInternalServerError).SetGin(c)
		return
	}
	respf().Ok().SetGin(c)
}

func Reload(c *gin.Context) {
	err := svc.Reload()
	if err != nil {
		logrus.Println(err)
		respf().Fail(response.HttpInternalServerError).SetGin(c)
		return
	}
	respf().Ok().SetGin(c)
}
