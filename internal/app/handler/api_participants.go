package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"volleyreg/internal/app/ds"
	"volleyreg/internal/app/serializer"
)

// ApiGetParticipants — GET /api/participants: список участников с фильтрацией по ФИО и амплуа
func (h *Handler) ApiGetParticipants(ctx *gin.Context) {
	participants, err := h.Repository.GetParticipants(ctx.Query("participant_name"), ctx.Query("position"))
	if err != nil {
		apiError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, serializer.Participants(participants))
}

// ApiGetParticipant — GET /api/participants/:id: один участник
func (h *Handler) ApiGetParticipant(ctx *gin.Context) {
	id, ok := uintParam(ctx, "id")
	if !ok {
		return
	}
	participant, err := h.Repository.GetParticipant(id)
	if err != nil {
		apiError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, serializer.Participant(participant))
}

// ApiCreateParticipant — POST /api/participants (multipart/form-data):
// новый участник + файлы image и video, которые сохраняются в Minio под латинскими именами
func (h *Handler) ApiCreateParticipant(ctx *gin.Context) {
	heightCm, err := strconv.Atoi(ctx.PostForm("height_cm"))
	if err != nil || heightCm <= 0 {
		badRequest(ctx, "height_cm должен быть положительным числом")
		return
	}
	birthDate, err := time.Parse("2006-01-02", ctx.PostForm("birth_date"))
	if err != nil {
		badRequest(ctx, "birth_date должен быть в формате YYYY-MM-DD")
		return
	}
	participant := ds.Participant{
		FullName:    strings.TrimSpace(ctx.PostForm("full_name")),
		Description: strings.TrimSpace(ctx.PostForm("description")),
		BirthDate:   birthDate,
		Position:    strings.TrimSpace(ctx.PostForm("position")),
		HeightCm:    heightCm,
		SportRank:   strings.TrimSpace(ctx.PostForm("sport_rank")),
		City:        strings.TrimSpace(ctx.PostForm("city")),
	}
	if participant.FullName == "" || participant.Description == "" || participant.Position == "" ||
		participant.SportRank == "" || participant.City == "" {
		badRequest(ctx, "заполните full_name, description, position, sport_rank, city")
		return
	}
	if err := h.Repository.CreateParticipant(&participant); err != nil {
		apiError(ctx, err)
		return
	}

	// имена файлов генерируются на латинице: participant_<id>.jpg, participant_<id>_video.mp4
	var imageKey, videoKey *string
	if file, err := ctx.FormFile("image"); err == nil {
		key := fmt.Sprintf("participant_%d%s", participant.ID, strings.ToLower(filepath.Ext(file.Filename)))
		if err := h.Storage.Upload(ctx, key, file); err != nil {
			apiError(ctx, err)
			return
		}
		imageKey = &key
	}
	if file, err := ctx.FormFile("video"); err == nil {
		key := fmt.Sprintf("participant_%d_video%s", participant.ID, strings.ToLower(filepath.Ext(file.Filename)))
		if err := h.Storage.Upload(ctx, key, file); err != nil {
			apiError(ctx, err)
			return
		}
		videoKey = &key
	}
	if err := h.Repository.SetParticipantMedia(participant.ID, imageKey, videoKey); err != nil {
		apiError(ctx, err)
		return
	}
	participant.ImageKey, participant.VideoKey = imageKey, videoKey
	ctx.JSON(http.StatusCreated, serializer.Participant(participant))
}
