<template>
  <div>
    <h2 style="font-size: 22px; font-weight: 800; margin-bottom: 20px;">동네가게 사업자 인증 관리</h2>

    <table class="data-table">
      <thead>
        <tr>
          <th>상호명</th>
          <th>업종 카테고리</th>
          <th>소재지 주소</th>
          <th>연락처</th>
          <th>인증 상태</th>
          <th>심사 승인</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="b in businesses" :key="b.id">
          <td style="font-weight: 700;">{{ b.name }}</td>
          <td>{{ b.category_name || '동네서비스' }}</td>
          <td>{{ b.public_address }}</td>
          <td>{{ b.phone_number || '-' }}</td>
          <td>
            <span
              :style="{
                padding: '2px 8px',
                borderRadius: '4px',
                fontSize: '11px',
                fontWeight: '700',
                backgroundColor: b.verification_status === 'VERIFIED' ? '#ECFDF5' : '#FFFBEB',
                color: b.verification_status === 'VERIFIED' ? '#16A34A' : '#D97706'
              }"
            >
              {{ b.verification_status }}
            </span>
          </td>
          <td>
            <div style="display: flex; gap: 6px;" v-if="b.verification_status !== 'VERIFIED'">
              <button @click="verify(b.id, true)" class="btn-admin btn-admin-primary">
                승인 (Badge)
              </button>
              <button @click="verify(b.id, false)" class="btn-admin btn-admin-danger">
                반려
              </button>
            </div>
            <span v-else style="font-size: 11px; color: #16A34A; font-weight: 700;">인증완료</span>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const businesses = ref<any[]>([])

const loadBusinesses = async () => {
  try {
    const res = await fetch('http://localhost:8080/api/v1/businesses?neighborhood_id=b0000000-0000-0000-0000-000000000001')
    const data = await res.json()
    if (data.data) {
      businesses.value = data.data
    }
  } catch (e) {
    console.error(e)
  }
}

const verify = async (id: string, approve: boolean) => {
  const token = localStorage.getItem('bn_admin_token')
  await fetch(`http://localhost:8080/api/v1/admin/businesses/${id}/verify`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({ approve })
  })
  loadBusinesses()
}

onMounted(() => {
  loadBusinesses()
})
</script>

