package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type RestClient struct {
	client  *http.Client
	include bool
}

func NewRestClient(include bool) RestClient {
	return RestClient{
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		include: include,
	}
}

func (rc RestClient) Get(url string, headers map[string]string) (string, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	for k, v := range headers {
		request.Header.Set(k, v)
	}

	return rc.doRequest(request)
}

func (rc RestClient) Put(url, body string) {

}

func (rc RestClient) Post(url, body string) {

}

func (rc RestClient) Delete(url string, body string) {

}

func (rc RestClient) Head(url string) {

}

func (rc RestClient) doRequest(request *http.Request) (string, error) {
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
	return string(dataBytes), nil
}
