package core

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func InitClient(proxyURL *url.URL, skipVerifySSL bool) http.Client {
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

	client := &http.Client{
		Transport: transport,
	}

	return *client
}

func initReq(url string, username string,
	password string, currentPage int) (*http.Request, error) {

	url = url + PAGE_SIZE_ARG + fmt.Sprintf(CURRENT_PAGE_ARG, currentPage)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(username, password)

	return req, nil

}

func executeReq(client http.Client, req *http.Request) ([]byte, error) {

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

func getPage(client http.Client, url string,
	username string, password string, currentPage int) ([]byte, error) {

	req, err := initReq(url, username, password, currentPage)
	if err != nil {
		return []byte{}, err
	}

	body, err := executeReq(client, req)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func Gather[T AnsibleType](client http.Client, target url.URL,
	username string, password string, endpoint string) ([]T, error) {

	var objects []T
	count := 0
	current := 0
	page := 1

	url := target.String() + endpoint

	body, err := getPage(client, url, username, password, page)
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
			body, err := getPage(client, url, username, password, page)
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

func GatherObject[T AnsibleType](installUUID string, client http.Client,
	target url.URL, username string, password string, endpoint string) (
	objectMap map[int]T, err error) {

	objectMap = make(map[int]T)

	objects, err := Gather[T](client, target, username, password, endpoint)
	if err != nil {
		return nil, err
	}

	for _, object := range objects {
		object.InitOID(installUUID)
		objectMap[object.GetID()] = object
	}

	return objectMap, nil
}

func GatherAnsibleInstance(client http.Client, target url.URL) (instance AnsibleInstance, err error) {

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

func AuthenticateOnAnsibleInstance(client http.Client,
	target url.URL, username string, password string, endpoint string) ([]byte, error) {

	url := target.String() + endpoint

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(username, password)

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
