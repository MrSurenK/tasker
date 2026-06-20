package task

//Fields that the task must have
type Task struct {
	ID int
	Text string
	Done bool
}

//To keep a key value store where the id represents the task number
//the map should store the pointer reference to the actual Task object so that it can be overwriten by calling it from the map
type TaskStore map[int]*Task

