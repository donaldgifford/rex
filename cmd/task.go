/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// taskCmd represents the task command
var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage tasks",
	Long: `Manage tasks for your project.

Tasks are stored in markdown files with YAML frontmatter in docs/tasks/.
The SQLite cache provides fast filtering and querying.

Available subcommands:
  create   - Create a new task with interactive prompts
  list     - List tasks with optional filtering
  complete - Mark a task as completed and move to completed/
  stats    - Show task statistics

Examples:
  rex task create
  rex task list --type core --status in_progress
  rex task complete TASK-001
  rex task stats`,
}

func init() {
	rootCmd.AddCommand(taskCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// taskCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// taskCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
