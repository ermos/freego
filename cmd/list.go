package cmd

import (
	"github.com/ermos/freego/internal/cli/command"

	"github.com/spf13/cobra"
)

var listCmdHandler = command.List{}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all active domains",
	Long:  `List all active domains`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := listCmdHandler.Execute(cmd, args); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmdHandler.Flags(listCmd)
}
