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

	"github.com/spf13/cobra"
)

var (
	records int
	filter  string
	adrType string
)

// adrListCmd represents the adrList command
var adrListCmd = &cobra.Command{
	Use:   "list",
	Short: "List ADR's",
	Long: `List the ADR's in console. For example:

Return list in JSON format:
rex adr list -t json 

Return list in markdown table:
rex adr list -t md

Return oldest 5 records:
rex adr list -r 5 -f oldest

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("adrList called")
	},
}

func init() {
	adrCmd.AddCommand(adrListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// adrListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// adrListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	adrListCmd.Flags().IntVarP(&records, "records", "r", 20, "Number of records returned")

	adrListCmd.Flags().StringVarP(&filter, "filter", "f", "", "Filter records - oldest, newest")
	adrListCmd.Flags().StringVarP(&adrType, "type", "t", "", "Output type - json, md")

	// TODO: eventually add in a sort or filter to return newest, oldest, or by ref number
}
