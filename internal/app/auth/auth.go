package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"volleyreg/internal/app/ds"
)

var ErrInvalidToken = errors.New("токен недействителен")

// Claims — содержимое JWT: id пользователя, роль и стандартные поля (exp, iat, jti)
type Claims struct {
	UserID      uint   `json:"user_id"`
	Login       string `json:"login"`
	IsOrganizer bool   `json:"is_organizer"`
	jwt.RegisteredClaims
}

// Auth — выдача и проверка JWT; отозванные при выходе токены хранятся в Redis (чёрный список)
type Auth struct {
	secret []byte
	ttl    time.Duration
	redis  *redis.Client
}

func New() (*Auth, error) {
	minutes, err := strconv.Atoi(os.Getenv("JWT_TTL_MINUTES"))
	if err != nil || minutes <= 0 {
		minutes = 60
	}
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis: %w", err)
	}
	return &Auth{secret: []byte(os.Getenv("JWT_SECRET")), ttl: time.Duration(minutes) * time.Minute, redis: client}, nil
}

// Issue — выдать JWT пользователю после успешного входа
func (a *Auth) Issue(user ds.User) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(a.ttl)
	claims := Claims{
		UserID:      user.ID,
		Login:       user.Login,
		IsOrganizer: user.IsOrganizer,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			Issuer:    "volleyreg",
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(a.secret)
	return token, exp, err
}

// Parse — проверить подпись и срок действия токена и что он не отозван (нет в чёрном списке Redis)
func (a *Auth) Parse(ctx context.Context, tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return a.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	revoked, err := a.redis.Exists(ctx, blacklistKey(claims.ID)).Result()
	if err != nil {
		return nil, err
	}
	if revoked > 0 {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

// Revoke — отозвать токен при выходе: jti попадает в чёрный список до истечения срока токена
func (a *Auth) Revoke(ctx context.Context, claims *Claims) error {
	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		return nil
	}
	value := fmt.Sprintf("user_id=%d login=%s", claims.UserID, claims.Login)
	return a.redis.Set(ctx, blacklistKey(claims.ID), value, ttl).Err()
}

func blacklistKey(jti string) string {
	return "blacklist:" + jti
}
