package model

import "time"

type UserINode struct {
	// Name: unique username
	Name         string
	RootBlockID  int64
	CreatedTime  time.Time
	ModifiedTime time.Time
}
