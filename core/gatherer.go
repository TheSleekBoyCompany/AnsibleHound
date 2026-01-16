package core

import (
	"ansible-hound/core/ansible"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func Gather[T ansible.AnsibleType](client AHClient, target url.URL,
	endpoint string) ([]T, error) {

	var objectList []T
	count := 0
	current := 0
	page := 1

	url := target.String() + endpoint

	body, err := getPage(client, url, page)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		r := ansible.Response[T]{}
		err = json.Unmarshal(body, &r)
		if err != nil {
			return nil, err
		}
		count = r.Count
		objectList = append(objectList, r.Results...)
	}
	current += PAGE_SIZE

	if count >= PAGE_SIZE {
		for {
			page += 1
			body, err := getPage(client, url, page)
			if err != nil {
				return nil, err
			}
			r := ansible.Response[T]{}
			err = json.Unmarshal(body, &r)
			if err != nil {
				return nil, err
			}
			objectList = append(objectList, r.Results...)
			current += PAGE_SIZE
			if current >= count {
				break
			}
		}
	}

	return objectList, nil

}

func GatherObject[T ansible.AnsibleType](installUUID string, client AHClient,
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

func GatherAnsibleInstance(client AHClient, target url.URL) (instance ansible.AnsibleInstance, err error) {

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

func HasAccessTo[T ansible.AnsibleType](objectMap map[int]T, ID int) (result bool) {
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
