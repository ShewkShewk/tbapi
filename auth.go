package tbapi

import (
	"errors"
	"io"
	"net/url"
	"regexp"
)

type CredentialsRetriever interface {
	retrieveCredentials() (error, string)
}

const (
	SALT = "salt"
	SHA  = "sha"
)

type TabroomCredentialsRetriever struct {
	url       url.URL
	username  string
	password  string
	requester httpRequester
}

func newTabroomCredentialsRetriever(
	url url.URL,
	username string,
	password string,
	requester httpRequester,
) *TabroomCredentialsRetriever {
	if requester == nil {
		requester = newDefaultHttpRequester(url)
	}
	return &TabroomCredentialsRetriever{
		url:       url,
		username:  username,
		password:  password,
		requester: requester,
	}
}

func (t *TabroomCredentialsRetriever) retrieveCredentials() (string, error) {
	loginParameters, err := t.getLoginParameters()
	if err != nil {
		return "", err
	}
	if loginParameters[SALT] == "" {
		return "", errors.New("unable to find salt in login parameters")
	}
	if loginParameters[SHA] == "" {
		return "", errors.New("unable to find sha in login parameters")
	}
	loginForm, err := t.getLoginForm(loginParameters[SALT], loginParameters[SHA])
	if err != nil {
		return "", err
	}
	tabroomToken, err := t.getTabroomToken(loginForm)
	if err != nil {
		return "", err
	}
	return tabroomToken, nil
}

func (t *TabroomCredentialsRetriever) getLoginParameters() (map[string]string, error) {
	result, err := t.requester.get("/index/index.mhtml")
	if err != nil {
		return nil, err
	}
	toReturn := make(map[string]string)
	body, err := io.ReadAll(result.Body)
	re := regexp.MustCompile(`<input[^>]*name\s*=\s*"([^"]+)"[^>]*value\s*=\s*"([^"]+)"`)
	matches := re.FindAllStringSubmatch(string(body), -1)
	for _, match := range matches {
		toReturn[match[1]] = match[2]
	}
	return toReturn, nil
}

func (t *TabroomCredentialsRetriever) getLoginForm(salt string, sha string) (url.Values, error) {
	loginForm := url.Values{}
	loginForm.Add("salt", salt)
	loginForm.Add("sha", sha)
	loginForm.Add("username", t.username)
	loginForm.Add("password", t.password)
	return loginForm, nil
}

func (t *TabroomCredentialsRetriever) getTabroomToken(loginForm url.Values) (string, error) {
	client := t.requester
	_, err := client.postForm("/user/login/login_save.mhtml", loginForm)
	if err != nil {
		return "", err
	}
	cookies := client.cookies()
	val, ok := cookies["TabroomToken"]
	if !ok {
		return "", errors.New("unable to find TabroomToken within cookies after login")
	}
	return val, nil
}
