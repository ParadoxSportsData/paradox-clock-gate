// src/api/client.ts
// Single API module. All fetch calls go through here.
// MOCK_MODE=true: returns fixture data from mock.ts — no backend required.
// MOCK_MODE=false: fetches from BASE_URL, validates with Zod schemas at boundary.

import { GameSummarySchema, GameTimelineResponseSchema, GameStateResponseSchema } from './schemas'
import type { GameSummary, GameTimelineResponse, GameStateResponse } from './schemas'
import { MOCK_GAMES, MOCK_TIMELINE } from './mock'

export const MOCK_MODE = import.meta.env.VITE_MOCK_MODE === 'true'
export const BASE_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'

// Binary search: find the largest play.tick <= targetTick.
// Returns the play at that index, or the first play if none found before targetTick.
function findNearestPlay(plays: GameTimelineResponse['plays'], targetTick: number) {
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

export async function listGames(): Promise<GameSummary[]> {
  if (MOCK_MODE) {
    return Promise.resolve(MOCK_GAMES)
  }
  const res = await fetch(`${BASE_URL}/games`)
  if (!res.ok) {
    throw new Error(`listGames: ${res.status} ${res.statusText}`)
  }
  const data = await res.json()
  return GameSummarySchema.array().parse(data)
}

export async function getTimeline(gameId: string): Promise<GameTimelineResponse> {
  if (MOCK_MODE) {
    return Promise.resolve(MOCK_TIMELINE)
  }
  const res = await fetch(`${BASE_URL}/games/${encodeURIComponent(gameId)}/timeline`)
  if (!res.ok) {
    throw new Error(`getTimeline(${gameId}): ${res.status} ${res.statusText}`)
  }
  const data = await res.json()
  return GameTimelineResponseSchema.parse(data)
}

export async function getState(gameId: string, tick: number): Promise<GameStateResponse> {
  if (MOCK_MODE) {
    const play = findNearestPlay(MOCK_TIMELINE.plays, tick)
    // Map PlaySnapshot to GameStateResponse shape (add defteam + has_state)
    return Promise.resolve({
      ...play,
      defteam: null,
      has_state: true,
    })
  }
  const res = await fetch(
    `${BASE_URL}/games/${encodeURIComponent(gameId)}/state?tick=${tick}`
  )
  if (!res.ok) {
    throw new Error(`getState(${gameId}, ${tick}): ${res.status} ${res.statusText}`)
  }
  const data = await res.json()
  return GameStateResponseSchema.parse(data)
}
