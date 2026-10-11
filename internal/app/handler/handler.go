package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"volleyreg/internal/app/auth"
	"volleyreg/internal/app/repository"
	"volleyreg/internal/app/singleton"
	"volleyreg/internal/app/storage"
)

type Handler struct {
	Repository *repository.Repository
	Storage    *storage.Storage
	Auth       *auth.Auth
}

func NewHandler(r *repository.Repository, s *storage.Storage, a *auth.Auth) *Handler {
	return &Handler{Repository: r, Storage: s, Auth: a}
}

func (h *Handler) errorPage(ctx *gin.Context, status int, message string) {
	ctx.HTML(status, "error.html", gin.H{"status": status, "message": message})
}

func idParam(ctx *gin.Context) (uint, bool) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		return 0, false
	}
	return uint(id), true
}

// GetParticipants — страница списка участников с поиском по ФИО
func (h *Handler) GetParticipants(ctx *gin.Context) {
	participantName := ctx.Query("participant_name")
	participants, err := h.Repository.GetParticipants(participantName, "")
	if err != nil {
		logrus.Error(err)
		h.errorPage(ctx, http.StatusInternalServerError, "Не удалось загрузить участников")
		return
	}

	// карточка текущей заявки: активна, только если у представителя есть черновик
	var draftID uint
	var draftCount int64
	if draft, count, err := h.Repository.GetDraftTeamApplication(singleton.GetCreatorID()); err == nil {
		draftID = draft.ID
		draftCount = count
	}

	ctx.HTML(http.StatusOK, "participants.html", gin.H{
		"participants":    participants,
		"participantName": participantName,
		"draftID":         draftID,
		"draftCount":      draftCount,
	})
}

// GetParticipant — страница одного участника
func (h *Handler) GetParticipant(ctx *gin.Context) {
	id, ok := idParam(ctx)
	if !ok {
		h.errorPage(ctx, http.StatusBadRequest, "Некорректный id участника")
		return
	}
	participant, err := h.Repository.GetParticipant(id)
	if err != nil {
		h.errorPage(ctx, http.StatusNotFound, "Участник не найден")
		return
	}
	ctx.HTML(http.StatusOK, "participant.html", gin.H{"participant": participant})
}

// AddParticipantToDraft — POST: добавить участника в заявку-черновик
func (h *Handler) AddParticipantToDraft(ctx *gin.Context) {
	id, ok := idParam(ctx)
	if !ok {
		h.errorPage(ctx, http.StatusBadRequest, "Некорректный id участника")
		return
	}
	if _, err := h.Repository.AddParticipantToDraft(singleton.GetCreatorID(), id); err != nil {
		logrus.Error(err)
		h.errorPage(ctx, http.StatusNotFound, "Не удалось добавить участника в заявку")
		return
	}
	ctx.Redirect(http.StatusFound, "/participants")
}

// GetTeamApplication — страница заявки команды с её составом
func (h *Handler) GetTeamApplication(ctx *gin.Context) {
	id, ok := idParam(ctx)
	if !ok {
		h.errorPage(ctx, http.StatusBadRequest, "Некорректный id заявки")
		return
	}
	app, items, err := h.Repository.GetTeamApplication(id)
	if err != nil {
		h.errorPage(ctx, http.StatusNotFound, "Заявка не найдена или удалена")
		return
	}
	ctx.HTML(http.StatusOK, "team_application.html", gin.H{
		"teamApplication": app,
		"items":           items,
	})
}

// DeleteTeamApplication — POST: логическое удаление заявки (SQL UPDATE)
func (h *Handler) DeleteTeamApplication(ctx *gin.Context) {
	id, ok := idParam(ctx)
	if !ok {
		h.errorPage(ctx, http.StatusBadRequest, "Некорректный id заявки")
		return
	}
	if err := h.Repository.DeleteTeamApplication(id, singleton.GetCreatorID()); err != nil {
		logrus.Error(err)
		h.errorPage(ctx, http.StatusInternalServerError, "Не удалось удалить заявку")
		return
	}
	ctx.Redirect(http.StatusFound, "/participants")
}
