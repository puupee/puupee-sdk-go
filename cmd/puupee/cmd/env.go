/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	local = "localhost:44355"
	dev   = "dev.api.puupee.com"
	prod  = "api.puupee.com"
)

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "切换环境",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			currentEnv := viper.GetString("env")
			if currentEnv == "" {
				currentEnv = "local"
			}
			fmt.Printf("Current env: %s\n\n", currentEnv)
			fmt.Println("Env:")
			fmt.Printf("local: %s\n", local)
			fmt.Printf("dev  : %s\n", dev)
			fmt.Printf("prod : %s\n", prod)
			return
		}

		env := args[0]
		if env == "local" {
			viper.Set("env", "local")
			viper.Set("host", local)
			fmt.Printf("Switched to local %s\n", local)
		}
		if env == "dev" {
			viper.Set("env", "dev")
			viper.Set("host", dev)
			fmt.Printf("Switched to dev %s\n", dev)
		}
		if env == "prod" {
			viper.Set("env", "prod")
			viper.Set("host", prod)
			fmt.Printf("Switched to prod %s\n", prod)
		}
		err := viper.WriteConfig()
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(envCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// envCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// envCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
