package tbapi

import (
	"encoding/json"
	"io"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

type TabroomApi struct {
	username string
	password string
	client   httpClient
}

func (t *TabroomApi) GetTournaments() ([]Tournament, error) {
	if err := t.ensureAuthenticated(); err != nil {
		return nil, err
	}
	resp, err := t.client.get("/user/tourn/all.mhtml")
	if err != nil {
		return nil, err
	}
	content, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	return parseTournaments(string(content))
}

func (t *TabroomApi) GetTournamentData(tournamentId int) (*TournamentData, error) {
	if err := t.ensureAuthenticated(); err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Add("tourn_id", strconv.Itoa(tournamentId))

	resp, err := t.client.post("/api/download_data.mhtml", "application/x-www-form-urlencoded", params.Encode())
	if err != nil {
		return nil, err
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	tabroomTournamentData := TournamentData{}
	err = json.Unmarshal(content, &tabroomTournamentData)
	if err != nil {
		return nil, err
	}
	return &tabroomTournamentData, nil
}

func parseTournaments(html string) ([]Tournament, error) {
	otherRe := regexp.MustCompile(`(?is)<td[^>]*>\s*([0-9]{4}-[0-9]{2}-[0-9]{2})\s*</td>.*?href="select\.mhtml\?tourn_id=([0-9]+)".*?>\s*([^<\s](?:[^<]*?[^<\s])?)\s*</a>`)

	matches := otherRe.FindAllStringSubmatch(html, -1)
	tournaments := make([]Tournament, 0, len(matches))
	for _, match := range matches {
		tournamentDate, err := time.Parse(time.DateOnly, match[1])
		if err != nil {
			log.Printf("Unable to convert %s to time.Time", matches[1])
			continue
		}
		tournamentId, err := strconv.Atoi(match[2])
		if err != nil {
			log.Printf("Unable to convert %s to int", match[2])
			continue
		}
		tournamentName := match[3]
		tournaments = append(tournaments, Tournament{
			Id:   tournamentId,
			Date: tournamentDate,
			Name: tournamentName,
		})
	}

	return tournaments, nil
}
