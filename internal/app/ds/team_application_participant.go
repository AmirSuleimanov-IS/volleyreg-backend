package ds

// TeamApplicationParticipant — участник в составе заявки команды (м-м)
type TeamApplicationParticipant struct {
	ID                uint `gorm:"primaryKey"`
	TeamApplicationID uint `gorm:"not null;uniqueIndex:idx_team_application_participant"`
	ParticipantID     uint `gorm:"not null;uniqueIndex:idx_team_application_participant"`
	JerseyNumber      *int // игровой номер участника в этой команде

	TeamApplication TeamApplication `gorm:"foreignKey:TeamApplicationID"`
	Participant     Participant     `gorm:"foreignKey:ParticipantID"`
}
