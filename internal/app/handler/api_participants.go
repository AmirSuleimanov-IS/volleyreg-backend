package handler

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"volleyreg/internal/app/ds"
	"volleyreg/internal/app/serializer"
)

// ApiGetParticipants godoc
// @Summary      Список участников
// @Description  Действующие участники с фильтрацией по ФИО и амплуа. Доступно гостю.
// @Tags         Участники
// @Produce      json
// @Param        participant_name  query     string  false  "фрагмент ФИО"
// @Param        position          query     string  false  "амплуа (Связующая, Либеро, ...)"
// @Success      200               {array}   serializer.ParticipantResponse
// @Router       /api/participants [get]
func (h *Handler) ApiGetParticipants(ctx *gin.Context) {
	participants, err := h.Repository.GetParticipants(ctx.Query("participant_name"), ctx.Query("position"))
	if err != nil {
		apiError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, serializer.Participants(participants))
}

// ApiGetParticipant godoc
// @Summary      Один участник
// @Tags         Участники
// @Produce      json
// @Param        id   path      int  true  "id участника"
// @Success      200  {object}  serializer.ParticipantResponse
// @Failure      404  {object}  serializer.ErrorResponse
// @Router       /api/participants/{id} [get]
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

// ApiCreateParticipant godoc
// @Summary      Добавить участника
// @Description  Новый участник с фото и коротким видео. Файлы сохраняются в Minio под латинскими именами. Только организатор.
// @Tags         Участники
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        full_name    formData  string  true   "ФИО"
// @Param        description  formData  string  true   "описание"
// @Param        birth_date   formData  string  true   "дата рождения YYYY-MM-DD"
// @Param        position     formData  string  true   "амплуа"
// @Param        height_cm    formData  int     true   "рост, см"
// @Param        sport_rank   formData  string  true   "спортивный разряд"
// @Param        city         formData  string  true   "город"
// @Param        image        formData  file    false  "фото"
// @Param        video        formData  file    false  "короткое видео"
// @Success      201          {object}  serializer.ParticipantResponse
// @Failure      400          {object}  serializer.ErrorResponse
// @Failure      401          {object}  serializer.ErrorResponse
// @Failure      403          {object}  serializer.ErrorResponse
// @Router       /api/participants [post]
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
	// файлы проверяются до создания участника
	image, imageExt, ok := mediaFile(ctx, "image", imageTypes)
	if !ok {
		return
	}
	video, videoExt, ok := mediaFile(ctx, "video", videoTypes)
	if !ok {
		return
	}
	if err := h.Repository.CreateParticipant(&participant); err != nil {
		apiError(ctx, err)
		return
	}

	// имена файлов генерируются на латинице: participant_<id>.jpg, participant_<id>_video.mp4
	var imageKey, videoKey *string
	if image != nil {
		key := fmt.Sprintf("participant_%d%s", participant.ID, imageExt)
		if err := h.Storage.Upload(ctx, key, imageTypes[imageExt], image); err != nil {
			apiError(ctx, err)
			return
		}
		imageKey = &key
	}
	if video != nil {
		key := fmt.Sprintf("participant_%d_video%s", participant.ID, videoExt)
		if err := h.Storage.Upload(ctx, key, videoTypes[videoExt], video); err != nil {
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

// Допустимые файлы участника: расширение → Content-Type. Тип от клиента не используется,
// чтобы в публичный бакет нельзя было положить, например, HTML или SVG со скриптом.
var (
	imageTypes = map[string]string{".jpg": "image/jpeg", ".jpeg": "image/jpeg", ".png": "image/png", ".webp": "image/webp"}
	videoTypes = map[string]string{".mp4": "video/mp4", ".webm": "video/webm"}
)

// mediaFile — необязательный файл из формы с проверкой расширения; ok = false, если ответ с ошибкой уже отправлен
func mediaFile(ctx *gin.Context, field string, allowed map[string]string) (*multipart.FileHeader, string, bool) {
	file, err := ctx.FormFile(field)
	if err != nil {
		return nil, "", true
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if _, ok := allowed[ext]; !ok {
		badRequest(ctx, "недопустимый тип файла "+field+": "+ext)
		return nil, "", false
	}
	return file, ext, true
}
