<template>
  <div class="container" style="padding-top: 24px; padding-bottom: 50px;">
    <!-- Breadcrumb & Header -->
    <div style="display: flex; justify-content: space-between; align-items: flex-end; margin-bottom: 24px; flex-wrap: wrap; gap: 16px;">
      <div>
        <h1 style="font-size: 26px; font-weight: 800; color: #0F172A;">중고거래 매물 탐색</h1>
        <p style="font-size: 14px; color: #64748B;">{{ auth.currentNeighborhood.name }} 및 인근 지역의 이웃 직거래 물품</p>
      </div>

      <NuxtLink v-if="auth.isAuthenticated" to="/marketplace/create" class="btn-primary">
        + 내 물건 판매하기
      </NuxtLink>
    </div>

    <!-- Filters Bar -->
    <div style="background: white; border: 1px solid #E2E8F0; border-radius: 12px; padding: 16px; margin-bottom: 24px; display: flex; flex-wrap: wrap; gap: 16px; align-items: center;">
      <!-- Radius Filter -->
      <div style="display: flex; align-items: center; gap: 8px;">
        <span style="font-size: 13px; font-weight: 600; color: #475569;">반경:</span>
        <select v-model="radius" @change="fetchListings" class="form-input" style="width: 110px; padding: 6px 10px;">
          <option value="1000">1 km</option>
          <option value="3000">3 km</option>
          <option value="5000">5 km</option>
          <option value="10000">10 km</option>
        </select>
      </div>

      <!-- Condition Filter -->
      <div style="display: flex; align-items: center; gap: 8px;">
        <span style="font-size: 13px; font-weight: 600; color: #475569;">상태:</span>
        <select v-model="condition" @change="fetchListings" class="form-input" style="width: 120px; padding: 6px 10px;">
          <option value="">전체 상태</option>
          <option value="NEW">미개봉/새상품</option>
          <option value="LIKE_NEW">거의 새것</option>
          <option value="GOOD">사용감 적음</option>
          <option value="FAIR">사용감 있음</option>
        </select>
      </div>

      <!-- Free Only Toggle -->
      <label style="display: flex; align-items: center; gap: 6px; font-size: 13px; font-weight: 600; cursor: pointer;">
        <input type="checkbox" v-model="isFree" @change="fetchListings" />
        무료나눔만 보기
      </label>

      <!-- Sort Filter -->
      <div style="margin-left: auto; display: flex; align-items: center; gap: 8px;">
        <span style="font-size: 13px; font-weight: 600; color: #475569;">정렬:</span>
        <select v-model="sortOrder" @change="fetchListings" class="form-input" style="width: 120px; padding: 6px 10px;">
          <option value="latest">최신순</option>
          <option value="distance">가까운순</option>
          <option value="price_asc">낮은가격순</option>
          <option value="price_desc">높은가격순</option>
        </select>
      </div>
    </div>

    <!-- Listings Grid -->
    <div v-if="listings.length > 0" style="display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)); gap: 20px;">
      <ListingCard
        v-for="listing in listings"
        :key="listing.id"
        :listing="listing"
      />
    </div>

    <!-- Empty State -->
    <div v-else style="text-align: center; padding: 80px 0; background: white; border-radius: 16px; border: 1px dashed #CBD5E1;">
      <div style="font-size: 40px; margin-bottom: 12px;">📦</div>
      <div style="font-size: 18px; font-weight: 700; color: #1E293B; margin-bottom: 6px;">등록된 물품이 없습니다</div>
      <p style="font-size: 14px; color: #64748B; margin-bottom: 16px;">검색 필터를 변경하거나 첫 번째 물품을 등록해보세요!</p>
      <NuxtLink v-if="auth.isAuthenticated" to="/marketplace/create" class="btn-primary">
        + 첫 물품 등록하기
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '~/stores/auth'
import { useApi } from '~/composables/useApi'

const auth = useAuthStore()
const { fetchAPI } = useApi()
const route = useRoute()

const listings = ref<any[]>([])
const radius = ref('10000')
const condition = ref('')
const isFree = ref(false)
const sortOrder = ref('latest')

const fetchListings = async () => {
  let url = `/listings?sort=${sortOrder.value}&radius_meters=${radius.value}`
  if (condition.value) url += `&condition=${condition.value}`
  if (isFree.value) url += `&is_free=true`
  if (route.query.category_id) url += `&category_id=${route.query.category_id}`
  if (route.query.q) url += `&q=${encodeURIComponent(String(route.query.q))}`

  const res = await fetchAPI(url)
  if (res.data) {
    listings.value = res.data
  }
}

watch(() => route.query, () => {
  fetchListings()
})

onMounted(() => {
  fetchListings()
})
</script>

