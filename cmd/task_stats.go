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

	"github.com/donaldgifford/rex/internal/db"
	"github.com/spf13/cobra"
)

var taskStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show task statistics",
	Long: `Show comprehensive statistics about tasks.

Displays:
- Total task count
- Breakdown by type (core, plugin, ui, other)
- Breakdown by status (planned, in_progress, blocked, completed, cancelled)
- Breakdown by priority (P0, P1, P2, P3)
- Time tracking (estimated vs actual hours)

Example:
  rex task stats`,
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

		// Get stats
		stats, err := database.GetTaskStats()
		if err != nil {
			return fmt.Errorf("get stats: %w", err)
		}

		// Print header
		fmt.Println("Task Statistics")
		fmt.Println("===============")
		fmt.Println()

		// Total
		fmt.Printf("Total Tasks: %d\n", stats.Total)
		fmt.Println()

		// By Type
		if len(stats.ByType) > 0 {
			fmt.Println("By Type:")
			typeOrder := []string{"core", "plugin", "ui", "other"}
			for _, taskType := range typeOrder {
				if count, ok := stats.ByType[taskType]; ok {
					percentage := float64(count) / float64(stats.Total) * 100
					bar := getProgressBar(count, stats.Total, 20)
					fmt.Printf("  %-8s : %3d tasks (%5.1f%%) %s\n", taskType, count, percentage, bar)
				}
			}
			fmt.Println()
		}

		// By Status
		if len(stats.ByStatus) > 0 {
			fmt.Println("By Status:")
			statusOrder := []string{"planned", "in_progress", "blocked", "completed", "cancelled"}
			for _, status := range statusOrder {
				if count, ok := stats.ByStatus[status]; ok {
					percentage := float64(count) / float64(stats.Total) * 100
					bar := getProgressBar(count, stats.Total, 20)
					emoji := getTaskStatusEmoji(status)
					fmt.Printf("  %s %-12s : %3d tasks (%5.1f%%) %s\n", emoji, status, count, percentage, bar)
				}
			}
			fmt.Println()
		}

		// By Priority
		if len(stats.ByPriority) > 0 {
			fmt.Println("By Priority:")
			priorityOrder := []string{"P0", "P1", "P2", "P3"}
			for _, priority := range priorityOrder {
				if count, ok := stats.ByPriority[priority]; ok {
					percentage := float64(count) / float64(stats.Total) * 100
					bar := getProgressBar(count, stats.Total, 20)
					label := getPriorityLabel(priority)
					fmt.Printf("  %-14s : %3d tasks (%5.1f%%) %s\n", label, count, percentage, bar)
				}
			}
			fmt.Println()
		}

		// Time Tracking
		fmt.Println("Time Tracking:")
		fmt.Printf("  Estimated : %.1f hours\n", stats.TotalHours.Estimated)
		fmt.Printf("  Actual    : %.1f hours\n", stats.TotalHours.Actual)
		if stats.TotalHours.Estimated > 0 {
			remaining := stats.TotalHours.Estimated - stats.TotalHours.Actual
			if remaining > 0 {
				fmt.Printf("  Remaining : %.1f hours\n", remaining)
			}
		}

		return nil
	},
}

func getTaskStatusEmoji(status string) string {
	switch status {
	case "planned":
		return "📋"
	case "in_progress":
		return "🔨"
	case "blocked":
		return "🚫"
	case "completed":
		return "✅"
	case "cancelled":
		return "❌"
	default:
		return "  "
	}
}

func getPriorityLabel(priority string) string {
	switch priority {
	case "P0":
		return "P0 (blocker)"
	case "P1":
		return "P1 (high)"
	case "P2":
		return "P2 (medium)"
	case "P3":
		return "P3 (low)"
	default:
		return priority
	}
}

func getProgressBar(value, total, width int) string {
	if total == 0 {
		return ""
	}

	filled := int(float64(value) / float64(total) * float64(width))
	if filled > width {
		filled = width
	}

	bar := "["
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	bar += "]"

	return bar
}

func init() {
	taskCmd.AddCommand(taskStatsCmd)
	taskStatsCmd.Flags().StringP("docs-dir", "d", "docs", "Documentation directory")
}
