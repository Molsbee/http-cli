package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"net/url"
	"os"
	"strings"
)

var headers map[string]string
var stringHeaders []string
var include bool

func init() {
	rootCmd.AddCommand(head, get, put, post, delete)
	rootCmd.PersistentFlags().StringSliceVarP(&stringHeaders, "headers", "H", []string{}, `-H "Content-Type: application/json"`)
	rootCmd.PersistentFlags().BoolVarP(&include, "include", "i", false, "Includes the response headers in the output")
}

var rootCmd = &cobra.Command{
	Use:   "http",
	Short: "http is a simplified command line client for performing http requests",
	Long: `http is a simplified command line client for performing http requests.
Default scheme is http:// and can be omitted from url parameter`,
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Validate the URL passed in and normalize to default scheme if needed
		if len(args) == 0 || len(args[0]) == 0 {
			return errors.New("please provide a url")
		}

		if !strings.Contains(args[0], "://") {
			args[0] = "http://" + args[0]
		}

		if _, err := url.Parse(args[0]); err != nil {
			return errors.New("please provide a valid url")
		}

		// Parse Headers into something better
		if len(stringHeaders) != 0 {
			headers = make(map[string]string)
			for _, v := range stringHeaders {
				parts := strings.Split(v, ":")
				headers[parts[0]] = strings.TrimSpace(parts[1])
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
