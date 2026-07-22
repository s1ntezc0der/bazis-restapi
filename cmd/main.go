package main

import (
	"log"

	"github.com/s1ntezc0der/bazis-restapi/internal/run"
)

// @title Task Management API
// @version 1.0
// @description Сервис управления задачами с командной работой и историей изменений
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите токен в формате: Bearer <token>

func main() {
	if err := run.Run(); err != nil {
		log.Fatal(err)
	}
}
