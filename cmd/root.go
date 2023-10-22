package cmd

import (
	"fmt"
	"github.com/ermos/freego/internal/pkg/config"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "freego",
	Short: "Reverse proxy for development environment",
	Long:  `Reverse proxy for development environment`,
}

func Execute() {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				fmt.Println(err.Error())
			} else {
				fmt.Println(r)
			}
			os.Exit(0)
		}
	}()

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func init() {
	cobra.OnInitialize(func() {
		if err := config.Init(); err != nil {
			panic(err)
		}
	})
}
