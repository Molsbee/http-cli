package cmd

import (
	"fmt"
	"github.com/Molsbee/http-cli/service"
	"github.com/spf13/cobra"
)

var head = &cobra.Command{
	Use:   "head",
	Short: "perform head request to supplied url",
	Run: func(cmd *cobra.Command, args []string) {
		client := service.NewRestClient(include, verbose, prettyPrint)
		response, err := client.Head(args[0], headers)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(response.Body)
	},
}
