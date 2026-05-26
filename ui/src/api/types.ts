// src/api/types.ts
// Copied verbatim from docs/api-contract.md — do NOT edit here directly.
// All changes must go through a contract amendment ticket (see api-contract.md).

export interface GameSummary {
  game_id: string;       // "2011_01_NO_GB"
  home_team: string;     // "GB"
  away_team: string;     // "NO"
  home_score: number;
  away_score: number;
  duration: number;      // maxTick in seconds
}

export interface GameStateResponse {
  tick: number;
  quarter: number;           // 1-6
  down: number | null;       // 1-4, null for kickoff/special
  yards_to_go: number | null;
  yard_line: number | null;  // 0-100, yards to opponent endzone
  home_score: number;
  away_score: number;
  posteam: string | null;    // "GB", "NO", etc.
  defteam: string | null;
  win_prob: number | null;   // 0.0-1.0, null if unknown
  play_type: string;         // "run"|"pass"|"punt"|"kickoff"|"no_play"|"other"|""
  description: string;
  has_state: boolean;        // false = before any play at this tick
}

export interface PlaySnapshot {
  tick: number;
  quarter: number;
  down: number | null;
  yards_to_go: number | null;
  yard_line: number | null;
  home_score: number;
  away_score: number;
  posteam: string | null;
  win_prob: number | null;
  play_type: string;
  description: string;
}

export interface GameTimelineResponse {
  game_id: string;
  home_team: string;
  away_team: string;
  max_tick: number;
  plays: PlaySnapshot[]; // only has_state=true ticks, sorted by tick ASC
}

export interface ApiError {
  error: string;
  max_tick?: number; // present on 422 only
}
