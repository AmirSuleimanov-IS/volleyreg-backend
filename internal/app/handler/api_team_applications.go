package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"volleyreg/internal/app/repository"
	"volleyreg/internal/app/serializer"
	"volleyreg/internal/app/singleton"
)

// ApiGetDraft — GET /api/team-applications/draft: иконка корзины (id черновика и количество участников)
func (h *Handler) ApiGetDraft(ctx *gin.Context) {
	app, count, err := h.Repository.GetDraftTeamApplication(singleton.GetCreatorID())
	if errors.Is(err, repository.ErrNotFound) {
		ctx.JSON(http.StatusOK, serializer.DraftResponse{TeamApplicationID: nil, ParticipantsCount: 0})
		return
	}
	if err != nil {
		apiError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, serializer.DraftResponse{TeamApplicationID: &app.ID, ParticipantsCount: count})
}

// ApiGetTeamApplications — GET /api/team-applications?status=&formed_from=&formed_to=:
// список заявок без удалённых и черновиков
func (h *Handler) ApiGetTeamApplications(ctx *gin.Context) {
	apps, err := h.Repository.GetTeamApplications(repository.TeamApplicationFilter{
		Status:     ctx.Query("status"),
		FormedFrom: ctx.Query("formed_from"),
		FormedTo:   ctx.Query("formed_to"),
	})
	if err != nil {
		apiError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, serializer.TeamApplications(apps))
}

// ApiGetTeamApplication — GET /api/team-applications/:id: заявка и её состав с картинками
func (h *Handler) ApiGetTeamApplication(ctx *gin.Context) {
	id, ok := uintParam(ctx, "id")
	if !ok {
		return
	}
	app, items, err := h.Repository.GetTeamApplication(id)
	if err != nil {
		apiError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, serializer.TeamApplicationDetail(app, items))
}

// ApiUpdateTeamApplication — PUT /api/team-applications/:id: изменение полей заявки по теме
func (h *Handler) ApiUpdateTeamApplication(ctx *gin.Context) {
	id, ok := uintParam(ctx, "id")
	if !ok {
		return
	}
	var req serializer.TeamApplicationUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequest(ctx, "некорректный JSON")
		return
	}
	err := h.Repository.UpdateTeamApplication(id, singleton.GetCreatorID(), req.TeamName, req.Institution, req.CoachName, req.CompetitionDate)
	if err != nil {
		apiError(ctx, err)
		return
	}
	h.respondTeamApplication(ctx, id)
}

// ApiFormTeamApplication — PUT /api/team-applications/:id/form: сформировать заявку (представитель).
// Проверяются обязательные поля и рассчитывается средний возраст команды.
func (h *Handler) ApiFormTeamApplication(ctx *gin.Context) {
	id, ok := uintParam(ctx, "id")
	if !ok {
		return
	}
	if err := h.Repository.FormTeamApplication(id, singleton.GetCreatorID()); err != nil {
		apiError(ctx, err)
		return
	}
	h.respondTeamApplication(ctx, id)
}

// ApiModerateTeamApplication — PUT /api/team-applications/:id/moderate: допустить или отклонить (организатор)
func (h *Handler) ApiModerateTeamApplication(ctx *gin.Context) {
	id, ok := uintParam(ctx, "id")
	if !ok {
		return
	}
	var req serializer.ModerateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequest(ctx, "action должен быть complete или reject")
		return
	}
	if err := h.Repository.ModerateTeamApplication(id, singleton.GetModeratorID(), req.Action); err != nil {
		apiError(ctx, err)
		return
	}
	h.respondTeamApplication(ctx, id)
}

// ApiDeleteTeamApplication — DELETE /api/team-applications/:id: удалить черновик (представитель)
func (h *Handler) ApiDeleteTeamApplication(ctx *gin.Context) {
	id, ok := uintParam(ctx, "id")
	if !ok {
		return
	}
	if err := h.Repository.DeleteTeamApplication(id, singleton.GetCreatorID()); err != nil {
		apiError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *Handler) respondTeamApplication(ctx *gin.Context, id uint) {
	app, items, err := h.Repository.GetTeamApplication(id)
	if err != nil {
		apiError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, serializer.TeamApplicationDetail(app, items))
}
