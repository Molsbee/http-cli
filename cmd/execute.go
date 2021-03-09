package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Molsbee/http-cli/model"
	"github.com/Molsbee/http-cli/service"
	"github.com/gookit/color"
	"github.com/savaki/jq"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

var executeCmd = &cobra.Command{
	Use:  "execute",
	Long: `executes a yaml formatted file`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 || len(args[0]) == 0 {
			return errors.New("please provide a path to a valid yaml file")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		b, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Printf("failed to read file (%s) - %s", filePath, err)
			os.Exit(1)
		}

		var requestFile model.RequestFile
		if err := yaml.Unmarshal(b, &requestFile); err != nil {
			fmt.Printf("failed to convert file to expected format - %s", err)
			os.Exit(1)
		}

		variables := make(map[string]string)
		client := service.NewRestClient(include, verbose, true)
		for i := 0; i < len(requestFile.Requests); i++ {
			request := requestFile.Requests[i]
			color.FgLightRed.Printf("##### %s #####\n", request.Name)

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

				yaml.Unmarshal(buffer.Bytes(), &requestFile)
			}

			fmt.Println(resp)
			color.FgLightRed.Println("##### END OF REQUEST #####")
		}
	},
}

func executeRequest(client service.RestClient, request model.Request) (string, error) {
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

func updateVariables(variables map[string]string, request model.Request, resp string) bool {
	for k, v := range request.Parse {
		op, _ := jq.Parse(v)
		value, _ := op.Apply([]byte(resp))
		variables[k] = strings.TrimSuffix(strings.TrimPrefix(string(value), "\""), "\"")
		return true
	}

	return false
}
