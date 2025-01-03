/*
Copyright © 2024 Donald Gifford <dgifford06@gmail.com>

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
	"github.com/spf13/cobra"

	"github.com/donaldgifford/rex/internal/rex"
)

var (
	directories bool
	index       bool
)

// configGenerateCmd represents the configGenerate command
var configGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate and install config files, directories, etc",
	Long: `Generate allos you to install directories and files used 
for rex. For example:

Create directories listed from .rex.yaml config file:
  rex config generate --directories

Create index listed from .rex.yaml config file:
  rex config generate --index`,
	Run: func(cmd *cobra.Command, args []string) {
		rex := rex.New()

		// if --directories and --index set
		if directories {
			err := rex.GenerateDirectories()
			if err != nil {
				cmd.Println(err.Error())
			}
		}

		// if just --index
		if index {
			err := rex.GenerateIndex()
			if err != nil {
				cmd.Println(err.Error())
			}
		}
	},
}

func init() {
	configCmd.AddCommand(configGenerateCmd)

	configGenerateCmd.Flags().BoolVarP(&force, "force", "f", false, "force overwritting config")
	configGenerateCmd.Flags().BoolVarP(&index, "index", "x", false, "create index for docs")
	configGenerateCmd.Flags().BoolVarP(&directories, "directories", "d", false, "create directories for docs")
}
