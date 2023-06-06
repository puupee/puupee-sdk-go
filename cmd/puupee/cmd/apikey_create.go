/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/puupee/puupee-api-go"
	puupeesdk "github.com/puupee/puupee-sdk-go"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "创建ApiKey",
	Run: func(cmd *cobra.Command, args []string) {
		dto := puupee.NewApiKeyCreateDto(cmd.Flag("name").Value.String(), cmd.Flag("key").Value.String())
		err := puupeesdk.NewSdk().ApiKey.Create(*dto)
		cobra.CheckErr(err)
	},
}

func init() {
	apikeyCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createCmd.Flags().StringP("name", "n", "", "name of the apikey")
	createCmd.Flags().StringP("key", "k", "", "key of the apikey")
}
