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

	// у каждого представителя не более одного черновика: частичный уникальный индекс
	err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_one_draft_per_creator ON team_applications (creator_id) WHERE status = 'draft'").Error
	if err != nil {
		logrus.Fatal("cant create index: ", err)
	}
	logrus.Info("migration done")
}
