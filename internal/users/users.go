package profile

import (
	"github.com/asaskevich/EventBus"
	"github.com/timdevlet/todo/internal/rbac"
)

type UserService struct {
	repo *UserRepository
	rbac *rbac.RBACService
	bus  *EventBus.Bus
}

func NewUserService(repo *UserRepository, bus *EventBus.Bus) *UserService {
	return &UserService{
		repo: repo,
		bus:  bus,
	}
}
