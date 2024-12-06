package cmd

import (
	"freb/utils"
	"github.com/spf13/cobra"
)

// vars below are set by '-X'
var (
	version = "dev"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "版本",
	Run: func(cmd *cobra.Command, args []string) {
		utils.SysInfof("version: %s", version)
	},
}
