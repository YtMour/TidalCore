import request from './request'

export interface CheckinRequest {
  duration: number
  cycles: number
}

export interface CheckinResponse {
  checkin: {
    id: number
    user_id: number
    duration: number
    cycles: number
    checked_at: string
  }
  current_streak: number
  max_streak: number
  total_checkin: number
}

export interface CheckinRecord {
  id: number
  user_id: number
  duration: number
  cycles: number
  checked_at: string
}

export function checkin(data: CheckinRequest): Promise<CheckinResponse> {
  return request.post('/checkin', data)
}

export function getHistory(limit = 30): Promise<CheckinRecord[]> {
  return request.get('/checkin/history', { params: { limit } })
}

export function getHeatmap(days = 365): Promise<CheckinRecord[]> {
  return request.get('/checkin/heatmap', { params: { days } })
}

export function getGlobalHeatmap(days = 365): Promise<Record<string, number>> {
  return request.get('/heatmap/global', { params: { days } })
}

export interface LeaderboardUser {
  id: number
  username: string
  display_name: string
  title: string
  streak: number
  max_streak: number
  total_checkin: number
}

export function getLeaderboard(limit = 20): Promise<LeaderboardUser[]> {
  return request.get('/leaderboard', { params: { limit } })
}

// 访问统计相关
export interface VisitStats {
  today_visits: number
  total_visits: number
}

export function recordVisit(): Promise<{ today_visits: number }> {
  return request.post('/visit')
}

export function getVisitStats(): Promise<VisitStats> {
  return request.get('/visit/stats')
}

export function getVisitHeatmap(days = 365): Promise<Record<string, number>> {
  return request.get('/visit/heatmap', { params: { days } })
}
