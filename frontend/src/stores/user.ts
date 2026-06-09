import { defineStore } from 'pinia'
import { login as loginApi, logout as logoutApi, type LoginResult } from '@/api/auth'
import { getCurrentUser, type UserInfo } from '@/api/user'

interface UserState {
  token: string
  userInfo: UserInfo | null
}

function loadToken(): string {
  return localStorage.getItem('token') || ''
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    token: loadToken(),
    userInfo: null,
  }),

  getters: {
    isLoggedIn: (state) => !!state.token,
    username: (state) => state.userInfo?.username ?? '',
    isAdmin: (state) => state.userInfo?.role === 'admin',
  },

  actions: {
    async login(username: string, password: string) {
      const result: LoginResult = await loginApi({ username, password })
      this.token = result.token
      this.userInfo = result.user_info
      localStorage.setItem('token', result.token)
    },

    async fetchUserInfo() {
      try {
        const user = await getCurrentUser()
        this.userInfo = user
      } catch {
        this.logout()
      }
    },

    async logout() {
      try {
        await logoutApi()
      } catch {
        // ignore
      }
      this.token = ''
      this.userInfo = null
      localStorage.removeItem('token')
    },
  },
})
