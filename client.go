package tbapi

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type httpClient interface {
	get(path string) (*http.Response, error)
	postForm(path string, values url.Values) (*http.Response, error)
	cookies() map[string]string
}

type defaultHttpClient struct {
	url    url.URL
	client http.Client
}

func newDefaultHttpClient(url url.URL) httpClient {
	jar, _ := cookiejar.New(nil)
	return &defaultHttpClient{
		url: url,
		client: http.Client{
			Timeout: 15 * time.Second,
			Jar:     jar,
		},
	}
}

func (h *defaultHttpClient) get(path string) (*http.Response, error) {
	return h.client.Get(h.url.JoinPath(path).String())
}

func (h *defaultHttpClient) postForm(path string, values url.Values) (*http.Response, error) {
	return h.client.PostForm(h.url.JoinPath(path).String(), values)
}

func (h *defaultHttpClient) cookies() map[string]string {
	toReturn := make(map[string]string)
	parsed, _ := url.Parse(fmt.Sprintf("%s://%s", h.url.Scheme, h.url.Hostname()))
	for _, cookie := range h.client.Jar.Cookies(parsed) {
		toReturn[cookie.Name] = cookie.Value
	}
	return toReturn
}
