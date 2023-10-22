package cmd

import (
	"github.com/ermos/freego/internal/cli/command"

	"github.com/spf13/cobra"
)

var serveCmdHandler = command.Serve{}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the reverse proxy server",
	Long:  `Start the reverse proxy server`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := serveCmdHandler.Execute(cmd, args); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmdHandler.Flags(serveCmd)
}
