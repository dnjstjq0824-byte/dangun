import { defineStore } from 'pinia'

interface User {
  id: string
  email: string
  nickname: string
  locale: string
  country_code: string
  status: string
  profile?: {
    bio?: string
    public_display_name?: string
    manner_temperature?: number
  }
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as User | null,
    token: '' as string,
    currentNeighborhood: {
      id: 'b0000000-0000-0000-0000-000000000001',
      name: '청라동 (Cheongna-dong)',
      isVerified: true
    }
  }),
  getters: {
    isAuthenticated: (state) => !!state.token
  },
  actions: {
    initAuth() {
      if (process.client) {
        const savedToken = localStorage.getItem('bn_near_token')
        const savedUser = localStorage.getItem('bn_near_user')
        if (savedToken && savedUser) {
          try {
            this.token = savedToken
            this.user = JSON.parse(savedUser)
          } catch (e) {
            this.logout()
          }
        }
      }
    },
    setAuth(user: User, token: string) {
      this.user = user
      this.token = token
      if (process.client) {
        localStorage.setItem('bn_near_token', token)
        localStorage.setItem('bn_near_user', JSON.stringify(user))
      }
    },
    logout() {
      this.user = null
      this.token = ''
      if (process.client) {
        localStorage.removeItem('bn_near_token')
        localStorage.removeItem('bn_near_user')
      }
    }
  }
})
