package vfsgo

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {
	path, err := getProjRoot()
	if err != nil {
		t.Error(err.Error())
		return
	}

	user, err := CreateUser(path+"/testdata/userplayground", "test")
	if err != nil {
		t.Error(err.Error())
		return
	}

	if _, err := os.Stat(user.GetUserPath()); err != nil {
		t.Error(err.Error())
		return
	}

	if _, err := os.Stat(user.GetUserINodePath()); err != nil {
		t.Error(err.Error())
		return
	}

	file, err := os.Open(user.GetUserINodePath())
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

	var data User
	if err := json.Unmarshal(b, &data); err != nil {
		t.Error(err.Error())
		return
	}

	if data.Name != "test" {
		t.Error("name not equal")
		return
	}

	if err := os.RemoveAll(data.GetUserPath()); err != nil {
		t.Error(err.Error())
		return
	}
}

func TestGetUser(t *testing.T) {
	path, err := getProjRoot()
	if err != nil {
		t.Error(err.Error())
		return
	}

	user, err := GetUser(path+"/testdata/userplayground", "testexist")
	if err != nil {
		t.Error(err.Error())
		return
	}

	file, err := os.Open(user.GetUserINodePath())
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

	var fileUser User
	if err := json.Unmarshal(b, &fileUser); err != nil {
		t.Error(err.Error())
		return
	}

	if user.Name != fileUser.Name {
		t.Error("name not equal")
		return
	}
}

func TestDeleteUser(t *testing.T) {
	path, err := getProjRoot()
	if err != nil {
		t.Error(err.Error())
		return
	}

	user := User{
		RootPath:    path + "/testdata/userplayground",
		Name:        "deleteuser",
		CreatedTime: time.Now(),
	}

	// create block folder
	if err := os.Mkdir(user.RootPath+"/"+user.Name, 0755); err != nil {
		t.Error(err.Error())
		return
	}

	if err := user.Save(); err != nil {
		t.Error(err.Error())
		return
	}

	if err := DeleteUser(user.RootPath, user.Name); err != nil {
		t.Error(err.Error())
		return
	}

	if _, err := os.Stat(user.RootPath + "/" + user.Name); err == nil {
		t.Error("user dir not deleted")
		return
	}
}
