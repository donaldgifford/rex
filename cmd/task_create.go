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
	"github.com/donaldgifford/rex/internal/generator"
	"github.com/donaldgifford/rex/internal/templates"
	"github.com/spf13/cobra"
)

var taskCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new task",
	Long: `Create a new task with interactive prompts.

This command will prompt you for:
- Title (required)
- Type (core, plugin, ui, other)
- Priority (P0-P3)
- Estimated hours
- Assignee (optional)
- Phase (optional)
- Tags (comma-separated)
- Relationships (blocked_by, blocks, related_to)

The task will be created in docs/tasks/{type}/active/ with a unique ID.

Example:
  rex task create`,
	RunE: func(cmd *cobra.Command, args []string) error {
		docsDir, err := cmd.Flags().GetString("docs-dir")
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

		reader := bufio.NewReader(os.Stdin)

		// Prompt for task details
		fmt.Println("Creating new task...")
		fmt.Println()

		// Title
		fmt.Print("Title: ")
		title, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read title: %w", err)
		}
		title = strings.TrimSpace(title)
		if title == "" {
			return fmt.Errorf("title is required")
		}

		// Type
		fmt.Print("Type (core/plugin/ui/other) [core]: ")
		taskType, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read type: %w", err)
		}
		taskType = strings.TrimSpace(taskType)
		if taskType == "" {
			taskType = "core"
		}
		validTypes := []string{"core", "plugin", "ui", "other"}
		if !contains(validTypes, taskType) {
			return fmt.Errorf("invalid type: %s (must be core, plugin, ui, or other)", taskType)
		}

		// Priority
		fmt.Print("Priority (P0/P1/P2/P3) [P2]: ")
		priority, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read priority: %w", err)
		}
		priority = strings.TrimSpace(priority)
		if priority == "" {
			priority = "P2"
		}
		validPriorities := []string{"P0", "P1", "P2", "P3"}
		if !contains(validPriorities, priority) {
			return fmt.Errorf("invalid priority: %s (must be P0, P1, P2, or P3)", priority)
		}

		// Estimated hours
		fmt.Print("Estimated hours [8]: ")
		estimatedStr, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read estimated hours: %w", err)
		}
		estimatedStr = strings.TrimSpace(estimatedStr)
		if estimatedStr == "" {
			estimatedStr = "8"
		}
		estimatedHours, err := strconv.ParseFloat(estimatedStr, 64)
		if err != nil {
			return fmt.Errorf("invalid estimated hours: %w", err)
		}
		if estimatedHours <= 0 {
			return fmt.Errorf("estimated hours must be positive")
		}

		// Assignee (optional)
		fmt.Print("Assignee (optional): ")
		assignee, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read assignee: %w", err)
		}
		assignee = strings.TrimSpace(assignee)
		var assigneePtr *string
		if assignee != "" {
			assigneePtr = &assignee
		}

		// Phase (optional)
		fmt.Print("Phase (optional): ")
		phaseStr, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read phase: %w", err)
		}
		phaseStr = strings.TrimSpace(phaseStr)
		var phasePtr *int
		if phaseStr != "" {
			phase, err := strconv.Atoi(phaseStr)
			if err != nil {
				return fmt.Errorf("invalid phase: %w", err)
			}
			phasePtr = &phase
		}

		// Tags
		fmt.Print("Tags (comma-separated, optional): ")
		tagsStr, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read tags: %w", err)
		}
		tagsStr = strings.TrimSpace(tagsStr)
		var tags []string
		if tagsStr != "" {
			tags = strings.Split(tagsStr, ",")
			for i, tag := range tags {
				tags[i] = strings.TrimSpace(tag)
			}
		}

		// Relationships (optional)
		fmt.Print("Blocked by (comma-separated task IDs, optional): ")
		blockedByStr, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read blocked_by: %w", err)
		}
		blockedByStr = strings.TrimSpace(blockedByStr)
		var blockedBy []string
		if blockedByStr != "" {
			blockedBy = strings.Split(blockedByStr, ",")
			for i, id := range blockedBy {
				blockedBy[i] = strings.TrimSpace(id)
			}
		}

		fmt.Print("Blocks (comma-separated task IDs, optional): ")
		blocksStr, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read blocks: %w", err)
		}
		blocksStr = strings.TrimSpace(blocksStr)
		var blocks []string
		if blocksStr != "" {
			blocks = strings.Split(blocksStr, ",")
			for i, id := range blocks {
				blocks[i] = strings.TrimSpace(id)
			}
		}

		fmt.Print("Related to (comma-separated task IDs, optional): ")
		relatedToStr, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read related_to: %w", err)
		}
		relatedToStr = strings.TrimSpace(relatedToStr)
		var relatedTo []string
		if relatedToStr != "" {
			relatedTo = strings.Split(relatedToStr, ",")
			for i, id := range relatedTo {
				relatedTo[i] = strings.TrimSpace(id)
			}
		}

		fmt.Println()

		// Get next task number
		number, err := generator.GetNextTaskNumber(database.Conn())
		if err != nil {
			return fmt.Errorf("get next task number: %w", err)
		}

		// Generate task ID
		taskID := generator.GenerateTaskID(number)

		// Generate filename
		filename := generator.GenerateTaskFilename(taskID, title)
		taskDir := filepath.Join(docsDir, "tasks", taskType, "active")
		filePath := filepath.Join(taskDir, filename)

		// Ensure directory exists
		if err := os.MkdirAll(taskDir, 0755); err != nil {
			return fmt.Errorf("create task directory: %w", err)
		}

		// Render template
		content, err := templates.RenderTask(templates.TaskData{
			ID:             taskID,
			Title:          title,
			Type:           taskType,
			Priority:       priority,
			EstimatedHours: estimatedHours,
			Assignee:       assigneePtr,
			Tags:           tags,
			Phase:          phasePtr,
		})
		if err != nil {
			return fmt.Errorf("render template: %w", err)
		}

		// Write file
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("write file: %w", err)
		}

		// Create task in database
		startedDate := time.Now().Format("2006-01-02")
		task := &db.Task{
			ID:             taskID,
			Title:          title,
			Type:           taskType,
			Status:         "planned",
			Priority:       priority,
			EstimatedHours: estimatedHours,
			StartedDate:    &startedDate,
			Assignee:       assigneePtr,
			Phase:          phasePtr,
			FilePath:       filePath,
			Content:        content,
			BlockedBy:      blockedBy,
			Blocks:         blocks,
			RelatedTo:      relatedTo,
			Tags:           tags,
		}

		if err := database.CreateTask(task); err != nil {
			return fmt.Errorf("create task in database: %w", err)
		}

		fmt.Printf("✓ Created task %s: %s\n", taskID, title)
		fmt.Printf("  Type: %s | Priority: %s | Estimated: %.1f hours\n", taskType, priority, estimatedHours)
		fmt.Printf("  File: %s\n", filePath)

		return nil
	},
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func init() {
	taskCmd.AddCommand(taskCreateCmd)
	taskCreateCmd.Flags().StringP("docs-dir", "d", "docs", "Documentation directory")
}
