package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var StatusTodo string = "todo"
var StatusInProgress string = "in-progress"
var StatusDone string = "done"

func NewTask(description string) (*Task, error) {
	meta, err := LoadMeta()
	if err != nil {
		return nil, err
	}
	id := meta.NextID
	meta.NextID++
	if err = SaveMeta(meta); err != nil {
		return nil, err
	}
	return &Task{
		// NOTE id need more change
		ID:          id,
		Description: description,
		Status:      StatusTodo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func AddTask(description string) (int, error) {
	task, err := NewTask(description)
	if err != nil {
		return 0, err
	}
	if err = AppendJSONL(task); err != nil {
		return 0, err
	}
	id := task.ID
	return id, nil
}

var ErrTaskNotFound = errors.New("task not found")

// func GetTaskByID(id int) (*Task, error) {
// 	f, err := os.Open(TasksFile)
// 	if err != nil {
// 		return nil, err
// 	}
// 	decoder := json.NewDecoder(f)
// 	for {
// 		var t Task
// 		err := decoder.Decode(&t)
// 		if errors.Is(err, io.EOF) {
// 			return nil, ErrTaskNotFound
// 		}
// 		if err := decoder.Decode(&t); err != nil {
// 			return nil, err
// 		}
// 		if t.ID == id {
// 			return &t, nil
// 		}
// 	}
// }

func UpdateTask(id int, newDiscription string) error {
	tasks, err := getAllTasks()
	if err != nil {
		return err
	}
	found := false
	for i, t := range tasks {
		if t.ID == id {
			tasks[i].Description = newDiscription
			tasks[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}
	if found == false {
		return ErrTaskNotFound
	}

	if err := RewriteJSON(tasks); err != nil {
		return err
	}
	return nil
}

func DeleteTask(id int) error {
	tasks, err := getAllTasks()
	if err != nil {
		return err
	}
	found := false
	var newTasks []Task
	for _, t := range tasks {
		if t.ID != id {
			newTasks = append(newTasks, t)
		} else {
			found = true
		}
	}
	if found == false {
		return ErrTaskNotFound
	}

	_, err = os.Create(TasksFile)
	if err != nil {
		return err
	}
	for _, t := range newTasks {
		if err := AppendJSONL(t); err != nil {
			return err
		}
	}

	return nil
}

func ChangeStatus(id int, newStatus string) error {
	tasks, err := getAllTasks()
	if err != nil {
		return err
	}

	found := false
	for i, t := range tasks {
		if t.ID == id {
			tasks[i].Status = newStatus
			found = true
		}
	}
	if found == false {
		return ErrTaskNotFound
	}
	if err := RewriteJSON(tasks); err != nil {
		return err
	}
	return nil
}

func PrintAllTasks() error {
	tasks, err := getAllTasks()
	if err != nil {
		return err
	}
	printTasks(tasks)
	return nil
}

func printTasks(tasks []Task) {
	if len(tasks) == 0 {
		fmt.Println("no tasks found")
		return
	}
	for _, t := range tasks {
		fmt.Printf("ID: %d\tStatus: %s\tCreatedAt: %s\tUpdatedAt: %s\tDescription: %s\n",
			t.ID,
			t.Status,
			t.CreatedAt.Format(time.RFC3339),
			t.UpdatedAt.Format(time.RFC3339),
			t.Description,
		)
	}
}

func PrintTasksByStatus(status string) error {
	tasks, err := getTasksByStatus(status)
	if err != nil {
		return err
	}
	printTasks(tasks)
	return nil
}

func getTasksByStatus(status string) ([]Task, error) {
	tasks, err := getAllTasks()
	if err != nil {
		return nil, err
	}
	var filtered []Task
	for _, t := range tasks {
		if t.Status == status {
			filtered = append(filtered, t)
		}
	}
	return filtered, nil
}

func getAllTasks() (tasks []Task, err error) {
	f, err := os.Open(TasksFile)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(f)
	for {
		var t Task
		if err := decoder.Decode(&t); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
