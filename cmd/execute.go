package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Molsbee/http-cli/service"
	"github.com/savaki/jq"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

type requests struct {
	Requests []request `json:"requests"`
}

type request struct {
	Name    string            `json:"name"`
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers []string          `json:"headers"`
	Data    string            `json:"data"`
	Parse   map[string]string `json:"parse"`
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
		b, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Printf("failed to read file (%s) - %s", filePath, err)
			os.Exit(1)
		}

		var requests requests
		if err := yaml.Unmarshal(b, &requests); err != nil {
			fmt.Printf("failed to convert file to expected format - %s", err)
			os.Exit(1)
		}

		variables := make(map[string]string)
		client := service.NewRestClient(true, true)
		for i := 0; i < len(requests.Requests); i++ {
			request := requests.Requests[i]
			fmt.Printf("##### %s #####", request.Name)
			resp, clientErr := executeRequest(client, request)
			if clientErr != nil {
				fmt.Println(clientErr)
				fmt.Println("##### END OF REQUEST #####")
				continue
			}

			if updated := updateVariables(variables, request, resp); updated {
				var buffer bytes.Buffer
				temp, _ := template.ParseFiles(filePath)
				temp.Execute(&buffer, variables)

				yaml.Unmarshal(buffer.Bytes(), &requests)
			}

			fmt.Println(resp)
			fmt.Println("##### END OF REQUEST #####")
		}
	},
}

func executeRequest(client service.RestClient, request request) (string, error) {
	switch {
	case strings.EqualFold(request.Method, "GET"):
		return client.Get(request.URL, parseHeaders(request.Headers))
	case strings.EqualFold(request.Method, "PUT"):
		return client.Put(request.URL, request.Data, parseHeaders(request.Headers))
	case strings.EqualFold(request.Method, "POST"):
		return client.Post(request.URL, request.Data, parseHeaders(request.Headers))
	case strings.EqualFold(request.Method, "DELETE"):
		return client.Delete(request.URL, request.Data, parseHeaders(request.Headers))
	}

	return "unsupported method", errors.New("unsupported method")
}

func updateVariables(variables map[string]string, request request, resp string) bool {
	for k, v := range request.Parse {
		op, _ := jq.Parse(v)
		value, _ := op.Apply([]byte(resp))
		variables[k] = strings.TrimSuffix(strings.TrimPrefix(string(value), "\""), "\"")
		return true
	}

	return false
}