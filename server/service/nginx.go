package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"nginx-gateway/pkg/config"
	"nginx-gateway/pkg/db"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// nginx 相关处理

func (s *Service) updateConfigAndReload() error {
	err := s.updateNginxBackend()
	if err != nil {
		return err
	}
	err = s.generateNginxConfig()
	if err != nil {
		return err
	}
	err = s.reloadNginx()
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) watchReload() {
	logrus.Infoln("watch reload start")
	var currentReloadTime int64
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		if currentReloadTime == 0 {
			currentReloadTime = time.Now().Unix()
			err := db.RDB.Set(s.ctx, config.ReloadKey, currentReloadTime, 0).Err()
			if err != nil {
				logrus.Errorln(err)
			}
			err = s.updateConfigAndReload()
			if err != nil {
				logrus.Errorf("reload nginx: %v", err)
			}
			continue
		}
		result, err := db.RDB.Get(s.ctx, config.ReloadKey).Result()
		if err != nil {
			logrus.Errorln(err)
		}
		latestReloadTime, err := strconv.ParseInt(result, 10, 64)
		if err != nil {
			logrus.Errorln(err)
			continue
		}

		if latestReloadTime > currentReloadTime {
			currentReloadTime = latestReloadTime
			err = s.updateConfigAndReload()
			if err != nil {
				logrus.Errorf("reload nginx: %v", err)
				continue
			}
			logrus.Infoln("nginx reload success")
		}

	}
}

func (s *Service) reloadNginx() error {
	t := exec.Command("./nginx", "-t")
	t.Dir = config.Conf.OpenrestyPath
	out, err := t.CombinedOutput()
	if err != nil {
		return nil
	}
	if !strings.Contains(string(out), "successful") {
		err = errors.New("nginx config test failed")
		return err
	}
	r := exec.Command("./nginx", "-s", "reload")
	r.Dir = config.Conf.OpenrestyPath
	err = r.Run()
	if err != nil {
		return err
	}
	return err
}

// 从redis 获取gateway 配置，根据go template 文件生成nginx.conf
func (s *Service) generateNginxConfig() error {
	gatewaysRes, err := db.RDB.HGetAll(s.ctx, config.GatewaysKey).Result()
	if err != nil {
		return err
	}
	if len(gatewaysRes) > 0 {
		var gatewayConfig config.GatewayConfig
		var gateways []config.Gateway
		for _, gatewayRes := range gatewaysRes {
			var gateway config.Gateway
			err = json.Unmarshal([]byte(gatewayRes), &gateway)
			if err != nil {
				return err
			}
			for _, routeName := range gateway.RouteNames {
				routeRes, err := db.RDB.HGet(s.ctx, config.RoutesKey, routeName).Result()
				if err != nil {
					return err
				}

				if len(routeRes) > 0 {
					var route config.Route
					err = json.Unmarshal([]byte(routeRes), &route)
					if err != nil {
						return err
					}
					gateway.Routes = append(gateway.Routes, route)
				}
			}
			gateways = append(gateways, gateway)
		}
		gatewayConfig.Gateways = gateways

		tmpl, err := ioutil.ReadFile(config.Conf.NginxTemplateFile)
		if err != nil {
			return err
		}

		p := template.New("config")
		t := template.Must(p.Parse(string(tmpl)))

		nginxConfigFile, err := os.Create(config.Conf.OpenrestyPath + "/conf/nginx.conf")
		if err != nil {
			return err
		}
		defer nginxConfigFile.Close()
		err = t.Execute(nginxConfigFile, gatewayConfig)
		if err != nil {
			return err
		}
	}
	return nil
}

// 更新openresty balancer 动态后端配置
func (s *Service) updateNginxBackend() error {
	backendsRes, err := db.RDB.HGetAll(s.ctx, config.BackendsKey).Result()
	if err != nil {
		return err
	}
	var gatewayConfig config.GatewayConfig
	var backends []config.Backend
	for _, backendRes := range backendsRes {
		var backend config.Backend
		err = json.Unmarshal([]byte(backendRes), &backend)
		if err != nil {
			return err
		}
		backends = append(backends, backend)
	}
	gatewayConfig.Backends = backends

	client := &http.Client{Timeout: time.Second * 3}
	body, err := json.Marshal(gatewayConfig)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", config.Conf.BalancerConfigAPI, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("x-save", "1")
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	go NewService().watchReload()
}
