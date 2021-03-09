package cmd

import (
	"fmt"
	"github.com/Molsbee/http-cli/service"
	"github.com/spf13/cobra"
)

var put = &cobra.Command{
	Use:   "put",
	Short: "perform put request to supplied url",
	Run: func(cmd *cobra.Command, args []string) {
		client := service.NewRestClient(include, verbose, prettyPrint)
		response, err := client.Put(args[0], data, headers)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(response)
	},
}
