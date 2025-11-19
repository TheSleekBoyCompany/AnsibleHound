package core

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
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

func Gather[T AnsibleType](client AHClient, target url.URL,
	endpoint string) ([]T, error) {

	var objects []T
	count := 0
	current := 0
	page := 1

	url := target.String() + endpoint

	body, err := getPage(client, url, page)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		r := Response[T]{}
		err = json.Unmarshal(body, &r)
		if err != nil {
			return nil, err
		}
		count = r.Count
		objects = append(objects, r.Results...)
	}
	current += PAGE_SIZE

	if count >= PAGE_SIZE {
		for {
			page += 1
			body, err := getPage(client, url, page)
			if err != nil {
				return nil, err
			}
			r := Response[T]{}
			err = json.Unmarshal(body, &r)
			if err != nil {
				return nil, err
			}
			objects = append(objects, r.Results...)
			current += PAGE_SIZE
			if current >= count {
				break
			}
		}
	}

	return objects, nil

}

func GatherObject[T AnsibleType](installUUID string, client AHClient,
	target url.URL, endpoint string) (
	objectMap map[int]T, err error) {

	objectMap = make(map[int]T)

	objects, err := Gather[T](client, target, endpoint)
	if err != nil {
		return nil, err
	}

	for _, object := range objects {
		object.InitOID(installUUID)
		objectMap[object.GetID()] = object
	}

	return objectMap, nil
}

func GatherAnsibleInstance(client AHClient, target url.URL) (instance AnsibleInstance, err error) {

	url := target.String() + PING_ENDPOINT

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return instance, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return instance, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return instance, err
	}

	err = json.Unmarshal(body, &instance)
	if err != nil {
		return instance, err
	}

	instance.InitOID(instance.InstallUUID)

	return instance, nil
}

func HasAccessTo[T AnsibleType](objectMap map[int]T, ID int) (result bool) {
	// NOTE: If ID = 0, then the resource is not bound to a resource of this type.
	// EX: A Credential can exist without being bound to an Organization.
	// This might also mean a resource is used, but your user cannot read it.
	if ID != 0 {
		if _, ok := objectMap[ID]; ok {
			result = true
		}
	}
	return result
}

func AuthenticateOnAnsibleInstance(client AHClient,
	target url.URL, endpoint string) ([]byte, error) {

	url := target.String() + endpoint

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return []byte{}, fmt.Errorf("HTTP error occurred: %s", resp.Status)
	}

	return []byte{}, nil
}
