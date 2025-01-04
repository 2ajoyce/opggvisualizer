package client

import (
	"encoding/json"
	"fmt"
	"log"
	"opggvisualizer/internal/db"
	"opggvisualizer/internal/models"
	"time"
)

func FetchAndStoreChampionData() error {
	// Check the last time the champion data was updated
	lastUpdated, err := db.GetDatabaseConnection().GetLastFetch("CHAMPIONS")
	if err != nil {
		log.Printf("error getting last fetch time: %v", err) // Log the error, but continue
	}

	log.Printf("Last champion data update: %v", lastUpdated)
	if time.Since(lastUpdated) < 24*time.Hour {
		log.Println("Champion data is up to date.")
		return nil
	}

	// Fetch the latest champion data version
	versionsBytes, err := FetchData(ChampionDataVersionURL)
	if err != nil {
		return fmt.Errorf("error fetching champion data versions: %w", err)
	}

	var versions []string
	if err := json.Unmarshal(versionsBytes, &versions); err != nil {
		return fmt.Errorf("error unmarshaling champion data versions: %w", err)
	}

	if len(versions) == 0 {
		return fmt.Errorf("no versions found in ChampionDataVersionURL")
	}

	latestVersion := versions[0]
	log.Printf("Latest Champion Data Version: %s", latestVersion)

	// Construct the ChampionDataURL with the latest version
	formattedChampionDataURL := fmt.Sprintf(ChampionDataURL, latestVersion)

	// Fetch champion data using the latest version
	championDataBytes, err := FetchData(formattedChampionDataURL)
	if err != nil {
		return fmt.Errorf("error fetching champion data: %w", err)
	}

	var championData models.ChampionData
	if err := json.Unmarshal(championDataBytes, &championData); err != nil {
		return fmt.Errorf("error unmarshaling champion data: %w", err)
	}

	log.Printf("Fetched %d champions.", len(championData.Data))

	database := db.GetDatabaseConnection()
	// Insert champions into the database
	for _, champ := range championData.Data {
		if err := database.InsertChampion(champ); err != nil {
			log.Printf("Error inserting champion %s: %v", champ.Name, err)
			continue
		}
	}

	// List and log all champion_ids after insertion
	championIDs, err := database.ListChampionIDs()
	if err != nil {
		log.Printf("Error listing champion IDs: %v", err)
	} else {
		log.Printf("Total champions in database: %d", len(championIDs))
		// Uncomment the next line if you want to log all champion IDs
		// log.Printf("Champion IDs: %v", championIDs)
	}

	// Update the last fetch time
	if err := database.SetLastFetch("CHAMPIONS", time.Now()); err != nil {
		return fmt.Errorf("error updating last fetch time for champions: %w", err)
	}

	newFetchTime, err := database.GetLastFetch("CHAMPIONS")
	if err != nil {
		log.Printf("error getting champions last fetch time: %v", err) // Log the error, but continue
	}
	log.Printf("Updated last fetch time: %v", newFetchTime)

	return nil
}
