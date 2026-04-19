package tbapi

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type httpRequester interface {
	get(path string) (*http.Response, error)
	postForm(path string, values url.Values) (*http.Response, error)
	cookies() map[string]string
}

type defaultHttpRequester struct {
	url    url.URL
	client http.Client
}

func newDefaultHttpRequester(url url.URL) httpRequester {
	jar, _ := cookiejar.New(nil)
	return &defaultHttpRequester{
		url: url,
		client: http.Client{
			Timeout: 15 * time.Second,
			Jar:     jar,
		},
	}
}

func (h *defaultHttpRequester) get(path string) (*http.Response, error) {
	return h.client.Get(h.url.JoinPath(path).String())
}

func (h *defaultHttpRequester) postForm(path string, values url.Values) (*http.Response, error) {
	return h.client.PostForm(h.url.JoinPath(path).String(), values)
}

func (h *defaultHttpRequester) cookies() map[string]string {
	toReturn := make(map[string]string)
	parsed, _ := url.Parse(fmt.Sprintf("%s://%s", h.url.Scheme, h.url.Hostname()))
	for _, cookie := range h.client.Jar.Cookies(parsed) {
		toReturn[cookie.Name] = cookie.Value
	}
	return toReturn
}
