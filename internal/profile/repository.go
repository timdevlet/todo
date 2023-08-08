package profile

import (
	"encoding/json"

	"github.com/davecgh/go-spew/spew"
	"github.com/timdevlet/todo/pkg/postgres"
	"gorm.io/datatypes"
)

type ProfileRepository struct {
	gorm *postgres.GDB
}

func NewProfileRepository(db *postgres.GDB) *ProfileRepository {
	return &ProfileRepository{
		gorm: db,
	}
}

func (repo *ProfileRepository) CreateUser(name, email string, meta *GUserMeta) (u GUser, err error) {

	jsonBytes, _ := json.Marshal(meta)
	if meta == nil {
		jsonBytes = []byte("{}")
	}

	u = GUser{
		Name:  name,
		Email: email,
		Meta:  datatypes.JSON(jsonBytes),
	}

	result := repo.gorm.DB.Table("users").Create(&u)
	err = result.Error

	spew.Dump("cccc", u.Meta)

	var m GUserMeta
	json.Unmarshal([]byte(u.Meta), &m)

	spew.Dump("cccc2", m.Fake)

	return
}

func (repo *ProfileRepository) Update(uuid, col string, val any) error {
	result := repo.gorm.DB.Table("users").Where("uuid = ?", uuid).Update(col, val)

	return result.Error
}
