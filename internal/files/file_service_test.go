package files

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

/*
Test for generating file name
*/
func TestFileNameGenerator(t *testing.T ){

	currDate := time.Now()
	date := currDate.Format("2006-01-02")
	expected_fileName := "ToDoList-"+date
	generated_fileName := generateFileName(currDate)

	if generated_fileName != expected_fileName{
		t.Errorf("Expected %v , but got %v", expected_fileName, generated_fileName)
	}
}

func TestFileCreationAndRetrieval(t *testing.T){
	//Create isolated temp directory
	tempDir := t.TempDir()

	//create a service instance 
	service := &FileService{
		root: tempDir,
	}

	//call method under test 
	file, err := service.GetOrCreateFile()

	//Check error
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	//Check file is not null 
	if file == nil {
		t.Fatal("expected file, got nil")
	}

	//verify file actually exists on disk
	info, err := os.Stat(file.Name())

	if err != nil {
		t.Fatalf("expected file to exist on disk, got error: %v", err)
	}

	//ensure its actually a file and not a directory 
	if info.IsDir(){
		t.Fatal("expected a file but got a directory")
	}
}


func TestCleanUp(t *testing.T){
	//isolated file system 
	tempDir := t.TempDir()

	//Create service
	service := FileService{
		root: tempDir,
	}

	//create Test files 
	file1, err := os.Create(filepath.Join(tempDir,"file1.txt"))
	if err != nil {
		t.Fatal(err)
	}

	file1.Close()

	file2, err := os.Create(filepath.Join(tempDir,"file2.txt"))

	if err != nil {
		t.Fatal(err)
	}

	file2.Close()

	//create a folder (should not be deleted)
	err = os.Mkdir(filepath.Join(tempDir, "folder1"), 0755)
	if err != nil{
		t.Fatal(err)
	}

	//run clean up 
	err = service.CleanUp(tempDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	//check that files are deleted 
	if _, err := os.Stat(filepath.Join(tempDir, "file1.txt")); !os.IsNotExist(err){

		t.Fatal("expected file1.txt to be deleted")
	}

	if _, err := os.Stat(filepath.Join(tempDir, "file2.txt")); !os.IsNotExist(err){

		t.Fatal("expected file2.txt to be deleted")
	}

	//check folder still exists
	if _, err := os.Stat(filepath.Join(tempDir, "folder1")); err != nil {
		t.Fatal("expected file2.txt to be deleted")
	}

}







