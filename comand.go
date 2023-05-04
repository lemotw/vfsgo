package vfsgo

type ICommandService interface {
	GetCurrentUser() *User
	GetCurrentBlock() *BlockINode

	Register(name string) (*User, error)
	Use(name string) error
	ChangeFolder(path string) error

	CreateFolder(oldName string) error
	DeleteFolder(oldName string) error
	ListFolder(dirName string) ([]string, error)
	RenameFolder(oldName string, newName string) error

	CreateFile(dirName, fileName, desc string) error
	ListFile(dirName string) ([]string, error)
}

func NewCommandService(root string) ICommandService {
	return &commandService{
		root:        root,
		currentUser: nil,
		userMap:     make(map[string]*User),
	}
}

type commandService struct {
	root         string
	currentUser  *User
	currentBlock *BlockINode

	userMap map[string]*User
}

func (cs *commandService) GetCurrentUser() *User {
	return cs.currentUser
}

func (cs *commandService) GetCurrentBlock() *BlockINode {
	return cs.currentBlock
}

func (cs *commandService) Register(name string) (*User, error) {
	panic("implement me")

}
func (cs *commandService) Use(name string) error {
	panic("implement me")
}

func (cs *commandService) ChangeFolder(path string) error {
	panic("implement me")
}

func (cs *commandService) CreateFolder(dirName string) error {
	panic("implement me")
}

func (cs *commandService) DeleteFolder(oldName string) error {
	panic("implement me")
}

func (cs *commandService) ListFolder(dirName string) ([]string, error) {
	panic("implement me")
}

func (cs *commandService) RenameFolder(oldName string, newName string) error {
	panic("implement me")
}

func (cs *commandService) CreateFile(dirName, fileName, desc string) error {
	panic("implement me")
}

func (cs *commandService) ListFile(dirName string) ([]string, error) {
	panic("implement me")
}
