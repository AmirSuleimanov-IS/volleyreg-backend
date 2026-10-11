package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"volleyreg/internal/app/serializer"
)

// ApiRegister — POST /api/users/register: регистрация пользователя
func (h *Handler) ApiRegister(ctx *gin.Context) {
	var req serializer.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequest(ctx, "login (3-50 символов) и password (от 6 символов) обязательны")
		return
	}
	user, err := h.Repository.CreateUser(req.Login, req.Password, req.FullName)
	if err != nil {
		apiError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, serializer.User(user))
}

// ApiLogin — POST /api/users/login: аутентификация (заглушка до лабораторной 4)
func (h *Handler) ApiLogin(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "аутентификация будет реализована в лабораторной работе 4"})
}

// ApiLogout — POST /api/users/logout: деавторизация (заглушка до лабораторной 4)
func (h *Handler) ApiLogout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "деавторизация будет реализована в лабораторной работе 4"})
}
