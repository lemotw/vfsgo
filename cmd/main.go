package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lemotw/vfsgo"
)

const FSROOTPATH = "/fs"

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

func sendCMD(serv vfsgo.ICommandService, command string) bool {
	cmdSplitSlice := strings.Split(strings.TrimSpace(command), " ")

	cmdSlice := make([]string, 0, len(cmdSplitSlice))
	for i := 0; i < len(cmdSplitSlice); i++ {
		if cmdSplitSlice[i] == "" {
			continue
		}
		cmdSlice = append(cmdSlice, cmdSplitSlice[i])
	}

	// check commandß
	if len(cmdSlice) <= 0 {
		return false
	}

	switch cmdSlice[0] {
	case "register":
		if len(cmdSlice) != 2 {
			log.Println("register command format: register username password")
			return false
		}
		log.Println("exec: register")
		if err := serv.Register(cmdSlice[1]); err != nil {
			log.Println(err.Error())
		} else {
			log.Println(fmt.Sprintf("Add [%s] successfully.", cmdSlice[1]))
		}
	case "use":
		if len(cmdSlice) != 2 {
			log.Println("use command format: use username")
			return false
		}
		log.Println("exec: use")
		if err := serv.Use(cmdSlice[1]); err != nil {
			log.Println(err.Error())
		} else {
			log.Println(fmt.Sprintf("You are using [%s].", cmdSlice[1]))
		}
	case "create-folder":
		if len(cmdSlice) != 2 {
			log.Println("create-folder command format: create-folder foldername path")
			return false
		}
		log.Println("exec: create-folder")
		if err := serv.CreateFolder(cmdSlice[1]); err != nil {
			log.Println(err.Error())
		} else {
			log.Println(fmt.Sprintf("Create [%s] successfully.", cmdSlice[1]))
		}
	case "delete-folder":
		if len(cmdSlice) != 2 {
			log.Println("delete-folder command format: delete-folder foldername path")
			return false
		}
		log.Println("exec: delete-folder")
		if err := serv.DeleteFolder(cmdSlice[1]); err != nil {
			log.Println(err.Error())
		} else {
			log.Println(fmt.Sprintf("Delete [%s] successfully.", cmdSlice[1]))
		}
	case "cd":
		if len(cmdSlice) != 2 {
			log.Println("cd command format: cd path")
			return false
		}
		log.Println("exec: cd")
		if err := serv.ChangeFolder(cmdSlice[1]); err != nil {
			log.Println(err.Error())
		}
	case "ls":

		var sortField *vfsgo.SortType
		var sortOrder *string

		if !(len(cmdSlice) == 2 || len(cmdSlice) == 4) {
			log.Println("ls command format: ls path [sortField sortOrder]")
			return false
		}

		if len(cmdSlice) == 4 {
			if cmdSlice[2] == "--sort-name" {
				f := vfsgo.SortByName
				sortField = &f
			}

			if cmdSlice[2] == "--sort-created" {
				f := vfsgo.SortByCreatedTime
				sortField = &f
			}

			if cmdSlice[3] == "asc" {
				order := vfsgo.ASC
				sortOrder = &order
			}

			if cmdSlice[3] == "desc" {
				order := vfsgo.DESC
				sortOrder = &order
			}
		}

		log.Println("exec: ls")
		if files, err := serv.List(cmdSlice[1], sortField, sortOrder); err != nil {
			log.Println(err.Error())
		} else {
			log.Println("files: ")
			for _, file := range files {
				log.Println(file)
			}
		}
	case "rename-folder":
		if len(cmdSlice) != 3 {
			log.Println("rename-folder command format: rename-folder oldname newname path")
			return false
		}
		log.Println("exec: rename-folder")
		if err := serv.RenameFolder(cmdSlice[1], cmdSlice[2]); err != nil {
			log.Println(err.Error())
		} else {
			log.Println(fmt.Sprintf("Rename [%s] to [%s] successfully.", cmdSlice[1], cmdSlice[2]))
		}
	case "create-file":
		if len(cmdSlice) != 3 {
			log.Println("create-file command format: create-file filename path")
			return false
		}
		log.Println("exec: create-file")
		if err := serv.CreateFile(cmdSlice[1], cmdSlice[2]); err != nil {
			log.Println(err.Error())
		} else {
			log.Println(fmt.Sprintf("Create [%s] successfully.", cmdSlice[1]))
		}
	case "delete-file":
		if len(cmdSlice) != 2 {
			log.Println("delete-file command format: delete-file filename path")
			return false
		}
		log.Println("exec: delete-file")
		if err := serv.DeleteFile(cmdSlice[1]); err != nil {
			log.Println(err.Error())
		} else {
			log.Println(fmt.Sprintf("Delete [%s] successfully.", cmdSlice[1]))
		}
	case "exit":
		return true
	}

	return false
}

func main() {
	path, err := getProjRoot()
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(os.Stdin)
	serv := vfsgo.NewCommandService(path + FSROOTPATH)

	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		if sendCMD(serv, command) {
			log.Println("goodbye!! ")
			break
		}
	}
}
