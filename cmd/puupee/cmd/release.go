/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	puupeesdk "github.com/puupee/puupee-sdk-go"
	"github.com/spf13/cobra"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "获取AppRelease列表",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := puupeesdk.NewSdk().Release.List(cmd.Flag("name").Value.String())
		cobra.CheckErr(err)
		if len(result.Items) > 0 {
			puupeesdk.PrintArray(result.Items)
		} else {
			fmt.Println("暂时没有版本发布")
		}
	},
}

func init() {
	appCmd.AddCommand(releaseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// releaseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// releaseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	releaseCmd.Flags().StringP("name", "n", "", "name of the app")
}
