package task

import (
	"io"
	"os"
	"strings"
	"errors"
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

// helper method to get the task and return error if no task found with that id
func (tasks TaskStore) getTask(number int)(*Task, error){
	task, ok := tasks[number]

	if !ok {
		return nil, ErrTaskNotFound
	}
	return task, nil
}









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
