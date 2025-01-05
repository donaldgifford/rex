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

// generateIndexCmd represents the generateIndex command
var generateIndexCmd = &cobra.Command{
	Use:   "index",
	Short: "index creates or regenerates the index for rex",
	Long: `index will regenerate the index file set for rex. If the 
index isn't found in the configured location, it will create it. 

Create index listed from .rex.yaml config file:
  rex config generate index

Regenerate index listed from .rex.yaml config file:
  rex config generate index --force

The index subcommand looks at the .rex.yaml config file to 
see where to save the index file, name, and what template to use.`,
	Run: func(cmd *cobra.Command, args []string) {
		rex := rex.New()

		err := rex.GenerateIndex(force)
		if err != nil {
			cmd.Println(err.Error())
		}
	},
}

func init() {
	configGenerateCmd.AddCommand(generateIndexCmd)

	generateIndexCmd.Flags().BoolVarP(&force, "force", "f", false, "force overwritting config")
}
