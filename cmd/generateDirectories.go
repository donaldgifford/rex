/*
Copyright © 2024-2025 Donald Gifford <dgifford06@gmail.com>

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
	"github.com/donaldgifford/rex/internal/rex"
	"github.com/spf13/cobra"
)

// generateDirectoriesCmd represents the generateDirectories command
var generateDirectoriesCmd = &cobra.Command{
	Use:   "directories",
	Short: "Create directories used by rex",
	Long: `directories subcommand will create the directories needed 
for rex based on your .rex.yaml config file. For example:

Create directories listed from .rex.yaml config file:
  rex config generate directories

If "templates.enabled: true" in your .rex.yaml config file then the 
directories subcommand will create the directories listed in your config 
at "templates.path".`,
	Run: func(cmd *cobra.Command, args []string) {
		rex := rex.New()

		err := rex.GenerateDirectories()
		if err != nil {
			cmd.Println(err.Error())
		}
	},
}

func init() {
	configGenerateCmd.AddCommand(generateDirectoriesCmd)
}
