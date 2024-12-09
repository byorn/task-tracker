package dao

import (
	"encoding/json"
	"fmt"
	"os"

	entity "github.com/byorn/task_tracker/entity"
)

// TaskFile defines the name of the file used to persist tasks.
var TaskFile = "tasks.json"

// TaskFile defines the name of the file used to persist tasks.
var tasksListMap map[string]entity.Task = make(map[string]entity.Task)

// SaveTask loads tasks from file, appends the new task and saves the list to the file.
func SaveTask(task entity.Task) error {

	tasksInFile, err := getTasksFromFile()

	if err != nil {
		return err
	}

	tasks := append(tasksInFile, task)
	data, err := json.MarshalIndent(&tasks, "", "   ")
	if err != nil {
		return fmt.Errorf("Could not marshal tasks to file %w", err)
	}

	return os.WriteFile(TaskFile, data, os.ModePerm)
}

// CompleteTask marks a task as completed by its UUID.
// It first checks the in-memory map, and if not found,
// loads tasks from the file and tries again.
func CompleteTask(uuId string) (string, error) {
	if len(tasksListMap) > 0 {
		task, found := tasksListMap[uuId]
		if found {
			task.Status = entity.StatusCompleted
			tasksListMap[uuId] = task
			saveMapToFile()
			return task.ID, nil
		} else {
			return "", fmt.Errorf("couldnt find task for id: %s", uuId)
		}
	}

	tasks := ListTasks()
	if len(tasks) > 0 {
		return CompleteTask(uuId)
	}
	return "", fmt.Errorf("no records in store found")
}

// FindTask retrieves a task by its UUID.
// It checks the in-memory map first, and if not found,
// loads tasks from the file into the map before searching again.
func FindTask(uuId string) (*entity.Task, error) {
	if len(tasksListMap) > 0 {
		task, found := tasksListMap[uuId]
		if found {
			return &task, nil
		} else {
			return nil, fmt.Errorf("task not found for uuId %s", uuId)
		}
	} else {

		tasksInFile, err := getTasksFromFile()
		if err != nil {
			fmt.Printf("Error loading tasks from file %v \n", err)
		}

		loadTasksToMap(tasksInFile)
		task, found := tasksListMap[uuId]
		if found {
			return &task, nil
		} else {
			return nil,
				fmt.Errorf("task not found for uuId %s", uuId)
		}
	}
}

// writes the in-memory map of tasks to the tasks file.
func saveMapToFile() error {
	var tasks []entity.Task
	for _, task := range tasksListMap {
		tasks = append(tasks, task)
	}
	data, err := json.MarshalIndent(&tasks, "", "   ")
	if err != nil {
		return fmt.Errorf("could not marshal tasks to file %w", err)
	}

	return os.WriteFile(TaskFile, data, os.ModePerm)
}

// DeleteTask removes a task by its ID from both the in-memory map and the file.
// Returns the ID of the deleted task if successful.
func DeleteTask(id string) (*string, error) {

	ListTasks()

	if len(tasksListMap) > 0 {
		delete(tasksListMap, id)
		saveMapToFile()

		return &id, nil
	}

	return nil, fmt.Errorf("task ID not found %s", id)
}

// ListTasks retrieves all tasks.
// If tasks exist in memory, it returns them directly.
// Otherwise, it loads tasks from the file and returns them.
func ListTasks() []entity.Task {

	if len(tasksListMap) > 0 {
		return getTasksFromMemory()
	}

	tasksInFile, err := getTasksFromFile()
	if err != nil {
		fmt.Printf("Error loading tasks from file %v \n", err)
	}

	loadTasksToMap(tasksInFile)
	return tasksInFile
}

// getTasksFromMemory retrieves all tasks from the in-memory map.
func getTasksFromMemory() []entity.Task {
	var tasks []entity.Task

	for _, v := range tasksListMap {
		tasks = append(tasks, v)
	}
	return tasks
}

// getTasksFromFile reads tasks from the tasks file.
// If the file does not exist, it returns an empty task list.
func getTasksFromFile() ([]entity.Task, error) {
	//check if file exists
	//os.Stat will return a fileInfo, and err. the err is passed to os.IsNotExist to check if file does not exist.
	if _, err := os.Stat(TaskFile); os.IsNotExist(err) {
		return []entity.Task{}, nil
	}
	var tasks []entity.Task
	data, err := os.ReadFile(TaskFile)
	if err == nil {
		err = json.Unmarshal(data, &tasks)
	}

	if err != nil {
		err = fmt.Errorf("failed to unmarchsal JSON %w", err)
	}
	return tasks, err
}

// loadTasksToMap populates the in-memory map with tasks loaded from the file.
func loadTasksToMap(tasksFromFile []entity.Task) {

	if len(tasksListMap) == 0 {
		for _, task := range tasksFromFile {
			tasksListMap[task.ID] = task
		}
	}
}
