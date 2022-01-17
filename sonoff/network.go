package sonoff

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func httpRequest(method string, url string, req []byte) ([]byte, error) {
	client := http.Client{}

	r, err := http.NewRequest(method, url, bytes.NewBufferString(string(req)))
	if err != nil {
		return nil, err
	}

	if req != nil {
		r.Header.Add("Content-type", "application/json")
	}

	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return nil, fmt.Errorf("http response error %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
