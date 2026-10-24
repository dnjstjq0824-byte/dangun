<template>
  <div>
    <h2 style="font-size: 22px; font-weight: 800; margin-bottom: 20px;">중고거래 매물 검수 및 관리</h2>

    <table class="data-table">
      <thead>
        <tr>
          <th>매물 제목</th>
          <th>카테고리</th>
          <th>동네</th>
          <th>가격</th>
          <th>상태</th>
          <th>게시상태</th>
          <th>조치</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="l in listings" :key="l.id">
          <td style="font-weight: 700; max-width: 200px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">
            {{ l.title }}
          </td>
          <td>{{ l.category_name || '디지털' }}</td>
          <td>{{ l.neighborhood_name }}</td>
          <td style="font-weight: 800; color: #0F766E;">{{ l.price_minor?.toLocaleString() }}원</td>
          <td><span class="badge-condition">{{ l.condition }}</span></td>
          <td>
            <span
              :style="{
                padding: '2px 8px',
                borderRadius: '4px',
                fontSize: '11px',
                fontWeight: '700',
                backgroundColor: l.listing_status === 'ACTIVE' ? '#ECFDF5' : '#FEF2F2',
                color: l.listing_status === 'ACTIVE' ? '#16A34A' : '#DC2626'
              }"
            >
              {{ l.listing_status }}
            </span>
          </td>
          <td>
            <button
              v-if="l.listing_status !== 'REMOVED'"
              @click="removeListing(l.id)"
              class="btn-admin btn-admin-danger"
            >
              게시중단
            </button>
            <span v-else style="font-size: 11px; color: #94A3B8;">조치완료</span>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const listings = ref<any[]>([])

const loadListings = async () => {
  const token = localStorage.getItem('bn_admin_token')
  try {
    const res = await fetch('http://localhost:8080/api/v1/listings', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.data) {
      listings.value = data.data
    }
  } catch (e) {
    console.error(e)
  }
}

const removeListing = async (id: string) => {
  if (!confirm('해당 매물을 게시 중단 조치하시겠습니까?')) return
  const token = localStorage.getItem('bn_admin_token')
  await fetch(`http://localhost:8080/api/v1/listings/${id}`, {
    method: 'DELETE',
    headers: { 'Authorization': `Bearer ${token}` }
  })
  loadListings()
}

onMounted(() => {
  loadListings()
})
</script>
