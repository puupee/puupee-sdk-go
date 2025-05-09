/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"path/filepath"

	puupeesdk "github.com/puupee/puupee-sdk-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfg     = puupeesdk.NewConfig()
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "puupee",
	Short: "小汪助理的命令行版本",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.puupee.yaml)")

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	viper.SetConfigFile(filepath.Join(home, ".puupee.yaml"))

	err = viper.Unmarshal(&cfg)
	cobra.CheckErr(err)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
			cobra.CheckErr(err)
		}
	}

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.PersistentFlags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVar(&cfg.Env, "env", "prod", "Api host")
	rootCmd.PersistentFlags().StringVar(&cfg.Host, "host", "api.puupee.com", "Api host")
	rootCmd.PersistentFlags().StringVar(&cfg.ApiKey, "api-key", "", "Api key")

	viper.BindPFlag("env", rootCmd.Flags().Lookup("env"))
	viper.BindPFlag("host", rootCmd.Flags().Lookup("host"))
	viper.BindPFlag("api-key", rootCmd.Flags().Lookup("api-key"))
}
