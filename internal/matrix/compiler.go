package matrix

import (
	"sort"

	"github.com/ParadoxSportsData/paradox-clock-gate/internal/ingestion"
)

// Compile builds a StateMatrix from parsed plays and game header.
// The resulting matrix supports O(1) temporal queries with no future-state bleed.
func Compile(plays []ingestion.RawPlay, header ingestion.GameHeader) StateMatrix {
	// Step 1: sort by (GameClockTotalSeconds ASC, PlayID ASC) so higher PlayID wins ties.
	sort.Slice(plays, func(i, j int) bool {
		if plays[i].GameClockTotalSeconds != plays[j].GameClockTotalSeconds {
			return plays[i].GameClockTotalSeconds < plays[j].GameClockTotalSeconds
		}
		return plays[i].PlayID < plays[j].PlayID
	})

	// Step 2: pre-compute total description length for arena allocation.
	totalLen := 0
	for _, p := range plays {
		totalLen += len(p.Description)
	}

	var m StateMatrix
	m.Arena = make([]byte, 0, totalLen)
	m.Meta.GameID = header.GameID
	m.Meta.HomeTeam = header.HomeTeam
	m.Meta.AwayTeam = header.AwayTeam
	m.PlayTicks = make([]uint16, 0, len(plays))

	var maxTick int
	lastTick := -1

	// Step 3: write each play into States[tick].
	for _, p := range plays {
		tick := p.GameClockTotalSeconds
		if tick < 0 || tick >= MaxTick {
			continue
		}
		// Record each distinct tick once (plays sorted ASC by tick; consecutive
		// same-tick entries are adjacent and the last overwrites States[tick]).
		if tick != lastTick {
			m.PlayTicks = append(m.PlayTicks, uint16(tick))
			lastTick = tick
		}

		gs := GameState{
			Elapsed:  uint16(tick),
			Quarter:  uint8(p.Quarter),
			HasState: true,
		}

		if p.Down != nil {
			gs.Down = uint8(*p.Down)
		}
		if p.YardsToGo != nil {
			v := *p.YardsToGo
			if v > 255 {
				v = 255
			}
			gs.YardsToGo = uint8(v)
		}
		if p.YardLine100 != nil {
			gs.YardLine = uint8(*p.YardLine100)
		}
		if p.PlayType != nil {
			gs.PlayType = mapPlayType(*p.PlayType)
		}
		if p.WP != nil {
			wp := *p.WP
			if wp < 0 {
				wp = 0
			} else if wp > 1 {
				wp = 1
			}
			// Raw wp is possession-team WP; normalize to home-team perspective.
			if p.Posteam != nil && *p.Posteam != header.HomeTeam {
				wp = 1.0 - wp
			}
			gs.WinProb = uint16(wp * float64(WinProbScale))
		} else {
			gs.WinProb = WinProbNull // null sentinel
		}

		if p.Posteam != nil {
			copyTeam(&gs.Posteam, *p.Posteam)
		}
		if p.Defteam != nil {
			copyTeam(&gs.Defteam, *p.Defteam)
		}

		// Score attribution: map possession scores to home/away based on posteam.
		if p.Posteam != nil && p.PosteamScore != nil && p.DefteamScore != nil {
			if *p.Posteam == header.HomeTeam {
				gs.HomeScore = uint8(*p.PosteamScore)
				gs.AwayScore = uint8(*p.DefteamScore)
			} else {
				gs.AwayScore = uint8(*p.PosteamScore)
				gs.HomeScore = uint8(*p.DefteamScore)
			}
		}

		// Arena: append description and record offset+length.
		gs.DescOffset = uint32(len(m.Arena))
		gs.DescLen = uint16(len(p.Description))
		m.Arena = append(m.Arena, p.Description...)

		m.States[tick] = gs

		if tick > maxTick {
			maxTick = tick
		}
	}

	// Step 4: forward-fill — copy States[t-1] into any empty States[t].
	for t := 1; t <= maxTick; t++ {
		if !m.States[t].HasState && m.States[t-1].HasState {
			m.States[t] = m.States[t-1]
			m.States[t].Elapsed = uint16(t)
		}
	}

	m.Meta.MaxTick = uint16(maxTick)
	return m
}

func mapPlayType(s string) PlayType {
	switch s {
	case "run":
		return PlayTypeRun
	case "pass":
		return PlayTypePass
	case "punt":
		return PlayTypePunt
	case "kickoff":
		return PlayTypeKickoff
	case "no_play":
		return PlayTypeNoPlay
	default:
		return PlayTypeOther
	}
}

func copyTeam(dst *[3]byte, src string) {
	for i := range dst {
		dst[i] = 0
	}
	for i := 0; i < len(src) && i < 3; i++ {
		dst[i] = src[i]
	}
}
