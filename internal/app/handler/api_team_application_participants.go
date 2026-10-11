package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"volleyreg/internal/app/serializer"
)

// ApiAddParticipantToDraft godoc
// @Summary      Добавить участника в черновик
// @Description  Добавляет участника в заявку-черновик текущего представителя. Черновик создаётся автоматически.
// @Tags         Состав заявки (м-м)
// @Produce      json
// @Security     BearerAuth
// @Param        participant_id  path      int  true  "id участника"
// @Success      201             {object}  serializer.DraftResponse
// @Failure      401             {object}  serializer.ErrorResponse
// @Failure      404             {object}  serializer.ErrorResponse
// @Router       /api/team-applications/draft/participants/{participant_id} [post]
func (h *Handler) ApiAddParticipantToDraft(ctx *gin.Context) {
	participantID, ok := uintParam(ctx, "participant_id")
	if !ok {
		return
	}
	userID, _ := currentUser(ctx)
	appID, err := h.Repository.AddParticipantToDraft(userID, participantID)
	if err != nil {
		apiError(ctx, err)
		return
	}
	_, count, err := h.Repository.GetDraftTeamApplication(userID)
	if err != nil {
		apiError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, serializer.DraftResponse{TeamApplicationID: &appID, ParticipantsCount: count})
}

// ApiUpdateJerseyNumber godoc
// @Summary      Изменить игровой номер
// @Description  Поле м-м: игровой номер участника в черновике заявки (без PK м-м).
// @Tags         Состав заявки (м-м)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id              path      int                             true  "id заявки"
// @Param        participant_id  path      int                             true  "id участника"
// @Param        body            body      serializer.JerseyNumberRequest  true  "игровой номер"
// @Success      200             {object}  map[string]int
// @Failure      400             {object}  serializer.ErrorResponse
// @Failure      403             {object}  serializer.ErrorResponse
// @Failure      409             {object}  serializer.ErrorResponse
// @Router       /api/team-applications/{id}/participants/{participant_id} [put]
func (h *Handler) ApiUpdateJerseyNumber(ctx *gin.Context) {
	appID, ok := uintParam(ctx, "id")
	if !ok {
		return
	}
	participantID, ok := uintParam(ctx, "participant_id")
	if !ok {
		return
	}
	var req serializer.JerseyNumberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil || *req.JerseyNumber < 0 || *req.JerseyNumber > 99 {
		badRequest(ctx, "jersey_number должен быть числом от 0 до 99")
		return
	}
	userID, _ := currentUser(ctx)
	if err := h.Repository.UpdateJerseyNumber(appID, participantID, userID, *req.JerseyNumber); err != nil {
		apiError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"team_application_id": appID, "participant_id": participantID, "jersey_number": *req.JerseyNumber})
}

// ApiRemoveParticipant godoc
// @Summary      Убрать участника из заявки
// @Tags         Состав заявки (м-м)
// @Security     BearerAuth
// @Param        id              path  int  true  "id заявки"
// @Param        participant_id  path  int  true  "id участника"
// @Success      204
// @Failure      403  {object}  serializer.ErrorResponse
// @Failure      404  {object}  serializer.ErrorResponse
// @Router       /api/team-applications/{id}/participants/{participant_id} [delete]
func (h *Handler) ApiRemoveParticipant(ctx *gin.Context) {
	appID, ok := uintParam(ctx, "id")
	if !ok {
		return
	}
	participantID, ok := uintParam(ctx, "participant_id")
	if !ok {
		return
	}
	userID, _ := currentUser(ctx)
	if err := h.Repository.RemoveParticipant(appID, participantID, userID); err != nil {
		apiError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
