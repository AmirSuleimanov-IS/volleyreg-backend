package repository

import (
	"errors"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"volleyreg/internal/app/ds"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

// GetParticipants — действующие участники, в ФИО которых встречается participantName (ORM)
func (r *Repository) GetParticipants(participantName string) ([]ds.Participant, error) {
	var participants []ds.Participant
	query := r.db.Where("is_deleted = ?", false)
	if participantName != "" {
		query = query.Where("full_name ILIKE ?", "%"+participantName+"%")
	}
	err := query.Order("id").Find(&participants).Error
	return participants, err
}

// GetParticipant — один действующий участник (ORM)
func (r *Repository) GetParticipant(id uint) (ds.Participant, error) {
	var participant ds.Participant
	err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&participant).Error
	return participant, err
}

// GetDraftTeamApplication — черновик заявки представителя и количество участников в нём (ORM)
func (r *Repository) GetDraftTeamApplication(creatorID uint) (ds.TeamApplication, int64, error) {
	var app ds.TeamApplication
	err := r.db.Where("creator_id = ? AND status = ?", creatorID, ds.StatusDraft).First(&app).Error
	if err != nil {
		return app, 0, err
	}
	var count int64
	err = r.db.Model(&ds.TeamApplicationParticipant{}).Where("team_application_id = ?", app.ID).Count(&count).Error
	return app, count, err
}

// AddParticipantToDraft — добавляет участника в черновик; черновик создаётся, если его ещё нет (ORM)
func (r *Repository) AddParticipantToDraft(creatorID, participantID uint) error {
	if _, err := r.GetParticipant(participantID); err != nil {
		return err
	}

	app, _, err := r.GetDraftTeamApplication(creatorID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		app = ds.TeamApplication{
			Status:    ds.StatusDraft,
			CreatedAt: time.Now(),
			CreatorID: creatorID,
		}
		err = r.db.Create(&app).Error
	}
	if err != nil {
		return err
	}

	item := ds.TeamApplicationParticipant{TeamApplicationID: app.ID, ParticipantID: participantID}
	return r.db.Where(&item).FirstOrCreate(&item).Error
}

// GetTeamApplication — заявка (кроме удалённых) и её состав с участниками (ORM)
func (r *Repository) GetTeamApplication(id uint) (ds.TeamApplication, []ds.TeamApplicationParticipant, error) {
	var app ds.TeamApplication
	err := r.db.Where("id = ? AND status <> ?", id, ds.StatusDeleted).First(&app).Error
	if err != nil {
		return app, nil, err
	}
	var items []ds.TeamApplicationParticipant
	err = r.db.Preload("Participant").Where("team_application_id = ?", id).Order("id").Find(&items).Error
	return app, items, err
}

// DeleteTeamApplication — логическое удаление черновика SQL-запросом UPDATE (без ORM)
func (r *Repository) DeleteTeamApplication(id, creatorID uint) error {
	return r.db.Exec(
		"UPDATE team_applications SET status = 'deleted', formed_at = NOW() WHERE id = ? AND creator_id = ? AND status = 'draft'",
		id, creatorID,
	).Error
}
