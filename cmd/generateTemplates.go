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
	"github.com/spf13/viper"
)

// generateTemplatesCmd represents the generateTemplates command
var generateTemplatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "create templates for rex",
	Long: `templates subcommand creates the default templates files for rex.
This subcommand only works if "templates.enabled: true" in your .rex.yaml 
config file.

Create templates listed from .rex.yaml config file:
  rex config generate templates 

Passing '--force, -f' will overwrite the templates if files 
are found.`,
	Run: func(cmd *cobra.Command, args []string) {
		// check to see if templates are enabled
		// if they aren't we exit since we dont need
		// template files unless this is enabled
		if !viper.GetBool("templates.enabled") {
			e := cmd.ErrOrStderr()
			_, err := e.Write([]byte("templates.enabled not set to true\n"))
			if err != nil {
				cmd.Println(err.Error())
			}
		}
		// init rex
		rex := rex.New()

		err := rex.GenerateTemplates(force)
		if err != nil {
			cmd.Println(err.Error())
		}
	},
}

func init() {
	configGenerateCmd.AddCommand(generateTemplatesCmd)

	generateTemplatesCmd.Flags().BoolVarP(&force, "force", "f", false, "force overwritting config")
}
