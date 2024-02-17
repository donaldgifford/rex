/*
Copyright © 2024 Donald Gifford <dgifford06@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/donaldgifford/rex/src"
	"github.com/spf13/cobra"
)

var forceOverwrite bool

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new ADR structure in the current project directory",
	Long: `This command creates a new ADR structure in your current project directory. 
  It will use the .rex.yaml file for it's init settings. If there is not a .rex.yaml
  file in the project, it will create one. Also, the default does not enable GitHub Pages:
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
		if forceOverwrite {
			fmt.Println("Overwriting existing files")
			src.Init(forceOverwrite)
		} else {
			fmt.Println("Creating new files")
			src.Init(forceOverwrite)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	initCmd.Flags().BoolVarP(&forceOverwrite, "force", "f", false, "Overwrite existing files create with command")
}
