package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lemotw/vfsgo"
)

const FSROOTPATH = "/cmdplayground"

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

func parseCommand(cmd string) (string, []string) {
	command := ""
	commands := strings.Split(cmd, " ")

	args := make([]string, len(commands)-1, 0)
	for i := 0; i < len(commands); i++ {
		if commands[i] == "" {
			continue
		}

		if command == "" {
			command = commands[i]
			continue
		}

		args = append(args, commands[i])
	}

	return commands[0], args
}

func sendCommand(serv vfsgo.ICommandService, command string, args []string) error {
	return nil
}

func main() {
	var command string

	path, err := getProjRoot()
	if err != nil {
		panic(err)
	}

	cmdService := vfsgo.NewCommandService(path + FSROOTPATH)

	for {
		if _, err := fmt.Scan(&command); err != nil {
			panic(err)
		}

		cmd, args := parseCommand(command)
		if err := sendCommand(cmdService, cmd, args); err != nil {
			// std err
			log.Println(err)
		}
	}
}
