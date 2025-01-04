// internal/client/client.go
package client

import (
	"fmt"
	"io"
	"net/http"
)

const (
	GameDataURL            = "https://lol-web-api.op.gg/api/v1.0/internal/bypass/games/na/summoners/%s?=&limit=20&hl=en_US&game_type=soloranked"
	ChampionDataURL        = "http://ddragon.leagueoflegends.com/cdn/%s/data/en_US/champion.json"
	ChampionDataVersionURL = "https://ddragon.leagueoflegends.com/api/versions.json"
)

func FetchData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}
