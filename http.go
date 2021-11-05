package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"git.sr.ht/~spc/go-log"
)

var (
	client    *http.Client
	userAgent string
	username  string
	password  string
)

func initHTTPClient(config *tls.Config, user, pass, ua string) {
	client = &http.Client{
		Transport: http.DefaultTransport.(*http.Transport).Clone(),
	}
	client.Transport.(*http.Transport).TLSClientConfig = config

	userAgent = ua
	username = user
	password = pass
}

func get(url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create HTTP request: %w", err)
	}

	for k, v := range headers {
		req.Header.Add(k, strings.TrimSpace(v))
	}
	req.Header.Add("User-Agent", userAgent)

	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	log.Debugf("sending HTTP request: %v %v", req.Method, req.URL)
	log.Tracef("request: %v", req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot get HTTP request: %w", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response body: %w", err)
	}
	log.Debugf("received HTTP response: %v %v", resp.Status, strings.TrimSpace(string(data)))
	log.Tracef("response: %v", resp)

	switch resp.StatusCode {
	case http.StatusOK:
		return data, nil
	default:
		return nil, fmt.Errorf("unexpected HTTP response: %v", resp.Status)
	}
}

func post(url string, headers map[string]string, body []byte) error {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("cannot create HTTP request: %w", err)
	}

	for k, v := range headers {
		req.Header.Add(k, strings.TrimSpace(v))
	}
	req.Header.Add("User-Agent", userAgent)

	log.Debugf("sending HTTP request: %v %v", req.Method, req.URL)
	log.Tracef("request: %v", req)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("cannot post HTTP requert: %w", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("cannot read response body: %w", err)
	}
	log.Debugf("received HTTP response: %v %v", resp.Status, strings.TrimSpace(string(data)))
	log.Tracef("response: %v", resp)

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		return nil
	default:
		return fmt.Errorf("unexpected HTTP response: %v", resp.Status)
	}
}
