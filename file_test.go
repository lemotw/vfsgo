package vfsgo

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func getProjRoot() (string, error) {
	projRoot, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(projRoot, "go.mod")); err == nil {
			break
		}

		if projRoot == filepath.Dir(projRoot) {
			return "", errors.New("can't find proj root")
		}

		projRoot = filepath.Dir(projRoot)
	}

	return projRoot, nil
}

func getHeader(path string) (FileHeader, error) {
	if _, err := os.Stat(path); err != nil {
		return FileHeader{}, err
	}

	file, err := os.Open(path)
	if err != nil {
		return FileHeader{}, err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return FileHeader{}, err
	}

	var data FileHeader
	if err := json.Unmarshal(b, &data); err != nil {
		return FileHeader{}, err
	}
	return data, nil
}

func TestCreateFile(t *testing.T) {
	// get package root
	projRoot, err := getProjRoot()
	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Log("userpath: " + projRoot + "/testdata/fileplayground")
	block := BlockINode{
		UserPath: projRoot + "/testdata/fileplayground",
		NodeID:   1,
		FileMap: map[string]*FileHeader{
			"existfile": {
				HashFileName: "existfile",
				Type:         File,
				DirNodeID:    nil,
				Description:  "this is desc",
				Name:         "existfile",
				CreatedTime:  time.Now(),
				ModifiedTime: time.Now(),
			},
		},
	}

	fileheader, err := CreateFile(&block, "testfile", "testfile desc")
	if err != nil {
		t.Error(err.Error())
		return
	}

	header, err := getHeader(block.GetPath() + "/" + fileheader.HashFileName)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if header.HashFileName != fileheader.HashFileName {
		t.Error("fileheader not equal")
		return
	}

	if header.Name != "testfile" {
		t.Error("fileheader name not equal")
		return
	}

	if header.Description != "testfile desc" {
		t.Error("fileheader desc not equal")
		return
	}
}

func TestGetFile(t *testing.T) {
	projRoot, err := getProjRoot()
	if err != nil {
		t.Error(err.Error())
		return
	}

	block := BlockINode{
		UserPath: projRoot + "/testdata/fileplayground",
		NodeID:   1,
		FileMap: map[string]*FileHeader{
			"fileexist": {
				HashFileName: "fileexist",
				Type:         File,
				DirNodeID:    nil,
				Name:         "fileexist",
				CreatedTime:  time.Now(),
				ModifiedTime: time.Now(),
			},
		},
	}

	fileHeader, err := GetFile(&block, "fileexist")
	if err != nil {
		t.Error(err.Error())
		return
	}

	header, err := getHeader(block.GetPath() + "/fileexist")
	if err != nil {
		t.Error(err.Error())
		return
	}

	if header.HashFileName != fileHeader.HashFileName {
		t.Error("get file header failed")
		return
	}

	if header.Type != fileHeader.Type {
		t.Error("get file header failed")
		return
	}

	if header.DirNodeID != fileHeader.DirNodeID {
		t.Error("get file header failed")
		return
	}

	if header.Name != fileHeader.Name {
		t.Error("get file header failed")
		return
	}

	if header.Description != fileHeader.Description {
		t.Error("get file header failed")
		return
	}

	if header.CreatedTime != fileHeader.CreatedTime {
		t.Error("get file header failed")
		return
	}

	if header.ModifiedTime != fileHeader.ModifiedTime {
		t.Error("get file header failed")
		return
	}
}

func TestUpdateFile(t *testing.T) {
	projRoot, err := getProjRoot()
	if err != nil {
		t.Error(err.Error())
		return
	}

	block := BlockINode{
		UserPath: projRoot + "/testdata/fileplayground",
		NodeID:   1,
		FileMap: map[string]*FileHeader{
			"fileexist": {
				HashFileName: "fileexist",
				Type:         File,
				DirNodeID:    nil,
				Name:         "fileexist",
				CreatedTime:  time.Now(),
				ModifiedTime: time.Now(),
			},
		},
	}

	updateDesc := "this is desc"
	updateDesc1 := "this is desc1"
	if err := UpdateFile(&block, "fileexist", updateDesc); err != nil {
		t.Error(err.Error())
		return
	}

	header, err := getHeader(block.GetPath() + "/fileexist")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if header.Description != updateDesc {
		t.Error("update file desc failed")
		return
	}

	if err := UpdateFile(&block, "fileexist", updateDesc1); err != nil {
		t.Error(err.Error())
		return
	}

	header, err = getHeader(block.GetPath() + "/fileexist")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if header.Description != updateDesc1 {
		t.Error("update file desc failed")
		return
	}
}

func TestDeleteFile(t *testing.T) {
	projRoot, err := getProjRoot()
	if err != nil {
		t.Error(err.Error())
		return
	}

	block := BlockINode{
		UserPath: projRoot + "/testdata/fileplayground",
		NodeID:   1,
		FileMap: map[string]*FileHeader{
			"deleteone": {
				HashFileName: "deleteone",
				Type:         File,
				DirNodeID:    nil,
				Name:         "deleteone",
				CreatedTime:  time.Now(),
				ModifiedTime: time.Now(),
			},
		},
	}

	deleteonePath := block.GetPath() + "/deleteone"
	if _, err := os.Stat(deleteonePath); err != nil {
		err := func() error {
			// prepare delete one
			header := FileHeader{
				HashFileName: "deleteone",
				Type:         File,
				DirNodeID:    nil,
				Description:  "delete one",
				Name:         "deleteone",
				CreatedTime:  time.Now(),
				ModifiedTime: time.Now(),
			}
			file, err := os.Create(deleteonePath)

			if err != nil {
				return err
			}
			defer file.Close()

			b, err := json.Marshal(header)
			if err != nil {
				return err
			}

			if _, err := file.Write(b); err != nil {
				return err
			}
			return nil
		}()
		if err != nil {
			t.Error(err.Error())
			return
		}
	}

	if err := DeleteFile(&block, "deleteone"); err != nil {
		t.Error(err.Error())
		return
	}
}
