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
	"fmt"

	"github.com/donaldgifford/rex/internal/config"
	"github.com/spf13/cobra"
)

var (
	directoriesCreate bool
	indexCreate       bool
)

// configCreateCmd represents the configCreate command
var configCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a rex.yaml config file",
	Long: `Create a default rex.yaml config file. For example:

rex config create

default config file created in the current directory. Ideally, you put it in 
your project root level directory. Creates .rex.yaml file.`,

	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("configCreate called")

		if install {
			configFile := config.NewRexConfInstall()
			err := configFile.GenerateConfig(force)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

		} else {
			fmt.Println("Please use `--install` when running this command")
		}
	},
}

func init() {
	configCmd.AddCommand(configCreateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCreateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCreateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	configCreateCmd.Flags().BoolVarP(&force, "force", "f", false, "force overwritting config")
}
