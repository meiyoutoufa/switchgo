package cmd

import (
	"github.com/spf13/cobra"
)

func GetRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "switchgo",
		Short: "switchgo",
		Long:  "Tools to quickly switch between go environments",
	}

	cmd.AddCommand(switchGoCmd)

	return cmd
}
