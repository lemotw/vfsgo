package vfsgo

type fileCommandService struct {
}

func (serv *fileCommandService) CreateFile(dirName, fileName, desc string) error {
	return nil
}

func (serv *fileCommandService) DeleteFolder(dirName, oldName string) error {
	return nil
}

func (serv *fileCommandService) ListFile(dirName string) ([]string, error) {
	return nil, nil
}
