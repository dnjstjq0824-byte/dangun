<template>
  <header class="navbar">
    <div class="container nav-container">
      <div class="logo-group">
        <NuxtLink to="/" style="display: flex; align-items: center; gap: 10px;">
          <span class="logo-badge">NEAR</span>
          <span class="logo-text">BN NEAR</span>
        </NuxtLink>

        <!-- Current Neighborhood Pill -->
        <div style="display: flex; align-items: center; gap: 6px; padding: 4px 10px; background: #F1F5F9; border-radius: 9999px; font-size: 13px; font-weight: 600; color: #0F766E;">
          <span>📍 {{ auth.currentNeighborhood.name }}</span>
          <span v-if="auth.currentNeighborhood.isVerified" style="font-size: 11px; background: #ECFDF5; color: #16A34A; padding: 1px 5px; border-radius: 4px;">인증됨</span>
        </div>
      </div>

      <!-- Center Search -->
      <div style="flex: 1; max-width: 420px; margin: 0 20px;">
        <input
          v-model="searchQuery"
          @keyup.enter="handleSearch"
          type="text"
          placeholder="청라동 근처에서 검색 (예: 맥북, 식탁, 유모차)"
          class="form-input"
          style="border-radius: 9999px; padding-left: 18px;"
        />
      </div>

      <!-- Right Navigation -->
      <nav class="nav-links">
        <NuxtLink to="/marketplace" class="nav-link">중고거래</NuxtLink>
        <NuxtLink to="/community" class="nav-link">동네생활</NuxtLink>
        <NuxtLink to="/businesses" class="nav-link">동네가게</NuxtLink>
        <NuxtLink to="/chats" class="nav-link">채팅</NuxtLink>

        <template v-if="auth.isAuthenticated">
          <NuxtLink to="/marketplace/create" class="btn-primary" style="padding: 6px 14px; font-size: 13px;">
            + 물품 등록
          </NuxtLink>
          <NuxtLink to="/me" class="nav-link" style="font-weight: 700; color: #0F766E;">
            {{ auth.user?.nickname }}님
          </NuxtLink>
          <button @click="auth.logout()" class="btn-secondary" style="padding: 4px 10px; font-size: 12px;">
            로그아웃
          </button>
        </template>
        <template v-else>
          <NuxtLink to="/login" class="btn-primary" style="padding: 6px 14px; font-size: 13px;">
            로그인 / 가입
          </NuxtLink>
        </template>
      </nav>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '~/stores/auth'

const auth = useAuthStore()
const router = useRouter()
const searchQuery = ref('')

const handleSearch = () => {
  if (searchQuery.value.trim()) {
    router.push({ path: '/marketplace', query: { q: searchQuery.value.trim() } })
  }
}
</script>
