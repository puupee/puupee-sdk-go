/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// apikeyCmd represents the apikey command
var apikeyCmd = &cobra.Command{
	Use:   "apikey",
	Short: "管理ApiKey",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("apikey called")
	},
}

func init() {
	rootCmd.AddCommand(apikeyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// apikeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// apikeyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
