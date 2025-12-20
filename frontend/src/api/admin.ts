import request from './request'
import type { UserInfo } from './auth'

export interface UsersResponse {
  users: UserInfo[]
  total: number
  page: number
  page_size: number
}

export interface SetAdminRequest {
  is_admin: boolean
}

export interface UpdateUserStatsRequest {
  streak?: number
  max_streak?: number
  total_checkin?: number
  title?: string
}

export function getUsers(page = 1, pageSize = 20): Promise<UsersResponse> {
  return request.get('/admin/users', { params: { page, page_size: pageSize } })
}

export function deleteUser(userId: number): Promise<void> {
  return request.delete(`/admin/users/${userId}`)
}

export function setUserAdmin(userId: number, isAdmin: boolean): Promise<void> {
  return request.put(`/admin/users/${userId}/admin`, { is_admin: isAdmin })
}

export function updateUserStats(userId: number, data: UpdateUserStatsRequest): Promise<UserInfo> {
  return request.put(`/admin/users/${userId}/stats`, data)
}
