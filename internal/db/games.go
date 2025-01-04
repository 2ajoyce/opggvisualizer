// internal/db/games.go
package db

import (
	"fmt"
	"log"
	"opggvisualizer/internal/models"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Insert a game into the games table
func (db *Database) InsertGame(game models.Game) error {
	insertGameSQL := `INSERT INTO games(
		game_id, created_at, game_length, tier, division, tier_image_url, border_image_url,
		is_remake, meta_version, game_type, is_opscore_active, is_recorded, version,
		first_game_created_at, last_game_created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	createdAt := game.CreatedAt.Format(time.RFC3339)
	firstGameCreatedAt := game.Meta.FirstGameCreatedAt.Format(time.RFC3339)
	lastGameCreatedAt := game.Meta.LastGameCreatedAt.Format(time.RFC3339)

	_, err := db.Conn.Exec(insertGameSQL,
		game.ID,
		createdAt,
		game.GameLengthSecond,
		game.AverageTierInfo.Tier,
		int(game.AverageTierInfo.Division),
		game.AverageTierInfo.TierImageURL,
		game.AverageTierInfo.BorderImageURL,
		game.IsRemake,
		game.MetaVersion,
		game.GameType,
		game.IsOpscoreActive,
		game.IsRecorded,
		game.Version,
		firstGameCreatedAt,
		lastGameCreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert game: %w", err)
	}
	return nil
}

// Insert a team into the teams table
func (db *Database) InsertTeam(gameID string, team models.Team) error {
	insertTeamSQL := `INSERT INTO teams(
		game_id, key, is_win, champion_first, inhibitor_first, rift_herald_first, death,
		champion_kill, inhibitor_kill, dragon_first, horde_first, rift_herald_kill,
		is_remake, gold_earned, kill, tower_first, horde_kill, assist, dragon_kill,
		baron_kill, baron_first, tower_kill
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	result, err := db.Conn.Exec(insertTeamSQL,
		gameID,
		team.Key,
		team.GameStat.IsWin,
		team.GameStat.ChampionFirst,
		team.GameStat.InhibitorFirst,
		team.GameStat.RiftHeraldFirst,
		team.GameStat.Death,
		team.GameStat.ChampionKill,
		team.GameStat.InhibitorKill,
		team.GameStat.DragonFirst,
		team.GameStat.HordeFirst,
		team.GameStat.RiftHeraldKill,
		team.GameStat.IsRemake,
		team.GameStat.GoldEarned,
		team.GameStat.Kill,
		team.GameStat.TowerFirst,
		team.GameStat.HordeKill,
		team.GameStat.Assist,
		team.GameStat.DragonKill,
		team.GameStat.BaronKill,
		team.GameStat.BaronFirst,
		team.GameStat.TowerKill,
	)
	if err != nil {
		return fmt.Errorf("failed to insert team: %w", err)
	}

	teamID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve last insert ID for team: %w", err)
	}

	// Insert banned champions
	for _, bannedChamp := range team.BannedChampions {
		if err := db.InsertTeamBannedChampion(int(teamID), bannedChamp); err != nil {
			log.Printf("Error inserting banned champion %v for team %d: %v", bannedChamp, teamID, err)
			continue
		}
	}

	return nil
}

// InsertTeamBannedChampion inserts a banned champion into the team_banned_champions table
func (db *Database) InsertTeamBannedChampion(teamID int, bannedChampion float64) error {
	insertBannedChampionSQL := `INSERT INTO team_banned_champions(
		team_id, banned_champion_id
	) VALUES (?, ?);`

	_, err := db.Conn.Exec(insertBannedChampionSQL,
		teamID,
		bannedChampion,
	)
	if err != nil {
		return fmt.Errorf("failed to insert banned champion: %w", err)
	}
	return nil
}

// InsertParticipant inserts a participant into the participants table
func (db *Database) InsertParticipant(gameID string, participant models.Participant) error {
	insertParticipantSQL := `INSERT INTO participants(
		game_id, participant_id, summoner_name, champion_id, position, role, kills,
		deaths, assists, gold_earned, damage_dealt, damage_taken, vision_score,
		primary_rune_id, secondary_rune_page_id, lane_score,
		team_key, result, ward_place, op_score_rank, barrack_kill, total_heal,
		game_type, is_remake
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	// Convert ChampionID from int to string
	championIDStr := strconv.Itoa(participant.ChampionID)

	result, err := db.Conn.Exec(insertParticipantSQL,
		gameID,
		participant.ParticipantID,
		participant.Summoner.Name,
		championIDStr, // Converted to string
		participant.Position,
		participant.Role,
		int(participant.Stats.Kill),
		int(participant.Stats.Death),
		int(participant.Stats.Assist),
		int(participant.Stats.GoldEarned),
		int(participant.Stats.TotalDamageDealtToChampions),
		int(participant.Stats.TotalDamageTaken),
		int(participant.Stats.VisionScore),
		int(participant.Rune.PrimaryRuneID),
		int(participant.Rune.SecondaryPageID),
		int(participant.Stats.LaneScore),
		participant.TeamKey,
		participant.Stats.Result,
		int(participant.Stats.WardPlace),
		int(participant.Stats.OpScoreRank),
		int(participant.Stats.BarrackKill),
		int(participant.Stats.TotalHeal),
		participant.GameType,
		participant.IsRemake,
	)
	if err != nil {
		return fmt.Errorf("failed to insert participant: %w", err)
	}

	participantDBID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve last insert ID for participant: %w", err)
	}

	// Insert items
	for _, item := range participant.Items {
		if err := db.InsertParticipantItem(int(participantDBID), int(item)); err != nil {
			log.Printf("Error inserting item %v for participant %d: %v", item, participantDBID, err)
			continue
		}
	}

	// Insert spells
	for _, spell := range participant.Spells {
		if err := db.InsertParticipantSpell(int(participantDBID), int(spell)); err != nil {
			log.Printf("Error inserting spell %v for participant %d: %v", spell, participantDBID, err)
			continue
		}
	}

	return nil
}

// InsertParticipantItem inserts an item into the participant_items table
func (db *Database) InsertParticipantItem(participantID int, itemID int) error {
	insertItemSQL := `INSERT INTO participant_items(
		participant_id, item_id
	) VALUES (?, ?);`

	_, err := db.Conn.Exec(insertItemSQL,
		participantID,
		itemID,
	)
	if err != nil {
		return fmt.Errorf("failed to insert participant item: %w", err)
	}
	return nil
}

// InsertParticipantSpell inserts a spell into the participant_spells table
func (db *Database) InsertParticipantSpell(participantID int, spellID int) error {
	insertSpellSQL := `INSERT INTO participant_spells(
		participant_id, spell_id
	) VALUES (?, ?);`

	_, err := db.Conn.Exec(insertSpellSQL,
		participantID,
		spellID,
	)
	if err != nil {
		return fmt.Errorf("failed to insert participant spell: %w", err)
	}
	return nil
}

// ClearGameData clears all data from the game-related tables
func (db *Database) ClearGameData() error {
	tables := []string{"games", "teams", "players", "stats"}
	for _, table := range tables {
		if _, err := db.Conn.Exec(fmt.Sprintf(`DELETE FROM %s;`, table)); err != nil {
			return fmt.Errorf("failed to clear data from table %s: %w", table, err)
		}
	}

	// Wipe last fetch time for games
	if _, err := db.Conn.Exec(`UPDATE fetch SET last_fetch = NULL WHERE fetch_type = "GAMES";`); err != nil {
		return fmt.Errorf("failed to clear last fetch time for games: %w", err)
	}
	return nil
}
