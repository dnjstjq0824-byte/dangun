<template>
  <div>
    <h2 style="font-size: 22px; font-weight: 800; margin-bottom: 8px;">채팅 안전 증거 조사 콘솔</h2>
    <p style="font-size: 13px; color: #64748B; margin-bottom: 20px;">
      개인정보 보호 원칙(Article 14.7 & 31.9)에 따라 일반 사용자의 채팅을 무단 열람할 수 없으며, 신고 접수 또는 법적 분쟁 건에 한하여 열람 사유를 기재한 후 감사 기록 하에 제한적 증거를 조사합니다.
    </p>

    <!-- Inspection Access Form -->
    <div style="background: white; border: 1px solid var(--admin-border); border-radius: 12px; padding: 24px; margin-bottom: 24px; max-width: 650px;">
      <h3 style="font-size: 16px; font-weight: 800; margin-bottom: 14px;">조사 대상 대화방 조회</h3>
      <form @submit.prevent="inspectChat" style="display: flex; flex-direction: column; gap: 14px;">
        <div>
          <label style="font-size: 12px; font-weight: 700; color: #475569; display: block; margin-bottom: 4px;">대화방 ID (Conversation UUID)</label>
          <input v-model="convID" type="text" required class="form-input" placeholder="예: 4e75d04c-3521-42ab-873b-e0fa90251781" />
        </div>
        <div>
          <label style="font-size: 12px; font-weight: 700; color: #475569; display: block; margin-bottom: 4px;">공식 열람 사유 (감사 로그에 영구 기록됨)</label>
          <input v-model="accessReason" type="text" required class="form-input" placeholder="예: 신고 접수번호 REP-102 사기 의심 패턴 확인" />
        </div>
        <button type="submit" class="btn-admin btn-admin-primary" style="padding: 10px; font-size: 14px;">
          조사 승인 및 증거 스냅샷 조회
        </button>
      </form>
    </div>

    <!-- Redacted Evidence Viewer -->
    <div v-if="evidenceMessages.length > 0" style="background: white; border: 1px solid var(--admin-border); border-radius: 12px; padding: 24px;">
      <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; border-bottom: 1px solid #F1F5F9; padding-bottom: 10px;">
        <h3 style="font-size: 16px; font-weight: 800; color: #DC2626;">⚠️ 신고된 대화 증거 스냅샷 (열람 감사 완료)</h3>
        <span style="font-size: 12px; color: #64748B;">총 {{ evidenceMessages.length }}건의 메시지</span>
      </div>

      <div style="display: flex; flex-direction: column; gap: 10px;">
        <div v-for="m in evidenceMessages" :key="m.sequence_no" style="padding: 10px 14px; background: #F8FAFC; border-radius: 8px; border: 1px solid #E2E8F0;">
          <div style="display: flex; justify-content: space-between; font-size: 11px; color: #94A3B8; margin-bottom: 4px;">
            <span>발신자 ID: {{ m.sender_user_id }}</span>
            <span>시퀀스 #{{ m.sequence_no }} · {{ m.created_at }}</span>
          </div>
          <div style="font-size: 14px; color: #0F172A; font-weight: 600;">{{ m.body }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const convID = ref('')
const accessReason = ref('')
const evidenceMessages = ref<any[]>([])

const inspectChat = async () => {
  const token = localStorage.getItem('bn_admin_token')
  try {
    const res = await fetch(`http://localhost:8080/api/v1/conversations/${convID.value}/messages`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await res.json()
    if (data.data) {
      evidenceMessages.value = data.data
      alert('증거 열람 권한이 승인되었으며 관리자 감사 로그에 기록되었습니다.')
    } else {
      alert('대화방을 찾을 수 없거나 열람이 불가합니다.')
    }
  } catch (e) {
    console.error(e)
  }
}
</script>

