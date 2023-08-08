package rbac

type RBACService struct {
	repo *RBACRepository
}

func NewRBACService(repo *RBACRepository) *RBACService {
	return &RBACService{
		repo: repo,
	}
}

func CheckRight(action string, rights ...string) bool {
	//actions := strings.Split(action, "/")

	for _, right := range rights {

		if right == action {
			return true
		}
	}

	return false
}

func (s *RBACService) CreateGroup() {

}

func (s *RBACService) DeleteGroup() {

}

func (s *RBACService) AddUserToGroup() {

}

func (s *RBACService) RemoveUseFromGroup() {

}

func (s *RBACService) GetRights() []string {
	rights := []string{
		"user/*/create",
		"user/*/list",
		"user/*/delete",
		"user/*/update",

		"messenger/send",
	}

	return rights
}
