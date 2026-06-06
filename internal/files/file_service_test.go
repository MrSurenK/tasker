package files

import (
	"testing"
	"time"
)

func TestFileNameGenerator(t *testing.T ){

	currDate := time.Now()
	date := currDate.Format("2006-01-02")
	expected_fileName := "ToDoList-"+date
	generated_fileName := generateFileName(currDate)

	if generated_fileName != expected_fileName{
		t.Errorf("Expected %v , but got %v", expected_fileName, generated_fileName)
	}
}





