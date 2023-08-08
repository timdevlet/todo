package notifications

import "github.com/timdevlet/todo/pkg/postgres"

type Email2DB struct {
	gorm *postgres.GDB
}

func NewEmail2DB(db *postgres.GDB) *Email2DB {
	return &Email2DB{
		gorm: db,
	}
}

func (s *Email2DB) Send() error {
	r := s.gorm.DB.Create(&Email{
		Text: "Hello world",
		To:   "aaaa@sss.sd",
	})

	if r.Error != nil {
		return r.Error
	}

	return nil
}
