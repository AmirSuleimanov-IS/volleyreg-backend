package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"volleyreg/internal/app/repository"
	"volleyreg/internal/app/serializer"
)

// apiError превращает ошибку бизнес-логики в HTTP-ответ с JSON
func apiError(ctx *gin.Context, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, repository.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, repository.ErrForbidden):
		status = http.StatusForbidden
	case errors.Is(err, repository.ErrStatus), errors.Is(err, repository.ErrConflict):
		status = http.StatusConflict
	case errors.Is(err, repository.ErrValidation):
		status = http.StatusBadRequest
	default:
		logrus.Error(err)
	}
	ctx.AbortWithStatusJSON(status, serializer.ErrorResponse{Error: err.Error()})
}

func badRequest(ctx *gin.Context, message string) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, serializer.ErrorResponse{Error: message})
}

// uintParam — целочисленный параметр пути
func uintParam(ctx *gin.Context, name string) (uint, bool) {
	id, err := strconv.Atoi(ctx.Param(name))
	if err != nil || id <= 0 {
		badRequest(ctx, "некорректный параметр "+name)
		return 0, false
	}
	return uint(id), true
}
