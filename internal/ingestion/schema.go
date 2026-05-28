// Package ingestion parses NFL play-by-play JSON files into typed structs for
// downstream compilation into a StateMatrix.
package ingestion

// GameHeader holds the top-level wrapper fields parsed before the plays array.
// home_team, away_team, home_score, away_score come from the JSON header —
// NOT from filename and NOT from play data (posteam/defteam change every drive).
type GameHeader struct {
	GameID    string `json:"game_id"`
	HomeTeam  string `json:"home_team"`
	AwayTeam  string `json:"away_team"`
	HomeScore int    `json:"home_score"`
	AwayScore int    `json:"away_score"`
	GameDate  string `json:"game_date"` // ISO-8601 e.g. "2011-09-08T00:00:00"
	Week      int    `json:"week"`
	Season    int    `json:"season"`
}

type RawPlay struct {
	PlayID                int      `json:"play_id"`
	Quarter               int      `json:"quarter"`
	GameClock             string   `json:"game_clock"`
	GameClockTotalSeconds int      `json:"game_clock_total_seconds"`
	Down                  *int     `json:"down"`
	YardsToGo             *int     `json:"ydstogo"`
	YardLine100           *int     `json:"yardline_100"`
	PlayType              *string  `json:"play_type"`
	YardsGained           *int     `json:"yards_gained"`
	Description           string   `json:"description"`
	Posteam               *string  `json:"posteam"`
	Defteam               *string  `json:"defteam"`
	PosteamScore          *int     `json:"posteam_score"`
	DefteamScore          *int     `json:"defteam_score"`
	WP                    *float64 `json:"wp"`
	EPA                   *float64 `json:"epa"`
}
