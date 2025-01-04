// internal/models/models.go
package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Config struct {
	SummonerID   string
	DatabasePath string
}

// GameData represents the structure of game data API response
type GameData struct {
	Meta GameDataMeta `json:"meta"`
	Data []GameEntry  `json:"data"`
}

// GameDataMeta holds the meta information in game data
type GameDataMeta struct {
	FirstGameCreatedAt string `json:"first_game_created_at"`
	LastGameCreatedAt  string `json:"last_game_created_at"`
}

type GameEntry struct {
	IsOpscoreActive  bool            `json:"is_opscore_active"`
	IsRecorded       sql.NullFloat64 `json:"is_recorded"` // Using sql.NullFloat64 to handle <nil>
	Teams            []Team          `json:"teams"`
	Memo             interface{}     `json:"memo"`
	ID               string          `json:"id"`
	Version          string          `json:"version"`
	GameLengthSecond float64         `json:"game_length_second"`
	AverageTierInfo  TierInfo        `json:"average_tier_info"`
	MyData           MyData          `json:"myData"`
	GameType         string          `json:"game_type"`
	IsRemake         bool            `json:"is_remake"`
	RecordInfo       interface{}     `json:"record_info"`
	CreatedAt        string          `json:"created_at"`
	Participants     []Participant   `json:"participants"`
	GameMap          string          `json:"game_map"`
	MetaVersion      string          `json:"meta_version"`
}

type TierInfo struct {
	Tier           string  `json:"tier"`
	Division       float64 `json:"division"`
	TierImageURL   string  `json:"tier_image_url"`
	BorderImageURL string  `json:"border_image_url"`
}

type MyData struct {
	Summoner      Summoner         `json:"summoner"`
	ParticipantID float64          `json:"participant_id"`
	ChampionID    int              `json:"champion_id"` // Changed to int
	Stats         ParticipantStats `json:"stats"`
	Rune          Rune             `json:"rune"`
	TeamKey       string           `json:"team_key"`
	// Add other fields as necessary
}

type Summoner struct {
	RevisionAt      string      `json:"revision_at"`
	ID              float64     `json:"id"`
	AcctID          string      `json:"acct_id"`
	Name            string      `json:"name"`
	InternalName    string      `json:"internal_name"`
	Level           float64     `json:"level"`
	UpdatedAt       string      `json:"updated_at"`
	RenewableAt     string      `json:"renewable_at"`
	Puuid           string      `json:"puuid"`
	GameName        string      `json:"game_name"`
	ProfileImageURL string      `json:"profile_image_url"`
	Player          interface{} `json:"player"`
	SummonerID      string      `json:"summoner_id"`
	Tagline         string      `json:"tagline"`
}

type ParticipantStats struct {
	WardPlace                      float64                 `json:"ward_place"`
	Result                         string                  `json:"result"`
	OpScoreRank                    float64                 `json:"op_score_rank"`
	ChampionLevel                  float64                 `json:"champion_level"`
	BarrackKill                    float64                 `json:"barrack_kill"`
	TotalHeal                      float64                 `json:"total_heal"`
	NeutralMinionKillEnemyJungle   sql.NullFloat64         `json:"neutral_minion_kill_enemy_jungle"`
	NeutralMinionKill              float64                 `json:"neutral_minion_kill"`
	OpScoreTimeline                []OpScoreTimeline       `json:"op_score_timeline"`
	TotalDamageTaken               float64                 `json:"total_damage_taken"`
	VisionWardsBoughtInGame        float64                 `json:"vision_wards_bought_in_game"`
	TurretKill                     float64                 `json:"turret_kill"`
	Assist                         float64                 `json:"assist"`
	NeutralMinionKillTeamJungle    sql.NullFloat64         `json:"neutral_minion_kill_team_jungle"`
	OpScore                        float64                 `json:"op_score"`
	LaneScore                      float64                 `json:"lane_score"`
	DamageSelfMitigated            float64                 `json:"damage_self_mitigated"`
	MagicDamageDealtPlayer         float64                 `json:"magic_damage_dealt_player"`
	TimeCcingOthers                float64                 `json:"time_ccing_others"`
	LargestKillingSpree            float64                 `json:"largest_killing_spree"`
	DamageDealtToTurrets           float64                 `json:"damage_dealt_to_turrets"`
	LargestCriticalStrike          float64                 `json:"largest_critical_strike"`
	Kill                           float64                 `json:"kill"`
	IsOpscoreMaxInTeam             bool                    `json:"is_opscore_max_in_team"`
	PhysicalDamageDealtToChampions float64                 `json:"physical_damage_dealt_to_champions"`
	TotalDamageDealtToChampions    float64                 `json:"total_damage_dealt_to_champions"`
	SightWardsBoughtInGame         float64                 `json:"sight_wards_bought_in_game"`
	LargestMultiKill               float64                 `json:"largest_multi_kill"`
	GoldEarned                     float64                 `json:"gold_earned"`
	OpScoreTimelineAnalysis        OpScoreTimelineAnalysis `json:"op_score_timeline_analysis"`
	Keyword                        string                  `json:"keyword"`
	DamageDealtToObjectives        float64                 `json:"damage_dealt_to_objectives"`
	TotalDamageDealt               float64                 `json:"total_damage_dealt"`
	WardKill                       float64                 `json:"ward_kill"`
	Death                          float64                 `json:"death"`
	PhysicalDamageTaken            float64                 `json:"physical_damage_taken"`
	VisionScore                    float64                 `json:"vision_score"`
	MinionKill                     float64                 `json:"minion_kill"`
}

type OpScoreTimeline struct {
	Second float64 `json:"second"`
	Score  float64 `json:"score"`
}

type OpScoreTimelineAnalysis struct {
	Left  string `json:"left"`
	Right string `json:"right"`
	Last  string `json:"last"`
}

type Rune struct {
	PrimaryPageID   float64 `json:"primary_page_id"`
	PrimaryRuneID   float64 `json:"primary_rune_id"`
	SecondaryPageID float64 `json:"secondary_page_id"`
}

type Participant struct {
	ParticipantID int              `json:"participant_id"` // Changed to int
	SummonerName  string           `json:"summoner_name"`
	ChampionID    int              `json:"champion_id"` // Changed to int
	Position      string           `json:"position"`
	Role          string           `json:"role"`
	TierInfo      TierInfo         `json:"tier_info"`
	Summoner      SummonerDetailed `json:"summoner"`
	Rune          Rune             `json:"rune"`
	Stats         ParticipantStats `json:"stats"`
	TeamKey       string           `json:"team_key"`
	GameType      string           `json:"game_type"`
	IsRemake      bool             `json:"is_remake"`
	RecordInfo    interface{}      `json:"record_info"`
	CreatedAt     string           `json:"created_at"`
	Items         []float64        `json:"items"`
	Spells        []float64        `json:"spells"`
	TrinketItem   float64          `json:"trinket_item"`
}

type SummonerDetailed struct {
	SummonerID      string      `json:"summoner_id"`
	GameName        string      `json:"game_name"`
	Tagline         string      `json:"tagline"`
	ProfileImageURL string      `json:"profile_image_url"`
	Player          interface{} `json:"player"`
	ID              float64     `json:"id"`
	AcctID          string      `json:"acct_id"`
	Name            string      `json:"name"`
	InternalName    string      `json:"internal_name"`
	Level           float64     `json:"level"`
	UpdatedAt       string      `json:"updated_at"`
	RenewableAt     string      `json:"renewable_at"`
	RevisionAt      string      `json:"revision_at"`
	Puuid           string      `json:"puuid"`
}

type ChampionData struct {
	Type    string              `json:"type"`
	Format  string              `json:"format"`
	Version string              `json:"version"`
	Data    map[string]Champion `json:"data"`
}

type Champion struct {
	ID      string             `json:"id"`
	Key     string             `json:"key"`
	Name    string             `json:"name"`
	Title   string             `json:"title"`
	Blurb   string             `json:"blurb"`
	Partype string             `json:"partype"`
	Info    ChampionInfo       `json:"info"`
	Stats   map[string]float64 `json:"stats"`
	Image   Image              `json:"image"`
	Tags    []string           `json:"tags"`
	Type    string             `json:"type"`
	Format  string             `json:"format"`
}

// ChampionInfo represents the attack, defense, magic, and difficulty stats of a champion
type ChampionInfo struct {
	Attack     float64 `json:"attack"`
	Defense    float64 `json:"defense"`
	Magic      float64 `json:"magic"`
	Difficulty float64 `json:"difficulty"`
}

// Image represents the image data for a champion
type Image struct {
	W       float64 `json:"w"`
	H       float64 `json:"h"`
	Full    string  `json:"full"`
	Sprite  string  `json:"sprite"`
	Group   string  `json:"group"`
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	Version string  `json:"version"`
}

// ImageURL constructs the full image URL for the champion
func (img Image) FullURL() string {
	return fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%s/img/champion/%s", img.Version, img.Full)
}

// Game represents a game entry to be inserted into the database
type Game struct {
	ID               string    `json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	GameLengthSecond int       `json:"game_length_second"`
	AverageTierInfo  TierInfo  `json:"average_tier_info"`
	IsRemake         bool      `json:"is_remake"`
	MetaVersion      string    `json:"meta_version"`
	GameType         string    `json:"game_type"`
	IsOpscoreActive  bool      `json:"is_opscore_active"`
	IsRecorded       bool      `json:"is_recorded"`
	Version          string    `json:"version"`
	Meta             GameMeta  `json:"meta"`
}

// GameMeta holds the meta information with time.Time fields
type GameMeta struct {
	FirstGameCreatedAt time.Time
	LastGameCreatedAt  time.Time
}

// Team represents a team entry to be inserted into the database
type Team struct {
	Key             string    `json:"key"`
	GameStat        TeamStat  `json:"game_stat"`
	BannedChampions []float64 `json:"banned_champions"`
}

// TeamStat holds the statistics for a team
type TeamStat struct {
	IsWin           bool `json:"is_win"`
	ChampionFirst   bool `json:"champion_first"`
	InhibitorFirst  bool `json:"inhibitor_first"`
	RiftHeraldFirst bool `json:"rift_herald_first"`
	Death           int  `json:"death"`
	ChampionKill    int  `json:"champion_kill"`
	InhibitorKill   int  `json:"inhibitor_kill"`
	DragonFirst     bool `json:"dragon_first"`
	HordeFirst      bool `json:"horde_first"`
	RiftHeraldKill  int  `json:"rift_herald_kill"`
	IsRemake        bool `json:"is_remake"`
	GoldEarned      int  `json:"gold_earned"`
	Kill            int  `json:"kill"`
	TowerFirst      bool `json:"tower_first"`
	HordeKill       int  `json:"horde_kill"`
	Assist          int  `json:"assist"`
	DragonKill      int  `json:"dragon_kill"`
	BaronKill       int  `json:"baron_kill"`
	BaronFirst      bool `json:"baron_first"`
	TowerKill       int  `json:"tower_kill"`
}

type FetchRecord struct {
	FetchType string // "GAMES" or "CHAMPIONS"
	LastFetch time.Time // The last time the records were fetched
}