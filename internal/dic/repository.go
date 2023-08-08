package dic

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"github.com/timdevlet/todo/internal/time"
	"github.com/timdevlet/todo/pkg/postgres"
)

type DicRepository struct {
	postgres *postgres.PDB
	gorm     *postgres.GDB
}

func NewDicGRepository(db *postgres.GDB) *DicRepository {
	return &DicRepository{
		gorm: db,
	}
}
func NewDicRepository(db *postgres.PDB) *DicRepository {
	return &DicRepository{
		postgres: db,
	}
}

// ----------------------------

func (repo *DicRepository) FetchUsers() ([]User, error) {
	tm := time.NewTime()

	users := []User{}

	result := repo.gorm.DB.Table("users").Scan(&users)

	if result.Error != nil {
		return []User{}, result.Error
	}

	spew.Dump(tm)
	logrus.WithField("mili", tm.Mili()).Debug("[repo:dic] FetchUsers finished")

	return users, nil

	//return helpers.FetchManyFromPostgres[User](repo.postgres, "users")
}
