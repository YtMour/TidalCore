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
  is_admin: boolean
}

export interface UpdateProfileRequest {
  display_name: string
}

export interface UpdateUsernameRequest {
  username: string
}

export interface UpdatePasswordRequest {
  old_password: string
  new_password: string
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

export function updateProfile(data: UpdateProfileRequest): Promise<UserInfo> {
  return request.put('/user/profile', data)
}

export function updateUsername(data: UpdateUsernameRequest): Promise<UserInfo> {
  return request.put('/user/username', data)
}

export function updatePassword(data: UpdatePasswordRequest): Promise<void> {
  return request.put('/user/password', data)
}
