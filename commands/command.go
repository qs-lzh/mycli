package commands

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/qs-lzh/mycli/internal/data"
)

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(markInProgressCmd)
	rootCmd.AddCommand(markDoneCmd)
	rootCmd.AddCommand(listCmd)
}

var rootCmd = &cobra.Command{
	Use:   "task-cli",
	Short: "task-cli is a simple cli tool",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run task-cli...")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println("error occur at rootCmd.Execute")
	}
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "adding a new task",
	Run: func(cmd *cobra.Command, args []string) {
		description := args[0]
		id, err := data.AddTask(description)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Task added successfully (Id: %d)\n", id)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "updating a task",
	Run: func(cmd *cobra.Command, args []string) {
		idStr, newDiscription := args[0], args[1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("id format incorrect!")
			log.Fatal(err)
		}
		if err := data.UpdateTask(id, newDiscription); err != nil {
			log.Fatal(err)
		}
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "deleting a task",
	Run: func(cmd *cobra.Command, args []string) {
		idStr := args[0]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("id format incorrect!")
			log.Fatal(err)
		}
		if err = data.DeleteTask(id); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("delete task id = %d successfully\n", id)
	},
}

var markInProgressCmd = &cobra.Command{
	Use:   "mark-in-progress",
	Short: "marking a task as in progress",
	Run: func(cmd *cobra.Command, args []string) {
		idStr := args[0]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("id format incorrect!")
			log.Fatal(err)
		}
		if err := data.ChangeStatus(id, data.StatusInProgress); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("mark task id = %d in progress successfully\n", id)
	},
}

var markDoneCmd = &cobra.Command{
	Use:   "mark-done",
	Short: "marking a task as done",
	Run: func(cmd *cobra.Command, args []string) {
		idStr := args[0]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("id format incorrect!")
			log.Fatal(err)
		}
		if err := data.ChangeStatus(id, data.StatusDone); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("mark task id = %d done successfully\n", id)
	},
}

var listCmd = &cobra.Command{
	Use:   "list [status]",
	Short: "listing tasks by status (or all tasks)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			if err := data.PrintAllTasks(); err != nil {
				log.Fatal(err)
			}
			return
		}

		statusArg := args[0]
		var status string
		switch statusArg {
		case "todo":
			status = data.StatusTodo
		case "in-progress":
			status = data.StatusInProgress
		case "done":
			status = data.StatusDone
		default:
			fmt.Println("unknown status, must be one of: todo, in-progress, done")
			return
		}

		if err := data.PrintTasksByStatus(status); err != nil {
			log.Fatal(err)
		}
	},
}
