package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"volleyreg/internal/api"
	"volleyreg/internal/app/auth"
	"volleyreg/internal/app/dsn"
	"volleyreg/internal/app/repository"
	"volleyreg/internal/app/storage"
)

// @title           VolleyReg API
// @version         1.0
// @description     Регистрация участников на спортивное соревнование (Первенство ЗГУ по волейболу 2026): участники, заявки команд, состав заявки, пользователи.
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 JWT из ответа /api/users/login в виде: Bearer <токен>
func main() {
	logrus.Info("Application start!")
	_ = godotenv.Load()

	repo, err := repository.New(dsn.FromEnv())
	if err != nil {
		logrus.Fatal("cant connect to db: ", err)
	}

	store, err := storage.New()
	if err != nil {
		logrus.Fatal("cant connect to minio: ", err)
	}

	a, err := auth.New()
	if err != nil {
		logrus.Fatal("cant connect to redis: ", err)
	}

	api.StartServer(repo, store, a)
	logrus.Info("Application terminated!")
}
