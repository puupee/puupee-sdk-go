/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/puupee/puupee-sdk-go/cli"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "获取单个App信息",
	Run: func(cmd *cobra.Command, args []string) {
		appDto, err := cli.NewpuupeeCli().AppOp.Get(cmd.Flag("name").Value.String())
		cobra.CheckErr(err)
		cli.PrintObject(appDto)
	},
}

func init() {
	appCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	getCmd.Flags().StringP("name", "n", "", "name of the app")
}
