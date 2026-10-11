package serializer

import (
	"time"

	"volleyreg/internal/app/ds"
)

const dateLayout = "2006-01-02"

// ---------- Участники ----------

// ParticipantResponse — участник в ответе API
type ParticipantResponse struct {
	ID          uint    `json:"id"`
	FullName    string  `json:"full_name"`
	Description string  `json:"description"`
	BirthDate   string  `json:"birth_date"`
	Position    string  `json:"position"`
	HeightCm    int     `json:"height_cm"`
	SportRank   string  `json:"sport_rank"`
	City        string  `json:"city"`
	ImageURL    *string `json:"image_url"`
	VideoURL    *string `json:"video_url"`
}

func Participant(p ds.Participant) ParticipantResponse {
	return ParticipantResponse{
		ID:          p.ID,
		FullName:    p.FullName,
		Description: p.Description,
		BirthDate:   p.BirthDate.Format(dateLayout),
		Position:    p.Position,
		HeightCm:    p.HeightCm,
		SportRank:   p.SportRank,
		City:        p.City,
		ImageURL:    minioURL(p.ImageKey),
		VideoURL:    minioURL(p.VideoKey),
	}
}

func Participants(list []ds.Participant) []ParticipantResponse {
	result := make([]ParticipantResponse, 0, len(list))
	for _, p := range list {
		result = append(result, Participant(p))
	}
	return result
}

func minioURL(key *string) *string {
	if key == nil || *key == "" {
		return nil
	}
	url := ds.MinioURL + *key
	return &url
}

// ---------- Заявки команд ----------

// DraftResponse — иконка корзины: черновик представителя и число участников в нём
type DraftResponse struct {
	TeamApplicationID *uint `json:"team_application_id"`
	ParticipantsCount int64 `json:"participants_count"`
}

// TeamApplicationResponse — заявка команды в списке
type TeamApplicationResponse struct {
	ID                uint     `json:"id"`
	Status            string   `json:"status"`
	CreatedAt         string   `json:"created_at"`
	FormedAt          *string  `json:"formed_at"`
	CompletedAt       *string  `json:"completed_at"`
	CreatorLogin      string   `json:"creator_login"`
	ModeratorLogin    *string  `json:"moderator_login"`
	TeamName          *string  `json:"team_name"`
	Institution       *string  `json:"institution"`
	CoachName         *string  `json:"coach_name"`
	CompetitionDate   *string  `json:"competition_date"`
	AverageAge        *float64 `json:"average_age"`
	ParticipantsCount int64    `json:"participants_count"`
}

// TeamApplicationParticipantResponse — участник в составе заявки с игровым номером
type TeamApplicationParticipantResponse struct {
	ParticipantResponse
	JerseyNumber *int `json:"jersey_number"`
}

// TeamApplicationDetailResponse — заявка команды вместе с составом
type TeamApplicationDetailResponse struct {
	TeamApplicationResponse
	Participants []TeamApplicationParticipantResponse `json:"participants"`
}

func TeamApplication(a ds.TeamApplication) TeamApplicationResponse {
	resp := TeamApplicationResponse{
		ID:                a.ID,
		Status:            a.Status,
		CreatedAt:         a.CreatedAt.Format(time.RFC3339),
		FormedAt:          timeString(a.FormedAt, time.RFC3339),
		CompletedAt:       timeString(a.CompletedAt, time.RFC3339),
		CreatorLogin:      a.Creator.Login,
		TeamName:          a.TeamName,
		Institution:       a.Institution,
		CoachName:         a.CoachName,
		CompetitionDate:   timeString(a.CompetitionDate, dateLayout),
		AverageAge:        a.AverageAge,
		ParticipantsCount: a.ParticipantsCount,
	}
	if a.Moderator != nil {
		resp.ModeratorLogin = &a.Moderator.Login
	}
	return resp
}

func TeamApplications(list []ds.TeamApplication) []TeamApplicationResponse {
	result := make([]TeamApplicationResponse, 0, len(list))
	for _, a := range list {
		result = append(result, TeamApplication(a))
	}
	return result
}

func TeamApplicationDetail(a ds.TeamApplication, items []ds.TeamApplicationParticipant) TeamApplicationDetailResponse {
	resp := TeamApplicationDetailResponse{TeamApplicationResponse: TeamApplication(a)}
	resp.ParticipantsCount = int64(len(items))
	resp.Participants = make([]TeamApplicationParticipantResponse, 0, len(items))
	for _, item := range items {
		resp.Participants = append(resp.Participants, TeamApplicationParticipantResponse{
			ParticipantResponse: Participant(item.Participant),
			JerseyNumber:        item.JerseyNumber,
		})
	}
	return resp
}

func timeString(t *time.Time, layout string) *string {
	if t == nil {
		return nil
	}
	s := t.Format(layout)
	return &s
}

// TeamApplicationUpdateRequest — поля заявки по теме, которые может менять представитель.
// Системные поля (id, статус, даты, создатель, модератор) с клиента не принимаются.
type TeamApplicationUpdateRequest struct {
	TeamName        *string `json:"team_name"`
	Institution     *string `json:"institution"`
	CoachName       *string `json:"coach_name"`
	CompetitionDate *string `json:"competition_date" example:"2026-11-20"`
}

// JerseyNumberRequest — изменение игрового номера участника в заявке (поле м-м)
type JerseyNumberRequest struct {
	JerseyNumber *int `json:"jersey_number" binding:"required"`
}

// ModerateRequest — решение организатора: complete (допустить) или reject (отклонить)
type ModerateRequest struct {
	Action string `json:"action" binding:"required,oneof=complete reject"`
}

// ---------- Пользователи ----------

type RegisterRequest struct {
	Login    string `json:"login" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name"`
}

type UserResponse struct {
	ID          uint   `json:"id"`
	Login       string `json:"login"`
	FullName    string `json:"full_name"`
	IsOrganizer bool   `json:"is_organizer"`
}

func User(u ds.User) UserResponse {
	return UserResponse{ID: u.ID, Login: u.Login, FullName: u.FullName, IsOrganizer: u.IsOrganizer}
}

// ErrorResponse — ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}
