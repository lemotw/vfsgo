package vfsgo

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"
)

const (
	UserINodeFileName = ".userINode"
)

type User struct {
	RootPath string `json:"root_path"`
	Name     string `json:"name"`

	BlockMap map[uint64]*BlockINode `json:"block_map"`

	CreatedTime time.Time `json:"created_time"`
}

func (u *User) GetUserPath() string {
	return u.RootPath + "/" + u.Name
}

func (u *User) GetUserINodePath() string {
	return u.GetUserPath() + "/" + UserINodeFileName
}

func (u *User) Save() error {
	if _, err := os.Stat(u.GetUserPath()); err != nil {
		return err
	}

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

func CreateUser(rootPath, name string) (User, error) {
	if _, err := os.Stat(rootPath); err != nil {
		return User{}, err
	}

	now := time.Now()
	user := User{
		RootPath:    rootPath,
		Name:        name,
		CreatedTime: now,
	}

	if _, err := os.Stat(user.GetUserPath()); err == nil {
		return User{}, errors.New("user already exist")
	}

	if err := os.Mkdir(user.GetUserPath(), 0755); err != nil {
		return User{}, err
	}

	if err := user.Save(); err != nil {
		return User{}, err
	}

	return user, nil
}

func AttemptUser(rootPath, name string) error {
	user := User{
		RootPath: rootPath,
		Name:     name,
	}

	if _, err := os.Stat(user.GetUserPath()); err != nil {
		return err
	}

	return nil
}

func GetUser(rootPath, name string) (User, error) {
	user := User{
		RootPath: rootPath,
		Name:     name,
	}

	file, err := os.Open(user.GetUserINodePath())
	if err != nil {
		return User{}, err
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return User{}, err
	}

	if err := json.Unmarshal(buf, &user); err != nil {
		return User{}, err
	}

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
