package vfsgo

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

// block test
func TestCreateBlock(t *testing.T) {
	path, err := getProjRoot()
	if err != nil {
		t.Error(err.Error())
		return
	}

	prevBlock := BlockINode{
		UserPath:   path + "/testdata/fileplayground",
		PrevNodeID: 0,
		NodeID:     0,
		FileMap:    make(map[string]FileHeader),
	}

	block, err := CreateBlock(&prevBlock, 2)
	if err != nil {
		t.Error(err.Error())
		return
	}

	log.Println(block.GetBlockPath())
	if _, err := os.Stat(block.GetBlockPath()); err != nil {
		t.Error(err.Error())
		return
	}

	if _, err := os.Stat(block.GetBlockINodePath()); err != nil {
		t.Error(err.Error())
		return
	}

	file, err := os.Open(block.GetBlockINodePath())
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		t.Error(err.Error())
		return
	}

	var data BlockINode
	if err := json.Unmarshal(b, &data); err != nil {
		t.Error(err.Error())
		return
	}

	if err := os.RemoveAll(block.GetBlockPath()); err != nil {
		t.Error(err.Error())
		return
	}
}

func TestGetBlock(t *testing.T) {
	path, err := getProjRoot()
	if err != nil {
		t.Error(err.Error())
		return
	}

	block, err := GetBlock(path+"/testdata/fileplayground", 3)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if block.NodeID != 3 {
		t.Error("block id error")
		return
	}

	if block.UserPath != path+"/testdata/fileplayground" {
		t.Error("block userpath error")
		return
	}

	// add file map validate
}

func TestDeleteBlock(t *testing.T) {
	path, err := getProjRoot()
	if err != nil {
		t.Error(err.Error())
		return
	}
	var id uint64 = 4

	block := BlockINode{
		UserPath: path,
		NodeID:   id,
		FileMap:  make(map[string]FileHeader),
	}

	if _, err := os.Stat(block.UserPath); err != nil {
		t.Error(err.Error())
		return
	}

	if _, err := os.Stat(block.GetBlockPath()); err == nil {
		t.Error(err.Error())
		return
	}

	// create block folder
	if err := os.Mkdir(block.GetBlockPath(), 0755); err != nil {
		t.Error(err.Error())
		return
	}

	// create block inode info
	if err := block.Save(); err != nil {
		t.Error(err.Error())
		return
	}

	if err := DelteBlock(&block); err != nil {
		t.Error(err.Error())
		return
	}

	if _, err := os.Stat(block.GetBlockPath()); err == nil {
		t.Error("block dir not deleted")
		return
	}
}
