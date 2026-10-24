<template>
  <div class="admin-layout">
    <!-- Left Sidebar -->
    <aside class="admin-sidebar">
      <div class="sidebar-header">
        <span style="background: #0F766E; color: white; padding: 4px 8px; border-radius: 4px; font-weight: 800; font-size: 14px;">ADMIN</span>
        <span style="font-weight: 800; color: white; font-size: 16px;">BN NEAR OPS</span>
      </div>

      <nav class="sidebar-menu">
        <RouterLink to="/" class="menu-item" exact-active-class="active">
          <span>📊</span> 대시보드
        </RouterLink>
        <RouterLink to="/geography" class="menu-item" active-class="active">
          <span>🗺️</span> 지역 & 동네 경계
        </RouterLink>
        <RouterLink to="/users" class="menu-item" active-class="active">
          <span>👥</span> 사용자 관리
        </RouterLink>
        <RouterLink to="/listings" class="menu-item" active-class="active">
          <span>🛍️</span> 중고거래 매물 검수
        </RouterLink>
        <RouterLink to="/moderation" class="menu-item" active-class="active">
          <span>🚨</span> 신고 & 트리아지
        </RouterLink>
        <RouterLink to="/businesses" class="menu-item" active-class="active">
          <span>🏪</span> 동네가게 인증
        </RouterLink>
        <RouterLink to="/chat-safety" class="menu-item" active-class="active">
          <span>🛡️</span> 채팅 안전 증거조사
        </RouterLink>
        <RouterLink to="/audit-logs" class="menu-item" active-class="active">
          <span>📜</span> 불변 감사 로그
        </RouterLink>
      </nav>

      <div style="padding: 16px; border-top: 1px solid #1E293B; font-size: 12px; color: #64748B;">
        Booking Network Platform
      </div>
    </aside>

    <!-- Main Section -->
    <div class="admin-main">
      <header class="admin-topbar">
        <div style="display: flex; align-items: center; gap: 12px;">
          <span style="background: #ECFDF5; color: #16A34A; border: 1px solid #A7F3D0; font-size: 11px; font-weight: 700; padding: 2px 8px; border-radius: 4px;">
            STAGING / PRODUCTION
          </span>
          <span style="font-size: 13px; color: #64748B;">API: http://localhost:8080/api/v1</span>
        </div>

        <div style="display: flex; align-items: center; gap: 14px;">
          <button @click="quickAdminLogin" class="btn-admin btn-admin-primary">
            🛡️ 최고 관리자 인증 토큰 갱신
          </button>
        </div>
      </header>

      <main class="admin-content">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import './assets/admin.css'

const quickAdminLogin = async () => {
  try {
    const res = await fetch('http://localhost:8080/api/v1/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        email: 'admin@bn-near.com',
        password: 'Password123!'
      })
    })
    const data = await res.json()
    if (data.data?.access_token) {
      localStorage.setItem('bn_admin_token', data.data.access_token)
      alert('관리자 토큰이 갱신되었습니다.')
    }
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  if (!localStorage.getItem('bn_admin_token')) {
    quickAdminLogin()
  }
})
</script>

