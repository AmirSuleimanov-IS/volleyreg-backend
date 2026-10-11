package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"volleyreg/internal/api"
	"volleyreg/internal/app/dsn"
	"volleyreg/internal/app/repository"
)

func main() {
	logrus.Info("Application start!")
	_ = godotenv.Load()

	repo, err := repository.New(dsn.FromEnv())
	if err != nil {
		logrus.Fatal("cant connect to db: ", err)
	}

	api.StartServer(repo)
	logrus.Info("Application terminated!")
}
