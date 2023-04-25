package vfsgo

type folderCommandService struct {
}

func (serv *folderCommandService) CreateFolder(user *User, oldName string) error {
	return nil
}

func (serv *folderCommandService) DeleteFolder(user *User, oldName string) error {
	return nil
}

func (serv *folderCommandService) ListFolder(user *User, dirName string) ([]string, error) {
	return nil, nil
}

func (serv *folderCommandService) RenameFolder(user *User, oldName string, newName string) error {
	return nil
}
