// src/components/TimelineScrubber.tsx
// PDX-24: Range slider over 0..maxTick with quarter markers.
// Full timeline loaded once via TanStack Query. All scrubbing is local binary search — no HTTP.

import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { getTimeline } from '../api/client'
import type { PlaySnapshot } from '../api/schemas'

interface TimelineScrubberProps {
  gameId: string
  onTickChange: (tick: number, play: PlaySnapshot | null) => void
}

// O(log N) binary search: largest play.tick <= targetTick
function findNearestPlay(plays: PlaySnapshot[], targetTick: number): PlaySnapshot | null {
  if (plays.length === 0) return null
  let lo = 0
  let hi = plays.length - 1
  let result = plays[0]
  while (lo <= hi) {
    const mid = (lo + hi) >> 1
    if (plays[mid].tick <= targetTick) {
      result = plays[mid]
      lo = mid + 1
    } else {
      hi = mid - 1
    }
  }
  return result
}

function tickToMMSS(tick: number): string {
  const m = Math.floor(tick / 60)
  const s = tick % 60
  return `${m}:${String(s).padStart(2, '0')}`
}

const QUARTER_TICKS = [900, 1800, 2700, 3600]
const QUARTER_LABELS = ['Q1', 'Q2', 'Q3', 'Q4']

export function TimelineScrubber({ gameId, onTickChange }: TimelineScrubberProps) {
  const [tick, setTick] = useState(0)

  const query = useQuery({
    queryKey: ['timeline', gameId],
    queryFn: () => getTimeline(gameId),
  })

  const maxTick = query.data?.max_tick ?? 3600

  function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    const newTick = Number(e.target.value)
    setTick(newTick)
    const play = query.data ? findNearestPlay(query.data.plays, newTick) : null
    onTickChange(newTick, play)
  }

  const nearestPlay = query.data ? findNearestPlay(query.data.plays, tick) : null
  const quarter = nearestPlay?.quarter ?? 1
  const displayTime = tickToMMSS(tick)

  if (query.isLoading) {
    return (
      <div className="px-6 py-4">
        <div className="h-8 bg-gray-700 rounded animate-pulse" />
        <div className="text-gray-500 text-sm mt-1">Loading timeline…</div>
      </div>
    )
  }

  if (query.isError) {
    return (
      <div className="px-6 py-2 text-red-400 text-sm">Timeline unavailable</div>
    )
  }

  return (
    <div className="px-6 py-4">
      <div className="flex justify-between text-xs text-gray-400 mb-1">
        <span>Q{quarter} — {displayTime}</span>
        <span>{tickToMMSS(maxTick)}</span>
      </div>

      {/* Slider + quarter markers */}
      <div className="relative">
        <input
          type="range"
          min={0}
          max={maxTick}
          value={tick}
          onChange={handleChange}
          className="w-full h-2 bg-gray-700 rounded appearance-none cursor-pointer accent-blue-500"
        />
        {/* Quarter marker lines */}
        {QUARTER_TICKS.filter((qt) => qt <= maxTick).map((qt, i) => {
          const pct = (qt / maxTick) * 100
          return (
            <div
              key={qt}
              className="absolute top-0 flex flex-col items-center"
              style={{ left: `${pct}%`, transform: 'translateX(-50%)' }}
            >
              <div className="w-px h-3 bg-gray-500 mt-0.5" />
              <span className="text-xs text-gray-500 mt-0.5">{QUARTER_LABELS[i]}</span>
            </div>
          )
        })}
        {/* OT marker */}
        {maxTick > 3600 && (
          <div
            className="absolute top-0 flex flex-col items-center"
            style={{ left: `${(4500 / maxTick) * 100}%`, transform: 'translateX(-50%)' }}
          >
            <div className="w-px h-3 bg-yellow-600 mt-0.5" />
            <span className="text-xs text-yellow-600 mt-0.5">OT</span>
          </div>
        )}
      </div>
    </div>
  )
}
