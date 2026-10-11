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

// CheckPassword — найти пользователя по логину и сверить пароль с bcrypt-хэшем
func (r *Repository) CheckPassword(login, password string) (ds.User, error) {
	var user ds.User
	if err := r.db.Where("login = ?", strings.TrimSpace(login)).First(&user).Error; err != nil {
		return user, notFound(err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return user, ErrForbidden
	}
	return user, nil
}

func (r *Repository) GetUser(id uint) (ds.User, error) {
	var user ds.User
	err := r.db.First(&user, id).Error
	return user, notFound(err)
}

// UpdateUser — изменение ФИО и пароля в личном кабинете
func (r *Repository) UpdateUser(id uint, fullName, password *string) (ds.User, error) {
	updates := map[string]interface{}{}
	if fullName != nil {
		updates["full_name"] = strings.TrimSpace(*fullName)
	}
	if password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			return ds.User{}, err
		}
		updates["password"] = string(hash)
	}
	if len(updates) > 0 {
		if err := r.db.Model(&ds.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return ds.User{}, err
		}
	}
	return r.GetUser(id)
}
