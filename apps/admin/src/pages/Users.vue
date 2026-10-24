<template>
  <div>
    <h2 style="font-size: 22px; font-weight: 800; margin-bottom: 20px;">사용자 및 계정 거버넌스</h2>

    <table class="data-table">
      <thead>
        <tr>
          <th>사용자 ID</th>
          <th>닉네임</th>
          <th>이메일</th>
          <th>매너온도</th>
          <th>계정상태</th>
          <th>작업 관리</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="u in users" :key="u.id">
          <td style="font-family: monospace; font-size: 11px;">{{ u.id.slice(0, 13) }}...</td>
          <td style="font-weight: 700;">{{ u.nickname }}</td>
          <td>{{ u.email }}</td>
          <td style="font-weight: 800; color: #0F766E;">{{ u.profile?.manner_temperature || 36.5 }}°C</td>
          <td>
            <span
              :style="{
                padding: '2px 8px',
                borderRadius: '4px',
                fontSize: '11px',
                fontWeight: '700',
                backgroundColor: u.status === 'ACTIVE' ? '#ECFDF5' : '#FEF2F2',
                color: u.status === 'ACTIVE' ? '#16A34A' : '#DC2626'
              }"
            >
              {{ u.status }}
            </span>
          </td>
          <td>
            <div style="display: flex; gap: 6px;">
              <button
                v-if="u.status === 'ACTIVE'"
                @click="updateStatus(u.id, 'RESTRICTED')"
                class="btn-admin"
                style="background: #F59E0B; color: white;"
              >
                이용제한
              </button>
              <button
                v-if="u.status === 'ACTIVE'"
                @click="updateStatus(u.id, 'SUSPENDED')"
                class="btn-admin btn-admin-danger"
              >
                영구정지
              </button>
              <button
                v-if="u.status !== 'ACTIVE'"
                @click="updateStatus(u.id, 'ACTIVE')"
                class="btn-admin btn-admin-primary"
              >
                정지해제
              </button>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const users = ref<any[]>([])

const loadUsers = async () => {
  const token = localStorage.getItem('bn_admin_token')
  try {
    const res = await fetch('http://localhost:8080/api/v1/admin/users', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.data) {
      users.value = data.data
    }
  } catch (e) {
    console.error(e)
  }
}

const updateStatus = async (userID: string, status: string) => {
  const token = localStorage.getItem('bn_admin_token')
  await fetch(`http://localhost:8080/api/v1/admin/users/${userID}/status`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({ status })
  })
  loadUsers()
}

onMounted(() => {
  loadUsers()
})
</script>
