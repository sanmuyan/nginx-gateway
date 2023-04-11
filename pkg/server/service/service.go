package service

import (
	"context"
	"encoding/json"
	"nginx-gateway/pkg/config"
	"nginx-gateway/pkg/db"
	"time"
)

// 接口逻辑

type Service struct {
	ctx context.Context
}

func NewService() *Service {
	return &Service{
		ctx: context.Background(),
	}
}

func (s *Service) UpdateBackend(backend *config.Backend) error {
	backendJson, err := json.Marshal(backend)
	if err != nil {
		return err
	}
	err = db.RDB.HSet(s.ctx, config.BackendsKey, backend.BackendName, backendJson).Err()
	if err != nil {
		return err
	}
	err = s.updateNginxBackend()
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateGateway(gateway *config.Gateway) error {
	gatewayJson, err := json.Marshal(gateway)
	if err != nil {
		return err
	}
	err = db.RDB.HSet(s.ctx, config.GatewaysKey, gateway.GatewayName, gatewayJson).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateRoute(route *config.Route) error {
	routeJson, err := json.Marshal(route)
	if err != nil {
		return err
	}
	err = db.RDB.HSet(s.ctx, config.RoutesKey, route.RouteName, routeJson).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Reload() error {
	err := db.RDB.Set(s.ctx, config.ReloadKey, time.Now().Unix(), 0).Err()
	if err != nil {
		return err
	}
	return nil
}
