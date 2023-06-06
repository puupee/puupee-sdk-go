/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	puupeesdk "github.com/puupee/puupee-sdk-go"
	"github.com/spf13/cobra"
)

var releaseNewPayload *puupeesdk.CreateReleasePayload

// releaseNewCmd represents the releaseNew command
var releaseNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new app release",
	Run: func(cmd *cobra.Command, args []string) {
		err := puupeesdk.NewSdk().
			Release.
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
	releaseNewPayload = &puupeesdk.CreateReleasePayload{}

	releaseNewCmd.Flags().StringVar(&releaseNewPayload.AppName, "app-name", "", "App name")
	releaseNewCmd.Flags().StringVar(&releaseNewPayload.Version, "version", "", "Release version")
	releaseNewCmd.Flags().StringVar(&releaseNewPayload.Notes, "notes", "", "Release notes")
	releaseNewCmd.Flags().StringVar(&releaseNewPayload.Platform, "platform", "", "Release platform")
	releaseNewCmd.Flags().StringVar(&releaseNewPayload.ProductType, "product-type", "", "Release product type")
	releaseNewCmd.Flags().StringVar(&releaseNewPayload.Channel, "channel", "", "Release channel")
	releaseNewCmd.Flags().StringVar(&releaseNewPayload.Environment, "env", "dev", "Release environment")
	releaseNewCmd.Flags().BoolVar(&releaseNewPayload.IsEnabled, "enabled", false, "Release enabled right now")
	releaseNewCmd.Flags().BoolVar(&releaseNewPayload.IsForceUpdate, "force-update", false, "Release force update")
	releaseNewCmd.Flags().StringVar(&releaseNewPayload.Filepath, "file", "", "Release product file path")

	releaseNewCmd.MarkFlagRequired("app-name")
	releaseNewCmd.MarkFlagRequired("version")
	releaseNewCmd.MarkFlagRequired("platform")
	releaseNewCmd.MarkFlagRequired("product-type")
	releaseNewCmd.MarkFlagRequired("file")
}
