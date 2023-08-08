package dic

import (
	"time"
)

type User struct {
	UUID  string `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name  string
	Email string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Citi struct {
	UUID string `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Country struct {
	UUID string `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
