package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"volleyreg/internal/app/serializer"
	"volleyreg/internal/app/singleton"
)

// ApiAddParticipantToDraft — POST /api/team-applications/draft/participants/:participant_id:
// добавить участника в заявку-черновик (черновик создаётся автоматически)
func (h *Handler) ApiAddParticipantToDraft(ctx *gin.Context) {
	participantID, ok := uintParam(ctx, "participant_id")
	if !ok {
		return
	}
	appID, err := h.Repository.AddParticipantToDraft(singleton.GetCreatorID(), participantID)
	if err != nil {
		apiError(ctx, err)
		return
	}
	_, count, err := h.Repository.GetDraftTeamApplication(singleton.GetCreatorID())
	if err != nil {
		apiError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, serializer.DraftResponse{TeamApplicationID: &appID, ParticipantsCount: count})
}

// ApiUpdateJerseyNumber — PUT /api/team-applications/:id/participants/:participant_id: игровой номер (м-м)
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
	if err := h.Repository.UpdateJerseyNumber(appID, participantID, singleton.GetCreatorID(), *req.JerseyNumber); err != nil {
		apiError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"team_application_id": appID, "participant_id": participantID, "jersey_number": *req.JerseyNumber})
}

// ApiRemoveParticipant — DELETE /api/team-applications/:id/participants/:participant_id: убрать участника из черновика
func (h *Handler) ApiRemoveParticipant(ctx *gin.Context) {
	appID, ok := uintParam(ctx, "id")
	if !ok {
		return
	}
	participantID, ok := uintParam(ctx, "participant_id")
	if !ok {
		return
	}
	if err := h.Repository.RemoveParticipant(appID, participantID, singleton.GetCreatorID()); err != nil {
		apiError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
