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
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/donaldgifford/rex/internal/db"
	"github.com/donaldgifford/rex/internal/parser"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var taskCompleteCmd = &cobra.Command{
	Use:   "complete <task-id>",
	Short: "Mark a task as completed",
	Long: `Mark a task as completed and move it to the completed/ directory.

This command will:
1. Prompt for actual hours spent
2. Update the task status to "completed"
3. Set the completed_date to today
4. Update the markdown file frontmatter
5. Move the file from {type}/active/ to {type}/completed/
6. Update the database cache

Example:
  rex task complete TASK-001
  rex task complete TASK-001 --actual-hours 10`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := strings.TrimSpace(args[0])
		if taskID == "" {
			return fmt.Errorf("task ID is required")
		}

		docsDir, err := cmd.Flags().GetString("docs-dir")
		if err != nil {
			return err
		}

		actualHoursFlag, err := cmd.Flags().GetFloat64("actual-hours")
		if err != nil {
			return err
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

		// Get task
		task, err := database.GetTask(taskID)
		if err != nil {
			return fmt.Errorf("get task: %w", err)
		}

		if task.Status == "completed" {
			return fmt.Errorf("task %s is already completed", taskID)
		}

		// Prompt for actual hours if not provided
		var actualHours float64
		if actualHoursFlag > 0 {
			actualHours = actualHoursFlag
		} else {
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("Task: %s\n", task.Title)
			fmt.Printf("Estimated: %.1f hours\n", task.EstimatedHours)
			fmt.Print("Actual hours spent: ")
			actualStr, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("read actual hours: %w", err)
			}
			actualStr = strings.TrimSpace(actualStr)
			actualHours, err = strconv.ParseFloat(actualStr, 64)
			if err != nil {
				return fmt.Errorf("invalid actual hours: %w", err)
			}
			if actualHours < 0 {
				return fmt.Errorf("actual hours must be non-negative")
			}
		}

		// Read current file
		content, err := os.ReadFile(task.FilePath)
		if err != nil {
			return fmt.Errorf("read file: %w", err)
		}

		// Parse frontmatter
		frontmatter, body, err := parser.ExtractFrontmatter(content)
		if err != nil {
			return fmt.Errorf("extract frontmatter: %w", err)
		}

		var taskMeta parser.TaskFrontmatter
		if err := yaml.Unmarshal(frontmatter, &taskMeta); err != nil {
			return fmt.Errorf("parse frontmatter: %w", err)
		}

		// Update frontmatter
		completedDate := time.Now().Format("2006-01-02")
		taskMeta.Status = "completed"
		taskMeta.CompletedDate = &completedDate
		taskMeta.ActualHours = &actualHours

		// Marshal updated frontmatter
		updatedFrontmatter, err := yaml.Marshal(&taskMeta)
		if err != nil {
			return fmt.Errorf("marshal frontmatter: %w", err)
		}

		// Reconstruct file content
		updatedContent := []byte(fmt.Sprintf("---\n%s---\n%s", string(updatedFrontmatter), string(body)))

		// Determine new file path (move from active/ to completed/)
		oldPath := task.FilePath
		newPath := strings.Replace(oldPath, "/active/", "/completed/", 1)

		// Ensure completed directory exists
		completedDir := filepath.Dir(newPath)
		if err := os.MkdirAll(completedDir, 0755); err != nil {
			return fmt.Errorf("create completed directory: %w", err)
		}

		// Write updated content to new location
		if err := os.WriteFile(newPath, updatedContent, 0644); err != nil {
			return fmt.Errorf("write file: %w", err)
		}

		// Remove old file
		if err := os.Remove(oldPath); err != nil {
			return fmt.Errorf("remove old file: %w", err)
		}

		// Update database
		task.Status = "completed"
		task.CompletedDate = &completedDate
		task.ActualHours = &actualHours
		task.FilePath = newPath
		task.Content = string(updatedContent)

		if err := database.UpdateTask(task); err != nil {
			return fmt.Errorf("update task in database: %w", err)
		}

		fmt.Printf("✓ Completed task %s: %s\n", taskID, task.Title)
		fmt.Printf("  Actual hours: %.1f (estimated: %.1f)\n", actualHours, task.EstimatedHours)
		if actualHours > task.EstimatedHours {
			diff := actualHours - task.EstimatedHours
			fmt.Printf("  ⚠️  Over estimate by %.1f hours\n", diff)
		} else if actualHours < task.EstimatedHours {
			diff := task.EstimatedHours - actualHours
			fmt.Printf("  ✓ Under estimate by %.1f hours\n", diff)
		}
		fmt.Printf("  Moved: %s -> %s\n", filepath.Base(oldPath), newPath)

		return nil
	},
}

func init() {
	taskCmd.AddCommand(taskCompleteCmd)
	taskCompleteCmd.Flags().StringP("docs-dir", "d", "docs", "Documentation directory")
	taskCompleteCmd.Flags().Float64("actual-hours", 0, "Actual hours spent (skip prompt)")
}
