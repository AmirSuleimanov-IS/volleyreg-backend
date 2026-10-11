package repository

import (
	"fmt"
	"strings"
	"time"

	"volleyreg/internal/app/ds"
)

// GetDraftTeamApplication — черновик представителя и количество участников в нём (ORM)
func (r *Repository) GetDraftTeamApplication(creatorID uint) (ds.TeamApplication, int64, error) {
	var app ds.TeamApplication
	err := r.db.Where("creator_id = ? AND status = ?", creatorID, ds.StatusDraft).First(&app).Error
	if err != nil {
		return app, 0, notFound(err)
	}
	count, err := r.countParticipants(r.db, app.ID)
	return app, count, err
}

// TeamApplicationFilter — фильтры списка заявок
type TeamApplicationFilter struct {
	Status     string // статус
	FormedFrom string // дата формирования с (YYYY-MM-DD)
	FormedTo   string // дата формирования по (YYYY-MM-DD)
	CreatorID  uint   // 0 — заявки всех представителей
}

// GetTeamApplications — список заявок без удалённых и черновиков, с логинами создателя и модератора
// и количеством участников (подзапрос)
func (r *Repository) GetTeamApplications(f TeamApplicationFilter) ([]ds.TeamApplication, error) {
	query := r.db.Model(&ds.TeamApplication{}).
		Select("team_applications.*, (SELECT COUNT(*) FROM team_application_participants tap WHERE tap.team_application_id = team_applications.id) AS participants_count").
		Preload("Creator").Preload("Moderator").
		Where("status NOT IN ?", []string{ds.StatusDeleted, ds.StatusDraft})
	if f.Status != "" {
		query = query.Where("status = ?", f.Status)
	}
	if f.FormedFrom != "" {
		from, err := time.Parse("2006-01-02", f.FormedFrom)
		if err != nil {
			return nil, fmt.Errorf("%w: formed_from должен быть в формате YYYY-MM-DD", ErrValidation)
		}
		query = query.Where("formed_at >= ?", from)
	}
	if f.FormedTo != "" {
		to, err := time.Parse("2006-01-02", f.FormedTo)
		if err != nil {
			return nil, fmt.Errorf("%w: formed_to должен быть в формате YYYY-MM-DD", ErrValidation)
		}
		query = query.Where("formed_at < ?", to.AddDate(0, 0, 1))
	}
	if f.CreatorID != 0 {
		query = query.Where("creator_id = ?", f.CreatorID)
	}
	var apps []ds.TeamApplication
	err := query.Order("formed_at DESC").Find(&apps).Error
	return apps, err
}

// GetTeamApplication — заявка (кроме удалённых) и её состав с участниками (ORM)
func (r *Repository) GetTeamApplication(id uint) (ds.TeamApplication, []ds.TeamApplicationParticipant, error) {
	var app ds.TeamApplication
	err := r.db.Preload("Creator").Preload("Moderator").
		Where("id = ? AND status <> ?", id, ds.StatusDeleted).First(&app).Error
	if err != nil {
		return app, nil, notFound(err)
	}
	var items []ds.TeamApplicationParticipant
	err = r.db.Preload("Participant").Where("team_application_id = ?", id).Order("id").Find(&items).Error
	return app, items, err
}

// UpdateTeamApplication — изменение полей заявки по теме (только черновик своего представителя)
func (r *Repository) UpdateTeamApplication(id, creatorID uint, teamName, institution, coachName, competitionDate *string) error {
	if _, err := r.draftOf(id, creatorID); err != nil {
		return err
	}
	updates := map[string]interface{}{}
	if teamName != nil {
		updates["team_name"] = strings.TrimSpace(*teamName)
	}
	if institution != nil {
		updates["institution"] = strings.TrimSpace(*institution)
	}
	if coachName != nil {
		updates["coach_name"] = strings.TrimSpace(*coachName)
	}
	if competitionDate != nil {
		date, err := time.Parse("2006-01-02", *competitionDate)
		if err != nil {
			return fmt.Errorf("%w: competition_date должен быть в формате YYYY-MM-DD", ErrValidation)
		}
		updates["competition_date"] = date
	}
	if len(updates) == 0 {
		return fmt.Errorf("%w: нет полей для изменения", ErrValidation)
	}
	return r.db.Model(&ds.TeamApplication{}).Where("id = ?", id).Updates(updates).Error
}

// FormTeamApplication — формирование заявки представителем: проверка обязательных полей
// и расчёт среднего возраста команды на дату соревнования
func (r *Repository) FormTeamApplication(id, creatorID uint) error {
	app, err := r.draftOf(id, creatorID)
	if err != nil {
		return err
	}
	if !ds.CanChangeStatus(app.Status, ds.StatusFormed, false) {
		return ErrStatus
	}
	if app.TeamName == nil || *app.TeamName == "" || app.CompetitionDate == nil {
		return fmt.Errorf("%w: заполните название команды и дату соревнования", ErrValidation)
	}

	var birthDates []time.Time
	err = r.db.Model(&ds.Participant{}).
		Joins("JOIN team_application_participants tap ON tap.participant_id = participants.id").
		Where("tap.team_application_id = ?", id).
		Pluck("participants.birth_date", &birthDates).Error
	if err != nil {
		return err
	}
	if len(birthDates) == 0 {
		return fmt.Errorf("%w: в заявке нет участников", ErrValidation)
	}

	averageAge := AverageAge(birthDates, *app.CompetitionDate)
	return r.db.Model(&ds.TeamApplication{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      ds.StatusFormed,
		"formed_at":   time.Now(),
		"average_age": averageAge,
	}).Error
}

// ModerateTeamApplication — организатор допускает (complete) или отклоняет (reject) сформированную заявку
func (r *Repository) ModerateTeamApplication(id, moderatorID uint, action string) error {
	target := ds.StatusCompleted
	if action == "reject" {
		target = ds.StatusRejected
	}
	var app ds.TeamApplication
	if err := r.db.Where("id = ? AND status <> ?", id, ds.StatusDeleted).First(&app).Error; err != nil {
		return notFound(err)
	}
	if !ds.CanChangeStatus(app.Status, target, true) {
		return fmt.Errorf("%w: заявку в статусе %s нельзя перевести в %s", ErrStatus, app.Status, target)
	}
	return r.db.Model(&ds.TeamApplication{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       target,
		"moderator_id": moderatorID,
		"completed_at": time.Now(),
	}).Error
}

// DeleteTeamApplication — логическое удаление черновика SQL-запросом UPDATE (без ORM),
// проставляется дата формирования
func (r *Repository) DeleteTeamApplication(id, creatorID uint) error {
	if _, err := r.draftOf(id, creatorID); err != nil {
		return err
	}
	return r.db.Exec(
		"UPDATE team_applications SET status = 'deleted', formed_at = NOW() WHERE id = ? AND creator_id = ? AND status = 'draft'",
		id, creatorID,
	).Error
}
