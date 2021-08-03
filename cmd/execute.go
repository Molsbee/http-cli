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
	"log"
	"strings"
	"text/template"
)

const (
	green = color.FgLightGreen
	gray  = color.FgGray
	red   = color.FgLightRed
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
			log.Panicf("failed to read file (%s) - %s", filePath, err)
		}

		var requestFile model.RequestFile
		if err := yaml.Unmarshal(b, &requestFile); err != nil {
			log.Panicf("failed to convert file to expected format - %s", err)
		}

		green.Printf("Executing: %s\n\n", requestFile.Name)
		variables := requestFile.Variables
		if len(variables) != 0 {
			updateTemplate(filePath, variables, &requestFile)
		}

		client := service.NewRestClient(include, verbose, true)
		for i := 0; i < len(requestFile.Requests); i++ {
			request := requestFile.Requests[i]
			green.Printf("→ %s\n", request.Name)
			resp, clientErr := client.Execute(request.Method, request.URL, parseHeaders(request.Headers), request.Data)
			if clientErr != nil {
				fmt.Println(clientErr)
				color.FgLightRed.Println("##### END OF REQUEST #####")
				if request.ExitOnError {
					break
				}

				continue
			}

			gray.Printf("%s %s [%s, %s, %s]\n", request.Method, request.URL, resp.Status, resp.ContentLength, resp.Duration)
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				if updated := updateVariables(variables, request, resp.Body); updated {
					updateTemplate(filePath, variables, &requestFile)
				}
			}
			printStatusCodeMessage(resp.StatusCode)

			//fmt.Println(resp.Body)
			color.FgLightRed.Println("##### END OF REQUEST #####")
		}
	},
}

func printStatusCodeMessage(statusCode int) {
	mark := red.Sprintf("✘")
	if statusCode >= 200 && statusCode < 300 {
		mark = green.Sprint("✓")
	}
	fmt.Printf("%s %s", mark, gray.Sprintf("Status Code is %d", statusCode))
}

func updateVariables(variables map[string]string, request model.Request, resp string) (updated bool) {
	for k, v := range request.Parse {
		op, _ := jq.Parse(v)
		value, _ := op.Apply([]byte(resp))
		variables[k] = strings.TrimSuffix(strings.TrimPrefix(string(value), "\""), "\"")
		updated = true
	}

	return
}

func updateTemplate(filePath string, variables map[string]string, requestFile *model.RequestFile) {
	var buffer bytes.Buffer
	temp, _ := template.ParseFiles(filePath)
	temp.Execute(&buffer, variables)
	yaml.Unmarshal(buffer.Bytes(), &requestFile)
}
