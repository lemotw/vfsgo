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

const (
	BlockINodeFileName = ".blockINode"
)

type BlockINode struct {
	// UserPath: user data pool path in real file system
	UserPath string `json:"user_path"`

	NodeID uint64 `json:"node_id"`
	// FileMap: file name -> file hash name
	FileMap map[string]*FileHeader `json:"file_map"`
}

func (b *BlockINode) GetBlockPath() string {
	return b.UserPath + "/" + strconv.FormatUint(b.NodeID, 10)
}

func (b *BlockINode) GetBlockINodePath() string {
	return b.UserPath + "/" + strconv.FormatUint(b.NodeID, 10) + "/" + BlockINodeFileName
}

func (b *BlockINode) Save() error {
	if _, err := os.Stat(b.GetBlockPath()); err != nil {
		return err
	}

	file, err := os.OpenFile(b.GetBlockINodePath(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	buf, err := json.Marshal(b)
	if err != nil {
		return err
	}

	if _, err := file.Write(buf); err != nil {
		return err
	}

	return nil
}

type FileType int8

const (
	Directory FileType = iota + 1
	File
)

type FileHeader struct {
	HashFileName string    `json:"hash_file_name"`
	Type         FileType  `json:"type"`
	DirNodeID    *uint64   `json:"dir_node_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	CreatedTime  time.Time `json:"created_time"`
	ModifiedTime time.Time `json:"modified_time"`
	// could added in the future
	// Content    []byte
	// Permission int32
	// Owner      string
	// ...
}

func (f *FileHeader) Save(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	buf, err := json.Marshal(f)
	if err != nil {
		return err
	}

	if _, err := file.Write(buf); err != nil {
		return err
	}

	return nil
}

func CreateBlock(path string, id uint64) (BlockINode, error) {
	if _, err := os.Stat(path); err != nil {
		return BlockINode{}, err
	}

	block := BlockINode{
		UserPath: path,
		NodeID:   id,
		FileMap:  make(map[string]*FileHeader),
	}

	// create block folder
	if err := os.Mkdir(block.GetBlockPath(), 0755); err != nil {
		return BlockINode{}, err
	}

	// create block inode info
	if err := block.Save(); err != nil {
		return BlockINode{}, err
	}

	return block, nil
}

func GetBlock(path string, id uint64) (BlockINode, error) {
	block := BlockINode{
		UserPath: path,
		NodeID:   id,
	}

	file, err := os.Open(block.GetBlockINodePath())
	if err != nil {
		return BlockINode{}, err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return BlockINode{}, err
	}

	if err := json.Unmarshal(b, &block); err != nil {
		return BlockINode{}, err
	}

	return block, nil
}

func DelteBlock(block *BlockINode) error {
	// remove folder
	if err := os.RemoveAll(block.GetBlockPath()); err != nil {
		return err
	}

	return nil
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

	if _, err := os.Stat(block.GetBlockPath()); err != nil {
		return FileHeader{}, err
	}

	if _, err := os.Stat(block.GetBlockPath() + "/" + filenameInFS); err == nil {
		return FileHeader{}, errors.New("file already exist")
	}

	file, err := os.Create(block.GetBlockPath() + "/" + filenameInFS)
	if err != nil {
		return FileHeader{}, err
	}
	defer file.Close()

	now := time.Now()
	header := FileHeader{
		HashFileName: filenameInFS,
		Type:         File,
		DirNodeID:    nil,
		Name:         filename,
		Description:  filedescription,
		CreatedTime:  now,
		ModifiedTime: now,
	}

	if err := header.Save(block.GetBlockPath() + "/" + filenameInFS); err != nil {
		return FileHeader{}, err
	}

	return header, nil
}

func GetFile(block *BlockINode, filename string) (FileHeader, error) {
	fileheader, ok := block.FileMap[filename]
	if !ok {
		return FileHeader{}, errors.New("file not found")
	}

	file, err := os.Open(block.GetBlockPath() + "/" + fileheader.HashFileName)
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

	file, err := os.OpenFile(block.GetBlockPath()+"/"+header.HashFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := header.Save(block.GetBlockPath() + "/" + header.HashFileName); err != nil {
		return err
	}

	return nil
}

func DeleteFile(block *BlockINode, filename string) error {
	header, ok := block.FileMap[filename]
	if !ok {
		return errors.New("file not found")
	}

	if err := os.Remove(block.GetBlockPath() + "/" + header.HashFileName); err != nil {
		return err
	}

	return nil
}
