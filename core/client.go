package core

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/charmbracelet/log"
)

type AHClient struct {
	Client  *http.Client
	Headers http.Header
}

func (ahc *AHClient) Do(req *http.Request) (*http.Response, error) {
	for k, vals := range ahc.Headers {
		for _, v := range vals {
			req.Header.Add(k, v)
		}
	}
	return ahc.Client.Do(req)
}

func executeReq(client AHClient, req *http.Request) ([]byte, error) {

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return []byte{}, fmt.Errorf("HTTP error occurred: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func getPage(client AHClient, url string, currentPage int) ([]byte, error) {

	req, err := initReq(url, currentPage)
	if err != nil {
		return []byte{}, err
	}

	body, err := executeReq(client, req)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func InitClient(proxyURL *url.URL, skipVerifySSL bool,
	username string, password string, token string) AHClient {
	// Returns a client configured for the gatherer.
	// For now, this only configures the Proxy, might be useful to handle SSL verification too.

	transport := &http.Transport{}

	if skipVerifySSL {
		TLSClientConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		transport.TLSClientConfig = TLSClientConfig
	}

	if proxyURL != nil {
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	headers := http.Header{}
	if username != "" && password != "" {
		raw := username + ":" + password
		encoded := base64.StdEncoding.EncodeToString([]byte(raw))
		headers.Add("Authorization", "Basic "+encoded)
		if token != "" {
			log.Warn("Prioritizing username/password authentication material, token was ignored.")
		}
	} else if token != "" {
		if password != "" || username != "" {
			log.Warn("Prioritizing token authentication material, username/password was ignored.")
		}
		headers.Add("Authorization", "Bearer "+token)
	} else {
		log.Fatal("No authentication material provided, exiting.")
	}

	client := AHClient{
		Client:  httpClient,
		Headers: headers,
	}

	return client
}

func initReq(url string, currentPage int) (*http.Request, error) {

	url = url + PAGE_SIZE_ARG + fmt.Sprintf(CURRENT_PAGE_ARG, currentPage)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return req, nil

}
