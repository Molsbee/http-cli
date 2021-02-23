package cmd

import "github.com/spf13/cobra"

var post = &cobra.Command{
	Use:   "post",
	Short: "perform post request to supplied url",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
