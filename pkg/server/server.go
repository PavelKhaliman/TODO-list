package server

import (
	"todo/pkg/api"

	"github.com/gofiber/fiber/v2"
)

// Создание и настройка сервера
func NewServer() *fiber.App {
	app := fiber.New()

	// Регистрация маршрутов API
	api.RegisterRoutes(app)

	return app
}
