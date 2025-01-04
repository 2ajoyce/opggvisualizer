package client

import (
	"encoding/json"
	"fmt"
	"log"
	"opggvisualizer/internal/config"
	"opggvisualizer/internal/db"
	"opggvisualizer/internal/models"
	"time"
)

func FetchAndStoreGameData() error {
	// Check the last time the game data was updated
	lastUpdated, err := db.GetDatabaseConnection().GetLastFetch("GAMES")
	if err != nil {
		log.Printf("error getting last fetch time: %v", err) // Log the error, but continue
	}

	log.Printf("Last game data update: %v", lastUpdated)
	if time.Since(lastUpdated) < 24*time.Hour {
		log.Println("Game data is up to date.")
		return nil
	}

	// Fetch the configuration
	cfg := config.GetConfig()

	// Fetch game data
	gameDataURL := fmt.Sprintf(GameDataURL, cfg.SummonerID)
	gameDataBytes, err := FetchData(gameDataURL)
	if err != nil {
		return fmt.Errorf("error fetching game data: %w", err)
	}

	var gameData models.GameData
	if err := json.Unmarshal(gameDataBytes, &gameData); err != nil {
		return fmt.Errorf("error unmarshalling game data: %w", err)
	}

	log.Printf("Fetched %d games.", len(gameData.Data))
	database := db.GetDatabaseConnection()
	// Insert games, teams, and participants into the database
	for _, gameEntry := range gameData.Data {
		// Parse time fields
		createdAt, err := time.Parse(time.RFC3339, gameEntry.CreatedAt)
		if err != nil {
			log.Printf("Error parsing created_at for game %s: %v", gameEntry.ID, err)
			continue
		}
		createdAt = createdAt.UTC()

		firstGameCreatedAt, err := time.Parse(time.RFC3339, gameData.Meta.FirstGameCreatedAt)
		if err != nil {
			log.Printf("Error parsing first_game_created_at for game %s: %v", gameEntry.ID, err)
			continue
		}
		firstGameCreatedAt = firstGameCreatedAt.UTC()

		lastGameCreatedAt, err := time.Parse(time.RFC3339, gameData.Meta.LastGameCreatedAt)
		if err != nil {
			log.Printf("Error parsing last_game_created_at for game %s: %v", gameEntry.ID, err)
			continue
		}
		lastGameCreatedAt = lastGameCreatedAt.UTC()

		// Insert game
		game := models.Game{
			ID:               gameEntry.ID,
			CreatedAt:        createdAt,
			GameLengthSecond: int(gameEntry.GameLengthSecond),
			AverageTierInfo:  gameEntry.AverageTierInfo,
			IsRemake:         gameEntry.IsRemake,
			MetaVersion:      gameEntry.MetaVersion,
			GameType:         gameEntry.GameType,
			IsOpscoreActive:  gameEntry.IsOpscoreActive,
			IsRecorded:       gameEntry.IsRecorded.Valid && gameEntry.IsRecorded.Float64 != 0,
			Version:          gameEntry.Version,
			Meta: models.GameMeta{
				FirstGameCreatedAt: firstGameCreatedAt,
				LastGameCreatedAt:  lastGameCreatedAt,
			},
		}

		if err := database.InsertGame(game); err != nil {
			log.Printf("Error inserting game %s: %v", game.ID, err)
			continue
		}

		// Insert teams
		for _, team := range gameEntry.Teams {
			if err := database.InsertTeam(game.ID, team); err != nil {
				log.Printf("Error inserting team for game %s: %v", game.ID, err)
				continue
			}
		}

		// Log a participant to verify the data types
		fmt.Printf("Participant: %+v\n", gameEntry.Participants[0])

		// Insert participants
		for _, participant := range gameEntry.Participants {
			if err := database.InsertParticipant(game.ID, participant); err != nil {
				log.Printf("Error inserting participant for game %s: %v", game.ID, err)
				continue
			}
		}
	}

	// Update the last fetch time
	if err := database.SetLastFetch("GAMES", time.Now()); err != nil {
		return fmt.Errorf("error updating last fetch time for games: %w", err)
	}

	newFetchTime, err := database.GetLastFetch("GAMES")
	if err != nil {
		log.Printf("error getting games last fetch time: %v", err) // Log the error, but continue
	}
	log.Printf("Updated last fetch time: %v", newFetchTime)

	return nil
}
