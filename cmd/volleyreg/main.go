package main

import (
	"github.com/sirupsen/logrus"

	"volleyreg/internal/api"
)

func main() {
	logrus.Info("Application start!")
	api.StartServer()
	logrus.Info("Application terminated!")
}
