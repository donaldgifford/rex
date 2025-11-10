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

var rfcListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all RFCs",
	Long: `List all RFCs (Requests for Comments) with optional filtering.

Examples:
  rex rfc list                     # List all RFCs
  rex rfc list --status draft      # List only draft RFCs
  rex rfc list --status approved   # List only approved RFCs`,
	RunE: func(cmd *cobra.Command, args []string) error {
		docsDir, err := cmd.Flags().GetString("docs-dir")
		if err != nil {
			return err
		}

		status, err := cmd.Flags().GetString("status")
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

		// List RFCs
		rfcs, err := database.ListRFCs(status)
		if err != nil {
			return fmt.Errorf("list RFCs: %w", err)
		}

		if len(rfcs) == 0 {
			if status != "" {
				fmt.Printf("No RFCs found with status: %s\n", status)
			} else {
				fmt.Println("No RFCs found")
			}
			return nil
		}

		// Get stats for header
		stats, err := database.GetRFCStats()
		if err != nil {
			return fmt.Errorf("get stats: %w", err)
		}

		// Print header
		fmt.Println("RFCs (Requests for Comments)")
		fmt.Println("=============================")
		fmt.Println()

		// Print stats
		total := 0
		for _, count := range stats {
			total += count
		}
		fmt.Printf("Total: %d RFCs\n", total)
		if len(stats) > 0 {
			fmt.Print("  ")
			first := true
			for status, count := range stats {
				if !first {
					fmt.Print(", ")
				}
				fmt.Printf("%s: %d", status, count)
				first = false
			}
			fmt.Println()
		}
		fmt.Println()

		// Print table
		fmt.Printf("%-6s %-14s %-20s %s\n", "Number", "Status", "Author", "Title")
		fmt.Println("────────────────────────────────────────────────────────────────────────────")

		for _, rfc := range rfcs {
			statusDisplay := getRFCStatusDisplay(rfc.Status)
			authorDisplay := rfc.Author
			if authorDisplay == "" {
				authorDisplay = "-"
			}
			if len(authorDisplay) > 20 {
				authorDisplay = authorDisplay[:17] + "..."
			}
			fmt.Printf("%-6d %-14s %-20s %s\n", rfc.Number, statusDisplay, authorDisplay, rfc.Title)
		}

		return nil
	},
}

func getRFCStatusDisplay(status string) string {
	switch status {
	case "draft":
		return "📝 draft"
	case "review":
		return "👀 review"
	case "approved":
		return "✅ approved"
	case "implemented":
		return "🚀 implemented"
	case "rejected":
		return "❌ rejected"
	default:
		return status
	}
}

func init() {
	rfcCmd.AddCommand(rfcListCmd)
	rfcListCmd.Flags().StringP("docs-dir", "d", "docs", "Documentation directory")
	rfcListCmd.Flags().StringP("status", "s", "", "Filter by status (draft, review, approved, implemented, rejected)")
}
