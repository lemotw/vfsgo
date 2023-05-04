package vfsgo

import (
	"encoding/json"
	"os"
	"time"
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
