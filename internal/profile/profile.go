package profile

import (
	"fmt"

	"github.com/asaskevich/EventBus"
	"github.com/timdevlet/todo/internal/helpers"
	"github.com/timdevlet/todo/internal/rbac"
)

type ProfileService struct {
	repo *ProfileRepository
	rbac *rbac.RBACService
	bus  *EventBus.Bus
}

func NewProfileService(repo *ProfileRepository, bus *EventBus.Bus) *ProfileService {
	return &ProfileService{
		repo: repo,
		bus:  bus,
	}
}

// ----------------------------

var ACTIONS = map[string]string{
	"profile/putch": "profile/putch",
}

// ----------------------------

func (s *ProfileService) CreateUser(name, email string) (GUser, error) {
	if err := helpers.ValidateEmail(email); err != nil {
		return GUser{}, err
	}

	return s.repo.CreateUser(name, email, nil)
}

func (s *ProfileService) CreateRandomUser() (GUser, error) {
	return s.repo.CreateUser(helpers.FakeName(), helpers.FakeEmail(), nil)
}

func (s *ProfileService) ChangeFIO(uuid string, name string) error {
	return s.repo.Update(uuid, "name", name)
}

// ----------------------------

func err(s string) error {
	return fmt.Errorf("[service:profile] %s", s)
}

func (s *ProfileService) pub(eventName string, date interface{}) {
	if (s.bus) != nil {
		b := *s.bus
		b.Publish(eventName, date)
	}
}

func (s *ProfileService) Go(action interface{}) error {
	actionName := helpers.GetType(action)

	if v, ok := action.(interface {
		Go() error
		SetService(s *ProfileService)
		Publish() (string, string)
	}); ok {
		// Set service to action
		v.SetService(s)

		// Run
		err := v.Go()
		if err != nil {
			return fmt.Errorf("[service:profile][action:%s] %s", actionName, err.Error())
		}

		//s.pub(helpers.GetType(action), v)
	} else {
		return err(fmt.Sprintf("action %v not found", helpers.GetType(action)))
	}

	return nil
}
