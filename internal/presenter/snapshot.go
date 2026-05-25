package presenter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ParadoxSportsData/paradox-clock-gate/internal/matrix"
)

// RenderText returns the box-drawing text view of a game state.
func RenderText(gs matrix.GameState, meta matrix.GameMeta, arena []byte) string {
	if !gs.HasState {
		return "┌─────────────────────────────────────┐\n│  No state available at this tick.  │\n└─────────────────────────────────────┘\n"
	}

	posteam := teamStr(gs.Posteam)
	defteam := teamStr(gs.Defteam)

	elapsed := gs.Elapsed
	mins := elapsed / 60
	secs := elapsed % 60
	elapsedStr := fmt.Sprintf("%ds (%d:%02d)", elapsed, mins, secs)

	qStr := fmt.Sprintf("Q%d", gs.Quarter)

	// Down/distance string.
	downStr := ""
	if gs.Down > 0 {
		downStr = fmt.Sprintf(" │  %s & %d  at %s %d", ordinal(gs.Down), gs.YardsToGo, posteam, gs.YardLine)
	}

	// Win probability — relative to posteam.
	wpStr := ""
	if gs.WinProb != 65535 {
		wpPct := float64(gs.WinProb) / 100.0
		wpStr = fmt.Sprintf("  Win Prob: %s %.1f%%", posteam, wpPct)
	}

	// Play description from arena.
	desc := ""
	if int(gs.DescOffset)+int(gs.DescLen) <= len(arena) {
		desc = string(arena[gs.DescOffset : gs.DescOffset+uint32(gs.DescLen)])
	}

	width := 60
	line := strings.Repeat("─", width)
	header := fmt.Sprintf("  %s @ %s   │  %s  │  Elapsed: %s", meta.AwayTeam, meta.HomeTeam, qStr, elapsedStr)
	score := fmt.Sprintf("  Score:  %s %d  –  %s %d", meta.HomeTeam, gs.HomeScore, meta.AwayTeam, gs.AwayScore)
	ball := fmt.Sprintf("  Ball:   %s possession%s", posteam, downStr)

	_ = defteam

	var sb strings.Builder
	sb.WriteString("┌" + line + "┐\n")
	sb.WriteString("│  " + padRight(strings.TrimPrefix(header, "  "), width-2) + "│\n")
	sb.WriteString("├" + line + "┤\n")
	sb.WriteString("│  " + padRight(strings.TrimPrefix(score, "  "), width-2) + "│\n")
	sb.WriteString("│  " + padRight(strings.TrimPrefix(ball, "  "), width-2) + "│\n")
	if wpStr != "" {
		sb.WriteString("│  " + padRight(strings.TrimPrefix(wpStr, "  "), width-2) + "│\n")
	}
	sb.WriteString("├" + line + "┤\n")
	if desc != "" {
		// Truncate description if too long for the box.
		maxDesc := width - 4
		if len(desc) > maxDesc {
			desc = desc[:maxDesc-1] + "…"
		}
		sb.WriteString("│  " + padRight(desc, width-2) + "│\n")
	}
	sb.WriteString("└" + line + "┘\n")
	return sb.String()
}

type jsonState struct {
	Elapsed   uint16  `json:"elapsed"`
	Quarter   uint8   `json:"quarter"`
	Down      uint8   `json:"down"`
	YardsToGo uint8   `json:"yards_to_go"`
	YardLine  uint8   `json:"yard_line"`
	HomeScore uint8   `json:"home_score"`
	AwayScore uint8   `json:"away_score"`
	Posteam   string  `json:"posteam"`
	Defteam   string  `json:"defteam"`
	WinProb   *float64 `json:"win_prob"`
	HasState  bool    `json:"has_state"`
	PlayDesc  string  `json:"play_description"`
}

// RenderJSON returns a JSON representation of the game state.
func RenderJSON(gs matrix.GameState, meta matrix.GameMeta, arena []byte) string {
	js := jsonState{
		Elapsed:   gs.Elapsed,
		Quarter:   gs.Quarter,
		Down:      gs.Down,
		YardsToGo: gs.YardsToGo,
		YardLine:  gs.YardLine,
		HomeScore: gs.HomeScore,
		AwayScore: gs.AwayScore,
		Posteam:   teamStr(gs.Posteam),
		Defteam:   teamStr(gs.Defteam),
		HasState:  gs.HasState,
	}
	if gs.WinProb != 65535 {
		v := float64(gs.WinProb) / 10000.0
		js.WinProb = &v
	}
	if gs.HasState && int(gs.DescOffset)+int(gs.DescLen) <= len(arena) {
		js.PlayDesc = string(arena[gs.DescOffset : gs.DescOffset+uint32(gs.DescLen)])
	}
	b, _ := json.MarshalIndent(js, "", "  ")
	return string(b)
}

func teamStr(b [3]byte) string {
	end := 0
	for end < 3 && b[end] != 0 {
		end++
	}
	return string(b[:end])
}

func padRight(s string, n int) string {
	if len(s) >= n {
		return s[:n]
	}
	return s + strings.Repeat(" ", n-len(s))
}

func ordinal(n uint8) string {
	switch n {
	case 1:
		return "1st"
	case 2:
		return "2nd"
	case 3:
		return "3rd"
	default:
		return fmt.Sprintf("%dth", n)
	}
}
