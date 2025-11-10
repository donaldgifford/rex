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
	"strings"
	"time"

	"github.com/donaldgifford/rex/internal/db"
	"github.com/donaldgifford/rex/internal/generator"
	"github.com/donaldgifford/rex/internal/templates"
	"github.com/spf13/cobra"
)

var planCreateCmd = &cobra.Command{
	Use:   "create [title]",
	Short: "Create a new plan",
	Long: `Create a new plan or roadmap document.

This command will:
1. Generate a unique plan number
2. Create a markdown file from the template
3. Update the SQLite cache

Plans are useful for:
- Project roadmaps
- Sprint plans
- Release plans
- Quarterly objectives

Example:
  rex plan create "Q4 2025 Roadmap"
  rex plan create  # Interactive mode - will prompt for title`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get title
		var title string
		if len(args) > 0 {
			title = strings.Join(args, " ")
		} else {
			// Interactive mode
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Plan Title: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("read title: %w", err)
			}
			title = strings.TrimSpace(input)
		}

		if title == "" {
			return fmt.Errorf("title is required")
		}

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

		// Get next plan number
		number, err := generator.GetNextPlanNumber(database.Conn())
		if err != nil {
			return fmt.Errorf("get next plan number: %w", err)
		}

		// Generate filename
		filename := generator.GeneratePlanFilename(number, title)
		planDir := filepath.Join(docsDir, "plans")
		filePath := filepath.Join(planDir, filename)

		// Ensure plan directory exists
		if err := os.MkdirAll(planDir, 0755); err != nil {
			return fmt.Errorf("create plans directory: %w", err)
		}

		// Render template
		date := time.Now().Format("2006-01-02")
		content, err := templates.RenderPlan(templates.PlanData{
			Number:      number,
			Title:       title,
			CreatedDate: date,
		})
		if err != nil {
			return fmt.Errorf("render template: %w", err)
		}

		// Write file
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("write file: %w", err)
		}

		// Insert into database
		plan := &db.Plan{
			Number:      number,
			Title:       title,
			Status:      "draft",
			CreatedDate: date,
			FilePath:    filePath,
			Content:     content,
		}
		if err := database.CreatePlan(plan); err != nil {
			return fmt.Errorf("create plan in database: %w", err)
		}

		fmt.Printf("✓ Created plan %04d: %s\n", number, title)
		fmt.Printf("  File: %s\n", filePath)

		return nil
	},
}

func init() {
	planCmd.AddCommand(planCreateCmd)
	planCreateCmd.Flags().StringP("docs-dir", "d", "docs", "Documentation directory")
}
