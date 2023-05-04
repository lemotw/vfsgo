package vfsgo

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"golang.org/x/xerrors"
)

const (
	UserINodeFileName = ".userINode"
)

type User struct {
	RootPath      string `json:"root_path"`
	Name          string `json:"name"`
	CurrentNodeID uint64 `json:"current_node_id"`

	BlockMap map[uint64]BlockINode `json:"block_map"`

	CreatedTime time.Time `json:"created_time"`
}

func (u *User) GetUserPath() string {
	return u.RootPath + "/" + u.Name
}

func (u *User) GetUserINodePath() string {
	return u.GetUserPath() + "/" + UserINodeFileName
}

func (u *User) Save() error {
	file, err := os.OpenFile(u.GetUserINodePath(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	buf, err := json.Marshal(u)
	if err != nil {
		return err
	}

	if _, err := file.Write(buf); err != nil {
		return err
	}

	return nil
}

func AttemptUser(rootPath, name string) error {
	user := User{
		RootPath: rootPath,
		Name:     name,
	}

	if _, err := os.Stat(user.GetUserPath()); err != nil {
		return xerrors.Errorf("error in os.Stat: %w", err)
	}

	return nil
}

func CreateUser(rootPath, name string) (User, error) {
	if _, err := os.Stat(rootPath); err != nil {
		return User{}, xerrors.Errorf("error in os.Stat: %w", err)
	}

	user := User{
		RootPath:    rootPath,
		Name:        name,
		CreatedTime: time.Now(),
		BlockMap:    make(map[uint64]BlockINode),
	}

	if _, err := os.Stat(user.GetUserPath()); err == nil {
		return User{}, xerrors.New("user already exist")
	}

	if err := os.Mkdir(user.GetUserPath(), 0755); err != nil {
		return User{}, xerrors.Errorf("error in os.Mkdir: %w", err)
	}

	// create root block
	b, err := CreateBlock(&BlockINode{
		UserPath:   user.GetUserPath(),
		NodeID:     0,
		PrevNodeID: 0,
	}, 0)
	if err != nil {
		return User{}, xerrors.Errorf("error in CreateBlock: %w", err)
	}
	user.BlockMap[0] = b

	if err := user.Save(); err != nil {
		return User{}, err
	}

	return user, nil
}

func GetUser(rootPath, name string) (User, error) {
	user := User{
		RootPath: rootPath,
		Name:     name,
	}

	file, err := os.Open(user.GetUserINodePath())
	if err != nil {
		return User{}, xerrors.Errorf("error in os.Open: %w", err)
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return User{}, xerrors.Errorf("error in ioutil.ReadAll: %w", err)
	}

	if err := json.Unmarshal(buf, &user); err != nil {
		return User{}, xerrors.Errorf("error in json.Unmarshal: %w", err)
	}

	blocks, err := ioutil.ReadDir(user.GetUserPath())
	if err != nil {
		return User{}, xerrors.Errorf("error in ioutil.ReadDir: %w", err)
	}

	// refetch max block id
	maxid := uint64(0)
	for _, block := range blocks {
		bid, err := strconv.ParseUint(block.Name(), 10, 64)
		if err != nil {
			continue
		}
		if bid > maxid {
			maxid = uint64(bid)
		}
	}
	user.CurrentNodeID = maxid

	return user, nil
}

func DeleteUser(rootPath, name string) error {
	user := User{
		RootPath: rootPath,
		Name:     name,
	}

	if _, err := os.Stat(user.GetUserPath()); err != nil {
		return err
	}

	if err := os.RemoveAll(user.GetUserPath()); err != nil {
		return err
	}

	return nil
}
