/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	puupeesdk "github.com/puupee/puupee-sdk-go"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "应用列表",
	Run: func(cmd *cobra.Command, args []string) {
		c := puupeesdk.NewSdk()
		resp, err := c.App.List()
		if err != nil {
			fmt.Println(err)
			return
		}
		if *resp.TotalCount > 0 {
			puupeesdk.PrintArray(resp.Items)
		}
	},
}

func init() {
	appCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
