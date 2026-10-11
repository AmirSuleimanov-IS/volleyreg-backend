package repository

import (
	"strings"

	"golang.org/x/crypto/bcrypt"

	"volleyreg/internal/app/ds"
)

// CreateUser — регистрация пользователя (пароль хранится в виде bcrypt-хэша)
func (r *Repository) CreateUser(login, password, fullName string) (ds.User, error) {
	login = strings.TrimSpace(login)
	var count int64
	if err := r.db.Model(&ds.User{}).Where("login = ?", login).Count(&count).Error; err != nil {
		return ds.User{}, err
	}
	if count > 0 {
		return ds.User{}, ErrConflict
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ds.User{}, err
	}
	user := ds.User{Login: login, Password: string(hash), FullName: strings.TrimSpace(fullName)}
	return user, r.db.Create(&user).Error
}
