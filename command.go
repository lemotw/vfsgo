package vfsgo

import (
	"errors"
	"strings"
)

type ICommandService interface {
	GetCurrentUser() *User

	Register(name string) (*User, error)
	Use(name string) error
	ChangeFolder(path string) error

	CreateFolder(user *User, oldName string) error
	DeleteFolder(user *User, oldName string) error
	ListFolder(user *User, dirName string) ([]string, error)
	RenameFolder(user *User, oldName string, newName string) error

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

func validRegister(name string) error {
	return nil
}

func (cs *commandService) Register(name string) (*User, error) {

	if err := validRegister(name); err != nil {
		return nil, err
	}

	u, err := CreateUser(cs.root, name)
	if err != nil {
		return nil, err
	}

	cs.userMap[name] = &u

	return &u, nil
}
func (cs *commandService) Use(name string) error {
	if u, ok := cs.userMap[name]; ok {
		if err := AttemptUser(u.RootPath, u.Name); err != nil {
			return err
		}
		cs.currentUser = u
		return nil
	}

	u, err := GetUser(cs.root, name)
	if err != nil {
		return err
	}
	cs.currentUser = &u
	cs.userMap[name] = &u

	return nil
}

func (cs *commandService) travelFolder(dirPath []string) (*BlockINode, error) {
	ret := cs.currentBlock

	for i := 0; i < len(dirPath); i++ {
		dirname := strings.TrimSpace(dirPath[i])
		if dirname == "." {
			continue
		}

		var nodeid uint64
		if dirname == ".." {
			nodeid = ret.PrevNodeID
		} else if ret.FileMap != nil {
			file, ok := ret.FileMap[dirname]
			if ok && file.Type == Directory && file.DirNodeID != nil {
				nodeid = *file.DirNodeID
			} else {
				return nil, errors.New("not a directory")
			}
		} else {
			return nil, errors.New("directory not exist")
		}

		if block, ok := cs.currentUser.BlockMap[nodeid]; ok {
			ret = block
		} else {
			return nil, errors.New("directory not exist")
		}
	}

	return ret, nil
}

func (cs *commandService) ChangeFolder(path string) error {
	if path == "." {
		return nil
	}

	dirPath := strings.Split(path, "/")

	block, err := cs.travelFolder(dirPath)
	if err != nil {
		return err
	}
	cs.currentBlock = block

	return nil
}

func (cs *commandService) CreateFolder(user *User, oldName string) error {
	return nil
}
func (cs *commandService) DeleteFolder(user *User, oldName string) error {
	return nil
}
func (cs *commandService) ListFolder(user *User, dirName string) ([]string, error) {
	return nil, nil
}
func (cs *commandService) RenameFolder(user *User, oldName string, newName string) error {
	return nil
}

func (cs *commandService) CreateFile(dirName, fileName, desc string) error {

	createBlock := cs.currentBlock
	dirname := strings.TrimSpace(dirName)
	if dirname == "." || dirname == "" {
		dirPath := strings.Split(dirName, "/")
		block, err := cs.travelFolder(dirPath)
		if err != nil {
			return err
		}

		createBlock = block
	}

	if createBlock != nil {
		if _, ok := createBlock.FileMap[fileName]; ok {
			return errors.New("file already exist")
		}

		file, err := CreateFile(createBlock, fileName, desc)
		if err != nil {
			return err
		}

		createBlock.FileMap[fileName] = &file
	}

	return nil
}
func (cs *commandService) ListFile(dirName string) ([]string, error) {
	if cs.currentBlock == nil {
		return nil, errors.New("please chose a folder")
	}

	ret := make([]string, len(cs.currentBlock.FileMap), 0)
	for fname, file := range cs.currentBlock.FileMap {
		if file.Type == Directory {
			fname += "/"
		}

		ret = append(ret, fname)
	}

	return ret, nil
}
