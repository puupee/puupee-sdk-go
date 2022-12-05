/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/puupee/puupee-sdk-go/cli"
	"github.com/spf13/cobra"
)

var releaseNewPayload *cli.CreateReleasePayload

// releaseNewCmd represents the releaseNew command
var releaseNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new app release",
	Run: func(cmd *cobra.Command, args []string) {
		err := cli.NewpuupeeCli().
			ReleaseOp.
			Create(releaseNewPayload)
		cobra.CheckErr(err)
	},
}

func init() {
	releaseCmd.AddCommand(releaseNewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// releaseNewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// releaseNewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	releaseNewPayload = &cli.CreateReleasePayload{}

	releaseNewCmd.Flags().StringVar(&releaseNewPayload.AppName, "app-name", "", "App name")
	releaseNewCmd.Flags().StringVar(&releaseNewPayload.Version, "version", "", "App id")
	releaseNewCmd.Flags().StringVar(&releaseNewPayload.Notes, "notes", "", "App id")
	releaseNewCmd.Flags().StringVar(&releaseNewPayload.Platform, "platform", "", "App id")
	releaseNewCmd.Flags().StringVar(&releaseNewPayload.ProductType, "product-type", "", "App id")
	releaseNewCmd.Flags().StringVar(&releaseNewPayload.Channel, "channel", "", "App id")
	releaseNewCmd.Flags().StringVar(&releaseNewPayload.Environment, "env", "", "App id")
	releaseNewCmd.Flags().BoolVar(&releaseNewPayload.IsEnabled, "enabled", true, "App id")
	releaseNewCmd.Flags().BoolVar(&releaseNewPayload.IsForceUpdate, "force-update", false, "App id")
	releaseNewCmd.Flags().StringVar(&releaseNewPayload.Filepath, "file", "", "App id")
}
