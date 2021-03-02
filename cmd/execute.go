package cmd

import (
	"errors"
	"fmt"
	"github.com/Molsbee/http-cli/service"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

type requests struct {
	Requests []request `json:"requests"`
}

type request struct {
	Name    string   `json:"name"`
	Method  string   `json:"method"`
	URL     string   `json:"url"`
	Headers []string `json:"headers"`
	Data    string   `json:"data"`
}

var executeCmd = &cobra.Command{
	Use:  "execute",
	Long: `executes a yaml formatted file`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || len(args[0]) == 0 {
			fmt.Println("please provide a path to a yaml file")
			os.Exit(1)
		}

		filePath := args[0]
		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Printf("failed to read file (%s) - %s", filePath, err)
			os.Exit(1)
		}

		var requests requests
		if err := yaml.Unmarshal(bytes, &requests); err != nil {
			fmt.Printf("failed to convert file to expected format - %s", err)
			os.Exit(1)
		}

		client := service.NewRestClient(true, true)
		for _, request := range requests.Requests {
			fmt.Printf("##### %s #####", request.Name)

			var resp string
			var clientErr error
			switch {
			case strings.EqualFold(request.Method, "GET"):
				resp, clientErr = client.Get(request.URL, parseHeaders(request.Headers))
			case strings.EqualFold(request.Method, "PUT"):
				resp, clientErr = client.Put(request.URL, request.Data, parseHeaders(request.Headers))
			case strings.EqualFold(request.Method, "POST"):
				resp, clientErr = client.Post(request.URL, request.Data, parseHeaders(request.Headers))
			case strings.EqualFold(request.Method, "DELETE"):
				resp, clientErr = client.Delete(request.URL, request.Data, parseHeaders(request.Headers))
			default:
				resp, clientErr = "unsupported method", errors.New("unsupported method")
			}

			if clientErr != nil {
				fmt.Println(clientErr)
			} else {
				fmt.Println(resp)
			}
			fmt.Println("##### END OF REQUEST #####")
		}
	},
}
