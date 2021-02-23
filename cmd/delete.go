package cmd

import "github.com/spf13/cobra"

var delete = &cobra.Command{
	Use:   "delete",
	Short: "perform delete request to supplied url",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
