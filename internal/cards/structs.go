package cards

import (
	"time"

	"gorm.io/datatypes"
)

type Card struct {
	UUID      string `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Payload datatypes.JSON `sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
}
