package vfsgo

type userCommandService struct {
}

func (serv *userCommandService) GetCurrentUser() *User {
	return nil
}

func (serv *userCommandService) Register(name string) error {
	return nil
}

func (serv *userCommandService) Use(name string) error {
	return nil
}

func (serv *userCommandService) ChangeFolder(path string) error {
	return nil
}
