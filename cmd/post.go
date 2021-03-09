package cmd

import (
	"fmt"
	"github.com/Molsbee/http-cli/service"
	"github.com/spf13/cobra"
)

var post = &cobra.Command{
	Use:   "post",
	Short: "perform post request to supplied url",
	Run: func(cmd *cobra.Command, args []string) {
		client := service.NewRestClient(include, verbose, prettyPrint)
		response, err := client.Post(args[0], data, headers)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(response)
	},
}
