package service

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/gookit/color"
	"github.com/yosssi/gohtml"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	lightBlue = color.FgLightBlue.Render
	teal      = color.FgCyan.Render
	purple    = color.Magenta.Render
)

type RestClient struct {
	client      *http.Client
	include     bool
	prettyPrint bool
	verbose     bool
}

func NewRestClient(include, verbose, prettyPrint bool) RestClient {
	return RestClient{
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		include:     include,
		verbose:     verbose,
		prettyPrint: prettyPrint,
	}
}

func (rc RestClient) Get(url string, headers map[string]string) (string, error) {
	return rc.doReq("GET", url, headers, "")
}

func (rc RestClient) Put(url, body string, headers map[string]string) (string, error) {
	return rc.doReq("PUT", url, headers, body)
}

func (rc RestClient) Post(url, body string, headers map[string]string) (string, error) {
	return rc.doReq("POST", url, headers, body)
}

func (rc RestClient) Delete(url string, body string, headers map[string]string) (string, error) {
	return rc.doReq("DELETE", url, headers, body)
}

func (rc RestClient) Head(url string, headers map[string]string) (string, error) {
	return rc.doReq("HEAD", url, headers, "")
}

func (rc RestClient) doReq(method, url string, headers map[string]string, body string) (string, error) {
	request, err := rc.createRequest(method, url, headers, body)
	if err != nil {
		return "", err
	}

	if rc.verbose {
		fmt.Printf("%s %s %s\n", teal(method), url, rc.sPrintProto(request.ProtoMajor, request.ProtoMinor))
		for k, v := range request.Header {
			fmt.Printf("%s: %s\n", lightBlue(k), strings.Join(v, ", "))
		}
		fmt.Printf("\n%s\n", body)
	}

	resp, err := rc.client.Do(request)
	if err != nil {
		return "", errors.New("error occurred performing request - " + err.Error())
	}
	defer resp.Body.Close()

	// Print Response Headers
	if rc.include {
		fmt.Printf("%s %s\n", rc.sPrintProto(resp.ProtoMajor, resp.ProtoMinor), purple(resp.Status))
		for k, v := range resp.Header {
			fmt.Printf("%s: %s\n", lightBlue(k), strings.Join(v, ", "))
		}
		fmt.Println()
	}

	dataBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("error occurred reading response body - " + err.Error())
	}

	if rc.prettyPrint {
		return prettyPrint(resp.Header.Get("Content-Type"), dataBytes)
	}

	return string(dataBytes), nil
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

func prettyPrint(contentType string, data []byte) (string, error) {
	switch {
	case strings.Contains(contentType, "application/json"):
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, data, "", "\t"); err == nil {
			return string(prettyJSON.Bytes()), nil
		}
	case strings.Contains(contentType, "text/html"):
		return string(gohtml.FormatBytes(data)), nil
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
					return string(prettyXML.Bytes()), nil
				}
				break
			}

			if err := encoder.EncodeToken(token); err != nil {
				break
			}
		}
	}
	return string(data), nil
}
