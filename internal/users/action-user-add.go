package profile

type ActionUserAdd struct {
	//IAM  IAM
	UUID string
	Name string
}

func (s *UserService) ActionUserAdd(a *ActionUserAdd) error {
	return nil
}
