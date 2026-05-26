// src/components/WinProbChart.tsx
// PDX-26: Recharts LineChart of win probability over game timeline.
// Shares ['timeline', gameId] query key with TimelineScrubber — zero extra fetches.
// Filters out plays with null win_prob to avoid line discontinuities.

import { useQuery } from '@tanstack/react-query'
import {
  ResponsiveContainer,
  LineChart,
  Line,
  XAxis,
  YAxis,
  ReferenceLine,
  Tooltip,
} from 'recharts'
import { getTimeline } from '../api/client'

interface WinProbChartProps {
  gameId: string
  homeTeam: string
  awayTeam: string
  currentTick: number
}

function tickToMMSS(tick: number): string {
  const m = Math.floor(tick / 60)
  const s = tick % 60
  return `${m}:${String(s).padStart(2, '0')}`
}

const QUARTER_TICKS = [900, 1800, 2700, 3600]
const QUARTER_LABELS = ['Q1', 'Q2', 'Q3', 'Q4']

export function WinProbChart({ gameId, homeTeam, currentTick }: WinProbChartProps) {
  const query = useQuery({
    queryKey: ['timeline', gameId],
    queryFn: () => getTimeline(gameId),
  })

  if (query.isLoading) {
    return (
      <div className="bg-gray-800 rounded-lg p-4 animate-pulse">
        <div className="h-48 bg-gray-700 rounded" />
      </div>
    )
  }

  if (query.isError || !query.data) {
    return null
  }

  const data = query.data.plays
    .filter((p) => p.win_prob !== null)
    .map((p) => ({ tick: p.tick, wp: p.win_prob as number }))

  const maxTick = query.data.max_tick

  return (
    <div className="bg-gray-800 rounded-lg p-4">
      <div className="text-xs text-gray-400 mb-2">{homeTeam} Win %</div>
      <ResponsiveContainer width="100%" height={200}>
        <LineChart data={data} margin={{ top: 4, right: 16, left: 0, bottom: 4 }}>
          <XAxis
            dataKey="tick"
            domain={[0, maxTick]}
            tickFormatter={tickToMMSS}
            tick={{ fill: '#9ca3af', fontSize: 10 }}
            tickLine={false}
            axisLine={{ stroke: '#374151' }}
          />
          <YAxis
            domain={[0, 1]}
            tickFormatter={(v: number) => `${Math.round(v * 100)}%`}
            tick={{ fill: '#9ca3af', fontSize: 10 }}
            tickLine={false}
            axisLine={false}
            width={40}
          />
          <Tooltip
            formatter={(value: unknown) => [`${((value as number) * 100).toFixed(1)}%`, `${homeTeam} Win`]}
            labelFormatter={(label: unknown) => tickToMMSS(label as number)}
            contentStyle={{ backgroundColor: '#1f2937', border: '1px solid #374151', borderRadius: 6 }}
            labelStyle={{ color: '#d1d5db' }}
            itemStyle={{ color: '#93c5fd' }}
          />
          {/* 50% baseline */}
          <ReferenceLine y={0.5} stroke="#6b7280" strokeDasharray="3 3" />
          {/* Quarter boundary lines */}
          {QUARTER_TICKS.filter((qt) => qt <= maxTick).map((qt, i) => (
            <ReferenceLine
              key={qt}
              x={qt}
              stroke="#374151"
              strokeDasharray="4 2"
              label={{ value: QUARTER_LABELS[i], fill: '#6b7280', fontSize: 10, position: 'top' }}
            />
          ))}
          {/* Live cursor */}
          <ReferenceLine x={currentTick} stroke="#ef4444" strokeDasharray="4 2" strokeWidth={1.5} />
          <Line
            type="monotone"
            dataKey="wp"
            stroke="#3b82f6"
            strokeWidth={2}
            dot={false}
            isAnimationActive={false}
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  )
}
