package vfsgo

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"golang.org/x/xerrors"
)

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
}

func (f *FileHeader) Save(path string) error {
	file, err := os.OpenFile(path+"/"+f.HashFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return xerrors.Errorf("error in os.OpenFile: %w", err)
	}
	defer file.Close()

	buf, err := json.Marshal(f)
	if err != nil {
		return xerrors.Errorf("error in json.Marshal: %w", err)
	}

	if _, err := file.Write(buf); err != nil {
		return xerrors.Errorf("error in file.Write: %w", err)
	}

	return nil
}

// randHash: sha256 with random
func randHash() (string, error) {
	data := make([]byte, 16)
	if _, err := rand.Read(data); err != nil {
		return "", xerrors.Errorf("error in rand.Read: %w", err)
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

func CreateFolder(block *BlockINode, nodeid uint64, foldername, desc string) (FileHeader, error) {
	if _, ok := block.FileMap[foldername]; ok {
		return FileHeader{}, xerrors.New("file already exist")
	}

	filenameInFS, err := randHash()
	if err != nil {
		return FileHeader{}, xerrors.Errorf("error in randHash: %w", err)
	}

	if _, err := os.Stat(block.GetBlockPath()); err != nil {
		return FileHeader{}, xerrors.New("block path not exis")
	}

	file, err := os.Create(block.GetBlockPath() + "/" + filenameInFS)
	if err != nil {
		return FileHeader{}, xerrors.Errorf("error in os.Create: %w", err)
	}
	defer file.Close()

	now := time.Now()
	header := FileHeader{
		HashFileName: filenameInFS,
		Type:         Directory,
		DirNodeID:    &nodeid,
		Name:         foldername,
		Description:  desc,
		CreatedTime:  now,
		ModifiedTime: now,
	}

	if err := header.Save(block.GetBlockPath()); err != nil {
		return FileHeader{}, xerrors.Errorf("error in header.Save: %w", err)
	}

	block.FileMap[foldername] = header
	if err := block.Save(); err != nil {
		return FileHeader{}, xerrors.Errorf("error in block.Save: %w", err)
	}

	return header, nil
}

func CreateFile(block *BlockINode, filename, filedescription string) (FileHeader, error) {
	if _, ok := block.FileMap[filename]; ok {
		return FileHeader{}, xerrors.New("file already exist")
	}

	filenameInFS, err := randHash()
	if err != nil {
		return FileHeader{}, xerrors.Errorf("error in randHash: %w", err)
	}

	if _, err := os.Stat(block.GetBlockPath()); err != nil {
		return FileHeader{}, xerrors.New("block path not exis")
	}

	file, err := os.Create(block.GetBlockPath() + "/" + filenameInFS)
	if err != nil {
		return FileHeader{}, xerrors.Errorf("error in os.Create: %w", err)
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

	if err := header.Save(block.GetBlockPath()); err != nil {
		return FileHeader{}, xerrors.Errorf("error in header.Save: %w", err)
	}

	block.FileMap[filename] = header
	if err := block.Save(); err != nil {
		return FileHeader{}, xerrors.Errorf("error in block.Save: %w", err)
	}

	return header, nil
}

func GetFile(block *BlockINode, filename string) (FileHeader, error) {
	fileheader, ok := block.FileMap[filename]
	if !ok {
		return FileHeader{}, xerrors.New("file not found")
	}

	file, err := os.Open(block.GetBlockPath() + "/" + fileheader.HashFileName)
	if err != nil {
		return FileHeader{}, xerrors.Errorf("error in os.Open: %w", err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return FileHeader{}, xerrors.Errorf("error in ioutil.ReadAll: %w", err)
	}

	var data FileHeader
	if err := json.Unmarshal(b, &data); err != nil {
		return FileHeader{}, xerrors.Errorf("error in json.Unmarshal: %w", err)
	}

	return data, nil
}

func UpdateFile(block *BlockINode, filename, filedescription string) error {
	header, ok := block.FileMap[filename]
	if !ok {
		return xerrors.New("file not found")
	}

	header.Description = filedescription

	if err := header.Save(block.GetBlockPath()); err != nil {
		return xerrors.Errorf("error in header.Save: %w", err)
	}

	return nil
}

func DeleteFile(block *BlockINode, filename string) error {
	header, ok := block.FileMap[filename]
	if !ok {
		return xerrors.New("file not found")
	}

	if err := os.Remove(block.GetBlockPath() + "/" + header.HashFileName); err != nil {
		return xerrors.Errorf("error in os.Remove: %w", err)
	}

	delete(block.FileMap, filename)
	if err := block.Save(); err != nil {
		return xerrors.Errorf("error in block.Save: %w", err)
	}

	return nil
}
