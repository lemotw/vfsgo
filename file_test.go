package vfsgo

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/xerrors"
)

/*
testdata:
fileplayground:
	1: -> for file test
	2: -> for block create test
	3: -> for block get test
	4: -> for block delete test
*/

// file test

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

func getRootBlock() (*BlockINode, error) {
	projRoot, err := getProjRoot()
	if err != nil {
		return nil, err
	}

	block := BlockINode{
		UserPath: projRoot + "/testdata/file",
		NodeID:   0,
	}

	file, err := os.Open(block.GetBlockINodePath())
	if err != nil {
		return nil, xerrors.Errorf("error in os.Open: %w", err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, xerrors.Errorf("error in ioutil.ReadAll: %w", err)
	}

	if err := json.Unmarshal(b, &block); err != nil {
		return nil, xerrors.Errorf("error in json.Unmarshal: %w", err)
	}

	return &block, nil
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

func readFromFile(block *BlockINode, header FileHeader) (*FileHeader, error) {
	file, err := os.Open(block.GetBlockPath() + "/" + header.HashFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, xerrors.Errorf("error in ioutil.ReadAll: %w", err)
	}

	var data FileHeader
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, xerrors.Errorf("error in json.Unmarshal: %w", err)
	}

	return &data, nil
}

func deleteFile(block *BlockINode, header FileHeader) error {
	if err := os.Remove(block.GetBlockPath() + "/" + header.HashFileName); err != nil {
		return err
	}

	delete(block.FileMap, header.Name)
	block.Save()

	return nil
}
func TestCreateFile(t *testing.T) {
	// get package root
	block, err := getRootBlock()
	if err != nil {
		t.Error(err.Error())
		return
	}

	fn := "createfile_test"
	desc := "file desc at create file case"
	header, err := CreateFile(block, fn, desc)
	if err != nil {
		t.Error(err.Error())
		return
	}

	// validate return header data
	if header.Name != fn {
		t.Error("created file name not equal")
		return
	}

	if header.Description != desc {
		t.Error("created file desc not equal")
		return
	}

	headerFromFile, err := readFromFile(block, header)
	if err != nil {
		t.Error(err.Error())
		return
	}

	// validate file header data
	if headerFromFile.Name != fn {
		t.Error("created file name not equal")
		return
	}

	if headerFromFile.Description != desc {
		t.Error("created file desc not equal")
		return
	}

	// final process test data
	if err := deleteFile(block, header); err != nil {
		t.Error(err.Error())
		return
	}
}

func TestGetFile(t *testing.T) {
	block, err := getRootBlock()
	if err != nil {
		t.Error(err.Error())
		return
	}

	header, err := GetFile(block, "getfile_test")
	if err != nil {
		t.Error(err.Error())
		return
	}

	headerFromFile, err := readFromFile(block, header)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if header.HashFileName != headerFromFile.HashFileName {
		t.Error("get file header failed")
		return
	}

	if header.Type != headerFromFile.Type {
		t.Error("get file header failed")
		return
	}

	if header.DirNodeID != headerFromFile.DirNodeID {
		t.Error("get file header failed")
		return
	}

	if header.Name != headerFromFile.Name {
		t.Error("get file header failed")
		return
	}

	if header.Description != headerFromFile.Description {
		t.Error("get file header failed")
		return
	}

	if header.CreatedTime != headerFromFile.CreatedTime {
		t.Error("get file header failed")
		return
	}

	if header.ModifiedTime != headerFromFile.ModifiedTime {
		t.Error("get file header failed")
		return
	}
}

func TestUpdateFile(t *testing.T) {
	block, err := getRootBlock()
	if err != nil {
		t.Error(err.Error())
		return
	}

	oriHeader, ok := block.FileMap["updatefile_test"]
	if !ok {
		t.Error("update file test data not exist")
		return
	}

	header, err := readFromFile(block, oriHeader)
	if err != nil {
		t.Error(err.Error())
		return
	}

	originalDesc := header.Description
	updateDesc := "this is desc"
	updateDesc1 := "this is desc1"

	// update case 1
	if err := UpdateFile(block, header.Name, updateDesc); err != nil {
		t.Error(err.Error())
		return
	}

	updateHeader, err := readFromFile(block, *header)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if updateHeader.Description != updateDesc {
		t.Error("file desc been updated not equal")
		return
	}

	// update case 2
	if err := UpdateFile(block, header.Name, updateDesc1); err != nil {
		t.Error(err.Error())
		return
	}

	updateHeader1, err := readFromFile(block, *header)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if updateHeader1.Description != updateDesc1 {
		t.Error("file desc been updated not equal")
		return
	}

	// rollback test data
	if err := UpdateFile(block, header.Name, originalDesc); err != nil {
		t.Error(err.Error())
		return
	}

	rollbackHeader, err := readFromFile(block, *header)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if rollbackHeader.Description != originalDesc {
		t.Error("file desc been updated not equal")
		return
	}
}

func TestDeleteFile(t *testing.T) {
	block, err := getRootBlock()
	if err != nil {
		t.Error(err.Error())
		return
	}

	// create delete file case
	if _, err := os.Stat(block.GetBlockPath() + "/deleteone"); err != nil {
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
			file, err := os.Create(block.GetBlockPath() + "/deleteone")

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

			block.FileMap["deleteone"] = header

			return nil
		}()
		if err != nil {
			t.Error(err.Error())
			return
		}
	}

	if err := DeleteFile(block, "deleteone"); err != nil {
		t.Error(err.Error())
		return
	}

	if _, err := os.Stat(block.GetBlockPath() + "/deleteone"); err == nil {
		t.Error("delete file failed")
		return
	}
}
