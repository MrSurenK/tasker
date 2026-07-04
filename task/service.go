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


//FIXME: When getting task and placing into Task object strip away the markdown formatting properly
//Get tasks from md file (Need to parse it first)
func GetSavedTasks(file *os.File) ([]string, error) {
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(content)),"\n"),nil
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

// TODO: Pull out the markdown processing logic to its own method 

//map tasks from map to md file (unordered)
//When method is called it will map the items to markdown file and 
func (tasks TaskStore) UpdateDoc(file *os.File)error{
	fmt.Printf("Updating tasks in file: %s " , file.Name())
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
		taskToWrite := ""

		if(currTask.Done){
			taskToWrite = "-[x] " + currTask.Text //prepend markdown checkbox with check
		} else {
			taskToWrite = "-[ ] " + currTask.Text //prepend markdown checkbox without check
		}
		
		//Write each task on a new line
		if _, err := writer.WriteString(taskToWrite + "\n"); err != nil {
			return err
		}
	}

	//write everything in buffer to file 
	if err := writer.Flush(); err != nil {
		return err
	}

	fmt.Print("Markdown file successfully updated!")
	return nil
}
