package files


import (
	"os"
	"fmt"
	"path/filepath"
)

/*


1. Get the respective OS directory and create a common Tasks folder to save their input tasks
3. Save input tasks with date time stamp
4. Method to check if folder exists in path. If yes then update the file in there. If not creaete a new folder to ensure gracefullness 
5. End of the day delete all task files

*/

const folderName = "Tasks"





//Function to start a fresh new todo list 
func GetOrCreateTaskFile(){
	/*
	1. Check if folder exists -> yes, proceed to next check | no, create new folder and skip next check 
	2. Check if file exists -> Look for substring for current date and if it exists then just return that file | no, create a new file and append current date to it
	*/


}


//function to remove todo list 
func DeleteFile(){

}



//Helper function to assist with file operations 

func createFolder(){
	// user home directory regardless of OS
	homeDir, err:=os.UserHomeDir()

	if err != nil {
		fmt.Println("Error getting home directory: ", err)
		return
	}
	
	//Create a new folder in home dir to save tasks in markdown file
	filePath := filepath.Join(homeDir, "Documents", folderName)

	file, err :=os.Create(filePath)
	
	if err != nil {
		fmt.Println("Error creating folder:", err)
		return 
	}

	defer file.Close()//Close file operation even if error thrown

	fmt.Println("File successfully creted at: ", filePath)
}

