package model

import "time"

type BlockINode struct {
	NodeID int64
	// file name -> file hash name
	FileMap map[string]string
}

type FileType int8

const (
	Directory FileType = iota + 1
	File
)

type FileHeader struct {
	Type         FileType
	Name         string
	CreatedTime  time.Time
	ModifiedTime time.Time

	// could added in the future
	// Content    []byte
	// Permission int32
	// Owner      string
	// ...
}
