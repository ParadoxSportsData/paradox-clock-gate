// src/components/DownDistance.tsx
// PDX-25: Pure display — quarter, down & distance, field position.
// No hooks. All nullable fields guarded.

import type { PlaySnapshot } from '../api/schemas'

interface DownDistanceProps {
  play: PlaySnapshot | null
}

function ordinal(n: number): string {
  const suffixes = ['th', 'st', 'nd', 'rd']
  const v = n % 100
  return n + (suffixes[(v - 20) % 10] ?? suffixes[v] ?? suffixes[0])
}

export function DownDistance({ play }: DownDistanceProps) {
  if (play === null) {
    return (
      <div className="bg-gray-800 rounded-lg p-4 flex items-center justify-center">
        <span className="text-gray-500 italic text-sm">Pre-game</span>
      </div>
    )
  }

  const { quarter, down, yards_to_go, yard_line, play_type } = play

  // Kickoff / special — no down
  if (down === null) {
    return (
      <div className="bg-gray-800 rounded-lg p-4">
        <div className="text-xs text-gray-500 uppercase tracking-wide mb-1">Situation</div>
        <div className="text-white font-semibold">
          Q{quarter} — <span className="capitalize">{play_type || 'special'}</span>
        </div>
      </div>
    )
  }

  const yardLineDisplay = yard_line !== null
    ? `at ${100 - yard_line} yd line`
    : ''

  return (
    <div className="bg-gray-800 rounded-lg p-4">
      <div className="text-xs text-gray-500 uppercase tracking-wide mb-1">Situation</div>
      <div className="text-white font-semibold">
        Q{quarter} &nbsp;|&nbsp; {ordinal(down)} &amp; {yards_to_go ?? '?'}
        {yardLineDisplay && <span className="text-gray-400 font-normal"> &nbsp;{yardLineDisplay}</span>}
      </div>
    </div>
  )
}
