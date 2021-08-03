package model

import (
	"github.com/Molsbee/http-cli/validator"
	"strings"
)

var validMethods = list{
	"GET",
	"PUT",
	"POST",
	"DELETE",
	"HEAD",
}

type list []string

func (v list) Contains(method string) bool {
	for _, val := range v {
		if val == strings.ToUpper(method) {
			return true
		}
	}
	return false
}

type RequestFile struct {
	Name      string            `json:"name"`
	Variables map[string]string `json:"variables"`
	Requests  []Request         `json:"requests"`
}

func (rf RequestFile) Validate() map[interface{}][]string {
	errors := map[interface{}][]string{}
	for i, req := range rf.Requests {
		validationErrors := req.Validate()
		if len(validationErrors) != 0 {
			if len(req.Name) != 0 {
				errors[req.Name] = validationErrors
			} else {
				errors[i] = validationErrors
			}
		}
	}
	return errors
}

type Request struct {
	Name        string            `json:"name"`
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	Headers     []string          `json:"headers"`
	Data        string            `json:"data"`
	Parse       map[string]string `json:"parse"`
	ExitOnError bool              `json:"exitOnError"`
}

func (r Request) Validate() []string {
	var errors []string
	if len(r.Name) == 0 {
		errors = append(errors, "please provide a name for your request")
	}
	if len(r.Method) == 0 || !validMethods.Contains(r.Method) {
		errors = append(errors, "please provide a valid method [GET, PUT, POST, DELETE, PATCH, HEAD]")
	}
	if len(r.URL) == 0 || !validator.IsValidURL(r.URL) {
		errors = append(errors, "please provide a valid url")
	}

	return errors
}
