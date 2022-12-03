/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// keyvalueCmd represents the keyvalue command
var keyvalueCmd = &cobra.Command{
	Use:   "keyvalue",
	Short: "键值对",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("keyvalue called")
	},
}

func init() {
	rootCmd.AddCommand(keyvalueCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// keyvalueCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// keyvalueCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
