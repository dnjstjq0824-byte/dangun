<template>
  <div>
    <h2 style="font-size: 22px; font-weight: 800; margin-bottom: 20px;">신고 및 콘텐츠 트리아지</h2>

    <table class="data-table">
      <thead>
        <tr>
          <th>신고 대상</th>
          <th>대상 ID</th>
          <th>사유 코드</th>
          <th>상세 설명</th>
          <th>우선순위</th>
          <th>상태</th>
          <th>조치 결정</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="r in reports" :key="r.id">
          <td><span class="badge-condition">{{ r.target_type }}</span></td>
          <td style="font-family: monospace; font-size: 11px;">{{ r.target_id.slice(0, 10) }}...</td>
          <td style="font-weight: 700;">{{ r.reason_code }}</td>
          <td>{{ r.description || '-' }}</td>
          <td>
            <span
              :style="{
                fontWeight: '700',
                fontSize: '11px',
                color: r.priority === 'URGENT' ? '#DC2626' : '#D97706'
              }"
            >
              {{ r.priority }}
            </span>
          </td>
          <td>{{ r.status }}</td>
          <td>
            <div style="display: flex; gap: 6px;" v-if="r.status === 'PENDING'">
              <button @click="resolve(r.id, 'CONTENT_REMOVED')" class="btn-admin btn-admin-danger">
                삭제조치
              </button>
              <button @click="resolve(r.id, 'NO_VIOLATION')" class="btn-admin" style="background: #E2E8F0; color: #475569;">
                기각
              </button>
            </div>
            <span v-else style="font-size: 11px; color: #16A34A; font-weight: 700;">처리완료</span>
          </td>
        </tr>
        <tr v-if="reports.length === 0">
          <td colspan="7" style="text-align: center; color: #94A3B8; padding: 30px;">
            대기 중인 신고 내역이 없습니다.
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const reports = ref<any[]>([])

const loadReports = async () => {
  const token = localStorage.getItem('bn_admin_token')
  try {
    const res = await fetch('http://localhost:8080/api/v1/admin/reports', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.data) {
      reports.value = data.data
    }
  } catch (e) {
    console.error(e)
  }
}

const resolve = async (reportID: string, decision: string) => {
  const token = localStorage.getItem('bn_admin_token')
  await fetch(`http://localhost:8080/api/v1/admin/moderation/cases/${reportID}/decide`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({ decision, reason: '관리자 심사 완료' })
  })
  loadReports()
}

onMounted(() => {
  loadReports()
})
</script>
