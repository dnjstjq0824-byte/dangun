// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2026-10-24',
  devtools: { enabled: false },
  modules: ['@pinia/nuxt'],
  app: {
    head: {
      title: 'BN NEAR — 우리 동네의 모든 연결',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: '가까운 이웃과 중고물품을 거래하고, 동네 소식을 나누며, 로컬 가게를 발견하는 글로벌 지역 생활 플랫폼 BN NEAR.' },
        { property: 'og:title', content: 'BN NEAR — Global Hyperlocal Marketplace & Community' },
        { property: 'og:description', content: 'Your neighborhood. Your marketplace.' },
        { property: 'og:type', content: 'website' }
      ],
      link: [
        { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' },
        { rel: 'stylesheet', href: 'https://fonts.googleapis.com/css2?family=Pretendard:wght@400;500;600;700&display=swap' }
      ]
    }
  },
  runtimeConfig: {
    public: {
      apiBase: process.env.API_BASE_URL || 'http://localhost:8080/api/v1',
      wsBase: process.env.WS_BASE_URL || 'ws://localhost:8080/api/v1/ws'
    }
  }
})

