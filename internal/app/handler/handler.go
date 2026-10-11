package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"volleyreg/internal/app/repository"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{Repository: r}
}

// GetParticipants — страница списка участников с поиском по ФИО
func (h *Handler) GetParticipants(ctx *gin.Context) {
	participantName := ctx.Query("participant_name")
	participants := h.Repository.GetParticipants(participantName)

	// карточка текущей заявки: id черновика и количество участников в нём
	draftID, draftCount := 0, 0
	if draft, ok := h.Repository.GetDraftTeamApplication(); ok {
		draftID = draft.ID
		draftCount = len(draft.Participants)
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
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusBadRequest, "некорректный id участника")
		return
	}
	participant, err := h.Repository.GetParticipant(id)
	if err != nil {
		ctx.String(http.StatusNotFound, err.Error())
		return
	}
	ctx.HTML(http.StatusOK, "participant.html", gin.H{
		"participant": participant,
	})
}

// teamApplicationRow — строка состава: участник и его игровой номер
type teamApplicationRow struct {
	Participant  repository.Participant
	JerseyNumber int
}

// GetTeamApplication — страница заявки команды с её составом
func (h *Handler) GetTeamApplication(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.String(http.StatusBadRequest, "некорректный id заявки")
		return
	}
	app, err := h.Repository.GetTeamApplication(id)
	if err != nil {
		ctx.String(http.StatusNotFound, err.Error())
		return
	}

	var rows []teamApplicationRow
	for _, item := range app.Participants {
		p, err := h.Repository.GetParticipant(item.ParticipantID)
		if err != nil {
			continue
		}
		rows = append(rows, teamApplicationRow{Participant: p, JerseyNumber: item.JerseyNumber})
	}

	ctx.HTML(http.StatusOK, "team_application.html", gin.H{
		"teamApplication": app,
		"rows":            rows,
	})
}
