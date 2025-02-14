package usecase

import (
	"gate-web/pkg/config"
	"gate-web/pkg/logger"
	"net/http/httputil"
	"net/url"
)

type RouterUsecase struct {
	Config *config.Config
	log    *logger.Logger
}

func NewRouterUsecase(cfg *config.Config, log *logger.Logger) *RouterUsecase {
	return &RouterUsecase{Config: cfg, log: log}
}

func (ru *RouterUsecase) GetReverseProxy(targetURL string) *httputil.ReverseProxy {
	target, err := url.Parse(targetURL)
	if err != nil {
		ru.log.Errorf("Ошибка парсинга URL: %v", err)
		return nil
	}
	return httputil.NewSingleHostReverseProxy(target)
}

func (ru *RouterUsecase) GetTarget(service string) string {
	for _, route := range ru.Config.Routes {
		if route.Path == service {
			return route.Target
		}
	}

	return ""
}
