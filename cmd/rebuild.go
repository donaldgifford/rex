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

	"github.com/donaldgifford/rex/internal/db"
	"github.com/spf13/cobra"
)

// rebuildCmd represents the rebuild command
var rebuildCmd = &cobra.Command{
	Use:   "rebuild",
	Short: "Rebuild the SQLite cache from markdown files",
	Long: `Rebuild scans all markdown files in the docs/ directory and rebuilds
the SQLite cache (.rex.db).

This command:
- Scans docs/adr/, docs/rfc/, docs/tasks/, docs/plans/
- Parses YAML frontmatter from each markdown file
- Updates the SQLite cache with all document metadata
- Sets the last rebuild timestamp

The cache is automatically rebuilt when stale, but this command can be used
to force a rebuild if needed.

Example:
  rex rebuild
  rex rebuild --docs-dir ./documentation`,
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

		// Check if docs directory exists
		if _, err := os.Stat(docsDir); os.IsNotExist(err) {
			return fmt.Errorf("docs directory does not exist: %s", docsDir)
		}

		// Open database
		dbPath := filepath.Join(docsDir, "..", ".rex.db")
		database, err := db.Open(dbPath)
		if err != nil {
			return fmt.Errorf("open database: %w", err)
		}
		defer database.Close()

		fmt.Printf("Rebuilding cache from %s...\n", docsDir)

		// Rebuild cache
		if err := database.Rebuild(docsDir); err != nil {
			return fmt.Errorf("rebuild cache: %w", err)
		}

		fmt.Println("✓ Cache rebuilt successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rebuildCmd)

	rebuildCmd.Flags().StringP("docs-dir", "d", "docs", "Documentation directory to scan")
}
