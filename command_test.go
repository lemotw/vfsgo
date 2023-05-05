package vfsgo

import (
	"os"
	"testing"
)

func getCmdService() (ICommandService, error) {
	root, err := getProjRoot()
	if err != nil {
		return nil, err
	}

	return NewCommandService(root + "/testdata/cmd"), nil
}

func TestRegister(t *testing.T) {
	root, err := getProjRoot()
	if err != nil {
		t.Error(err.Error())
		return
	}
	cmdService := NewCommandService(root + "/testdata/cmd")

	registCase := "testRegisterCase"
	if err := cmdService.Register(registCase); err != nil {
		t.Error(err.Error())
		return
	}

	userFromFile, err := GetUser(root+"/testdata/cmd", registCase)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if userFromFile.Name != registCase {
		t.Error("userFromFile.Name != user.Name")
		return
	}

	if err := os.RemoveAll(userFromFile.GetUserPath()); err != nil {
		t.Error(err.Error())
		return
	}
}

func TestUse(t *testing.T) {
	root, err := getProjRoot()
	if err != nil {
		t.Error(err.Error())
		return
	}
	cmdRoot := root + "/testdata/cmd"

	cmdService := NewCommandService(cmdRoot)

	if err := cmdService.Register("testUse"); err != nil {
		t.Error(err.Error())
		return
	}
	defer func() {
		if err := os.RemoveAll(cmdService.GetCurrentUser().GetUserPath()); err != nil {
			t.Error(err.Error())
			return
		}
	}()

	if err := cmdService.Use("testUse"); err != nil {
		t.Error(err.Error())
		return
	}
	curUser := cmdService.GetCurrentUser()

	if curUser.Name != cmdService.GetCurrentUser().Name {
		t.Error("curUser.Name != userPreCreate.Name")
		return
	}
}

func TestCreateFolder(t *testing.T) {
	cmdService, err := getCmdService()
	if err != nil {
		t.Error(err.Error())
		return
	}

	if err := cmdService.Register("testCreateFolder"); err != nil {
		t.Error(err.Error())
		return
	}
	defer func() {
		if err := os.RemoveAll(cmdService.GetCurrentUser().GetUserPath()); err != nil {
			t.Error(err.Error())
			return
		}
	}()

	if err := cmdService.Use("testCreateFolder"); err != nil {
		t.Error(err.Error())
		return
	}

	if err := cmdService.CreateFolder("cFolder"); err != nil {
		t.Error(err.Error())
		return
	}

	createFolderUser := cmdService.GetCurrentUser()
	rootBlock, ok := createFolderUser.BlockMap[0]
	if !ok {
		t.Error("rootBlock not found")
		return
	}

	headerInMem, ok := rootBlock.FileMap["cFolder"]
	if !ok {
		t.Error("headerInMem not found")
		return
	}

	headerInFile, err := getHeader(rootBlock.GetBlockPath() + "/" + headerInMem.HashFileName)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if headerInMem.HashFileName != headerInFile.HashFileName {
		t.Error("headerInMem.HashFileName != headerInFile.HashFileName")
		return
	}

	if headerInMem.Type != Directory || headerInMem.Type != headerInFile.Type {
		t.Error("headerInMem.Type != Directory || headerInMem.Type != headerInFile.Type")
		return
	}

	if headerInMem.Name != "cFolder" || headerInMem.Name != headerInFile.Name {
		t.Error("headerInMem.Name != \"cFolder\" || headerInMem.Name != headerInFile.Name")
		return
	}
}

func TestDeleteFolder(t *testing.T) {
	cmdService, err := getCmdService()
	if err != nil {
		t.Error(err.Error())
		return
	}

	if err := cmdService.Register("testDeleteFolder"); err != nil {
		t.Error(err.Error())
		return
	}
	defer func() {
		if err := os.RemoveAll(cmdService.GetCurrentUser().GetUserPath()); err != nil {
			t.Error(err.Error())
			return
		}
	}()

	if err := cmdService.Use("testDeleteFolder"); err != nil {
		t.Error(err.Error())
		return
	}

	if err := cmdService.CreateFolder("dFolder"); err != nil {
		t.Error(err.Error())
		return
	}

	deleteFolderUser := cmdService.GetCurrentUser()
	rootBlock, ok := deleteFolderUser.BlockMap[0]
	if !ok {
		t.Error("rootBlock not found")
		return
	}

	headerInMem, ok := rootBlock.FileMap["dFolder"]
	if !ok {
		t.Error("headerInMem not found")
		return
	}

	if _, err := os.Stat(rootBlock.GetBlockPath() + "/" + headerInMem.HashFileName); err != nil {
		t.Error(err.Error())
		return
	}

	if err := cmdService.DeleteFolder("dFolder"); err != nil {
		t.Error(err.Error())
		return
	}

	if _, err := os.Stat(deleteFolderUser.GetUserPath() + "/" + headerInMem.HashFileName); err == nil {
		t.Error(err.Error())
		return
	}
}

func TestRenameFolder(t *testing.T) {
	cmdService, err := getCmdService()
	if err != nil {
		t.Error(err.Error())
		return
	}

	if err := cmdService.Register("testRenameFolder"); err != nil {
		t.Error(err.Error())
		return
	}
	defer func() {
		if err := os.RemoveAll(cmdService.GetCurrentUser().GetUserPath()); err != nil {
			t.Error(err.Error())
			return
		}
	}()

	if err := cmdService.Use("testRenameFolder"); err != nil {
		t.Error(err.Error())
		return
	}

	if err := cmdService.CreateFolder("rnFolder"); err != nil {
		t.Error(err.Error())
		return
	}

	renameFolderUser := cmdService.GetCurrentUser()
	rootBlock, ok := renameFolderUser.BlockMap[0]
	if !ok {
		t.Error("rootBlock not found")
		return
	}

	headerInMem, ok := rootBlock.FileMap["rnFolder"]
	if !ok {
		t.Error("headerInMem not found")
		return
	}

	if err := cmdService.RenameFolder("rnFolder", "rnFolder2"); err != nil {
		t.Error(err.Error())
		return
	}

	headerInFile, err := getHeader(rootBlock.GetBlockPath() + "/" + headerInMem.HashFileName)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if headerInFile.Name != "rnFolder2" {
		t.Error("headerInFile.Name != \"rnFolder2\"")
		return
	}
}

func TestCreateFileCMD(t *testing.T) {
	cmdService, err := getCmdService()
	if err != nil {
		t.Error(err.Error())
		return
	}

	if err := cmdService.Register("testCreateFile"); err != nil {
		t.Error(err.Error())
		return
	}
	defer func() {
		if err := os.RemoveAll(cmdService.GetCurrentUser().GetUserPath()); err != nil {
			t.Error(err.Error())
			return
		}
	}()

	if err := cmdService.Use("testCreateFile"); err != nil {
		t.Error(err.Error())
		return
	}

	if err := cmdService.CreateFile(".", "createdFile", "this is created file desc"); err != nil {
		t.Error(err.Error())
		return
	}

	// check mem block and header
	createFileUser := cmdService.GetCurrentUser()
	rootBlock, ok := createFileUser.BlockMap[0]
	if !ok {
		t.Error("rootBlock not found")
		return
	}

	headerInMem, ok := rootBlock.FileMap["createdFile"]
	if !ok {
		t.Error("createdFile not found")
		return
	}

	if headerInMem.Name != "createdFile" {
		t.Error("headerInMem.Name != \"createdFile\"")
		return
	}

	if headerInMem.Type != File {
		t.Error("headerInMem.Type != File")
		return
	}

	if headerInMem.Description != "this is created file desc" {
		t.Error("headerInMem.Description != \"this is created file desc\"")
		return
	}

	// check in file
	headerInFile, err := GetFile(&rootBlock, "createdFile")
	if err != nil {
		t.Error(err.Error())
		return
	}

	if headerInFile.Name != "createdFile" {
		t.Error("headerInFile.Name != \"createdFile\"")
		return
	}

	if headerInFile.Type != File {
		t.Error("headerInFile.Type != File")
		return
	}

	if headerInFile.Description != "this is created file desc" {
		t.Error("headerInFile.Description != \"this is created file desc\"")
		return
	}
}

func TestChangeFolder(t *testing.T) {
	cmdService, err := getCmdService()
	if err != nil {
		t.Error(err.Error())
		return
	}

	if err := cmdService.Register("testChangeFolder"); err != nil {
		t.Error(err.Error())
		return
	}
	defer func() {
		if err := os.RemoveAll(cmdService.GetCurrentUser().GetUserPath()); err != nil {
			t.Error(err.Error())
			return
		}
	}()

	if err := cmdService.Use("testChangeFolder"); err != nil {
		t.Error(err.Error())
		return
	}

	if err := cmdService.CreateFolder("testFolder"); err != nil {
		t.Error(err.Error())
		return
	}

	if err := cmdService.ChangeFolder("testFolder"); err != nil {
		t.Error(err.Error())
		return
	}

	curentBlock := cmdService.GetCurrentBlock()
	prevBlock, err := GetBlock(cmdService.GetCurrentUser(), curentBlock.PrevNodeID)
	if err != nil {
		t.Error(err.Error())
		return
	}

	curHeader, ok := prevBlock.FileMap["testFolder"]
	if !ok {
		t.Error("curHeader not found")
		return
	}

	if curHeader.Type != Directory {
		t.Error("curHeader.Type != Directory")
		return
	}

	if curHeader.DirNodeID == nil {
		t.Error("curHeader.DirNodeID != nil")
		return
	}

	if *curHeader.DirNodeID != cmdService.GetCurrentBlock().NodeID {
		t.Error("curHeader.DirNodeID !=  &cmdService.GetCurrentBlock().NodeID()")
		return
	}
}

func TestList(t *testing.T) {
	cmdService, err := getCmdService()
	if err != nil {
		t.Error(err.Error())
		return
	}

	if err := cmdService.Register("testChangeFolder"); err != nil {
		t.Error(err.Error())
		return
	}
	defer func() {
		if err := os.RemoveAll(cmdService.GetCurrentUser().GetUserPath()); err != nil {
			t.Error(err.Error())
			return
		}
	}()

	if err := cmdService.Use("testChangeFolder"); err != nil {
		t.Error(err.Error())
		return
	}

	if err := cmdService.CreateFolder("testFolder"); err != nil {
		t.Error(err.Error())
		return
	}

	if err := cmdService.CreateFolder("testFolder1"); err != nil {
		t.Error(err.Error())
		return
	}

	if err := cmdService.CreateFolder("testFolder2"); err != nil {
		t.Error(err.Error())
		return
	}

	if err := cmdService.CreateFile(".", "testListFile", "thisis the file Desc"); err != nil {
		t.Error(err.Error())
		return
	}

	files, err := cmdService.List(".")
	if err != nil {
		t.Error(err.Error())
		return
	}

	if len(files) != 4 {
		t.Error("len(files) != 4")
		return
	}

	// testFolder1s
	f := false
	for i := 0; i < len(files); i++ {
		if files[i] == "testFolder/" {
			f = true
		}
	}
	if !f {
		t.Error("testFolder not found")
		return
	}

	// testFolder1
	f = false
	for i := 0; i < len(files); i++ {
		if files[i] == "testFolder1/" {
			f = true
		}
	}
	if !f {
		t.Error("testFolder1 not found")
		return
	}

	// testFolder2
	f = false
	for i := 0; i < len(files); i++ {
		if files[i] == "testFolder2/" {
			f = true
		}
	}
	if !f {
		t.Error("testFolder2 not found")
		return
	}

	// testListFile
	f = false
	for i := 0; i < len(files); i++ {
		if files[i] == "testListFile" {
			f = true
		}
	}
	if !f {
		t.Error("testListFile not found")
		return
	}
}
