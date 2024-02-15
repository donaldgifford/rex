/*
Copyright © 2024 Donald Gifford <dgifford06@gmail.com>
*/
package cmd

import (
	"time"

	"github.com/donaldgifford/rex/src"
	"github.com/spf13/cobra"
)

var (
	title  string
	author string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new ADR record",
	Long: `Create a new ADR in the path specifed in the .rex.yaml config. For example:

  rex create -t "My ADR Title" -a "Donald Gifford"
`,
	Run: func(cmd *cobra.Command, args []string) {
		myADR := src.ADR{
			Title:  title,
			Author: author,
			Status: "Draft",
			Date:   time.Now().Format(time.DateOnly),
		}
		src.CreateADR(myADR)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createCmd.Flags().StringVarP(&title, "title", "t", "", "Title for ADR.")
	createCmd.Flags().StringVarP(&author, "author", "a", "", "Author for ADR.")
}
