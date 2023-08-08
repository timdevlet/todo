package profile

import (
	"time"

	"gorm.io/datatypes"
)

type User struct {
	UUID  string `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name  string
	Email string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Meta datatypes.JSON `sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
}
