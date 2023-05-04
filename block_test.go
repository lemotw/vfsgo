package vfsgo

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"golang.org/x/xerrors"
)

func getUser() (*User, error) {
	projRoot, err := getProjRoot()
	if err != nil {
		return nil, err
	}

	user := User{
		RootPath:      projRoot + "/testdata/block",
		Name:          "user1",
		CurrentNodeID: 0,
		BlockMap:      make(map[uint64]BlockINode),
		CreatedTime:   time.Now(),
	}

	file, err := os.Open(user.GetUserINodePath())
	if err != nil {
		return nil, xerrors.Errorf("error in os.Open: %w", err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, xerrors.Errorf("error in ioutil.ReadAll: %w", err)
	}

	if err := json.Unmarshal(b, &user); err != nil {
		return nil, xerrors.Errorf("error in json.Unmarshal: %w", err)
	}

	return &user, nil
}

func getBlock(path string) (BlockINode, error) {
	if _, err := os.Stat(path); err != nil {
		return BlockINode{}, err
	}

	file, err := os.Open(path)
	if err != nil {
		return BlockINode{}, err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return BlockINode{}, err
	}

	var data BlockINode
	if err := json.Unmarshal(b, &data); err != nil {
		return BlockINode{}, err
	}
	return data, nil
}

// block test
func TestCreateBlock(t *testing.T) {
	user, err := getUser()
	if err != nil {
		t.Error(err.Error())
		return
	}

	block, err := CreateBlock(&BlockINode{
		UserPath: user.GetUserPath(),
		NodeID:   0,
	}, 1)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if _, err := os.Stat(block.GetBlockPath()); err != nil {
		t.Error(err.Error())
		return
	}

	if err := os.RemoveAll(block.GetBlockPath()); err != nil {
		t.Error(err.Error())
		return
	}
}

func TestGetBlock(t *testing.T) {
	user, err := getUser()
	if err != nil {
		t.Error(err.Error())
		return
	}

	block, err := GetBlock(user, 0)
	if err != nil {
		t.Error(err.Error())
		return
	}

	blockFromFile, err := getBlock(block.GetBlockINodePath())
	if err != nil {
		t.Error(err.Error())
		return
	}

	if block.PrevNodeID != blockFromFile.PrevNodeID {
		t.Error("block.PrevNodeID != blockFromFile.PrevNodeID")
		return
	}

	if block.NodeID != blockFromFile.NodeID {
		t.Error("block.NodeID != blockFromFile.NodeID")
		return
	}

	if block.UserPath != blockFromFile.UserPath {
		t.Error("block.UserPath != blockFromFile.UserPath")
		return
	}
}

func TestDeleteBlock(t *testing.T) {
	user, err := getUser()
	if err != nil {
		t.Error(err.Error())
		return
	}

	deleteBlock := BlockINode{
		UserPath:   user.GetUserPath(),
		PrevNodeID: 0,
		NodeID:     99,
		FileMap:    make(map[string]FileHeader),
	}

	// create delete file case
	if _, err := os.Stat(user.GetUserPath() + "/block/99"); err != nil {
		err := func() error {
			// create block folder
			if err := os.Mkdir(deleteBlock.GetBlockPath(), 0755); err != nil {
				return xerrors.Errorf("error in os.Mkdir: %w", err)
			}

			if err := deleteBlock.Save(); err != nil {
				return err
			}

			user.BlockMap[deleteBlock.NodeID] = deleteBlock

			return nil
		}()
		if err != nil {
			t.Error(err.Error())
			return
		}
	}

	if err := DeleteBlock(user, 99); err != nil {
		t.Error(err.Error())
		return
	}

	if _, err := os.Stat(deleteBlock.GetBlockPath()); err == nil {
		t.Error("delete file failed")
		return
	}
}
