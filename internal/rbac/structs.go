package rbac

import "time"

type RBAC struct {
	UUID        string `gorm:"type:uuid;default:uuid_generate_v4()"`
	Rights      string
	Name        string
	Description string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type RBAC_USER struct {
	UUID        string `gorm:"type:uuid;default:uuid_generate_v4()"`
	UserUuid    string
	RbacUuid    string
	Description string

	CreatedAt time.Time
}
