package main

import (
	"gate-web/internal/delivery"
	"gate-web/internal/usecase"
	"gate-web/pkg/config"
	"gate-web/pkg/logger"
	"log"
	"net/http"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig("config/routes.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Настраиваем логирования
	l := logger.NewLogger()

	// Создаем сервер маршрутизации
	routerUsecase := usecase.NewRouterUsecase(cfg, l)

	// Запускаем HTTP-сервер
	handler := delivery.NewHttpHandler(routerUsecase, l)
	l.Info("API Gatewab запущен на порту 8081")

	err = http.ListenAndServe(":8081", handler)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
