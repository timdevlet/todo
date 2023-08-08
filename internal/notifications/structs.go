package notifications

import "time"

type SMS struct {
	UUID string `gorm:"type:uuid;default:gen_random_uuid()"`
	Text string `gorm:"type:varchar(200);default:''" validate:"min=3,max=200"`
	To   int64  `validate:"gte=9999,lte=99999999999"`

	CreatedAt time.Time `gorm:"index"`
}

type Email struct {
	UUID    string `gorm:"type:uuid;default:gen_random_uuid()"`
	To      string `gorm:"type:varchar(200);default:''" validate:"email"`
	From    string `gorm:"type: varchar(50); default:''" validate:"email"`
	Subject string `gorm:"type:varchar(200);default:''" validate:"min=3,max=200"`
	Text    string `gorm:"type:text;default:''"`

	CreatedAt time.Time `gorm:"index"`
}
