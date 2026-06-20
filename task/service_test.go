package task

import (
	"testing"
)


func TestAddTask(t *testing.T) {
	store := TaskStore{}

	store.AddTask("Buy groceries")

	// Check we have exactly 1 task
	if len(store) != 1 {
		t.Fatalf("expected 1 task, got %d", len(store))
	}

	// Check the task has correct fields
	task := store[1]
	if task.ID != 1 {
		t.Errorf("expected ID 1, got %d", task.ID)
	}
	if task.Text != "Buy groceries" {
		t.Errorf("unexpected text: %q", task.Text)
	}
	if task.Done != false {
		t.Errorf("expected Done=false, got %v", task.Done)
	}
}


func TestGetTask(t *testing.T){
	//Arrange
	task1 := Task{
		ID: 1, 
		Text: "First Task",
		Done: false,
	}

	task2 := Task{
		ID:2, 
		Text: "Second Task",
		Done: false,
	}

	store := TaskStore{
		1:&task1,
		2:&task2,
	}

	//Act
	taskGotten, err := store.getTask(2)

	if err!=nil{
		t.Fatal(err)
	}

	//Assert
	//Check that correct task was gotten
	if taskGotten.ID != 2 {
		t.Errorf("Expected second task but got %d task", taskGotten.ID)
	}

	//Chek format of text gotten back is correct(no formatting in the task here)
	if taskGotten.Text != "Second Task" {
		t.Errorf("Unexpected text %q", taskGotten.Text)
	}
}

