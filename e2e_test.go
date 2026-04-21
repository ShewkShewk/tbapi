package tbapi

import (
	"reflect"
	"testing"
)

func TestEndToEnd(t *testing.T) {
	tests := []struct {
		name          string
		scenario      string
		username      string
		password      string
		toCall        func(api TabroomApi) (interface{}, error)
		expected      interface{}
		expectedError error
	}{
		{
			name:     "End to End GetTournaments",
			scenario: "test_scenarios/end_to_end/get_tournaments.json",
			username: "right_username",
			password: "right_password",
			toCall:   getTournamentsFunc(),
			expected: []Tournament{
				{
					Id:   27632,
					Date: getTimeForDate("2026-12-31"),
					Name: "dElo Test Tournament",
				},
				{
					Id:   39832,
					Date: getTimeForDate("2026-04-25"),
					Name: "dElo Test Tournament 2",
				},
				{
					Id:   39722,
					Date: getTimeForDate("2026-04-18"),
					Name: "dElo Test Tournament 3",
				},
			},
			expectedError: nil,
		},
		{
			name:     "End to End GetTournamentData",
			scenario: "test_scenarios/end_to_end/get_tournament_data.json",
			username: "right_username",
			password: "right_password",
			toCall:   getTournamentData(27632),
			expected: &TournamentData{
				Timezone: "America/Chicago",
				Start:    "2026-12-31 09:00:00",
				Name:     "dElo Test Tournament",
				Webname:  "delo",
				City:     "Dallas",
			},
			expectedError: nil,
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
			tabroomApi, err := NewBuilder().WithUsername(test.username).WithPassword(test.password).WithHostname(testServer.URL).Build()
			if err != nil {
				t.Errorf("unable to create TabroomApi for scenario %s", test.name)
			}
			got, err := test.toCall(*tabroomApi)
			if !reflect.DeepEqual(got, test.expected) {
				t.Errorf("outcome: got: %v, expected: %v", got, test.expected)
			}
			if !reflect.DeepEqual(err, test.expectedError) {
				t.Errorf("error: got: %v, expected: %v", err, test.expectedError)
			}
		})
	}
}

func getTournamentsFunc() func(api TabroomApi) (interface{}, error) {
	return func(api TabroomApi) (interface{}, error) {
		return api.GetTournaments()
	}
}

func getTournamentData(tournamentId int) func(api TabroomApi) (interface{}, error) {
	return func(api TabroomApi) (interface{}, error) {
		return api.GetTournamentData(tournamentId)
	}
}
