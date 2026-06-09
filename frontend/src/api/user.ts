import request from './request'

export interface UserInfo {
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

export interface CreateUserParams {
  username: string
  password: string
  nickname?: string
  email?: string
  role?: string
  status?: number
}

export interface UpdateUserParams {
  nickname?: string
  email?: string
  role?: string
  status?: number
  password?: string
}

export function getCurrentUser() {
  return request.get<any, UserInfo>('/user/me')
}

export function getUserList(params: { page: number; page_size: number }) {
  return request.get<any, { list: UserInfo[]; total: number; page: number; page_size: number }>('/users', { params })
}

export function createUser(params: CreateUserParams) {
  return request.post<any, UserInfo>('/users', params)
}

export function updateUser(id: number, params: UpdateUserParams) {
  return request.put<any, UserInfo>(`/users/${id}`, params)
}

export function deleteUser(id: number) {
  return request.delete(`/users/${id}`)
}
