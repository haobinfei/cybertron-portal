import request from './request'

export interface LoginParams {
  username: string
  password: string
}

export interface LoginResult {
  token: string
  user_info: {
    id: number
    username: string
    nickname: string
    email: string
    avatar: string
    role: string
    status: number
    last_login_at: string
    created_at: string
    updated_at: string
  }
}

export function login(params: LoginParams) {
  return request.post<any, LoginResult>('/auth/login', params)
}

export function logout() {
  return request.post('/auth/logout')
}
