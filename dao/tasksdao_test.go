package dao

import (
	"os"
	"testing"

	entity "github.com/byorn/task_tracker/entity"
)

func setUpMockTestFile() func() {
	TaskFile = "tasks_test.json"
	return func() {
		os.Remove(TaskFile)
	}
}

func TestSaveTask(t *testing.T) {
	cleanUp := setUpMockTestFile()
	defer cleanUp()

	expectedRecordCount := 1

	task := entity.Task{
		ID:          "123",
		Description: "Test Description",
		Status:      entity.StatusPending,
	}

	SaveTask(task)

	tasks, err := getTasksFromFile()
	if err == nil {
		if len(tasks) != expectedRecordCount {
			t.Errorf("Received count %d expecting %d", len(tasks), expectedRecordCount)
		}
		taskInFile := tasks[0]
		if task.ID != taskInFile.ID {
			t.Errorf("Task id created was %s expecting %s", taskInFile.ID, task.ID)
		}

		if task.Description != taskInFile.Description {
			t.Errorf("Task description created was %s expecting %s", taskInFile.Description, task.Description)
		}
		if string(task.Status) != string(taskInFile.Status) {
			t.Errorf("Task status created was %s expecting %s", taskInFile.Status, task.Status)
		}

	} else {
		t.Errorf("Error geting tasks from file %v", err)
	}
}
func populateDBWithTwoRecords() {

	task1 := entity.Task{
		ID:          "123",
		Description: "Test Description 123",
		Status:      entity.StatusPending,
	}

	task2 := entity.Task{
		ID:          "456",
		Description: "Test Description 456",
		Status:      entity.StatusPending,
	}

	SaveTask(task1)
	SaveTask(task2)
}
func TestDeleteTask(t *testing.T) {
	expectedRecordsInFile := 1
	cleanUp := setUpMockTestFile()
	defer cleanUp()

	populateDBWithTwoRecords()

	DeleteTask("456")

	tasks, err := getTasksFromFile()

	if err == nil {
		if len(tasks) != expectedRecordsInFile {
			t.Errorf("Found %d reocrds in file, expected %d", len(tasks), expectedRecordsInFile)
		}
	} else {
		t.Error(err)
	}

}

func TestFindTask(t *testing.T) {
	cleanUp := setUpMockTestFile()
	defer cleanUp()

	task := entity.Task{
		ID:          "123",
		Description: "Test Description",
		Status:      entity.StatusPending,
	}

	SaveTask(task)

	taskInStore, error := FindTask(task.ID)

	if error == nil {
		if taskInStore.ID != task.ID {
			t.Errorf("Task ID do not match create: %s expected %s", taskInStore.ID, task.ID)
		}
		if taskInStore.Description != task.Description {
			t.Errorf("Task Description not not match. created :%s, expected: %s", taskInStore.Description, task.Description)
		}

		if string(taskInStore.Status) != string(task.Status) {
			t.Errorf("Task Status not not match. created :%s, expected: %s", taskInStore.Status, task.Status)
		}
	} else {
		t.Error(error)
	}
}

func TestListTask(t *testing.T) {

	expected := 2
	cleanUp := setUpMockTestFile()
	defer cleanUp()

	populateDBWithTwoRecords()

	tasks := ListTasks()

	if len(tasks) != expected {
		t.Errorf("List tasks record count was : %d, expected : %d", len(tasks), expected)
	}
}

func TestCompleteTask(t *testing.T) {

	expected := entity.StatusCompleted
	cleanUp := setUpMockTestFile()
	defer cleanUp()

	populateDBWithTwoRecords()

	id, err := CompleteTask("456")

	if err != nil {
		t.Errorf("Error occured %v", err)
	}
	task, error := FindTask(id)
	if error == nil {
		if task.Status != expected {
			t.Errorf("Expected completed task to be: %s ,  but was : %s", expected, task.Status)
		}
	}
}
