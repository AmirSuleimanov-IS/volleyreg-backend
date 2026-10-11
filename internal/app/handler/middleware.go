package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"volleyreg/internal/app/serializer"
)

// ключи контекста запроса с данными пользователя из JWT
const (
	ctxUserID      = "user_id"
	ctxIsOrganizer = "is_organizer"
	ctxClaims      = "claims"
)

// RequireAuth — пропускает только запросы с действующим JWT в заголовке Authorization: Bearer <токен>
func (h *Handler) RequireAuth(ctx *gin.Context) {
	header := ctx.GetHeader("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
	if header == "" || token == "" || token == header {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, serializer.ErrorResponse{Error: "требуется авторизация"})
		return
	}
	claims, err := h.Auth.Parse(ctx, token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, serializer.ErrorResponse{Error: "токен недействителен или отозван"})
		return
	}
	ctx.Set(ctxUserID, claims.UserID)
	ctx.Set(ctxIsOrganizer, claims.IsOrganizer)
	ctx.Set(ctxClaims, claims)
	ctx.Next()
}

// RequireOrganizer — пропускает только организатора соревнования (модератора); ставится после RequireAuth
func (h *Handler) RequireOrganizer(ctx *gin.Context) {
	if !ctx.GetBool(ctxIsOrganizer) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, serializer.ErrorResponse{Error: "доступно только организатору соревнования"})
		return
	}
	ctx.Next()
}

// currentUser — id и роль пользователя из токена
func currentUser(ctx *gin.Context) (uint, bool) {
	return ctx.GetUint(ctxUserID), ctx.GetBool(ctxIsOrganizer)
}
