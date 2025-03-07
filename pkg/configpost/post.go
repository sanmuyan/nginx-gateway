package configpost

import (
	"context"
	"nginx-gateway/pkg/config"
	"nginx-gateway/pkg/db"
	"nginx-gateway/server/controller"
	"nginx-gateway/server/service"
)

func PostInit(ctx context.Context) {
	db.InitRedis()
	service.NewService().WatchReload()
	controller.RunServer(config.Conf.ServerBind)
}
