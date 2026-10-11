package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"volleyreg/internal/app/ds"
)

// AddParticipantToDraft — добавляет участника в черновик представителя.
// Если черновика нет, он создаётся пустым: статус, создатель и дата создания проставляются автоматически.
func (r *Repository) AddParticipantToDraft(creatorID, participantID uint) (uint, error) {
	if _, err := r.GetParticipant(participantID); err != nil {
		return 0, err
	}

	app, _, err := r.GetDraftTeamApplication(creatorID)
	if errors.Is(err, ErrNotFound) {
		app = ds.TeamApplication{Status: ds.StatusDraft, CreatedAt: time.Now(), CreatorID: creatorID}
		err = r.db.Create(&app).Error
	}
	if err != nil {
		return 0, err
	}

	item := ds.TeamApplicationParticipant{TeamApplicationID: app.ID, ParticipantID: participantID}
	return app.ID, r.db.Where(&item).FirstOrCreate(&item).Error
}

// draftOf — черновик с указанным id, принадлежащий представителю
func (r *Repository) draftOf(id, creatorID uint) (ds.TeamApplication, error) {
	var app ds.TeamApplication
	if err := r.db.Where("id = ? AND status <> ?", id, ds.StatusDeleted).First(&app).Error; err != nil {
		return app, notFound(err)
	}
	if app.CreatorID != creatorID {
		return app, ErrForbidden
	}
	if app.Status != ds.StatusDraft {
		return app, ErrStatus
	}
	return app, nil
}

// UpdateJerseyNumber — изменение игрового номера участника в черновике (без PK м-м)
func (r *Repository) UpdateJerseyNumber(appID, participantID, creatorID uint, number int) error {
	if _, err := r.draftOf(appID, creatorID); err != nil {
		return err
	}
	res := r.db.Model(&ds.TeamApplicationParticipant{}).
		Where("team_application_id = ? AND participant_id = ?", appID, participantID).
		Update("jersey_number", number)
	if res.Error == nil && res.RowsAffected == 0 {
		return ErrNotFound
	}
	return res.Error
}

// RemoveParticipant — удаление участника из черновика (без PK м-м)
func (r *Repository) RemoveParticipant(appID, participantID, creatorID uint) error {
	if _, err := r.draftOf(appID, creatorID); err != nil {
		return err
	}
	res := r.db.Where("team_application_id = ? AND participant_id = ?", appID, participantID).
		Delete(&ds.TeamApplicationParticipant{})
	if res.Error == nil && res.RowsAffected == 0 {
		return ErrNotFound
	}
	return res.Error
}

// countParticipants — количество участников в заявке
func (r *Repository) countParticipants(db *gorm.DB, appID uint) (int64, error) {
	var count int64
	err := db.Model(&ds.TeamApplicationParticipant{}).Where("team_application_id = ?", appID).Count(&count).Error
	return count, err
}
