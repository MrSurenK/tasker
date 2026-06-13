package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/MrSurenK/tasker/internal/files"
	"github.com/MrSurenK/tasker/internal/terminal"
)



func main() {

	// ---- User input to start application flow ------ // 
	//Put terminal in raw mode 
	raw_mode := terminal.SetUpTerminal()
	defer terminal.Restore(raw_mode)

	
	//Welcome msg
	terminal.ShowWelcome()

	//Get user confirmation to start making to do list
	if !terminal.StartToDo(){
		return
	}

	/* 
	----- Set up environment to add to do list ------
	*/
	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Fatalf("Unable to find your home directory in your environment! %v\r\n", err) //Log error and exit the program 
	}

	requiredPath := path.Join(homeDir,"Documents","ToDoList")
	
	//Initialize file path into file serivice struct
	fileservice := files.NewFileService(requiredPath)

	file, err := fileservice.GetOrCreateFile()

	if err != nil {
		log.Fatalf("Encountered error: %v", err)
	}

	fmt.Println(file.Name())

	/*
	
	---- Perform crud operations on todolist---

	
	*/

}