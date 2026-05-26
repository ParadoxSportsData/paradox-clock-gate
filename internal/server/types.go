package server

type GameSummary struct {
	GameID    string `json:"game_id"`
	HomeTeam  string `json:"home_team"`
	AwayTeam  string `json:"away_team"`
	HomeScore uint8  `json:"home_score"`
	AwayScore uint8  `json:"away_score"`
	Duration  uint16 `json:"duration"` // maxTick in seconds
}

type GameStateResponse struct {
	Tick        uint16   `json:"tick"`
	Quarter     uint8    `json:"quarter"`
	Down        *uint8   `json:"down"`
	YardsToGo   *uint8   `json:"yards_to_go"`
	YardLine    *uint8   `json:"yard_line"`
	HomeScore   uint8    `json:"home_score"`
	AwayScore   uint8    `json:"away_score"`
	Posteam     *string  `json:"posteam"`
	Defteam     *string  `json:"defteam"`
	WinProb     *float64 `json:"win_prob"`
	PlayType    string   `json:"play_type"`
	Description string   `json:"description"`
	HasState    bool     `json:"has_state"`
}

type PlaySnapshot struct {
	Tick        uint16   `json:"tick"`
	Quarter     uint8    `json:"quarter"`
	Down        *uint8   `json:"down"`
	YardsToGo   *uint8   `json:"yards_to_go"`
	YardLine    *uint8   `json:"yard_line"`
	HomeScore   uint8    `json:"home_score"`
	AwayScore   uint8    `json:"away_score"`
	Posteam     *string  `json:"posteam"`
	WinProb     *float64 `json:"win_prob"`
	PlayType    string   `json:"play_type"`
	Description string   `json:"description"`
}

type GameTimelineResponse struct {
	GameID   string         `json:"game_id"`
	HomeTeam string         `json:"home_team"`
	AwayTeam string         `json:"away_team"`
	MaxTick  uint16         `json:"max_tick"`
	Plays    []PlaySnapshot `json:"plays"`
}

type ErrorResponse struct {
	Error   string  `json:"error"`
	MaxTick *uint16 `json:"max_tick,omitempty"`
}
