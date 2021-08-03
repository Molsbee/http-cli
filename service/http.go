package service

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/Molsbee/http-cli/model"
	"github.com/gookit/color"
	"github.com/yosssi/gohtml"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	lightBlue = color.FgLightBlue.Render
	teal      = color.FgCyan.Render
	purple    = color.Magenta.Render
)

type RestClient struct {
	client                 *http.Client
	includeResponseHeaders bool
	prettyPrint            bool
	verbose                bool
}

func NewRestClient(include, verbose, prettyPrint bool) RestClient {
	return RestClient{
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		includeResponseHeaders: include,
		verbose:                verbose,
		prettyPrint:            prettyPrint,
	}
}

func (rc RestClient) Get(url string, headers map[string]string) (model.HttpResponse, error) {
	return rc.Execute("GET", url, headers, "")
}

func (rc RestClient) Put(url, body string, headers map[string]string) (model.HttpResponse, error) {
	return rc.Execute("PUT", url, headers, body)
}

func (rc RestClient) Post(url, body string, headers map[string]string) (model.HttpResponse, error) {
	return rc.Execute("POST", url, headers, body)
}

func (rc RestClient) Delete(url string, body string, headers map[string]string) (model.HttpResponse, error) {
	return rc.Execute("DELETE", url, headers, body)
}

func (rc RestClient) Head(url string, headers map[string]string) (model.HttpResponse, error) {
	return rc.Execute("HEAD", url, headers, "")
}

func (rc RestClient) Execute(method, url string, headers map[string]string, body string) (httpResponse model.HttpResponse, err error) {
	request, requestError := rc.createRequest(method, url, headers, body)
	if requestError != nil {
		err = requestError
		return
	}

	if rc.verbose {
		fmt.Printf("%s %s %s\n", teal(method), url, rc.sPrintProto(request.ProtoMajor, request.ProtoMinor))
		for k, v := range request.Header {
			fmt.Printf("%s: %s\n", lightBlue(k), strings.Join(v, ", "))
		}
		if len(body) != 0 {
			fmt.Printf("\n%s\n\n", body)
		}
	}

	start := time.Now()
	resp, executeError := rc.client.Do(request)
	if executeError != nil {
		err = errors.New("error occurred performing request - " + executeError.Error())
		return
	}
	defer resp.Body.Close()
	dur := time.Now().Sub(start).Round(time.Millisecond)

	// Print Response Headers
	if rc.includeResponseHeaders {
		fmt.Printf("%s %s\n", rc.sPrintProto(resp.ProtoMajor, resp.ProtoMinor), purple(resp.Status))
		for k, v := range resp.Header {
			fmt.Printf("%s: %s\n", lightBlue(k), strings.Join(v, ", "))
		}
		fmt.Println()
	}

	dataBytes, readBodyError := ioutil.ReadAll(resp.Body)
	if readBodyError != nil {
		err = errors.New("error occurred reading response body - " + readBodyError.Error())
		return
	}

	httpResponse.StatusCode = resp.StatusCode
	httpResponse.Status = resp.Status
	httpResponse.ContentLength = FormatBytes(resp.ContentLength)
	httpResponse.Duration = dur.String()

	httpResponse.Body = formatResponse(resp.Header.Get("Content-Type"), dataBytes, rc.prettyPrint)
	return
}

func (rc RestClient) createRequest(method, url string, headers map[string]string, body string) (req *http.Request, err error) {
	if len(body) == 0 {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, strings.NewReader(body))
	}

	if err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return
}

func (rc RestClient) sPrintProto(protoMajor, protoMinor int) string {
	proto := color.Magenta.Sprintf("%d.%d", protoMajor, protoMinor)
	return fmt.Sprintf("%s/%s", teal("HTTP"), proto)
}

func formatResponse(contentType string, data []byte, prettyFormat bool) string {
	if prettyFormat {
		switch {
		case strings.Contains(contentType, "application/ld+json"):
		case strings.Contains(contentType, "application/json"):
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, data, "", "\t"); err == nil {
				return string(prettyJSON.Bytes())
			}
		case strings.Contains(contentType, "text/html"):
			xmlBytes := gohtml.FormatBytes(data)
			if len(xmlBytes) != 0 {
				return string(xmlBytes)
			}
		case strings.Contains(contentType, "text/xml"):
		case strings.Contains(contentType, "application/xml"):
			decoder := xml.NewDecoder(bytes.NewReader(data))

			var prettyXML bytes.Buffer
			encoder := xml.NewEncoder(&prettyXML)
			encoder.Indent("", "\t")
			for {
				token, err := decoder.Token()
				if err != nil {
					if err == io.EOF {
						encoder.Flush()
						return string(prettyXML.Bytes())
					}
					break
				}

				if err := encoder.EncodeToken(token); err != nil {
					break
				}
			}
		}
	}

	return string(data)
}
