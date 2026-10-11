package ds

import "time"

// MinioURL — публичный адрес бакета Minio с фото и видео участников.
// Меняется переменной MINIO_PUBLIC_URL (например, на IP компьютера в локальной сети для Tauri).
var MinioURL = "http://localhost:9000/volleyreg/"

// Participant — участник соревнования (услуга)
type Participant struct {
	ID          uint      `gorm:"primaryKey"`
	FullName    string    `gorm:"type:varchar(150);not null"`
	Description string    `gorm:"type:text;not null"`
	IsDeleted   bool      `gorm:"not null;default:false"` // статус: удалён / действует
	ImageKey    *string   `gorm:"type:varchar(255)"`
	VideoKey    *string   `gorm:"type:varchar(255)"`
	BirthDate   time.Time `gorm:"type:date;not null"`
	Position    string    `gorm:"type:varchar(50);not null"`
	HeightCm    int       `gorm:"not null"`
	SportRank   string    `gorm:"type:varchar(30);not null"`
	City        string    `gorm:"type:varchar(100);not null"`
}

// ImageURL — ссылка на фото в Minio (или фото по умолчанию)
func (p Participant) ImageURL() string {
	if p.ImageKey == nil || *p.ImageKey == "" {
		return "/styles/default_participant.svg"
	}
	return MinioURL + *p.ImageKey
}

// VideoURL — ссылка на видео в Minio
func (p Participant) VideoURL() string {
	if p.VideoKey == nil || *p.VideoKey == "" {
		return ""
	}
	return MinioURL + *p.VideoKey
}
