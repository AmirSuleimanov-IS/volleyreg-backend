package ds

import (
	"strconv"
	"time"
)

// Статусы заявки команды
const (
	StatusDraft     = "draft"     // черновик
	StatusDeleted   = "deleted"   // удалена
	StatusFormed    = "formed"    // сформирована
	StatusCompleted = "completed" // завершена (допущена)
	StatusRejected  = "rejected"  // отклонена
)

// TeamApplication — заявка команды на участие в соревновании
type TeamApplication struct {
	ID          uint       `gorm:"primaryKey"`
	Status      string     `gorm:"type:varchar(20);not null"`
	CreatedAt   time.Time  `gorm:"not null"`
	CreatorID   uint       `gorm:"not null"`
	FormedAt    *time.Time // дата формирования (удаление или формирование представителем)
	CompletedAt *time.Time // дата завершения (допуск или отклонение организатором)
	ModeratorID *uint

	TeamName        *string    `gorm:"type:varchar(100)"`
	Institution     *string    `gorm:"type:varchar(200)"`
	CoachName       *string    `gorm:"type:varchar(150)"`
	CompetitionDate *time.Time `gorm:"type:date"`
	AverageAge      *float64   `gorm:"type:numeric(4,1)"` // средний возраст команды, вычисляется

	Creator   User  `gorm:"foreignKey:CreatorID"`
	Moderator *User `gorm:"foreignKey:ModeratorID"`
}

// StatusTitle — название статуса для интерфейса
func (a TeamApplication) StatusTitle() string {
	switch a.Status {
	case StatusDraft:
		return "Черновик"
	case StatusDeleted:
		return "Удалена"
	case StatusFormed:
		return "Сформирована"
	case StatusCompleted:
		return "Допущена"
	case StatusRejected:
		return "Отклонена"
	}
	return a.Status
}

// CompetitionDateText — дата соревнования для интерфейса
func (a TeamApplication) CompetitionDateText() string {
	if a.CompetitionDate == nil {
		return "—"
	}
	return a.CompetitionDate.Format("02.01.2006")
}

// AverageAgeText — средний возраст команды для интерфейса
func (a TeamApplication) AverageAgeText() string {
	if a.AverageAge == nil {
		return "—"
	}
	return strconv.FormatFloat(*a.AverageAge, 'f', 1, 64)
}
