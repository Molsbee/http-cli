package cmd

import (
	"fmt"
	"github.com/Molsbee/http-cli/service"
	"github.com/spf13/cobra"
)

var delete = &cobra.Command{
	Use:   "delete",
	Short: "perform delete request to supplied url",
	Run: func(cmd *cobra.Command, args []string) {
		client := service.NewRestClient(include, prettyPrint)
		response, err := client.Delete(args[0], data, headers)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(response)
	},
}
