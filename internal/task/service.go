package task

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

//Sentinel error in Go
var ErrTaskNotFound = errors.New("task not found!")

// CRUD functionality for tasks
//Add task
func (tasks TaskStore) AddTask(task string){
	lastIdx := getLastId(tasks)

	newTask := Task{
		ID: lastIdx+1,
		Text: task,
		Done: false,
	}
	
	tasks[newTask.ID] = &newTask //match the id key with the task id for easier retrieval
}

/*
----- Edit task option ---------
Functional options pattern to pass optional edit parameters that can be modified 
*/
type EditOption func(*Task)

func WithText(text string) EditOption{
	return func (t *Task)  {
		t.Text = text
	}
}

func WithCompletionStatus(done bool) EditOption {
	return func (t *Task){
		t.Done = done
	}
}

func (tasks TaskStore)EditTask(id int, opts ...EditOption)(*Task){
	task, err := tasks.getTask(id)

	if errors.Is(err, ErrTaskNotFound){
		fmt.Print(err)
		return nil
	}

	for _, opt := range opts{
		opt(task) //each option mutates task in place
	}

	return task
}

// helper method to get the task and return error if no task found with that id
func (tasks TaskStore) getTask(number int)(*Task, error){
	task, ok := tasks[number]

	if !ok {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

// ----------- Delete service --------------- // 

//method to delete a task completely from list. (not the same as marking as done)
 func (tasks TaskStore) DeleteTask(taskId int){
	delete(tasks, taskId) //remove the task from TaskStore 
 }

// --------- Services that interact with File data ----------- //

func (store TaskStore) getAllTasks()(TaskStore){
	return store
}



//Get tasks from md file (Need to parse it first)
func getSavedTasks(file *os.File) ([]string, error) {
	//Reset pointer in markdown file to the first line to ensure all tasks are picked up
	_, err := file.Seek(0,0) //rewind to start
	if err != nil {
		return nil, err
	}
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(content)),"\n"),nil
}

//Function to check the status of the tasks written in the markdown file
func checkStatusInFile(line string) (bool){
	prefix := "- [x] "
	//Check if 3rd character of the line contains 'x'
	if (strings.HasPrefix(line, prefix)){
		return true
	}else{
		return false
	}
}

//Helper method to strip the checkboxes from string
func stripCheckboxes(line string)(string){
	task := line[6:]
	return task
}

//function to map markdown tasks to TaskStore map and Task struct
//When calling this function remember to delete any other existing instance of TaskStore present to prevent bugs and inconsistency
func (store TaskStore) MapMdToStore(file *os.File)(TaskStore,error){
	 tasks,err:= getSavedTasks(file)
	 if err != nil {
		return nil, err
	 }

	 for i,task := range tasks {
		done := checkStatusInFile(task)
		cleaned_task := stripCheckboxes(task)

		store[i+1] = &Task{
			ID: i + 1, 
			Text: cleaned_task,
			Done: done,
		}
		
	 }
	 return store, nil
}


func getLastId(store TaskStore) int{
	max := 0
	for id := range store{
		if id > max {
			max = id
		}
	}
	return max
}


//map tasks from map to md file (unordered)
//When method is called it will map the items to markdown file and 
func (tasks TaskStore) UpdateDoc(file *os.File)error{
	fmt.Printf("Updating tasks in file: %s\r\n " , file.Name())
	//Prepare the file for the tasks

	//1. Truncate the file (clear all its contents)
	if err := file.Truncate(0); err != nil {
		return err
	}

	//2. Move the cursor back to the first line in file
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	//3. Open a new buffer to store all the lines and write it one shot to the file when the buffer is full
	writer := bufio.NewWriter(file)
	
	/*
	- Loop through the map and write the item to the file
	- Want markdown formats: so check if status is updated to done. If yes, then update the string with the markdown done
	to reflect the change in the markdown file itself
	*/
	for _, currTask := range tasks {
		//check and format string
		taskToWrite := currTask.PrepareTask()

		//Write each task on a new line
		if _, err := writer.WriteString(taskToWrite + "\n"); err != nil {
			return err
		}
	}

	//write everything in buffer to file 
	if err := writer.Flush(); err != nil {
		return err
	}

	fmt.Print("\r\nMarkdown file successfully updated!\r\n")
	return nil
}

//Method to prepare task data to be written onto markdown
func (task Task) PrepareTask() (string){
	if task.Done {
		return "- [x] " + task.Text
	}
	return "- [ ] " + task.Text
}
