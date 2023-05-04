package vfsgo

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"golang.org/x/xerrors"
)

func getUserRoot() (string, error) {
	path, err := getProjRoot()
	if err != nil {
		return "", err
	}

	return path + "/testdata/user", nil
}

func getUserINode(path string) (User, error) {
	if _, err := os.Stat(path); err != nil {
		return User{}, err
	}

	file, err := os.Open(path)
	if err != nil {
		return User{}, err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return User{}, err
	}

	var data User
	if err := json.Unmarshal(b, &data); err != nil {
		return User{}, err
	}
	return data, nil
}

func TestAttemptUser(t *testing.T) {
	rootPath, err := getUserRoot()
	if err != nil {
		t.Error(err.Error())
		return
	}

	if err := AttemptUser(rootPath, "test"); err != nil {
		t.Error(err.Error())
		return
	}

	if err := AttemptUser(rootPath, "wqenwqkjdncsaucqbweqwejkqnk"); err == nil {
		t.Error("user attempt failed")
		return
	}
}

func TestCreateUser(t *testing.T) {
	root, err := getUserRoot()
	if err != nil {
		t.Error(err.Error())
		return
	}

	user, err := CreateUser(root, "test_create")
	if err != nil {
		t.Error(err.Error())
		return
	}

	if _, err := os.Stat(user.GetUserPath()); err != nil {
		t.Error(err.Error())
		return
	}

	userFromFile, err := getUserINode(user.GetUserINodePath())
	if err != nil {
		t.Error(err.Error())
		return
	}

	if userFromFile.Name != "test_create" {
		t.Error("userFromFile.Name != \"test_create\"")
		return
	}

	if userFromFile.Name != user.Name {
		t.Error("userFromFile.Name != user.Name")
		return
	}

	if err := os.RemoveAll(user.GetUserPath()); err != nil {
		t.Error(err.Error())
		return
	}
}

func TestGetUser(t *testing.T) {
	root, err := getUserRoot()
	if err != nil {
		t.Error(err.Error())
		return
	}

	user, err := GetUser(root, "test_get")
	if err != nil {
		t.Error(err.Error())
		return
	}

	userFromFile, err := getUserINode(user.GetUserINodePath())
	if err != nil {
		t.Error(err.Error())
		return
	}

	if userFromFile.Name != user.Name {
		t.Error("userFromFile.Name != user.Name")
		return
	}
}

func TestDeleteUser(t *testing.T) {
	root, err := getUserRoot()
	if err != nil {
		t.Error(err.Error())
		return
	}

	deleteUser := User{
		RootPath:      root,
		Name:          "test_delete",
		CurrentNodeID: 0,
		BlockMap:      make(map[uint64]BlockINode),
		CreatedTime:   time.Now(),
	}

	// create delete file case
	if _, err := os.Stat(deleteUser.GetUserPath()); err != nil {
		err := func() error {
			// create block folder
			if err := os.Mkdir(deleteUser.GetUserPath(), 0755); err != nil {
				return xerrors.Errorf("error in os.Mkdir: %w", err)
			}

			if err := deleteUser.Save(); err != nil {
				return err
			}

			return nil
		}()
		if err != nil {
			t.Error(err.Error())
			return
		}
	}

	if err := DeleteUser(root, deleteUser.Name); err != nil {
		t.Error(err.Error())
		return
	}

	if _, err := os.Stat(deleteUser.GetUserPath()); err == nil {
		t.Error("delete user failed")
		return
	}
}
