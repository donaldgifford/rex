/*
Copyright © 2024 Donald Gifford <dgifford06@gmail.com>
*/
package cmd

import (
	"github.com/donaldgifford/rex/src"
	"github.com/spf13/cobra"
)

// generateConfigCmd represents the generateConfig command
var generateConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Create a rex.yaml config file",
	Long: `Create a default rex.yaml config file. For example:

rex generate config

default config file created in the current directory. Ideally, you put it in 
your project root level directory. Creates .rex.yaml file.`,
	Run: func(cmd *cobra.Command, args []string) {
		src.CreateConfig()
	},
}

func init() {
	generateCmd.AddCommand(generateConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateConfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateConfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
