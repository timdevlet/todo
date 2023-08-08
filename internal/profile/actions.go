package profile

import (
	"fmt"

	"github.com/timdevlet/todo/internal/rbac"
)

var actions = map[string]string{
	"user_add":       "users/add",
	"user_patch_fio": "users/patch/fio",
	"user_group_add": "users/group/add",
	"user_delete":    "users/delete",
}

func (s *ProfileService) AddUser(iam IAM, email string) error {
	if ok := rbac.CheckRight(actions["user_add"], iam.Actions...); !ok {
		return err(fmt.Sprintf("[iam:%s][action:%s] you don't have permission to add user", iam.Email, actions["user_add"]))
	}

	s.pub("user_add", email)

	return nil
}

func (s *ProfileService) DeleteUser(iam IAM, email string) error {
	if ok := rbac.CheckRight(actions["user_delete"], iam.Actions...); !ok {
		return err(fmt.Sprintf("[iam:%s][email:%s] you don't have permission to delete user", iam.Email, email))
	}

	if iam.Email == email {
		return err(fmt.Sprintf("[iam:%s][email:%s] you can't delete yourself", iam.Email, email))
	}

	return nil
}

func (s *ProfileService) AddToGroup(iam IAM, email string, groupUuid string) error {
	if ok := rbac.CheckRight(actions["user_group_add"], iam.Actions...); !ok {
		return err(fmt.Sprintf("[iam:%s][email:%s][gp:%s] you don't have permission to delete user", iam.Email, email, groupUuid))
	}

	return nil
}
