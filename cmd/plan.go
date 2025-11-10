/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Manage plans and roadmaps",
	Long: `Manage plans and roadmaps for your project.

Plans are useful for documenting:
- Project roadmaps
- Sprint plans
- Release plans
- Quarterly objectives

Available subcommands:
  create  - Create a new plan
  list    - List all plans
  update  - Update plan README files

Examples:
  rex plan create "Q4 2025 Roadmap"
  rex plan list --status active
  rex plan update`,
}

func init() {
	rootCmd.AddCommand(planCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// planCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// planCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
