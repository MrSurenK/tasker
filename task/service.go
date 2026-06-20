package task

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/MrSurenK/tasker/internal/terminal"
	"github.com/MrSurenK/tasker/task"
)

//Sentinel error in Go
var ErrTaskNotFound = errors.New("task not found!")

// CRUD functionality for tasks
//Add task and throw error
func (tasks TaskStore) AddTask(task string){
	lastIdx := getLastId(tasks)

	newTask := Task{
		ID: lastIdx+1,
		Text: task,
		Done: false,
	}
	
	tasks[lastIdx + 1] = &newTask
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

func withCompletionStatus(done bool) EditOption {
	return func (t *Task){
		t.Done = done
	}
}

func (tasks TaskStore)  EditTask(id int, opts ...EditOption)(*Task){
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


// ------- End of Edit Task Service ----------- //








func (store TaskStore) getAllTasks()(TaskStore){
	return store
}


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
