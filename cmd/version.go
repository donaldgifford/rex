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
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// Version is set by goreleaser at build time
	Version = "dev"
	// Commit is set by goreleaser at build time
	Commit = "none"
	// Date is set by goreleaser at build time
	Date = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long: `Display version information for rex.

Shows the version, git commit, build date, and Go version.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("rex version %s\n", Version)
		fmt.Printf("  commit: %s\n", Commit)
		fmt.Printf("  built: %s\n", Date)
		fmt.Printf("  go: %s\n", runtime.Version())
		fmt.Printf("  platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
