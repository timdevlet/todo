package rbac

import (
	"github.com/timdevlet/todo/pkg/postgres"
)

type RBACRepository struct {
	postgres *postgres.PDB
	gorm     *postgres.GDB
}

func NewRBACRepository(db *postgres.GDB) *RBACRepository {
	return &RBACRepository{
		gorm: db,
	}
}
