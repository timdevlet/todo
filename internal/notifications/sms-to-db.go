package notifications

import (
	"fmt"

	"github.com/timdevlet/todo/internal/helpers"
	"github.com/timdevlet/todo/pkg/postgres"
)

type SMS2DB struct {
	gorm *postgres.GDB
	sms  *SMS

	before func(s SMS)
	after  func(s SMS)
}

func NewSms2DB(db *postgres.GDB) *SMS2DB {
	return &SMS2DB{
		gorm: db,
	}
}

func (s *SMS2DB) Set(sms *SMS) *SMS2DB {
	s.sms = sms

	return s
}

func (s *SMS2DB) SetBefore(f func(s SMS)) {
	s.before = f
}

func (s *SMS2DB) SetAfter(f func(s SMS)) {
	s.after = f
}

func (s *SMS2DB) Send() error {
	if s.before != nil {
		s.before(*s.sms)
	}

	if err, ok := helpers.ValidationStruct(s.sms); !ok {
		return fmt.Errorf("[validation:sms][errors:%v] %s", len(err), err[0])
	}

	r := s.gorm.DB.Create(s.sms)

	if r.Error != nil {
		return r.Error
	}

	if s.after != nil {
		s.after(*s.sms)
	}

	return nil
}
