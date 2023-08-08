package dic

import "github.com/sirupsen/logrus"

type DicService struct {
	repo *DicRepository

	usersByUuid  map[string]User
	citiesByUuid map[string]Citi
}

func NewDicService(repo *DicRepository) *DicService {
	return &DicService{
		repo:        repo,
		usersByUuid: make(map[string]User),
	}
}

func (s *DicService) Sync() error {
	users, err := s.repo.FetchUsers()

	//cities, err := s.repo.FetchCities()

	logrus.WithField("total", len(users)).Debug("fetched users")

	for _, user := range users {
		s.usersByUuid[user.UUID] = user
	}

	// for _, item := range cities {
	// 	s.citiesByUuid[item.UUID] = item
	// }

	return err
}

func (s *DicService) InvalidateUser(uuid string) {
	delete(s.usersByUuid, uuid)
}
