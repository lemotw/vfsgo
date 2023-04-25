package vfsgo

type IUserCommandService interface {
	GetCurrentUser() *User

	Register(name string) error
	Use(name string) error
	ChangeFolder(path string) error
}

type IFolderCommandService interface {
	CreateFolder(user *User, oldName string) error
	DeleteFolder(user *User, oldName string) error
	ListFolder(user *User, dirName string) ([]string, error)
	RenameFolder(user *User, oldName string, newName string) error
}

type IFileCommandService interface {
	CreateFile(dirName, fileName, desc string) error
	DeleteFolder(dirName, oldName string) error
	ListFile(dirName string) ([]string, error)
}
