package tools

import (
	"io/ioutil"
	"net/http"
)

func GetHttpContent(url string) (string, error) {
	resp, err := http.Get(url)

	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
