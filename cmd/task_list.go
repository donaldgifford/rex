/*
Copyright © 2025 Donald Gifford

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/donaldgifford/rex/internal/db"
	"github.com/spf13/cobra"
)

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Long: `List all tasks with optional filtering.

Filters:
  --type      Filter by type (core, plugin, ui, other)
  --status    Filter by status (planned, in_progress, blocked, completed, cancelled)
  --priority  Filter by priority (P0, P1, P2, P3)
  --tags      Filter by tags (comma-separated)

Examples:
  rex task list                                    # List all tasks
  rex task list --type core                        # List only core tasks
  rex task list --status in_progress               # List only in-progress tasks
  rex task list --priority P1                      # List only P1 tasks
  rex task list --type core --status in_progress   # Combined filters
  rex task list --tags database,cache              # Filter by tags`,
	RunE: func(cmd *cobra.Command, args []string) error {
		docsDir, err := cmd.Flags().GetString("docs-dir")
		if err != nil {
			return err
		}

		taskType, err := cmd.Flags().GetString("type")
		if err != nil {
			return err
		}

		status, err := cmd.Flags().GetString("status")
		if err != nil {
			return err
		}

		priority, err := cmd.Flags().GetString("priority")
		if err != nil {
			return err
		}

		tagsStr, err := cmd.Flags().GetString("tags")
		if err != nil {
			return err
		}

		var tags []string
		if tagsStr != "" {
			tags = strings.Split(tagsStr, ",")
			for i, tag := range tags {
				tags[i] = strings.TrimSpace(tag)
			}
		}

		// Make docs dir absolute
		docsDir, err = filepath.Abs(docsDir)
		if err != nil {
			return fmt.Errorf("resolve docs dir: %w", err)
		}

		// Open database
		dbPath := filepath.Join(docsDir, "..", ".rex.db")
		database, err := db.Open(dbPath)
		if err != nil {
			return fmt.Errorf("open database: %w", err)
		}
		defer database.Close()

		// Auto-rebuild if needed
		needsRebuild, err := database.NeedsRebuild(docsDir)
		if err != nil {
			return fmt.Errorf("check rebuild: %w", err)
		}
		if needsRebuild {
			fmt.Println("Cache is stale, rebuilding...")
			if err := database.Rebuild(docsDir); err != nil {
				return fmt.Errorf("rebuild cache: %w", err)
			}
		}

		// List tasks
		tasks, err := database.ListTasks(db.TaskListOptions{
			Type:     taskType,
			Status:   status,
			Priority: priority,
			Tags:     tags,
		})
		if err != nil {
			return fmt.Errorf("list tasks: %w", err)
		}

		if len(tasks) == 0 {
			fmt.Println("No tasks found")
			return nil
		}

		// Get stats for header
		stats, err := database.GetTaskStats()
		if err != nil {
			return fmt.Errorf("get stats: %w", err)
		}

		// Print header
		fmt.Println("Tasks")
		fmt.Println("=====")
		fmt.Println()

		// Print stats
		fmt.Printf("Total: %d tasks\n", stats.Total)
		if len(stats.ByStatus) > 0 {
			fmt.Print("  Status: ")
			first := true
			statusOrder := []string{"planned", "in_progress", "blocked", "completed", "cancelled"}
			for _, st := range statusOrder {
				if count, ok := stats.ByStatus[st]; ok {
					if !first {
						fmt.Print(", ")
					}
					fmt.Printf("%s: %d", st, count)
					first = false
				}
			}
			fmt.Println()
		}
		if len(stats.ByType) > 0 {
			fmt.Print("  Types: ")
			first := true
			for taskType, count := range stats.ByType {
				if !first {
					fmt.Print(", ")
				}
				fmt.Printf("%s: %d", taskType, count)
				first = false
			}
			fmt.Println()
		}
		fmt.Println()

		// Print table
		fmt.Printf("%-10s %-6s %-12s %-8s %s\n", "ID", "Type", "Status", "Priority", "Title")
		fmt.Println("────────────────────────────────────────────────────────────────────────────")

		for _, task := range tasks {
			statusDisplay := getTaskStatusDisplay(task.Status)
			fmt.Printf("%-10s %-6s %-12s %-8s %s\n",
				task.ID,
				task.Type,
				statusDisplay,
				task.Priority,
				task.Title,
			)
		}

		return nil
	},
}

func getTaskStatusDisplay(status string) string {
	switch status {
	case "planned":
		return "📋 planned"
	case "in_progress":
		return "🔨 progress"
	case "blocked":
		return "🚫 blocked"
	case "completed":
		return "✅ done"
	case "cancelled":
		return "❌ cancelled"
	default:
		return status
	}
}

func init() {
	taskCmd.AddCommand(taskListCmd)
	taskListCmd.Flags().StringP("docs-dir", "d", "docs", "Documentation directory")
	taskListCmd.Flags().StringP("type", "t", "", "Filter by type (core, plugin, ui, other)")
	taskListCmd.Flags().StringP("status", "s", "", "Filter by status (planned, in_progress, blocked, completed, cancelled)")
	taskListCmd.Flags().StringP("priority", "p", "", "Filter by priority (P0, P1, P2, P3)")
	taskListCmd.Flags().String("tags", "", "Filter by tags (comma-separated)")
}
