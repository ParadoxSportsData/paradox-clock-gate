// src/api/schemas.ts
// Zod schemas mirroring every TypeScript interface in types.ts.
// Zod validates at the fetch boundary — if the Go server changes a response shape
// without a contract amendment ticket, parse() throws immediately.
// Do NOT add manual error catch blocks here — let Zod throw natively.

import { z } from 'zod'

export const GameSummarySchema = z.object({
  game_id: z.string(),
  home_team: z.string(),
  away_team: z.string(),
  home_score: z.number(),
  away_score: z.number(),
  duration: z.number(),
})

export const GameStateResponseSchema = z.object({
  tick: z.number(),
  quarter: z.number(),
  down: z.number().nullable(),
  yards_to_go: z.number().nullable(),
  yard_line: z.number().nullable(),
  home_score: z.number(),
  away_score: z.number(),
  posteam: z.string().nullable(),
  defteam: z.string().nullable(),
  win_prob: z.number().nullable(),
  play_type: z.string(),
  description: z.string(),
  has_state: z.boolean(),
})

export const PlaySnapshotSchema = z.object({
  tick: z.number(),
  quarter: z.number(),
  down: z.number().nullable(),
  yards_to_go: z.number().nullable(),
  yard_line: z.number().nullable(),
  home_score: z.number(),
  away_score: z.number(),
  posteam: z.string().nullable(),
  win_prob: z.number().nullable(),
  play_type: z.string(),
  description: z.string(),
})

export const GameTimelineResponseSchema = z.object({
  game_id: z.string(),
  home_team: z.string(),
  away_team: z.string(),
  max_tick: z.number(),
  plays: z.array(PlaySnapshotSchema),
})

export const ApiErrorSchema = z.object({
  error: z.string(),
  max_tick: z.number().optional(),
})

// Inferred types — should structurally match types.ts; compile error if not.
export type GameSummary = z.infer<typeof GameSummarySchema>
export type GameStateResponse = z.infer<typeof GameStateResponseSchema>
export type PlaySnapshot = z.infer<typeof PlaySnapshotSchema>
export type GameTimelineResponse = z.infer<typeof GameTimelineResponseSchema>
export type ApiError = z.infer<typeof ApiErrorSchema>
