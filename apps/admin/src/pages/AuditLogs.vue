<template>
  <div>
    <h2 style="font-size: 22px; font-weight: 800; margin-bottom: 20px;">불변 감사 로그 (Tamper-Proof Audit Logs)</h2>
    <p style="font-size: 13px; color: #64748B; margin-bottom: 20px;">
      플랫폼 내의 모든 관리자 작업, 권한 변경, 상태 업데이트, 개인정보 및 증거 열람 내역은 삭제 불가능한 추가 전용(Append-Only) 감사 로그로 영구 기록됩니다.
    </p>

    <table class="data-table">
      <thead>
        <tr>
          <th>타임스탬프 (UTC)</th>
          <th>수행자 (Actor)</th>
          <th>작업 (Action)</th>
          <th>대상 리소스</th>
          <th>결과</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="log in logs" :key="log.id">
          <td style="font-family: monospace; font-size: 11px;">{{ log.created_at }}</td>
          <td><span style="font-weight: 700;">{{ log.actor_type }}</span> ({{ log.actor_id.slice(0, 8) }}...)</td>
          <td style="font-weight: 700; color: #0F766E;">{{ log.action }}</td>
          <td>{{ log.resource_type }} : {{ log.resource_id.slice(0, 10) }}...</td>
          <td>
            <span style="font-size: 11px; font-weight: 700; background: #ECFDF5; color: #16A34A; padding: 2px 6px; border-radius: 4px;">
              {{ log.result }}
            </span>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const logs = ref<any[]>([])

const loadLogs = async () => {
  const token = localStorage.getItem('bn_admin_token')
  try {
    const res = await fetch('http://localhost:8080/api/v1/admin/audit-logs', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.data) {
      logs.value = data.data
    }
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadLogs()
})
</script>
