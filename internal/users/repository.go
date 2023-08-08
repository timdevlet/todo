package profile

import (
	"github.com/timdevlet/todo/pkg/postgres"
)

type UserRepository struct {
	gorm *postgres.GDB
}

func NewUserRepository(db *postgres.GDB) *UserRepository {
	return &UserRepository{
		gorm: db,
	}
}
