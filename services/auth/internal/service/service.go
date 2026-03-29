package service

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) Login() error {
	return nil
}

func (s *Service) Register() error {
	return nil
}

func (s *Service) ActivateUser() error {
	return nil
}

func (s *Service) RefreshToken() error {
	return nil
}

func (s *Service) ResetPassword() error {
	return nil
}
func (s *Service) Logout() error {
	return nil
}
