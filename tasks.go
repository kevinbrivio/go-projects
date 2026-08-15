package main

import (
	"fmt"
	"sort"
	"strconv"
	"time"
)

type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in-progress"
	StatusDone       Status = "done"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskList struct {
	Tasks []Task
}

/// Add task into TaskList struct
func (tl *TaskList) addTask(description string) (int, error) {
	num := fmt.Sprintf("%d", len(tl.Tasks)+1)
	id, err := strconv.Atoi(num)
	if err != nil {
		fmt.Println("Cannot convert ID: ", err)
		return -1, err
	}

	newTask := Task{
		ID:          id,
		Description: description,
		Status:      StatusTodo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tl.Tasks = append(tl.Tasks, newTask)

	return id, nil
}

/// Show tasks
func (tl *TaskList) showTasks(status Status) error {
	if (len(tl.Tasks) == 0) {
		fmt.Println("No tasks was found.")
		return nil
	}
	// Sort by UpdatedAt
	sort.Slice(tl.Tasks, func(i, j int) bool {
		return tl.Tasks[i].UpdatedAt.After(tl.Tasks[j].UpdatedAt)
	})

	found := false

	fmt.Println("\n========================================================")
	for _, t := range tl.Tasks {
		if status != "" && t.Status != status {
			continue
		}

		found = true

		fmt.Printf(
			"ID: %d | %s | %s | updated: %s\n",
			t.ID,
			t.Status,
			t.Description,
			t.UpdatedAt.Format("2006-01-02 15:04"),
		)
	}
	if !found {
		fmt.Printf("no %s tasks was found...\n", status)
	}
	fmt.Println("========================================================")

	return nil
}

// GetTaskByID returns the task with given ID
func (tl *TaskList) GetTaskByID(ID int) (*Task, error) {
	for i := range tl.Tasks {
		if tl.Tasks[i].ID == ID {
			return &tl.Tasks[i], nil
		}
	}

	return nil, fmt.Errorf("task with id %d not found", ID)
}

// MarkTask return specific task by target status
func (tl *TaskList) MarkTask(ID string, status string) (error) {
	targetID, err := strconv.Atoi(ID)
	if err != nil {
		return fmt.Errorf("wrong id")
	}

	task, err := tl.GetTaskByID(targetID)
	if err != nil {
		return err
	}

	task.Status = Status(status)
	task.UpdatedAt = time.Now()

	return nil
}

// UpdateTask only updates the description and updated_at
func (tl *TaskList) UpdateTask(ID string, description string) (error) {
	targetID, err := strconv.Atoi(ID)
	if err != nil {
		return fmt.Errorf("wrong id")
	}

	task, err := tl.GetTaskByID(targetID)
	if err != nil {
		return err
	}

	task.Description = description
	task.UpdatedAt = time.Now()

	return nil
}

// DeleteTask will delete task permanently
func (tl *TaskList) DeleteTask(ID string) (error) {
	targetID, err := strconv.Atoi(ID)
	if err != nil {
		return fmt.Errorf("wrong id")
	}
	
	for i, task := range tl.Tasks {
		if task.ID == targetID {
			tl.Tasks = append(tl.Tasks[:i], tl.Tasks[i+1:]...)
		}
	}

	return nil
}

