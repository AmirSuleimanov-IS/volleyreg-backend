package ds

// User — пользователь системы: представитель команды или организатор соревнования
type User struct {
	ID          uint   `gorm:"primaryKey"`
	Login       string `gorm:"type:varchar(50);unique;not null"`
	Password    string `gorm:"type:varchar(100);not null"`
	FullName    string `gorm:"type:varchar(150)"`
	IsOrganizer bool   `gorm:"not null;default:false"`
}
