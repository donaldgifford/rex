/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// rfcCmd represents the rfc command
var rfcCmd = &cobra.Command{
	Use:   "rfc",
	Short: "Manage RFCs (Requests for Comments)",
	Long: `Manage RFCs (Requests for Comments) for your project.

RFCs document proposed changes, features, or architectural decisions that
require discussion and approval.

Available subcommands:
  create  - Create a new RFC
  list    - List all RFCs
  update  - Update RFC README files

Examples:
  rex rfc create "Add Plugin System"
  rex rfc list --status draft
  rex rfc update`,
}

func init() {
	rootCmd.AddCommand(rfcCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rfcCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rfcCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
