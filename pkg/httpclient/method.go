package httpclient

import "net/http"

func Get(url string) (Response, error) {
	return handleResponse(client.Get(url))
}

func Post(url string, param any) (Response, error) {
	return jsonDo(http.MethodPost, url, param)
}
