package vfsgo

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type BlockINode struct {
	// UserPath: user data pool path in real file system
	UserPath string

	NodeID uint64
	// FileMap: file name -> file hash name
	FileMap map[string]*FileHeader
}

func (b *BlockINode) GetPath() string {
	return b.UserPath + "/" + strconv.FormatUint(b.NodeID, 10)
}

type FileType int8

const (
	Directory FileType = iota + 1
	File
)

type FileHeader struct {
	HashFileName string
	Type         FileType
	DirNodeID    *uint64
	Name         string
	Description  string
	CreatedTime  time.Time
	ModifiedTime time.Time

	// could added in the future
	// Content    []byte
	// Permission int32
	// Owner      string
	// ...
}

func CreateBlock() {
	panic("implement me")
}

func GetBlock() {
	panic("implement me")
}

func UpdateBlock() {
	panic("implement me")
}

func DelteBlock() {
	panic("implement me")
}

// randHash: sha256 with random
func randHash() (string, error) {
	data := make([]byte, 16)
	if _, err := rand.Read(data); err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

func CreateFile(block *BlockINode, filename, filedescription string) (FileHeader, error) {
	filenameInFS, err := randHash()
	if err != nil {
		return FileHeader{}, err
	}

	if _, err := os.Stat(block.GetPath()); err != nil {
		return FileHeader{}, err
	}

	if _, err := os.Stat(block.GetPath() + "/" + filenameInFS); err == nil {
		return FileHeader{}, errors.New("file already exist")
	}

	file, err := os.Create(block.GetPath() + "/" + filenameInFS)
	if err != nil {
		return FileHeader{}, err
	}
	defer file.Close()

	now := time.Now()
	fileheader := FileHeader{
		HashFileName: filenameInFS,
		Type:         File,
		DirNodeID:    nil,
		Name:         filename,
		Description:  filedescription,
		CreatedTime:  now,
		ModifiedTime: now,
	}

	b, err := json.Marshal(fileheader)
	if err != nil {
		return FileHeader{}, err
	}

	if _, err := file.Write(b); err != nil {
		return FileHeader{}, err
	}

	return fileheader, nil
}

func GetFile(block *BlockINode, filename string) (FileHeader, error) {
	fileheader, ok := block.FileMap[filename]
	if !ok {
		return FileHeader{}, errors.New("file not found")
	}

	if _, err := os.Stat(block.GetPath() + "/" + fileheader.HashFileName); err != nil {
		return FileHeader{}, err
	}

	file, err := os.Open(block.GetPath() + "/" + fileheader.HashFileName)
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

func UpdateFile(block *BlockINode, filename, filedescription string) error {
	header, ok := block.FileMap[filename]
	if !ok {
		return errors.New("file not found")
	}

	header.Description = filedescription

	file, err := os.OpenFile(block.GetPath()+"/"+header.HashFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	// 寫入數據
	b, err := json.Marshal(header)
	if err != nil {
		return err
	}

	if _, err := file.Write(b); err != nil {
		return err
	}

	return nil
}

func DeleteFile(block *BlockINode, filename string) error {
	header, ok := block.FileMap[filename]
	if !ok {
		return errors.New("file not found")
	}

	if _, err := os.Stat(block.GetPath() + "/" + header.HashFileName); err != nil {
		return err
	}

	if err := os.Remove(block.GetPath() + "/" + header.HashFileName); err != nil {
		return err
	}

	return nil
}
