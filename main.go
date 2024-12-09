package main

import (
	"fmt"
	"os"

	"github.com/byorn/task_tracker/dao"
	"github.com/byorn/task_tracker/entity"
	"github.com/google/uuid"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tasktracker [options]")
		fmt.Println("options --> add, list, delete, complete")
		return
	}
	command := os.Args[1]

	switch command {
	case "add":
		{
			if len(os.Args) < 3 {
				fmt.Println("Usage: tasktracker add [argument]")
				fmt.Println("example: tasktracker add 'clean the laptop'")
				return
			}
			AddTask(os.Args[2])
		}
	case "list":
		{
			ListTasks()
		}
	case "delete":
		{
			if len(os.Args) < 3 {
				fmt.Println("Usage: tasktracker delete [id]")
				fmt.Println("example: tasktracker delete 12342")
				return
			}

			DeleteTask(os.Args[2])
		}

	case "complete":
		{
			if len(os.Args) < 3 {
				fmt.Println("Usage: tasktracker complete [id]")
				fmt.Println("example: tasktracker complete 12342")
				return
			}

			CompleteTask(os.Args[2])
		}

	default:
		{
			fmt.Println("Unknown command")
		}
	}

}

func AddTask(taskDescription string) {
	task := entity.Task{
		ID:          uuid.NewString(),
		Description: taskDescription,
		Status:      entity.StatusPending,
	}
	err := dao.SaveTask(task)
	if err != nil {
		fmt.Printf("Error saving task %v \n", err)
	}
	fmt.Printf("%s created", taskDescription)
}

func ListTasks() {
	fmt.Println("Listing all tasks")
	for _, task := range dao.ListTasks() {
		fmt.Println("Task ID: " + task.ID + " Task Description: " + task.Description + " Task Status: " + string(task.Status))
	}
}

func DeleteTask(uuId string) string {
	fmt.Println("Going to delete task ID :" + uuId)
	id, err := dao.DeleteTask(uuId)
	if err != nil {
		fmt.Printf("Failed to delete task %v ", err)
	}
	fmt.Printf("\n Deleted task %s", *id)
	return *id
}

func CompleteTask(uuId string) string {
	fmt.Println("Going to complete task ID :" + uuId)
	id, err := dao.CompleteTask(uuId)
	if err != nil {
		fmt.Printf("Failed to complete task %v ", err)
	}

	fmt.Printf("\n Completed task %s", id)
	return id
}
