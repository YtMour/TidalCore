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

export function getUsers(page = 1, pageSize = 20): Promise<UsersResponse> {
  return request.get('/admin/users', { params: { page, page_size: pageSize } })
}

export function deleteUser(userId: number): Promise<void> {
  return request.delete(`/admin/users/${userId}`)
}

export function setUserAdmin(userId: number, isAdmin: boolean): Promise<void> {
  return request.put(`/admin/users/${userId}/admin`, { is_admin: isAdmin })
}
