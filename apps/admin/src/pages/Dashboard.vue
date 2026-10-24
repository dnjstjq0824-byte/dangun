<template>
  <div>
    <h2 style="font-size: 22px; font-weight: 800; margin-bottom: 20px;">운영 대시보드</h2>

    <!-- Stat Grid -->
    <div class="stat-grid" v-if="stats">
      <div class="stat-card">
        <div style="font-size: 12px; font-weight: 700; color: #64748B; margin-bottom: 6px;">전체 등록 사용자</div>
        <div style="font-size: 28px; font-weight: 900; color: #0F172A;">{{ stats.total_users }}명</div>
        <div style="font-size: 11px; color: #16A34A; margin-top: 4px;">활성: {{ stats.active_users }}명</div>
      </div>

      <div class="stat-card">
        <div style="font-size: 12px; font-weight: 700; color: #64748B; margin-bottom: 6px;">등록된 중고물품</div>
        <div style="font-size: 28px; font-weight: 900; color: #0F766E;">{{ stats.total_listings }}건</div>
        <div style="font-size: 11px; color: #16A34A; margin-top: 4px;">판매중: {{ stats.active_listings }}건</div>
      </div>

      <div class="stat-card">
        <div style="font-size: 12px; font-weight: 700; color: #64748B; margin-bottom: 6px;">완료된 직거래</div>
        <div style="font-size: 28px; font-weight: 900; color: #2563EB;">{{ stats.completed_deals }}건</div>
        <div style="font-size: 11px; color: #64748B; margin-top: 4px;">실제 양측 확인 완료</div>
      </div>

      <div class="stat-card">
        <div style="font-size: 12px; font-weight: 700; color: #64748B; margin-bottom: 6px;">실시간 웹소켓 세션</div>
        <div style="font-size: 28px; font-weight: 900; color: #D97706;">{{ stats.live_ws_sessions }}개</div>
        <div style="font-size: 11px; color: #16A34A; margin-top: 4px;">동시 접속 유지 중</div>
      </div>

      <div class="stat-card">
        <div style="font-size: 12px; font-weight: 700; color: #64748B; margin-bottom: 6px;">미처리 신고 건수</div>
        <div style="font-size: 28px; font-weight: 900; color: #DC2626;">{{ stats.pending_reports }}건</div>
        <div style="font-size: 11px; color: #DC2626; margin-top: 4px;">트리아지 대기</div>
      </div>
    </div>

    <!-- System Telemetry Box -->
    <div style="background: white; border: 1px solid #E2E8F0; border-radius: 12px; padding: 24px; margin-top: 24px;">
      <h3 style="font-size: 16px; font-weight: 800; margin-bottom: 12px;">플랫폼 핵심 서브시스템 상태</h3>
      <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 16px;">
        <div style="padding: 12px; border: 1px solid #F1F5F9; border-radius: 8px;">
          <div style="font-size: 12px; color: #64748B;">PostGIS 지리 공간 엔진</div>
          <div style="font-size: 14px; font-weight: 700; color: #16A34A;">● 정상 가동 (Polygon Checks)</div>
        </div>
        <div style="padding: 12px; border: 1px solid #F1F5F9; border-radius: 8px;">
          <div style="font-size: 12px; color: #64748B;">WebSocket 실시간 허브</div>
          <div style="font-size: 14px; font-weight: 700; color: #16A34A;">● 정상 가동 (Sequenced Duplex)</div>
        </div>
        <div style="padding: 12px; border: 1px solid #F1F5F9; border-radius: 8px;">
          <div style="font-size: 12px; color: #64748B;">거래 예약 잠금 FSM</div>
          <div style="font-size: 14px; font-weight: 700; color: #16A34A;">● 충돌 방지 활성 (Row Lock)</div>
        </div>
        <div style="padding: 12px; border: 1px solid #F1F5F9; border-radius: 8px;">
          <div style="font-size: 12px; color: #64748B;">감사 로깅 (Audit Log)</div>
          <div style="font-size: 14px; font-weight: 700; color: #16A34A;">● 불변 추가 기록 중</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const stats = ref<any>(null)

const loadStats = async () => {
  const token = localStorage.getItem('bn_admin_token')
  try {
    const res = await fetch('http://localhost:8080/api/v1/admin/dashboard', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.data) {
      stats.value = data.data
    }
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadStats()
})
</script>

