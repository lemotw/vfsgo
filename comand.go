package vfsgo

import (
	"os"
	"strings"

	"golang.org/x/xerrors"
)

type ICommandService interface {
	GetCurrentUser() *User
	GetCurrentBlock() *BlockINode

	Register(name string) error
	Use(name string) error
	ChangeFolder(path string) error

	CreateFolder(oldName string) error
	DeleteFolder(oldName string) error
	RenameFolder(oldName string, newName string) error

	CreateFile(fileName, desc string) error
	DeleteFile(fileName string) error
	RenameFile(oldName, newName string, newDesc string) error

	List(dirName string) ([]string, error)
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

func (cs *commandService) validRegister(name string) error {
	if _, ok := cs.userMap[name]; ok {
		return xerrors.Errorf("user %s already exists", name)
	}

	return nil
}

func (cs *commandService) travelFolder(path string) (*BlockINode, error) {
	if cs.currentBlock == nil {
		return nil, xerrors.New("current block is nil")
	}

	if path == "" {
		return cs.currentBlock, nil
	}

	blockRet := cs.currentBlock
	directories := strings.Split(strings.TrimSpace(path), "/")

	for i := 0; i < len(directories); i++ {
		var nodeid uint64

		switch directories[i] {
		case ".":
			nodeid = blockRet.NodeID
		case "..":
			nodeid = blockRet.PrevNodeID
		default:
			if b, ok := blockRet.FileMap[directories[i]]; ok && b.Type == Directory && b.DirNodeID != nil {
				nodeid = *b.DirNodeID
			} else {
				return nil, xerrors.Errorf("path %s not exist", path)
			}
		}

		if b, ok := cs.currentUser.BlockMap[nodeid]; !ok {
			return nil, xerrors.Errorf("path %s not exist", path)
		} else {
			blockRet = &b
		}
	}

	return blockRet, nil
}

func (cs *commandService) Register(name string) error {
	if err := cs.validRegister(name); err != nil {
		return xerrors.Errorf("error in validRegister: %w", err)
	}

	u, err := CreateUser(cs.root, name)
	if err != nil {
		return xerrors.Errorf("error in CreateUser: %w", err)
	}

	cs.userMap[name] = &u

	return nil
}

func (cs *commandService) Use(name string) error {
	if u, ok := cs.userMap[name]; ok {
		// hit cached in memory
		if err := AttemptUser(u.RootPath, u.Name); err != nil {
			return xerrors.Errorf("error in AttemptUser: %w", err)
		}

		b, ok := u.BlockMap[0]
		if !ok {
			return xerrors.Errorf("user not has root path")
		}

		cs.currentUser = u
		cs.currentBlock = &b

		return nil
	}

	u, err := GetUser(cs.root, name)
	if err != nil {
		return xerrors.Errorf("error in GetUser: %w", err)
	}

	// get root block
	b, ok := u.BlockMap[0]
	if !ok {
		return xerrors.Errorf("user not has root path")
	}

	cs.userMap[name] = &u
	cs.currentUser = &u
	cs.currentBlock = &b

	return nil
}

func (cs *commandService) ChangeFolder(path string) error {
	block, err := cs.travelFolder(path)
	if err != nil {
		return xerrors.Errorf("err in travelFolder: %w", err)
	}

	if block == nil {
		return xerrors.New("block is nil")
	}

	cs.currentBlock = block

	return nil
}

func (cs *commandService) CreateFolder(dirName string) error {
	if cs.currentBlock == nil {
		return xerrors.New("current block is nil")
	}

	dirName = strings.TrimSpace(dirName)
	if strings.Index(dirName, "/") != -1 {
		return xerrors.New("invalid directory name")
	}

	if _, ok := cs.currentBlock.FileMap[dirName]; ok {
		return xerrors.New("directory already exist")
	}

	block, err := CreateBlock(cs.currentBlock, cs.currentUser.CurrentNodeID+1)
	if err != nil {
		return xerrors.Errorf("create folder: %w", err)
	}

	_, err = CreateFolder(cs.currentBlock, block.NodeID, dirName, "dir")
	if err != nil {
		return xerrors.Errorf("create folder: %w", err)
	}

	cs.currentUser.BlockMap[block.NodeID] = block
	cs.currentUser.CurrentNodeID++
	cs.currentUser.Save()

	return nil
}

func (cs *commandService) DeleteFolder(oldName string) error {
	header, ok := cs.currentBlock.FileMap[oldName]
	if !ok {
		return xerrors.New("directory not exist")
	}

	if header.Type != Directory || header.DirNodeID == nil {
		return xerrors.New("not a directory")
	}

	if err := os.Remove(cs.currentBlock.GetBlockPath() + "/" + header.HashFileName); err != nil {
		return xerrors.Errorf("err in os.Remove: %w", err)
	}

	dirBloc := BlockINode{UserPath: cs.currentUser.GetUserPath(), NodeID: *header.DirNodeID}
	if err := os.RemoveAll(dirBloc.GetBlockPath()); err != nil {
		return xerrors.Errorf("err in os.Remove: %w", err)
	}

	delete(cs.currentBlock.FileMap, oldName)
	delete(cs.currentUser.BlockMap, *header.DirNodeID)
	cs.currentUser.BlockMap[cs.currentBlock.NodeID] = *cs.currentBlock

	if err := cs.currentBlock.Save(); err != nil {
		return xerrors.Errorf("err in currentBlock.Save: %w", err)
	}

	if err := cs.currentUser.Save(); err != nil {
		return xerrors.Errorf("err in currentUser.Save: %w", err)
	}

	return nil
}

func (cs *commandService) RenameFolder(oldName string, newName string) error {
	header, ok := cs.currentBlock.FileMap[oldName]
	if !ok {
		return xerrors.New("directory not exist")
	}

	if header.Type != Directory || header.DirNodeID == nil {
		return xerrors.New("not a directory")
	}

	header.Name = newName
	if err := header.Save(cs.currentBlock.GetBlockPath()); err != nil {
		return xerrors.Errorf("err in header.Save: %w", err)
	}

	delete(cs.currentBlock.FileMap, oldName)
	delete(cs.currentUser.BlockMap, *header.DirNodeID)

	cs.currentBlock.FileMap[newName] = header
	if err := cs.currentBlock.Save(); err != nil {
		return xerrors.Errorf("err in currentBlock.Save: %w", err)
	}

	cs.currentUser.BlockMap[cs.currentBlock.NodeID] = *cs.currentBlock
	if err := cs.currentUser.Save(); err != nil {
		return xerrors.Errorf("err in currentUser.Save: %w", err)
	}

	return nil
}

func (cs *commandService) CreateFile(fileName, desc string) error {
	if cs.currentBlock == nil {
		return xerrors.New("block is nil")
	}

	if _, ok := cs.currentBlock.FileMap[fileName]; ok {
		return xerrors.New("file already exist")
	}

	file, err := CreateFile(cs.currentBlock, fileName, desc)
	if err != nil {
		return xerrors.Errorf("err in CreateFile: %w", err)
	}

	cs.currentBlock.FileMap[fileName] = file
	cs.currentUser.BlockMap[cs.currentBlock.NodeID] = *cs.currentBlock

	return nil
}

func (cs *commandService) DeleteFile(oldName string) error {
	header, ok := cs.currentBlock.FileMap[oldName]
	if !ok {
		return xerrors.New("file not exist")
	}

	if header.Type != File {
		return xerrors.New("not a file")
	}

	if err := os.Remove(cs.currentBlock.GetBlockPath() + "/" + header.HashFileName); err != nil {
		return xerrors.Errorf("err in os.Remove: %w", err)
	}

	delete(cs.currentBlock.FileMap, oldName)
	cs.currentUser.BlockMap[cs.currentBlock.NodeID] = *cs.currentBlock

	if err := cs.currentBlock.Save(); err != nil {
		return xerrors.Errorf("err in currentBlock.Save: %w", err)
	}

	if err := cs.currentUser.Save(); err != nil {
		return xerrors.Errorf("err in currentUser.Save: %w", err)
	}

	return nil
}

func (cs *commandService) RenameFile(oldName string, newName string, newDesc string) error {
	header, ok := cs.currentBlock.FileMap[oldName]
	if !ok {
		return xerrors.New("file not exist")
	}

	if header.Type != File {
		return xerrors.New("not a file")
	}

	header.Name = newName
	header.Description = newDesc

	delete(cs.currentBlock.FileMap, oldName)
	cs.currentBlock.FileMap[newName] = header
	cs.currentUser.BlockMap[cs.currentBlock.NodeID] = *cs.currentBlock

	if err := header.Save(cs.currentBlock.GetBlockPath()); err != nil {
		return xerrors.Errorf("err in header.Save: %w", err)
	}

	if err := cs.currentBlock.Save(); err != nil {
		return xerrors.Errorf("err in currentBlock.Save: %w", err)
	}

	if err := cs.currentUser.Save(); err != nil {
		return xerrors.Errorf("err in currentUser.Save: %w", err)
	}

	return nil
}

func (cs *commandService) List(dirName string) ([]string, error) {
	block, err := cs.travelFolder(dirName)
	if err != nil {
		return nil, xerrors.Errorf("err in travelFolder: %w", err)
	}

	ret := make([]string, 0, len(cs.currentBlock.FileMap))
	for fname, file := range block.FileMap {
		if file.Type == Directory {
			fname += "/"
		}

		ret = append(ret, fname)
	}

	return ret, nil
}
