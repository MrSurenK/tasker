package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/MrSurenK/tasker/internal/files"
)

func main() {

	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Fatalf("Unable to find your home directory in your environment! %v\n", err) //Log error and exit the program 
	}

	requiredPath := path.Join(homeDir,"Documents","ToDoList")
	
	//Initialize file path into file serivice struct
	fileservice := files.NewFileService(requiredPath)

	file, err := fileservice.GetOrCreateFile()

	if err != nil {
		log.Fatalf("Encountered error: %v", err)
	}

	fmt.Println(file.Name())
}