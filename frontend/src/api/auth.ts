import request from './request'

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  display_name: string
  password: string
}

export interface UserInfo {
  id: number
  username: string
  display_name: string
  streak: number
  max_streak: number
  total_checkin: number
}

export interface AuthResponse {
  token: string
  user: UserInfo
}

export function login(data: LoginRequest): Promise<AuthResponse> {
  return request.post('/auth/login', data)
}

export function register(data: RegisterRequest): Promise<AuthResponse> {
  return request.post('/auth/register', data)
}

export function getProfile(): Promise<UserInfo> {
  return request.get('/user/profile')
}
