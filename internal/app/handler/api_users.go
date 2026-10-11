package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"volleyreg/internal/app/auth"
	"volleyreg/internal/app/serializer"
)

// ApiRegister godoc
// @Summary      Регистрация пользователя
// @Description  Создаёт представителя команды. Пароль хранится в виде bcrypt-хэша.
// @Tags         Пользователи
// @Accept       json
// @Produce      json
// @Param        body  body      serializer.RegisterRequest  true  "логин, пароль, ФИО"
// @Success      201   {object}  serializer.UserResponse
// @Failure      400   {object}  serializer.ErrorResponse
// @Failure      409   {object}  serializer.ErrorResponse  "логин занят"
// @Router       /api/users/register [post]
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

// ApiLogin godoc
// @Summary      Аутентификация
// @Description  Проверяет логин и пароль и выдаёт JWT. Токен передаётся в заголовке Authorization: Bearer <токен>.
// @Tags         Пользователи
// @Accept       json
// @Produce      json
// @Param        body  body      serializer.LoginRequest  true  "логин и пароль"
// @Success      200   {object}  serializer.LoginResponse
// @Failure      401   {object}  serializer.ErrorResponse
// @Router       /api/users/login [post]
func (h *Handler) ApiLogin(ctx *gin.Context) {
	var req serializer.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequest(ctx, "login и password обязательны")
		return
	}
	user, err := h.Repository.CheckPassword(req.Login, req.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, serializer.ErrorResponse{Error: "неверный логин или пароль"})
		return
	}
	token, exp, err := h.Auth.Issue(user)
	if err != nil {
		apiError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, serializer.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresAt:   exp.Format(time.RFC3339),
		User:        serializer.User(user),
	})
}

// ApiLogout godoc
// @Summary      Деавторизация
// @Description  Отзывает текущий JWT: его идентификатор (jti) помещается в чёрный список Redis до конца срока действия.
// @Tags         Пользователи
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  serializer.MessageResponse
// @Failure      401  {object}  serializer.ErrorResponse
// @Router       /api/users/logout [post]
func (h *Handler) ApiLogout(ctx *gin.Context) {
	claims := ctx.MustGet(ctxClaims).(*auth.Claims)
	if err := h.Auth.Revoke(ctx, claims); err != nil {
		apiError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, serializer.MessageResponse{Message: "выход выполнен"})
}

// ApiGetMe godoc
// @Summary      Профиль текущего пользователя
// @Tags         Пользователи
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  serializer.UserResponse
// @Failure      401  {object}  serializer.ErrorResponse
// @Router       /api/users/me [get]
func (h *Handler) ApiGetMe(ctx *gin.Context) {
	userID, _ := currentUser(ctx)
	user, err := h.Repository.GetUser(userID)
	if err != nil {
		apiError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, serializer.User(user))
}

// ApiUpdateMe godoc
// @Summary      Изменение профиля
// @Description  Меняет ФИО и/или пароль текущего пользователя.
// @Tags         Пользователи
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      serializer.UpdateMeRequest  true  "новые ФИО и пароль"
// @Success      200   {object}  serializer.UserResponse
// @Failure      400   {object}  serializer.ErrorResponse
// @Failure      401   {object}  serializer.ErrorResponse
// @Router       /api/users/me [put]
func (h *Handler) ApiUpdateMe(ctx *gin.Context) {
	var req serializer.UpdateMeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequest(ctx, "некорректный JSON")
		return
	}
	if req.Password != nil && len(*req.Password) < 6 {
		badRequest(ctx, "пароль должен быть не короче 6 символов")
		return
	}
	userID, _ := currentUser(ctx)
	user, err := h.Repository.UpdateUser(userID, req.FullName, req.Password)
	if err != nil {
		apiError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, serializer.User(user))
}
