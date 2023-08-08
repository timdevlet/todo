package profile

import (
	"fmt"

	"github.com/timdevlet/todo/internal/eventbus"
	"github.com/timdevlet/todo/internal/rbac"
)

type ActionChangeName struct {
	IAM  IAM
	UUID string
	Name string
}

func (s *ProfileService) ActionChangeName(a *ActionChangeName) error {
	if ok := rbac.CheckRight(actions["user_patch_fio"], a.IAM.Actions...); !ok {
		return err(fmt.Sprintf("[iam:%s][action:%s] you don't have permission to patch fio", a.IAM.Email, actions["user_patch_fio"]))
	}

	event := eventbus.NewUpdateProfileEvent(a.Name)
	s.pub(event.Channel, *event)

	return s.ChangeFIO(a.UUID, a.Name)
}
