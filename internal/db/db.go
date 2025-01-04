// internal/db/db.go
package db

import (
	"database/sql"
	"fmt"
	"log"
	"opggvisualizer/internal/config"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *Database // Constant to hold the DB connection

type Database struct {
	Conn *sql.DB // At some point in the future we might want to make this an array
}

func newDatabase(dbPath string) (*Database, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db = &Database{Conn: conn}

	// Enable foreign key constraints
	if _, err := db.Conn.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	if err := db.initializeSchema(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return db, nil
}

func GetDatabaseConnection() *Database {
	if db == nil {
		cfg := config.GetConfig()
		newDb, err := newDatabase(cfg.DatabasePath)
		if err != nil {
			log.Fatalf("Error initializing database: %v", err)
		}
		db = newDb
	}
	return db
}

func (db *Database) Close() error {
	return db.Conn.Close()
}

// Initialize the database schema with enhanced structure
func (db *Database) initializeSchema() error {
	schemaStatements := []string{
		// Games Table
		`CREATE TABLE IF NOT EXISTS games (
			game_id TEXT PRIMARY KEY,
			created_at TEXT,
			game_length INTEGER,
			tier TEXT,
			division INTEGER,
			tier_image_url TEXT,
			border_image_url TEXT,
			is_remake BOOLEAN,
			meta_version TEXT,
			game_type TEXT,
			is_opscore_active BOOLEAN,
			is_recorded BOOLEAN,
			version TEXT,
			first_game_created_at TEXT,
			last_game_created_at TEXT
		);`,

		// Teams Table
		`CREATE TABLE IF NOT EXISTS teams (
			team_id INTEGER PRIMARY KEY AUTOINCREMENT,
			game_id TEXT,
			key TEXT,
			is_win BOOLEAN,
			champion_first BOOLEAN,
			inhibitor_first BOOLEAN,
			rift_herald_first BOOLEAN,
			death INTEGER,
			champion_kill INTEGER,
			inhibitor_kill INTEGER,
			dragon_first BOOLEAN,
			horde_first BOOLEAN,
			rift_herald_kill INTEGER,
			is_remake BOOLEAN,
			gold_earned INTEGER,
			kill INTEGER,
			tower_first BOOLEAN,
			horde_kill INTEGER,
			assist INTEGER,
			dragon_kill INTEGER,
			baron_kill INTEGER,
			baron_first BOOLEAN,
			tower_kill INTEGER,
			FOREIGN KEY(game_id) REFERENCES games(game_id)
		);`,

		// Participants Table
		`CREATE TABLE IF NOT EXISTS participants (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			game_id TEXT,
			participant_id REAL,
			summoner_name TEXT,
			champion_id TEXT,
			position TEXT,
			role TEXT,
			kills INTEGER,
			deaths INTEGER,
			assists INTEGER,
			gold_earned INTEGER,
			damage_dealt INTEGER,
			damage_taken INTEGER,
			vision_score INTEGER,
			primary_rune_id INTEGER,
			secondary_rune_page_id INTEGER,
			lane_score INTEGER,
			team_key TEXT,
			result TEXT,
			ward_place INTEGER,
			op_score_rank INTEGER,
			barrack_kill INTEGER,
			total_heal INTEGER,
			game_type TEXT,
			is_remake BOOLEAN,
			FOREIGN KEY(game_id) REFERENCES games(game_id),
			FOREIGN KEY(champion_id) REFERENCES champions(champion_id)
		);`,

		// Participant Items Table
		`CREATE TABLE IF NOT EXISTS participant_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			participant_id INTEGER,
			item_id INTEGER,
			FOREIGN KEY(participant_id) REFERENCES participants(id)
		);`,

		// Participant Spells Table
		`CREATE TABLE IF NOT EXISTS participant_spells (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			participant_id INTEGER,
			spell_id INTEGER,
			FOREIGN KEY(participant_id) REFERENCES participants(id)
		);`,

		// Team Banned Champions Table
		`CREATE TABLE IF NOT EXISTS team_banned_champions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			team_id INTEGER,
			banned_champion_id REAL,
			FOREIGN KEY(team_id) REFERENCES teams(team_id)
		);`,

		// Champions Table
		`CREATE TABLE IF NOT EXISTS champions (
			champion_id TEXT PRIMARY KEY,
			name TEXT,
			title TEXT,
			tags TEXT, -- Serialized JSON array
			type TEXT,
			format TEXT,
			blurb TEXT,
			partype TEXT,
			attack REAL,
			defense REAL,
			magic REAL,
			difficulty REAL,
			stats TEXT, -- Serialized JSON object
			image_url TEXT
		);`,

		// Fetch Table
		`CREATE TABLE IF NOT EXISTS fetch (
			fetch_type TEXT PRIMARY KEY, -- Type of fetch (CHAMPIONS, GAMES)
			last_fetch TEXT -- Last fetch timestamp
		);`,
	}

	for _, stmt := range schemaStatements {
		if _, err := db.Conn.Exec(stmt); err != nil {
			return fmt.Errorf("failed to execute schema statement: %w", err)
		}
	}

	return nil
}

// GetLastFetch returns the last fetch timestamp for fetchType="CHAMPIONS" or "GAMES"
func (db *Database) GetLastFetch(fetchType string) (time.Time, error) {
	var lastFetch string
	err := db.Conn.QueryRow("SELECT last_fetch FROM fetch WHERE fetch_type = ?;", fetchType).Scan(&lastFetch)
	if err != nil {
		return time.Time{}, err
	}

	lastFetchTime, err := time.Parse(time.RFC3339, lastFetch)
	if err != nil {
		return time.Time{}, err
	}

	return lastFetchTime, nil
}

// SetLastFetch sets the last fetch timestamp for fetchType="CHAMPIONS" or "GAMES"
func (db *Database) SetLastFetch(fetchType string, lastFetch time.Time) error {
	_, err := db.Conn.Exec("INSERT OR REPLACE INTO fetch (fetch_type, last_fetch) VALUES (?, ?);", fetchType, lastFetch.Format(time.RFC3339))
	return err
}
