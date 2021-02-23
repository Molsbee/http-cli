package cmd

import "github.com/spf13/cobra"

var head = &cobra.Command{
	Use:   "head",
	Short: "perform head request to supplied url",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
