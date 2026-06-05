import { defineStore } from 'pinia'

interface UserInfo {
  id: number
  username: string
  nickname: string
  avatar: string
  roles: string[]
}

interface UserState {
  token: string
  userInfo: UserInfo | null
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    token: '',
    userInfo: null,
  }),

  getters: {
    isLoggedIn: (state) => !!state.token,
    username: (state) => state.userInfo?.username ?? '',
  },

  actions: {
    setToken(token: string) {
      this.token = token
    },
    setUserInfo(info: UserInfo) {
      this.userInfo = info
    },
    logout() {
      this.token = ''
      this.userInfo = null
    },
  },
})
