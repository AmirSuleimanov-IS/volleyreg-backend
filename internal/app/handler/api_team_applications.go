package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"volleyreg/internal/app/repository"
	"volleyreg/internal/app/serializer"
)

// ApiGetDraft godoc
// @Summary      Иконка корзины
// @Description  id заявки-черновика текущего представителя и количество участников в ней (без входных параметров).
// @Tags         Заявки команд
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  serializer.DraftResponse
// @Failure      401  {object}  serializer.ErrorResponse
// @Router       /api/team-applications/draft [get]
func (h *Handler) ApiGetDraft(ctx *gin.Context) {
	userID, _ := currentUser(ctx)
	app, count, err := h.Repository.GetDraftTeamApplication(userID)
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

// ApiGetTeamApplications godoc
// @Summary      Список заявок
// @Description  Без удалённых и черновиков. Представитель видит только свои заявки, организатор — все.
// @Tags         Заявки команд
// @Produce      json
// @Security     BearerAuth
// @Param        status       query     string  false  "статус: formed, completed, rejected"
// @Param        formed_from  query     string  false  "дата формирования с (YYYY-MM-DD)"
// @Param        formed_to    query     string  false  "дата формирования по (YYYY-MM-DD)"
// @Success      200          {array}   serializer.TeamApplicationResponse
// @Failure      401          {object}  serializer.ErrorResponse
// @Router       /api/team-applications [get]
func (h *Handler) ApiGetTeamApplications(ctx *gin.Context) {
	userID, isOrganizer := currentUser(ctx)
	filter := repository.TeamApplicationFilter{
		Status:     ctx.Query("status"),
		FormedFrom: ctx.Query("formed_from"),
		FormedTo:   ctx.Query("formed_to"),
	}
	if !isOrganizer {
		filter.CreatorID = userID
	}
	apps, err := h.Repository.GetTeamApplications(filter)
	if err != nil {
		apiError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, serializer.TeamApplications(apps))
}

// ApiGetTeamApplication godoc
// @Summary      Одна заявка
// @Description  Поля заявки и её состав с картинками участников. Представителю доступны только свои заявки.
// @Tags         Заявки команд
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "id заявки"
// @Success      200  {object}  serializer.TeamApplicationDetailResponse
// @Failure      403  {object}  serializer.ErrorResponse
// @Failure      404  {object}  serializer.ErrorResponse
// @Router       /api/team-applications/{id} [get]
func (h *Handler) ApiGetTeamApplication(ctx *gin.Context) {
	id, ok := uintParam(ctx, "id")
	if !ok {
		return
	}
	h.respondTeamApplication(ctx, id)
}

// ApiUpdateTeamApplication godoc
// @Summary      Изменить поля заявки
// @Description  Поля по теме (команда, учебное заведение, тренер, дата соревнования). Только черновик своего представителя.
// @Tags         Заявки команд
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int                                      true  "id заявки"
// @Param        body  body      serializer.TeamApplicationUpdateRequest  true  "поля заявки"
// @Success      200   {object}  serializer.TeamApplicationDetailResponse
// @Failure      400   {object}  serializer.ErrorResponse
// @Failure      403   {object}  serializer.ErrorResponse
// @Failure      409   {object}  serializer.ErrorResponse
// @Router       /api/team-applications/{id} [put]
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
	userID, _ := currentUser(ctx)
	err := h.Repository.UpdateTeamApplication(id, userID, req.TeamName, req.Institution, req.CoachName, req.CompetitionDate)
	if err != nil {
		apiError(ctx, err)
		return
	}
	h.respondTeamApplication(ctx, id)
}

// ApiFormTeamApplication godoc
// @Summary      Сформировать заявку
// @Description  Представитель формирует черновик: проверяются обязательные поля и рассчитывается средний возраст команды.
// @Tags         Заявки команд
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "id заявки"
// @Success      200  {object}  serializer.TeamApplicationDetailResponse
// @Failure      400  {object}  serializer.ErrorResponse  "не заполнены обязательные поля"
// @Failure      403  {object}  serializer.ErrorResponse
// @Failure      409  {object}  serializer.ErrorResponse  "заявка не черновик"
// @Router       /api/team-applications/{id}/form [put]
func (h *Handler) ApiFormTeamApplication(ctx *gin.Context) {
	id, ok := uintParam(ctx, "id")
	if !ok {
		return
	}
	userID, _ := currentUser(ctx)
	if err := h.Repository.FormTeamApplication(id, userID); err != nil {
		apiError(ctx, err)
		return
	}
	h.respondTeamApplication(ctx, id)
}

// ApiModerateTeamApplication godoc
// @Summary      Допустить или отклонить заявку
// @Description  Организатор соревнования завершает (complete — допуск) или отклоняет (reject) сформированную заявку.
// @Tags         Заявки команд
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int                         true  "id заявки"
// @Param        body  body      serializer.ModerateRequest  true  "решение: complete или reject"
// @Success      200   {object}  serializer.TeamApplicationDetailResponse
// @Failure      401   {object}  serializer.ErrorResponse
// @Failure      403   {object}  serializer.ErrorResponse  "не организатор"
// @Failure      409   {object}  serializer.ErrorResponse  "заявка не сформирована"
// @Router       /api/team-applications/{id}/moderate [put]
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
	userID, _ := currentUser(ctx)
	if err := h.Repository.ModerateTeamApplication(id, userID, req.Action); err != nil {
		apiError(ctx, err)
		return
	}
	h.respondTeamApplication(ctx, id)
}

// ApiDeleteTeamApplication godoc
// @Summary      Удалить черновик
// @Description  Логическое удаление: статус deleted, проставляется дата формирования.
// @Tags         Заявки команд
// @Security     BearerAuth
// @Param        id   path  int  true  "id заявки"
// @Success      204
// @Failure      403  {object}  serializer.ErrorResponse
// @Failure      409  {object}  serializer.ErrorResponse
// @Router       /api/team-applications/{id} [delete]
func (h *Handler) ApiDeleteTeamApplication(ctx *gin.Context) {
	id, ok := uintParam(ctx, "id")
	if !ok {
		return
	}
	userID, _ := currentUser(ctx)
	if err := h.Repository.DeleteTeamApplication(id, userID); err != nil {
		apiError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// respondTeamApplication — заявка с составом; представителю — только своя
func (h *Handler) respondTeamApplication(ctx *gin.Context, id uint) {
	app, items, err := h.Repository.GetTeamApplication(id)
	if err != nil {
		apiError(ctx, err)
		return
	}
	userID, isOrganizer := currentUser(ctx)
	if !isOrganizer && app.CreatorID != userID {
		apiError(ctx, repository.ErrForbidden)
		return
	}
	ctx.JSON(http.StatusOK, serializer.TeamApplicationDetail(app, items))
}
