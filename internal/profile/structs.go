package profile

import (
	"time"

	"gorm.io/datatypes"
)

type User struct {
	UUID  string `sql:"uuid"`
	Name  string `sql:"name"`
	Email string `sql:"email"`
}

type GUser struct {
	UUID  string `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name  string
	Email string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Meta datatypes.JSON `sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
}

type GUserMeta struct {
	Fake *bool `json:"fake"`
}

//

type IAM struct {
	UUID    string
	Email   string
	Groups  []string
	Actions []string
}
