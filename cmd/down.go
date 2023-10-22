package cmd

import (
	"github.com/ermos/freego/internal/cli/command"

	"github.com/spf13/cobra"
)

var downCmdHandler command.Down

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and remove reverse proxy for domains inside configuration file",
	Long:  `Stop and remove reverse proxy for domains inside configuration file`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := downCmdHandler.Execute(cmd, args); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
	downCmdHandler.Flags(downCmd)
}
