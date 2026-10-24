import { defineStore } from 'pinia'

export const useAdminAuthStore = defineStore('adminAuth', {
  state: () => ({
    token: (typeof localStorage !== 'undefined' ? localStorage.getItem('bn_admin_token') : '') || '',
    user: (typeof localStorage !== 'undefined' && localStorage.getItem('bn_admin_user') ? JSON.parse(localStorage.getItem('bn_admin_user') || '{}') : null) as any
  }),
  getters: {
    isAuthenticated: (state) => !!state.token
  },
  actions: {
    setAuth(user: any, token: string) {
      this.user = user
      this.token = token
      localStorage.setItem('bn_admin_token', token)
      localStorage.setItem('bn_admin_user', JSON.stringify(user))
    },
    logout() {
      this.token = ''
      this.user = null
      localStorage.removeItem('bn_admin_token')
      localStorage.removeItem('bn_admin_user')
    }
  }
})
