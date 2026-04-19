package tbapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
)

type TestServerConfiguration struct {
	Configurations []struct {
		Method    string `json:"method"`
		Uri       string `json:"uri"`
		Responses []struct {
			Code    int               `json:"code"`
			Body    string            `json:"body"`
			Headers map[string]string `json:"headers"`
			Cookies []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
				Path  string `json:"path"`
			} `json:"cookies"`
		} `json:"responses"`
	} `json:"configurations"`
}

type TestCookie struct {
	name  string
	value string
	path  string
}
type TestHttpResponse struct {
	code    int
	body    string
	headers map[string]string
	cookies []TestCookie
}
type TestHttpServer struct {
	Responses map[string][]TestHttpResponse
}

func getTestServerConfiguration(scenarioDirectory string) (*httptest.Server, error) {
	content, err := os.ReadFile(scenarioDirectory)
	if err != nil {
		return nil, err
	}
	var testServerConfiguration TestServerConfiguration

	if err := json.Unmarshal(content, &testServerConfiguration); err != nil {
		return nil, err
	}
	httpServer, err := testServerConfiguration.toHttpServer()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to create new test http server from configuration: %s", scenarioDirectory))
	}
	handlerFunc := http.HandlerFunc(httpServer.handle)
	return httptest.NewServer(&handlerFunc), nil
}

func (t *TestServerConfiguration) toHttpServer() (*TestHttpServer, error) {
	httpServer := TestHttpServer{
		Responses: make(map[string][]TestHttpResponse),
	}

	for _, configuration := range t.Configurations {
		method := configuration.Method
		uri := configuration.Uri
		key := fmt.Sprintf("%s %s", method, uri)
		for _, response := range configuration.Responses {
			cookies := make([]TestCookie, len(response.Cookies))
			for i, cookie := range response.Cookies {
				cookies[i] = TestCookie{
					name:  cookie.Name,
					value: cookie.Value,
					path:  cookie.Path,
				}
			}
			converted := TestHttpResponse{
				code:    response.Code,
				body:    response.Body,
				headers: response.Headers,
				cookies: cookies,
			}
			httpServer.Responses[key] = append(httpServer.Responses[key], converted)
		}
	}
	return &httpServer, nil
}

func (t *TestHttpServer) handle(writer http.ResponseWriter, requests *http.Request) {
	responses, ok := t.Responses[fmt.Sprintf("%v %v", requests.Method, requests.URL.Path)]
	if !ok || len(responses) == 0 {
		writer.WriteHeader(404)
		return
	}

	response := responses[0]
	responses = responses[1:]

	for key, value := range response.headers {
		writer.Header().Add(key, value)
	}
	for _, cookie := range response.cookies {
		http.SetCookie(writer, &http.Cookie{
			Name:  cookie.name,
			Value: cookie.value,
			Path:  cookie.path,
		})
	}
	writer.WriteHeader(response.code)
	writer.Write([]byte(response.body))
}
