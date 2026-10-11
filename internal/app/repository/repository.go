package repository

import (
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Ошибки бизнес-логики (обработчики превращают их в HTTP-коды)
var (
	ErrNotFound   = errors.New("не найдено")
	ErrForbidden  = errors.New("недостаточно прав")
	ErrStatus     = errors.New("недопустимая смена статуса")
	ErrValidation = errors.New("ошибка в данных")
	ErrConflict   = errors.New("запись уже существует")
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

func notFound(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return err
}
