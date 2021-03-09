package cmd

import (
	"fmt"
	"github.com/Molsbee/http-cli/service"
	"github.com/spf13/cobra"
)

var get = &cobra.Command{
	Use:   "get [url]",
	Short: "perform get request to supplied url",
	Run: func(cmd *cobra.Command, args []string) {
		client := service.NewRestClient(include, verbose, prettyPrint)
		response, err := client.Get(args[0], headers)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Print(response)
	},
}
