package xhttp

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const timeout = time.Second * 5

func Request(method string, reqUrl string, postData string, header http.Header) (string, error) {
	var res string
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
	var ioreader io.Reader
	if postData != "" {
		ioreader = strings.NewReader(postData)
	}
	req, err := http.NewRequest(method, reqUrl, ioreader)
	if err != nil {
		return res, err
	}
	if header != nil {
		req.Header = header
	}

	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		byt, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return res, err
		}
		return fmt.Sprintf("%s", byt), nil
	} else {
		return res, fmt.Errorf("code:%d message:%s", resp.StatusCode, resp.Status)
	}
}
