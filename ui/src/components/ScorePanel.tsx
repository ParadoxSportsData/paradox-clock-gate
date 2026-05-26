// src/components/ScorePanel.tsx
// PDX-25: Pure display — score and possession indicator.
// No hooks, no fetches. All null-safe.

import type { PlaySnapshot } from '../api/schemas'

interface ScorePanelProps {
  homeTeam: string
  awayTeam: string
  play: PlaySnapshot | null
}

export function ScorePanel({ homeTeam, awayTeam, play }: ScorePanelProps) {
  const homeScore = play?.home_score ?? 0
  const awayScore = play?.away_score ?? 0
  const posteam = play?.posteam ?? null

  const homeLeads = homeScore > awayScore
  const awayLeads = awayScore > homeScore

  return (
    <div className="bg-gray-800 rounded-lg p-4 flex items-center justify-between gap-4">
      {/* Away team */}
      <div className="flex flex-col items-center flex-1">
        <div className="flex items-center gap-1">
          {posteam === awayTeam && (
            <span className="text-yellow-400 text-xs">▶</span>
          )}
          <span className="text-sm text-gray-400 font-mono">{awayTeam}</span>
        </div>
        <span className={`text-3xl font-bold ${awayLeads ? 'text-white' : 'text-gray-400'}`}>
          {awayScore}
        </span>
      </div>

      <div className="text-gray-600 text-lg font-light">–</div>

      {/* Home team */}
      <div className="flex flex-col items-center flex-1">
        <div className="flex items-center gap-1">
          <span className="text-sm text-gray-400 font-mono">{homeTeam}</span>
          {posteam === homeTeam && (
            <span className="text-yellow-400 text-xs">◀</span>
          )}
        </div>
        <span className={`text-3xl font-bold ${homeLeads ? 'text-white' : 'text-gray-400'}`}>
          {homeScore}
        </span>
      </div>
    </div>
  )
}
