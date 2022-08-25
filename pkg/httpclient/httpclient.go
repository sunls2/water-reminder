package httpclient

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"net"
	"net/http"
	"time"
)

var client *http.Client

func init() {
	client = &http.Client{Timeout: 15 * time.Second}
	client.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   15 * time.Second,
			KeepAlive: time.Minute,
		}).DialContext,
		MaxIdleConns:    50,
		MaxConnsPerHost: 10,
		IdleConnTimeout: 3 * time.Hour,
	}
}

func jsonDo(method, url string, param any) (Response, error) {
	var reader io.Reader
	if param != nil {
		body, err := json.Marshal(param)
		if err != nil {
			return nil, errors.Wrap(err, "json.Marshal param")
		}
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, errors.Wrap(err, "new request")
	}
	req.Header.Set("Content-Type", "application/json")
	return handleResponse(client.Do(req))
}

func handleResponse(resp *http.Response, err error) (Response, error) {
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return nil, errors.Wrap(err, "read response.body")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("code: %d, body: %s", resp.StatusCode, body)
	}
	return newResponse(body), nil
}
