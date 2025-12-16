import request from './request'

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  password: string
}

export interface AuthResponse {
  token: string
  user: {
    id: number
    username: string
    streak: number
    max_streak: number
    total_checkin: number
  }
}

export function login(data: LoginRequest): Promise<AuthResponse> {
  return request.post('/auth/login', data)
}

export function register(data: RegisterRequest): Promise<AuthResponse> {
  return request.post('/auth/register', data)
}

export function getProfile() {
  return request.get('/user/profile')
}
