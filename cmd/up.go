package cmd

import (
	"github.com/ermos/progo/internal/cli/command"
	"github.com/spf13/cobra"
)

var upCmdHandler command.Up

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Create and start reverse proxy for domains inside configuration file",
	Long:  `Create and start reverse proxy for domains inside configuration file`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := upCmdHandler.Execute(cmd, args); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
	upCmdHandler.Flags(upCmd)
}
