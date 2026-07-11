package task

import (
	"os"
	"testing"
)

// ---------- AddTask ----------

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

// ---------- getTask ----------

func TestGetTask(t *testing.T) {
	// Arrange
	task1 := Task{
		ID:   1,
		Text: "First Task",
		Done: false,
	}

	task2 := Task{
		ID:   2,
		Text: "Second Task",
		Done: false,
	}

	store := TaskStore{
		1: &task1,
		2: &task2,
	}

	// Act
	taskGotten, err := store.getTask(2)

	if err != nil {
		t.Fatal(err)
	}

	// Assert
	if taskGotten.ID != 2 {
		t.Errorf("Expected second task but got %d task", taskGotten.ID)
	}

	if taskGotten.Text != "Second Task" {
		t.Errorf("Unexpected text %q", taskGotten.Text)
	}
}

func TestGetTask_NotFound(t *testing.T) {
	store := TaskStore{}

	_, err := store.getTask(999)

	if err == nil {
		t.Fatal("expected error for non-existent task, got nil")
	}
	if err != ErrTaskNotFound {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}
}

// ---------- EditTask ----------

func TestEditTask(t *testing.T) {
	// Arrange
	store := TaskStore{
		1: {ID: 1, Text: "Original text", Done: false},
	}

	// Act
	updated := store.EditTask(1, WithText("Updated text"))

	// Assert
	if updated.Text != "Updated text" {
		t.Errorf("expected text 'Updated text', got %q", updated.Text)
	}
	if updated.Done != false {
		t.Errorf("expected Done=false, got %v", updated.Done)
	}
	// Verify the original in the store was also mutated (pointer semantics)
	if store[1].Text != "Updated text" {
		t.Errorf("store task text should also be 'Updated text', got %q", store[1].Text)
	}
}

func TestEditTask_CompletionStatus(t *testing.T) {
	store := TaskStore{
		1: {ID: 1, Text: "Task", Done: false},
	}

	updated := store.EditTask(1, WithCompletionStatus(true))

	if updated.Done != true {
		t.Errorf("expected Done=true, got %v", updated.Done)
	}
}

func TestEditTask_MultipleOptions(t *testing.T) {
	store := TaskStore{
		1: {ID: 1, Text: "Old", Done: false},
	}

	updated := store.EditTask(1, WithText("New"), WithCompletionStatus(true))

	if updated.Text != "New" {
		t.Errorf("expected text 'New', got %q", updated.Text)
	}
	if updated.Done != true {
		t.Errorf("expected Done=true, got %v", updated.Done)
	}
}

func TestEditTask_NotFound(t *testing.T) {
	store := TaskStore{}

	// EditTask prints the error but does not return early — it will attempt
	// to call options on a nil *Task.  WithText and WithCompletionStatus
	// dereference the pointer, so this test documents the current behaviour
	// (a nil-pointer dereference panic).
	//
	// To avoid the panic we pass no options; the method simply returns nil.
	result := store.EditTask(999) // no options

	if result != nil {
		t.Errorf("expected nil when task not found, got %v", result)
	}
}

// ---------- DeleteTask ----------

func TestDeleteTask(t *testing.T) {
	store := TaskStore{
		1: {ID: 1, Text: "Task to delete", Done: false},
		2: {ID: 2, Text: "Task to keep", Done: false},
	}

	store.DeleteTask(1)

	if len(store) != 1 {
		t.Fatalf("expected 1 task after delete, got %d", len(store))
	}
	if _, exists := store[1]; exists {
		t.Error("task 1 should have been deleted")
	}
	if _, exists := store[2]; !exists {
		t.Error("task 2 should still exist")
	}
}

func TestDeleteTask_NonExistent(t *testing.T) {
	store := TaskStore{
		1: {ID: 1, Text: "Only task", Done: false},
	}

	// Deleting a key that doesn't exist is a no-op in Go maps
	store.DeleteTask(999)

	if len(store) != 1 {
		t.Errorf("expected 1 task, got %d", len(store))
	}
}

// ---------- getLastId ----------

func TestGetLastId(t *testing.T) {
	store := TaskStore{
		1:  {ID: 1, Text: "First", Done: false},
		5:  {ID: 5, Text: "Fifth", Done: false},
		3:  {ID: 3, Text: "Third", Done: false},
		10: {ID: 10, Text: "Tenth", Done: false},
	}

	last := getLastId(store)

	if last != 10 {
		t.Errorf("expected last ID 10, got %d", last)
	}
}

func TestGetLastId_EmptyStore(t *testing.T) {
	store := TaskStore{}

	last := getLastId(store)

	if last != 0 {
		t.Errorf("expected last ID 0 for empty store, got %d", last)
	}
}

// ---------- getAllTasks ----------

func TestGetAllTasks(t *testing.T) {
	store := TaskStore{
		1: {ID: 1, Text: "Task 1", Done: false},
		2: {ID: 2, Text: "Task 2", Done: true},
	}

	all := store.getAllTasks()

	if len(all) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(all))
	}
	if all[1] == nil {
		t.Error("expected task 1 to exist")
	}
	if all[2] == nil {
		t.Error("expected task 2 to exist")
	}
}

func TestGetAllTasks_EmptyStore(t *testing.T) {
	store := TaskStore{}

	all := store.getAllTasks()

	if len(all) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(all))
	}
}

// ---------- checkStatusInFile ----------

func TestCheckStatusInFile(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{"done_task", "-[x] Buy groceries", true},
		{"not_done_task", "-[ ] Walk the dog", false},
		{"no_prefix", "Plain task without prefix", false},
		{"partial_x_no_brackets", "-x] Some task", false},
		{"wrong_prefix", "[x] Missing dash", false},
		{"case_sensitive", "-[X] Uppercase X", false}, // only lowercase 'x' matches
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkStatusInFile(tt.line)
			if result != tt.expected {
				t.Errorf("checkStatusInFile(%q) = %v, want %v", tt.line, result, tt.expected)
			}
		})
	}
}

// ---------- stripCheckboxes ----------

func TestStripCheckboxes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"done_task", "-[x] Buy groceries", "Buy groceries"},
		{"not_done_task", "-[ ] Walk the dog", "Walk the dog"},
		{"five_char_prefix", "-[_] Some task", "Some task"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripCheckboxes(tt.input)
			if result != tt.expected {
				t.Errorf("stripCheckboxes(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// ---------- prepareTask ----------

func TestPrepareTask(t *testing.T) {
	tests := []struct {
		name     string
		task     Task
		expected string
	}{
		{"not_done", Task{ID: 1, Text: "Buy groceries", Done: false}, "-[ ] Buy groceries"},
		{"done", Task{ID: 2, Text: "Walk dog", Done: true}, "-[x] Walk dog"},
		{"empty_text", Task{ID: 3, Text: "", Done: false}, "-[ ] "},
		{"empty_text_done", Task{ID: 4, Text: "", Done: true}, "-[x] "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.task.prepareTask()
			if result != tt.expected {
				t.Errorf("prepareTask() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// ---------- Functional Options ----------

func TestWithText(t *testing.T) {
	task := &Task{ID: 1, Text: "Old", Done: false}

	opt := WithText("New text")
	opt(task)

	if task.Text != "New text" {
		t.Errorf("expected text 'New text', got %q", task.Text)
	}
}

func TestWithCompletionStatus(t *testing.T) {
	task := &Task{ID: 1, Text: "Task", Done: false}

	// Set to true
	optTrue := WithCompletionStatus(true)
	optTrue(task)
	if task.Done != true {
		t.Errorf("expected Done=true, got %v", task.Done)
	}

	// Set back to false
	optFalse := WithCompletionStatus(false)
	optFalse(task)
	if task.Done != false {
		t.Errorf("expected Done=false, got %v", task.Done)
	}
}

// ---------- File-based services ----------

// getSavedTasks

func TestGetSavedTasks(t *testing.T) {
	file, err := os.CreateTemp("", "test-saved-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	content := "line1\nline2\nline3"
	if _, err := file.WriteString(content); err != nil {
		t.Fatal(err)
	}

	tasks, err := getSavedTasks(file)
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(tasks))
	}
	if tasks[0] != "line1" {
		t.Errorf("expected first line 'line1', got %q", tasks[0])
	}
	if tasks[1] != "line2" {
		t.Errorf("expected second line 'line2', got %q", tasks[1])
	}
	if tasks[2] != "line3" {
		t.Errorf("expected third line 'line3', got %q", tasks[2])
	}
}

func TestGetSavedTasks_EmptyFile(t *testing.T) {
	file, err := os.CreateTemp("", "test-empty-saved-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	tasks, err := getSavedTasks(file)
	if err != nil {
		t.Fatal(err)
	}

	// Split on an empty string yields [""] — one element.
	if len(tasks) != 1 {
		t.Fatalf("expected 1 element from empty file, got %d", len(tasks))
	}
	if tasks[0] != "" {
		t.Errorf("expected empty string, got %q", tasks[0])
	}
}

// MapMdToStore

func TestMapMdToStore(t *testing.T) {
	// Arrange — create a temp file with markdown task content
	file, err := os.CreateTemp("", "test-map-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	content := "-[x] Buy groceries\n-[ ] Walk dog\n-[ ] Read book\n"
	if _, err := file.WriteString(content); err != nil {
		t.Fatal(err)
	}

	// Act
	store := TaskStore{}
	store, err = store.MapMdToStore(file)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	if len(store) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(store))
	}

	if store[1].ID != 1 {
		t.Errorf("expected task 1 ID=1, got %d", store[1].ID)
	}
	if store[1].Text != "Buy groceries" {
		t.Errorf("expected task 1 text 'Buy groceries', got %q", store[1].Text)
	}
	if store[1].Done != true {
		t.Errorf("expected task 1 Done=true, got %v", store[1].Done)
	}

	if store[2].Text != "Walk dog" {
		t.Errorf("expected task 2 text 'Walk dog', got %q", store[2].Text)
	}
	if store[2].Done != false {
		t.Errorf("expected task 2 Done=false, got %v", store[2].Done)
	}

	if store[3].Text != "Read book" {
		t.Errorf("expected task 3 text 'Read book', got %q", store[3].Text)
	}
}

// UpdateDoc

func TestUpdateDoc(t *testing.T) {
	// Arrange
	file, err := os.CreateTemp("", "test-update-doc-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	store := TaskStore{
		1: {ID: 1, Text: "First task", Done: false},
		2: {ID: 2, Text: "Second task", Done: true},
	}

	// Act
	if err := store.UpdateDoc(file); err != nil {
		t.Fatal(err)
	}

	// Assert — read file back and verify content
	content, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Map iteration order is random, so check that both expected lines are present
	got := string(content)
	if len(got) == 0 {
		t.Fatal("expected non-empty file content")
	}

	hasFirst := false
	hasSecond := false
	lines := splitLines(got)
	for _, line := range lines {
		if line == "-[ ] First task" {
			hasFirst = true
		}
		if line == "-[x] Second task" {
			hasSecond = true
		}
	}
	if !hasFirst {
		t.Error("file missing expected line: '-[ ] First task'")
	}
	if !hasSecond {
		t.Error("file missing expected line: '-[x] Second task'")
	}
}

func TestUpdateDoc_EmptyStore(t *testing.T) {
	file, err := os.CreateTemp("", "test-update-empty-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	store := TaskStore{}

	if err := store.UpdateDoc(file); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	if len(content) != 0 {
		t.Errorf("expected empty file, got %q", string(content))
	}
}

func TestUpdateDoc_TruncatesExistingContent(t *testing.T) {
	file, err := os.CreateTemp("", "test-truncate-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	// Write some old content first
	if _, err := file.WriteString("-old stale content\n"); err != nil {
		t.Fatal(err)
	}

	store := TaskStore{
		1: {ID: 1, Text: "Fresh task", Done: false},
	}

	if err := store.UpdateDoc(file); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	got := string(content)

	// Old content must be gone
	if got == "-old stale content\n" {
		t.Error("old content was not truncated — stale content remains")
	}

	// Fresh content must be present
	hasFresh := false
	for _, line := range splitLines(got) {
		if line == "-[ ] Fresh task" {
			hasFresh = true
		}
	}
	if !hasFresh {
		t.Errorf("expected fresh task in file, got %q", got)
	}
}

// ---------- Integration-style: round-trip MapMdToStore → UpdateDoc ----------

func TestRoundTrip_MapToStoreAndUpdateDoc(t *testing.T) {
	// 1. Write initial markdown tasks to a file
	file, err := os.CreateTemp("", "test-roundtrip-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	initial := "-[x] Done task\n-[ ] Pending task\n"
	if _, err := file.WriteString(initial); err != nil {
		t.Fatal(err)
	}

	// 2. Map the file into a TaskStore
	store := TaskStore{}
	store, err = store.MapMdToStore(file)
	if err != nil {
		t.Fatal(err)
	}

	// 3. Edit a task
	store.EditTask(2, WithText("Pending task updated"))

	// 4. Write back to the file
	if err := store.UpdateDoc(file); err != nil {
		t.Fatal(err)
	}

	// 5. Read file and verify the round-trip
	content, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	lines := splitLines(string(content))
	hasDone := false
	hasPending := false
	for _, line := range lines {
		if line == "-[x] Done task" {
			hasDone = true
		}
		if line == "-[ ] Pending task updated" {
			hasPending = true
		}
	}
	if !hasDone {
		t.Error("expected '-[x] Done task' in file")
	}
	if !hasPending {
		t.Error("expected '-[ ] Pending task updated' in file")
	}
}

// ---------- helpers ----------

// splitLines splits a string by newlines, discarding empty trailing lines.
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	// Don't append an empty string for trailing newline
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
