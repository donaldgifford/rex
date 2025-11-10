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
	"os"
	"path/filepath"
	"time"

	"github.com/donaldgifford/rex/internal/db"
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database management commands",
	Long: `Manage the SQLite database cache.

Available subcommands:
  info  - Show database statistics and cache info

Examples:
  rex db info`,
}

var dbInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show database cache information",
	Long: `Display information about the SQLite cache database.

Shows:
- Database file path and size
- Last rebuild timestamp
- Document counts by type
- Schema version

Example:
  rex db info`,
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

		// Database path
		dbPath := filepath.Join(docsDir, "..", ".rex.db")

		// Check if database exists
		dbInfo, err := os.Stat(dbPath)
		if os.IsNotExist(err) {
			fmt.Println("Database cache does not exist yet.")
			fmt.Printf("Run 'rex rebuild' to create it.\n")
			return nil
		}
		if err != nil {
			return fmt.Errorf("stat database: %w", err)
		}

		// Open database
		database, err := db.Open(dbPath)
		if err != nil {
			return fmt.Errorf("open database: %w", err)
		}
		defer database.Close()

		fmt.Println("Database Cache Information")
		fmt.Println("==========================")
		fmt.Println()

		// File info
		fmt.Printf("Path: %s\n", dbPath)
		fmt.Printf("Size: %.2f KB\n", float64(dbInfo.Size())/1024)
		fmt.Printf("Modified: %s\n", dbInfo.ModTime().Format("2006-01-02 15:04:05"))
		fmt.Println()

		// Last rebuild
		lastRebuild, err := database.GetLastRebuild()
		if err != nil {
			return fmt.Errorf("get last rebuild: %w", err)
		}
		if !lastRebuild.IsZero() {
			fmt.Printf("Last Rebuild: %s\n", lastRebuild.Format("2006-01-02 15:04:05"))
			fmt.Printf("Time Since: %s ago\n", time.Since(lastRebuild).Round(time.Second))
		} else {
			fmt.Println("Last Rebuild: Never")
		}
		fmt.Println()

		// Schema version
		schemaVersion, err := database.GetMetadata("schema_version")
		if err == nil && schemaVersion != "" {
			fmt.Printf("Schema Version: %s\n", schemaVersion)
			fmt.Println()
		}

		// Document counts
		fmt.Println("Document Counts:")

		// ADRs
		adrStats, err := database.GetADRStats()
		if err != nil {
			return fmt.Errorf("get ADR stats: %w", err)
		}
		adrTotal := 0
		for _, count := range adrStats {
			adrTotal += count
		}
		fmt.Printf("  ADRs  : %d\n", adrTotal)
		if len(adrStats) > 0 {
			for status, count := range adrStats {
				fmt.Printf("    - %s: %d\n", status, count)
			}
		}

		// RFCs
		rfcStats, err := database.GetRFCStats()
		if err != nil {
			return fmt.Errorf("get RFC stats: %w", err)
		}
		rfcTotal := 0
		for _, count := range rfcStats {
			rfcTotal += count
		}
		fmt.Printf("  RFCs  : %d\n", rfcTotal)
		if len(rfcStats) > 0 {
			for status, count := range rfcStats {
				fmt.Printf("    - %s: %d\n", status, count)
			}
		}

		// Tasks
		taskStats, err := database.GetTaskStats()
		if err != nil {
			return fmt.Errorf("get task stats: %w", err)
		}
		fmt.Printf("  Tasks : %d\n", taskStats.Total)
		if len(taskStats.ByStatus) > 0 {
			for status, count := range taskStats.ByStatus {
				fmt.Printf("    - %s: %d\n", status, count)
			}
		}

		// Plans
		planStats, err := database.GetPlanStats()
		if err != nil {
			return fmt.Errorf("get plan stats: %w", err)
		}
		planTotal := 0
		for _, count := range planStats {
			planTotal += count
		}
		fmt.Printf("  Plans : %d\n", planTotal)
		if len(planStats) > 0 {
			for status, count := range planStats {
				fmt.Printf("    - %s: %d\n", status, count)
			}
		}

		fmt.Println()
		totalDocs := adrTotal + rfcTotal + taskStats.Total + planTotal
		fmt.Printf("Total Documents: %d\n", totalDocs)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(dbInfoCmd)
	dbInfoCmd.Flags().StringP("docs-dir", "d", "docs", "Documentation directory")
}
