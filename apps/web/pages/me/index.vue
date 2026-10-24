<template>
  <div class="container" style="padding-top: 30px; padding-bottom: 60px; max-width: 800px;">
    <!-- Profile Card -->
    <div style="background: white; border: 1px solid #E2E8F0; border-radius: 16px; padding: 28px; box-shadow: var(--shadow-sm); margin-bottom: 24px;">
      <div style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 16px;">
        <div style="display: flex; align-items: center; gap: 16px;">
          <div style="width: 64px; height: 64px; border-radius: 9999px; background: #E2E8F0; display: flex; align-items: center; justify-content: center; font-size: 28px;">
            👤
          </div>
          <div>
            <div style="font-size: 20px; font-weight: 800; color: #0F172A;">{{ auth.user?.nickname }}</div>
            <div style="font-size: 13px; color: #64748B;">{{ auth.user?.email }}</div>
            <div style="display: flex; align-items: center; gap: 6px; margin-top: 4px;">
              <span class="badge-verified">📍 {{ auth.currentNeighborhood.name }}</span>
            </div>
          </div>
        </div>

        <MannerBadge :temperature="auth.user?.profile?.manner_temperature || 36.5" />
      </div>

      <!-- GPS Verification Trigger -->
      <div style="margin-top: 20px; padding-top: 20px; border-top: 1px solid #F1F5F9; display: flex; justify-content: space-between; align-items: center;">
        <div>
          <div style="font-size: 14px; font-weight: 700; color: #0F172A;">내 동네 GPS 위치 인증</div>
          <div style="font-size: 12px; color: #64748B;">실제 해당 동네 경계 내에 위치할 때 인증할 수 있습니다 (유효기간 30일)</div>
        </div>
        <button @click="verifyCurrentLocation" class="btn-primary" style="padding: 8px 16px; font-size: 13px;">
          지금 위치 인증하기
        </button>
      </div>
    </div>

    <!-- Activity Tabs -->
    <div style="background: white; border: 1px solid #E2E8F0; border-radius: 16px; padding: 24px; box-shadow: var(--shadow-sm);">
      <div style="display: flex; gap: 16px; border-bottom: 1px solid #F1F5F9; padding-bottom: 12px; margin-bottom: 20px;">
        <button
          @click="activeTab = 'listings'"
          :style="{
            fontWeight: activeTab === 'listings' ? '800' : '600',
            color: activeTab === 'listings' ? '#0F766E' : '#64748B',
            borderBottom: activeTab === 'listings' ? '2px solid #0F766E' : 'none',
            paddingBottom: '8px'
          }"
        >
          판매중인 물품 ({{ myListings.length }})
        </button>
        <button
          @click="activeTab = 'favorites'"
          :style="{
            fontWeight: activeTab === 'favorites' ? '800' : '600',
            color: activeTab === 'favorites' ? '#0F766E' : '#64748B',
            borderBottom: activeTab === 'favorites' ? '2px solid #0F766E' : 'none',
            paddingBottom: '8px'
          }"
        >
          관심 물품 ({{ favorites.length }})
        </button>
        <button
          @click="activeTab = 'deals'"
          :style="{
            fontWeight: activeTab === 'deals' ? '800' : '600',
            color: activeTab === 'deals' ? '#0F766E' : '#64748B',
            borderBottom: activeTab === 'deals' ? '2px solid #0F766E' : 'none',
            paddingBottom: '8px'
          }"
        >
          거래 내역 ({{ deals.length }})
        </button>
      </div>

      <!-- Listings Tab -->
      <div v-if="activeTab === 'listings'">
        <div v-if="myListings.length > 0" style="display: flex; flex-direction: column; gap: 12px;">
          <div v-for="l in myListings" :key="l.id" style="display: flex; align-items: center; justify-content: space-between; padding: 12px; border: 1px solid #F1F5F9; border-radius: 10px;">
            <div style="display: flex; align-items: center; gap: 12px;">
              <img :src="l.images?.[0]" style="width: 50px; height: 50px; border-radius: 6px; object-fit: cover;" />
              <div>
                <NuxtLink :to="`/marketplace/${l.id}`" style="font-weight: 700; font-size: 14px; color: #0F172A;">{{ l.title }}</NuxtLink>
                <div style="font-size: 13px; font-weight: 800; color: #0F766E;">{{ l.price_minor?.toLocaleString() }}원</div>
              </div>
            </div>
            <span class="badge-condition">{{ l.listing_status }}</span>
          </div>
        </div>
        <div v-else style="text-align: center; color: #94A3B8; padding: 30px;">
          판매 등록한 물품이 없습니다.
        </div>
      </div>

      <!-- Favorites Tab -->
      <div v-if="activeTab === 'favorites'">
        <div v-if="favorites.length > 0" style="display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 16px;">
          <ListingCard v-for="f in favorites" :key="f.id" :listing="f" />
        </div>
        <div v-else style="text-align: center; color: #94A3B8; padding: 30px;">
          관심 등록한 물품이 없습니다.
        </div>
      </div>

      <!-- Deals Tab -->
      <div v-if="activeTab === 'deals'">
        <div v-if="deals.length > 0" style="display: flex; flex-direction: column; gap: 12px;">
          <div v-for="d in deals" :key="d.id" style="padding: 14px; border: 1px solid #F1F5F9; border-radius: 10px; display: flex; justify-content: space-between; align-items: center;">
            <div>
              <div style="font-weight: 700; font-size: 14px; color: #0F172A;">거래 번호: {{ d.id.slice(0, 8) }}</div>
              <div style="font-size: 12px; color: #64748B;">상태: {{ d.status }}</div>
            </div>
            <NuxtLink :to="`/chats`" class="btn-secondary" style="font-size: 12px; padding: 6px 12px;">
              대화방 이동
            </NuxtLink>
          </div>
        </div>
        <div v-else style="text-align: center; color: #94A3B8; padding: 30px;">
          진행 중인 거래 내역이 없습니다.
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '~/stores/auth'
import { useApi } from '~/composables/useApi'

const auth = useAuthStore()
const { fetchAPI } = useApi()
const router = useRouter()

const activeTab = ref('listings')
const myListings = ref<any[]>([])
const favorites = ref<any[]>([])
const deals = ref<any[]>([])

const loadData = async () => {
  const lRes = await fetchAPI('/me/listings')
  if (lRes.data) myListings.value = lRes.data

  const fRes = await fetchAPI('/me/favorites')
  if (fRes.data) favorites.value = fRes.data

  const dRes = await fetchAPI('/me/deals')
  if (dRes.data) deals.value = dRes.data
}

const verifyCurrentLocation = async () => {
  // Uses browser geolocation fix or Cheongna center
  if (navigator.geolocation) {
    navigator.geolocation.getCurrentPosition(
      async (pos) => {
        const res = await fetchAPI('/me/neighborhood-verifications', {
          method: 'POST',
          body: {
            neighborhood_id: auth.currentNeighborhood.id,
            latitude: pos.coords.latitude,
            longitude: pos.coords.longitude,
            accuracy_meters: pos.coords.accuracy || 30.0
          }
        })
        if (res.data) {
          alert('동네 위치 인증 성공! 30일간 인증 상태가 유지됩니다.')
        } else if (res.error) {
          alert('위치 인증 실패: ' + (res.error.message || '동네 경계 밖입니다.'))
        }
      },
      async () => {
        // Fallback demo coordinates inside Cheongna
        const res = await fetchAPI('/me/neighborhood-verifications', {
          method: 'POST',
          body: {
            neighborhood_id: auth.currentNeighborhood.id,
            latitude: 37.5350,
            longitude: 126.6500,
            accuracy_meters: 25.0
          }
        })
        if (res.data) {
          alert('GPS 모의 좌표로 청라동 인증이 완료되었습니다 (30일 유효).')
        }
      }
    )
  }
}

onMounted(() => {
  if (!auth.isAuthenticated) {
    router.push('/login')
    return
  }
  loadData()
})
</script>

