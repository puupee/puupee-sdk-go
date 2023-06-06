/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "设置",
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, err := cmd.Flags().GetString("api-key")
		if err == nil {
			viper.Set("apiKey", apiKey)
			viper.Set("apiKeys."+viper.GetString("env"), apiKey)
			err := viper.WriteConfig()
			cobra.CheckErr(err)
		}
	},
}

func init() {
	configCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	setCmd.Flags().String("apiKey", "", "apiKey")
}
