package cmd

import (
	"freb/models"
	"freb/utils/stdout"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "版本",
	Run: func(cmd *cobra.Command, args []string) {
		stdout.SysInfofln("version: %s", models.Version)
	},
}
