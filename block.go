package vfsgo

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	"golang.org/x/xerrors"
)

const (
	BlockINodeFileName = ".blockINode"
)

type BlockINode struct {
	// UserPath: user data pool path in real file system
	UserPath string `json:"user_path"`

	PrevNodeID uint64 `json:"prev_node_id"`
	NodeID     uint64 `json:"node_id"`
	// FileMap: file name -> file hash name
	FileMap map[string]FileHeader `json:"file_map"`
}

func (b *BlockINode) GetBlockPath() string {
	return b.UserPath + "/" + strconv.FormatUint(b.NodeID, 10)
}

func (b *BlockINode) GetBlockINodePath() string {
	return b.UserPath + "/" + strconv.FormatUint(b.NodeID, 10) + "/" + BlockINodeFileName
}

func (b *BlockINode) Save() error {
	if _, err := os.Stat(b.GetBlockPath()); err != nil {
		return xerrors.Errorf("error in os.Stat: %w", os.Create)
	}

	file, err := os.OpenFile(b.GetBlockINodePath(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return xerrors.Errorf("error in os.OpenFile: %w", os.Create)
	}
	defer file.Close()

	buf, err := json.Marshal(b)
	if err != nil {
		return xerrors.Errorf("error in json.Marshal: %w", os.Create)
	}

	if _, err := file.Write(buf); err != nil {
		return xerrors.Errorf("error in file.Write: %w", os.Create)
	}

	return nil
}

func CreateBlock(block *BlockINode, id uint64) (BlockINode, error) {
	if id != 0 {
		_, err := os.Stat(block.GetBlockPath())
		if err != nil {
			return BlockINode{}, xerrors.New(("parent block path not exist"))
		}
	}

	newBlock := BlockINode{
		UserPath:   block.UserPath,
		PrevNodeID: block.NodeID,
		NodeID:     id,
		FileMap:    make(map[string]FileHeader),
	}

	// create block folder
	if err := os.Mkdir(newBlock.GetBlockPath(), 0755); err != nil {
		return BlockINode{}, xerrors.Errorf("error in os.Mkdir: %w", os.Create)
	}

	// create block inode info
	if err := newBlock.Save(); err != nil {
		return BlockINode{}, xerrors.Errorf("error in block.Save: %w", os.Create)
	}

	return newBlock, nil
}

func GetBlock(path string, id uint64) (BlockINode, error) {
	block := BlockINode{
		UserPath: path,
		NodeID:   id,
	}

	file, err := os.Open(block.GetBlockINodePath())
	if err != nil {
		return BlockINode{}, xerrors.Errorf("error in os.Open: %w", os.Create)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return BlockINode{}, xerrors.Errorf("error in ioutil.ReadAll: %w", os.Create)
	}

	if err := json.Unmarshal(b, &block); err != nil {
		return BlockINode{}, xerrors.Errorf("error in json.Unmarshal: %w", os.Create)
	}

	return block, nil
}

func DelteBlock(block *BlockINode) error {
	// remove folder
	if err := os.RemoveAll(block.GetBlockPath()); err != nil {
		return xerrors.Errorf("error in os.RemoveAll: %w", os.Create)
	}

	return nil
}
