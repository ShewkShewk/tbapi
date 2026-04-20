package tbapi

import (
	"errors"
	url2 "net/url"
	"reflect"
	"testing"
)

func Test(t *testing.T) {
	tests := []struct {
		name          string
		scenario      string
		expectedToken string
		expectedErr   error
	}{
		{
			name:          "Token successfully retrieved",
			scenario:      "test_scenarios/complete.json",
			expectedToken: "tabroom_token_value",
			expectedErr:   nil,
		},
		{
			name:          "Missing salt",
			scenario:      "test_scenarios/missing_salt.json",
			expectedToken: "",
			expectedErr:   errors.New("unable to find salt in login parameters"),
		},
		{
			name:          "Missing sha",
			scenario:      "test_scenarios/missing_sha.json",
			expectedToken: "",
			expectedErr:   errors.New("unable to find sha in login parameters"),
		},
		{
			name:          "Wrong Login",
			scenario:      "test_scenarios/wrong_login.json",
			expectedToken: "",
			expectedErr:   errors.New("unable to find TabroomToken within cookies after login"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testServer, err := getTestServerConfiguration(test.scenario)
			if err != nil {
				t.Errorf("unable to start up test server %v", err)
				t.Fail()
			}
			defer testServer.Close()
			url, _ := url2.Parse(testServer.URL)
			api := TabroomApi{
				username: "",
				password: "",
				client:   newDefaultHttpRequester(*url),
			}
			token, err := api.retrieveCredentials()
			if !reflect.DeepEqual(token, test.expectedToken) {
				t.Errorf("retrieveCredentials token got = %v, want %v", token, test.expectedToken)
			}
			if !reflect.DeepEqual(err, test.expectedErr) {
				t.Errorf("retrieveCredentials error got = %v, want %v", err, test.expectedErr)
			}
		})
	}
}
