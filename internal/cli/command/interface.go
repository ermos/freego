package command

import "github.com/spf13/cobra"

type Command interface {
	Flags(cmd *cobra.Command)
	Execute(cmd *cobra.Command, args []string) error
}
