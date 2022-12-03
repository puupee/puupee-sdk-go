/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/puupee/puupee-api-go"
	"github.com/puupee/puupee-sdk-go/cli"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "新建App",
	Run: func(cmd *cobra.Command, args []string) {
		dto := puupee.NewCreateOrUpdateAppDto()
		dto.SetName(cmd.Flag("name").Value.String())
		dto.SetDisplayName(cmd.Flag("displayName").Value.String())
		dto.SetFramework(cmd.Flag("framework").Value.String())
		dto.SetAppType(cmd.Flag("appType").Value.String())
		dto.SetDescription(cmd.Flag("description").Value.String())
		dto.SetIcon(cmd.Flag("icon").Value.String())
		dto.SetGitRepository(cmd.Flag("git-repo").Value.String())
		dto.SetGitRepository(cmd.Flag("git-repo-type").Value.String())
		err := cli.NewpuupeeCli().AppOp.Create(*dto)
		cobra.CheckErr(err)
	},
}

func init() {
	appCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	newCmd.Flags().StringP("name", "n", "", "Name of app")
	newCmd.Flags().StringP("displayName", "d", "", "Display name of app")
	newCmd.Flags().StringP("framework", "f", "", "Framework of app")
	newCmd.Flags().StringP("appType", "t", "", "AppType of app")
	newCmd.Flags().StringP("description", "e", "", "Description of app")
	newCmd.Flags().StringP("icon", "i", "", "Icon of app")
	newCmd.Flags().StringP("git-repo", "g", "", "Git Repository of app")
	newCmd.Flags().StringP("git-repo-type", "b", "", "Git Repository Type of app")
}
