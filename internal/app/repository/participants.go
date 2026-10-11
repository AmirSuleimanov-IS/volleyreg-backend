package repository

import "volleyreg/internal/app/ds"

// GetParticipants — действующие участники с фильтрами по ФИО и амплуа (ORM)
func (r *Repository) GetParticipants(participantName, position string) ([]ds.Participant, error) {
	var participants []ds.Participant
	query := r.db.Where("is_deleted = ?", false)
	if participantName != "" {
		query = query.Where("full_name ILIKE ?", "%"+participantName+"%")
	}
	if position != "" {
		query = query.Where("position = ?", position)
	}
	err := query.Order("id").Find(&participants).Error
	return participants, err
}

// GetParticipant — один действующий участник (ORM)
func (r *Repository) GetParticipant(id uint) (ds.Participant, error) {
	var participant ds.Participant
	err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&participant).Error
	return participant, notFound(err)
}

// CreateParticipant — добавление нового участника
func (r *Repository) CreateParticipant(p *ds.Participant) error {
	return r.db.Create(p).Error
}

// SetParticipantMedia — сохранение ключей фото и видео из Minio
func (r *Repository) SetParticipantMedia(id uint, imageKey, videoKey *string) error {
	return r.db.Model(&ds.Participant{}).Where("id = ?", id).
		Updates(map[string]interface{}{"image_key": imageKey, "video_key": videoKey}).Error
}
