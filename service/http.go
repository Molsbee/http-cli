package service

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/yosssi/gohtml"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type RestClient struct {
	client      *http.Client
	include     bool
	prettyPrint bool
}

func NewRestClient(include, prettyPrint bool) RestClient {
	return RestClient{
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		include:     include,
		prettyPrint: prettyPrint,
	}
}

func (rc RestClient) Get(url string, headers map[string]string) (string, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	return rc.doRequest(request, headers)
}

func (rc RestClient) Put(url, body string, headers map[string]string) (string, error) {
	request, err := http.NewRequest("PUT", url, strings.NewReader(body))
	if err != nil {
		return "", err
	}

	return rc.doRequest(request, headers)
}

func (rc RestClient) Post(url, body string, headers map[string]string) (string, error) {
	request, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return "", err
	}

	return rc.doRequest(request, headers)
}

func (rc RestClient) Delete(url string, body string, headers map[string]string) (string, error) {
	request, err := http.NewRequest("DELETE", url, strings.NewReader(body))
	if err != nil {
		return "", err
	}

	return rc.doRequest(request, headers)
}

func (rc RestClient) Head(url string, headers map[string]string) (string, error) {
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return "", err
	}

	return rc.doRequest(request, headers)
}

func (rc RestClient) doRequest(request *http.Request, headers map[string]string) (string, error) {
	for k, v := range headers {
		request.Header.Set(k, v)
	}

	resp, err := rc.client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	dataBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if rc.include {
		fmt.Printf("%s %s\n", resp.Proto, resp.Status)
		for k, v := range resp.Header {
			fmt.Printf("%s: %s\n", k, strings.Join(v, ", "))
		}
		fmt.Println()
	}

	// do something with response
	if rc.prettyPrint {
		return prettyPrint(resp.Header.Get("Content-Type"), dataBytes)
	}

	return string(dataBytes), nil
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
