/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	puupeesdk "github.com/puupee/puupee-sdk-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// infoCmd represents the version command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "命令行信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\n", puupeesdk.Version)
		fmt.Printf("BuildTime: %s\n", puupeesdk.BuildTime)

		fmt.Printf("Env: %s\n", viper.GetString("env"))
		fmt.Printf("Host: %s\n", viper.GetString("host"))
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
