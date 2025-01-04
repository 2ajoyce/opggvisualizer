// internal/db/champions.go
package db

import (
	"encoding/json"
	"fmt"
	"opggvisualizer/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

// InsertChampion inserts a champion into the champions table
func (db *Database) InsertChampion(champion models.Champion) error {
	tagsJSON, err := json.Marshal(champion.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	statsJSON, err := json.Marshal(champion.Stats)
	if err != nil {
		return fmt.Errorf("failed to marshal stats: %w", err)
	}

	insertChampionSQL := `INSERT INTO champions(
		champion_id, name, title, tags, type, format, blurb, partype,
		attack, defense, magic, difficulty, stats, image_url
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(champion_id) DO UPDATE SET
		name=excluded.name,
		title=excluded.title,
		tags=excluded.tags,
		type=excluded.type,
		format=excluded.format,
		blurb=excluded.blurb,
		partype=excluded.partype,
		attack=excluded.attack,
		defense=excluded.defense,
		magic=excluded.magic,
		difficulty=excluded.difficulty,
		stats=excluded.stats,
		image_url=excluded.image_url;`

	// Use champion.Key as champion_id to match participant's champion_id
	_, err = db.Conn.Exec(insertChampionSQL,
		champion.Key, // Changed from champion.ID to champion.Key
		champion.Name,
		champion.Title,
		string(tagsJSON),
		champion.Type,
		champion.Format,
		champion.Blurb,
		champion.Partype,
		champion.Info.Attack,
		champion.Info.Defense,
		champion.Info.Magic,
		champion.Info.Difficulty,
		string(statsJSON),
		champion.Image.FullURL(),
	)
	if err != nil {
		return fmt.Errorf("failed to insert champion: %w", err)
	}
	return nil
}

// Function to list all champion_ids
func (db *Database) ListChampionIDs() ([]string, error) {
	rows, err := db.Conn.Query("SELECT champion_id FROM champions;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var championIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		championIDs = append(championIDs, id)
	}
	return championIDs, nil
}

// ClearChampionData clears all data from the champions table
func (db *Database) ClearChampionData() error {
	_, err := db.Conn.Exec(`DELETE FROM champions;`)
	if err != nil {
		return fmt.Errorf("failed to clear champion data: %w", err)
	}

	// Wipe last fetch time for champions
	_, err = db.Conn.Exec(`UPDATE fetch SET last_fetch = NULL WHERE fetch_type = "CHAMPIONS";`)
	if err != nil {
		return fmt.Errorf("failed to clear last fetch time for champions: %w", err)
	}
	return nil
}
