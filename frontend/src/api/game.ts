import request from './request'
import type { ApiResponse } from './types'

export interface GameScoreUpdate {
  tool_id: number
  player_id: string
  score: number
}

export interface GameScoreResponse {
  best_score: number
}

export interface LeaderboardEntry {
  rank: number
  player: string
  best_score: number
}

export const gameApi = {
  // 更新分数
  updateScore(data: GameScoreUpdate): Promise<ApiResponse<GameScoreResponse>> {
    return request.post('/games/score', data)
  },

  // 获取最高分
  getBestScore(toolId: number, playerId: string): Promise<ApiResponse<GameScoreResponse>> {
    return request.get('/games/score', {
      params: { tool_id: toolId, player_id: playerId }
    })
  },

  // 获取排行榜
  getLeaderboard(toolId: number, limit = 10): Promise<ApiResponse<LeaderboardEntry[]>> {
    return request.get('/games/leaderboard', {
      params: { tool_id: toolId, limit }
    })
  }
}
