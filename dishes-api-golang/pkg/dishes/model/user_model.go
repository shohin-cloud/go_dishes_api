package model

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"size:25;not null"`
	Firstname string `gorm:"size:25;not null"`
	Lastname  string `gorm:"size:25;not null"`
	Email     string `gorm:"size:25;not null"`
	Password  string `gorm:"size:25;not null"`
}
