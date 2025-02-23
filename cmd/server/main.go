package main

import (
	"gate-web/internal/delivery"
	"gate-web/internal/usecase"
	"gate-web/pkg/config"
	"gate-web/pkg/logger"
	"gate-web/pkg/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		DisableColors:   false,
		ForceQuote:      true,
		PadLevelText:    true,
	})
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig("config/routes.yaml")
	if err != nil {
		logrus.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Настраиваем логирования
	l := logger.NewLogger()

	// Создаем сервер маршрутизации
	routerUsecase := usecase.NewRouterUsecase(cfg, l)

	// Запускаем HTTP-сервер
	handler := delivery.NewHttpHandler(routerUsecase, l)

	// Оборачиваем handler в TraceMiddleware
	wrappedHandler := middleware.TraceMeddleWare(handler)

	logrus.Print("API Gatewab запущен на порту 8081")

	err = http.ListenAndServe(":8081", wrappedHandler)
	if err != nil {
		logrus.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
