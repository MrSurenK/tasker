package files

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

/*
1. Get the respective OS directory and create a common Tasks folder to save their input tasks
3. Save input tasks with date time stamp
4. Method to check if folder exists in path. If yes then update the file in there. If not creaete a new folder to ensure gracefullness
5. End of the day delete all task files
*/

//In main method must specify file path
type FileService struct {
	root string 
	filename string
}

//constructor to populate FileService struct
func NewFileService(root string) *FileService{
	return &FileService{
		root: root,
	}
}

//Function to start a fresh new todo list 
	/*
	1. Check if folder exists -> yes, proceed to next check | no, create new folder and skip next check 
	2. Check if file exists -> Look for substring for current date and if it exists then just return that file | no, create a new file and append current date to it
	*/

	/*
	Helper function to generate a dynamic todo list file name for the day
	*/
func generateFileName(now time.Time) string{
	//File name generated should follow the following pattern 
	// ToDoList-DateStamp
	date := now.Format("2006-01-02")
	fileName := "ToDoList-" + date
	return fileName
}

/*
Method to get today's todo list or create a new one all in one simple clean method
Caller has the reponsibility to give the root path to save todo list
*/
func (f *FileService) GetOrCreateFile()(*os.File, error){
	//Create file
	//today date 
	date := time.Now()
	f.filename = generateFileName(date) + ".md" //get today's file name 
	path := filepath.Join(f.root,f.filename) //creates OS specific paths

	//Check for existance and log it
	exists:= true

	if _, err := os.Stat(path); err != nil {
		if(os.IsNotExist(err)){
			exists = false
		} else {
			return nil, err
		}
	}

	if exists {
		log.Printf("Using existing todo file: %s" , f.filename)
	}else{
		log.Printf("Creating new todoFile in this path: %s", path)
	}


	//if path does not exists then createa the path too
	err := os.MkdirAll(f.root, 0755)
	if err != nil {
		return nil, err 
	}

	//if file does not exists it will create the file in the specified path
	file, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE, 0666) //if file already exists then it will just return that file and not truncate it (wipe out)

	if err != nil {
		return nil , err
	}

	log.Printf("Retrieved todo file successfully @ %s", path)
	return file, nil 
}

/*
Function to clear all todolists -- if user says exit then do it, else leave the todolist. 
!!! Do not perform a clean up by default when program exits !!! -- User has to ask for it or if its a new day then run it to clear all the old to do lists (ask user if they want to as well)
*/
func (f FileService)CleanUp(dir string) error{

	entries, err := os.ReadDir(f.root) //get all files in my root dir 

	if(err != nil){
		return err
	}

	for _, entry := range entries{
		if entry.IsDir(){
			continue //if entry is not an actual file skip it. Skips folders
		}

		path := filepath.Join(f.root, entry.Name())

		if err := os.Remove(path); err != nil {
			return err
		}
	}
	return nil
}






