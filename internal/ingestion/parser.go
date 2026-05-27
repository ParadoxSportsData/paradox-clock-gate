package ingestion

import (
	"encoding/json"
	"fmt"
	"os"
)

// ParseFile reads a game JSON file and returns the header and all plays.
// Each file is a wrapper object: {game_id, home_team, away_team, home_score, away_score, plays: [...]}
func ParseFile(path string) (GameHeader, []RawPlay, error) {
	f, err := os.Open(path)
	if err != nil {
		return GameHeader{}, nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)

	// Consume opening '{' of the wrapper object.
	if _, err := dec.Token(); err != nil {
		return GameHeader{}, nil, fmt.Errorf("parse %s: expected '{': %w", path, err)
	}

	var header GameHeader
	var plays []RawPlay

	// Read key-value pairs until the closing '}'.
	for dec.More() {
		keyToken, err := dec.Token()
		if err != nil {
			return GameHeader{}, nil, fmt.Errorf("parse %s key: %w", path, err)
		}
		key, ok := keyToken.(string)
		if !ok {
			return GameHeader{}, nil, fmt.Errorf("parse %s: expected string key, got %T", path, keyToken)
		}

		switch key {
		case "game_id":
			if err := dec.Decode(&header.GameID); err != nil {
				return GameHeader{}, nil, fmt.Errorf("parse %s game_id: %w", path, err)
			}
		case "home_team":
			if err := dec.Decode(&header.HomeTeam); err != nil {
				return GameHeader{}, nil, fmt.Errorf("parse %s home_team: %w", path, err)
			}
		case "away_team":
			if err := dec.Decode(&header.AwayTeam); err != nil {
				return GameHeader{}, nil, fmt.Errorf("parse %s away_team: %w", path, err)
			}
		case "home_score":
			if err := dec.Decode(&header.HomeScore); err != nil {
				return GameHeader{}, nil, fmt.Errorf("parse %s home_score: %w", path, err)
			}
		case "away_score":
			if err := dec.Decode(&header.AwayScore); err != nil {
				return GameHeader{}, nil, fmt.Errorf("parse %s away_score: %w", path, err)
			}
		case "game_date":
			if err := dec.Decode(&header.GameDate); err != nil {
				return GameHeader{}, nil, fmt.Errorf("parse %s game_date: %w", path, err)
			}
		case "week":
			if err := dec.Decode(&header.Week); err != nil {
				return GameHeader{}, nil, fmt.Errorf("parse %s week: %w", path, err)
			}
		case "season":
			if err := dec.Decode(&header.Season); err != nil {
				return GameHeader{}, nil, fmt.Errorf("parse %s season: %w", path, err)
			}
		case "plays":
			// Consume '['.
			if _, err := dec.Token(); err != nil {
				return GameHeader{}, nil, fmt.Errorf("parse %s plays '[': %w", path, err)
			}
			for dec.More() {
				var p RawPlay
				if err := dec.Decode(&p); err != nil {
					return GameHeader{}, nil, fmt.Errorf("parse %s play: %w", path, err)
				}
				plays = append(plays, p)
			}
			// Consume ']'.
			if _, err := dec.Token(); err != nil {
				return GameHeader{}, nil, fmt.Errorf("parse %s plays ']': %w", path, err)
			}
		default:
			// Skip unknown fields by consuming their value token(s).
			if err := skipValue(dec); err != nil {
				return GameHeader{}, nil, fmt.Errorf("parse %s skip %q: %w", path, key, err)
			}
		}
	}

	return header, plays, nil
}

// skipValue discards one JSON value from the decoder, regardless of depth.
func skipValue(dec *json.Decoder) error {
	t, err := dec.Token()
	if err != nil {
		return err
	}
	switch t {
	case json.Delim('{'), json.Delim('['):
		for dec.More() {
			if err := skipValue(dec); err != nil {
				return err
			}
		}
		// Consume closing delimiter.
		if _, err := dec.Token(); err != nil {
			return err
		}
	}
	return nil
}
