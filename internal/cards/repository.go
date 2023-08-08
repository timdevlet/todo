package cards

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/timdevlet/todo/pkg/postgres"
	"gorm.io/datatypes"
)

type CardsRepository struct {
	gorm *postgres.GDB
}

func NewCardsRepository(db *postgres.GDB) *CardsRepository {
	return &CardsRepository{
		gorm: db,
	}
}

func (s *CardsRepository) Search() {
	cards := []Card{}

	//rand int
	// https://postgrespro.ru/docs/postgresql/15/functions-json
	// select '{"a":[{"v":3}, {"v":4}]}'::jsonb,  '{"a":[{"v":3}, {"v":4}]}'::jsonb @? '$.a[*].v ? (@ == 4)';

	err := s.gorm.DB.Table("cards").Create(&Card{
		Payload: datatypes.JSON([]byte(`{"name": "test", "foo": "bar", "fields": [{"aaa": "bbb"},{"zzz": "vvv"}]}`)),
		Title:   "test",
	})

	if err.Error != nil {
		log.Fatal(err.Error)
	}

	err = s.gorm.DB.Where("payload ->> 'fields' @@  '$.a[*].v ? (@ == 4)'").Find(&cards)

	if err.Error != nil {
		log.Fatal(err.Error)
	}

	spew.Dump("cards", cards)
}
