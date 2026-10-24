<template>
  <div class="container" style="padding-top: 24px; padding-bottom: 60px;">
    <div style="margin-bottom: 24px;">
      <h1 style="font-size: 26px; font-weight: 800; color: #0F172A;">동네가게</h1>
      <p style="font-size: 14px; color: #64748B;">{{ auth.currentNeighborhood.name }}의 신뢰할 수 있는 지역 소상공인과 상점</p>
    </div>

    <!-- Businesses Grid -->
    <div style="display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 20px;">
      <div
        v-for="biz in businesses"
        :key="biz.id"
        class="card"
        style="padding: 24px; display: flex; flex-direction: column;"
      >
        <div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 8px;">
          <h3 style="font-size: 17px; font-weight: 800; color: #0F172A;">{{ biz.name }}</h3>
          <span class="badge-verified" v-if="biz.verification_status === 'VERIFIED'">인증완료</span>
        </div>

        <div style="font-size: 12px; font-weight: 700; color: #0F766E; margin-bottom: 10px;">
          {{ biz.category_name || '동네서비스' }}
        </div>

        <p style="font-size: 14px; color: #475569; line-height: 1.6; margin-bottom: 16px; flex: 1;">
          {{ biz.description }}
        </p>

        <div style="font-size: 12px; color: #64748B; border-top: 1px solid #F1F5F9; padding-top: 12px; display: flex; flex-direction: column; gap: 4px;">
          <div>📍 {{ biz.public_address }}</div>
          <div v-if="biz.phone_number">📞 {{ biz.phone_number }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '~/stores/auth'
import { useApi } from '~/composables/useApi'

const auth = useAuthStore()
const { fetchAPI } = useApi()
const businesses = ref<any[]>([])

const loadBusinesses = async () => {
  const res = await fetchAPI(`/businesses?neighborhood_id=${auth.currentNeighborhood.id}`)
  if (res.data) {
    businesses.value = res.data
  }
}

onMounted(() => {
  loadBusinesses()
})
</script>

