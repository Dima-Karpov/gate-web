package delivery

import (
	"fmt"
	"gate-web/internal/usecase"
	"gate-web/pkg/logger"
	"gate-web/pkg/middleware"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type HttpHandler struct {
	routerUsecase *usecase.RouterUsecase
	log           *logger.Logger
}

func NewHttpHandler(routerUsecase *usecase.RouterUsecase, log *logger.Logger) http.Handler {
	handler := &HttpHandler{
		routerUsecase: routerUsecase,
		log:           log,
	}
	router := mux.NewRouter()

	// Проксируем запросы вида /{service}/{rest}
	router.HandleFunc(
		"/{service}/{rest:.*}",
		handler.ProxyRequest,
	).Methods(
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodOptions,
	)

	return router
}

// Прокси запроса к целевому сервису
func (h *HttpHandler) ProxyRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"] // Например, "news-api"

	fmt.Print("service: ", service)

	// Получаем базовый URL сервиса
	targetURL := h.routerUsecase.GetTarget(service)
	if targetURL == "" {
		http.Error(w, "Сервис не найден", http.StatusNotFound)
		return
	}

	// Получаем trace_id из контекста
	traceID := middleware.GetTraceID(r.Context())

	// Логируем входящий запрос
	h.log.WithFields(logrus.Fields{
		"trace_id": traceID,
		"service":  service,
		"path":     r.URL.Path,
		"method":   r.Method,
	}).Info("ProxyRequest")

	// Формируем правильный путь (убираем news-api)
	restOfPath := strings.TrimPrefix(r.URL.Path, "/"+service)
	if !strings.HasPrefix(restOfPath, "/") {
		restOfPath = "/" + restOfPath //
	}
	finalURL := targetURL + restOfPath

	// Добавляем query параметры
	if r.URL.RawQuery != "" {
		finalURL += "?" + r.URL.RawQuery
	}

	r.URL.Path = restOfPath

	// Отправляем запрос через reverse proxy
	proxy := h.routerUsecase.GetReverseProxy(targetURL)

	// Передаем trace_id в заголовок
	director := proxy.Director
	proxy.Director = func(req *http.Request) {
		director(req)
		req.Header.Set("X-Forwarded-Host", req.Host)

		traceID := middleware.GetTraceID(r.Context())
		req.Header.Set("X-Request-ID", traceID)
	}

	// Логируем перед отправкой
	h.log.WithFields(logrus.Fields{
		"trace_id": traceID,
		"target":   finalURL,
		"method":   r.Method,
	}).Info("Forwarding request to service")

	// Проксируем запрос
	proxy.ServeHTTP(w, r)
}
