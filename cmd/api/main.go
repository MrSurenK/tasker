package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/MrSurenK/tasker/internal/files"
	"github.com/MrSurenK/tasker/internal/task"
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

	defer file.Close() //close file at the end of programe


	// ------------- Prepare existing tasks for applications to use ------------- // 

	//Read file and check if empty
	info, err := file.Stat()
	if err != nil{
		log.Fatalf("Could not stat file: %v", err)
	}

	//set up global struct store for app to use
	store := make(task.TaskStore)

	if info.Size() == 0 {
		os.Stdout.Write([]byte("\r\nLooks like we have no tasks for the day. Let's add some!\r\n"))
	}else{
		store.MapMdToStore(file)
		for id, t := range store {
			fmt.Printf("%d. %s\r\n", id, t.PrepareTask())
		}
	}

	//--------------------- Main application loop ---------------------//

	for {
		choice, ok := terminal.TellMeWhatToDo()

		if !ok {
			return //quit - defers handle file.Close() and terminal.Restore()
		}

		switch choice {
		case '1': //Add task API 
			text := terminal.ReadLine("Enter your tasks:")
			if text == ""{
				continue
			}
			store.AddTask(text)
		case '2':
			//get the task id of task that user wants to edit
			taskid := terminal.ReadLine("Enter the task id that you want to edit:")
			id, err := strconv.Atoi(taskid)
			if err != nil {
				os.Stdout.Write([]byte("\r\nInvalid number\r\n"))
				continue
			}
			//Sub menu after picking a task
			os.Stdout.Write([]byte("\r\nWhat do you want to update?\r\n1.Text\r\n2.Mark as done/undone\r\n3.Both\r\n"))
			key := make([]byte, 8)
			os.Stdin.Read(key) //Read stdin into key
			switch key[0]{
			case '1': //update text of task
				newText := terminal.ReadLine("New Text:")
				if newText == "" {
					continue //invalid input. User should not input empty task
				}
				store.EditTask(id, task.WithText(newText))
			case '2': //Update task completion status
				doneStr := terminal.ReadLine("Mark as done? (y/n):")
				done := len(doneStr) > 0 && (doneStr[0] == 'y' || doneStr[0] == 'Y')
				store.EditTask(id, task.WithCompletionStatus(done))
			
			case '3': //update both task and completion status
				newText := terminal.ReadLine("New Text:")
				if newText == ""{
					continue
				}
				doneStr := terminal.ReadLine("Mark as done? (y/n):")
				done := len(doneStr) > 0 && (doneStr[0] == 'Y' || doneStr[0] == 'y')
				store.EditTask(id, task.WithText(newText), task.WithCompletionStatus(done))
			default:
				os.Stdout.Write([]byte("\r\nInvalid option\r\n"))
			}
		//Back to main menu
		case '3':
			if len(store) == 0{
				os.Stdout.Write([]byte("\r\nNothing to delete"))
				continue
			}
			showTasks(store)
			taskId := terminal.ReadLine("Which task would you like to delete?")
			id, err := strconv.Atoi(taskId)
			if err != nil {
				os.Stdout.Write([]byte("\r\nInvalid number\r\n"))
				continue
			}
			store.DeleteTask(id)
			os.Stdout.Write([]byte("\r\nTask deleted!\r\n"))
		case '4':
			showTasks(store)
		case '5':
			//write store to markdown file and close app
			store.UpdateDoc(file)
			os.Stdout.Write([]byte("\r\nFile saved sucessfully!"))
			return
		}
	}
}


//NOTE: Named function has to live outside of Main in golang//
// -- Helper method to show TaskStore object --- // 
func showTasks(store task.TaskStore){
	if len(store) == 0 {
		os.Stdout.Write([]byte("\r\nNo tasks at the moment\r\n"))
	}
	for id, t := range store {
		fmt.Printf("%d. %s\r\n", id, t.PrepareTask())
	}
}