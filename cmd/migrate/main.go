package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"volleyreg/internal/app/ds"
	"volleyreg/internal/app/dsn"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		logrus.Fatal("cant connect to db: ", err)
	}

	err = db.AutoMigrate(
		&ds.User{},
		&ds.Participant{},
		&ds.TeamApplication{},
		&ds.TeamApplicationParticipant{},
	)
	if err != nil {
		logrus.Fatal("cant migrate db: ", err)
	}
	logrus.Info("migration done")
}
