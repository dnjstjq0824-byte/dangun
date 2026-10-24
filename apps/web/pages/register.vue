<template>
  <div class="container" style="padding-top: 50px; padding-bottom: 70px; max-width: 440px;">
    <div style="background: white; border: 1px solid #E2E8F0; border-radius: 16px; padding: 32px; box-shadow: var(--shadow-sm);">
      <div style="text-align: center; margin-bottom: 24px;">
        <span class="logo-badge" style="font-size: 16px; padding: 6px 12px;">NEAR</span>
        <h1 style="font-size: 24px; font-weight: 800; color: #0F172A; margin-top: 12px;">
          BN NEAR 회원가입
        </h1>
        <p style="font-size: 13px; color: #64748B; margin-top: 4px;">우리 동네 주민들과 함께 시작하세요</p>
      </div>

      <form @submit.prevent="handleRegister" style="display: flex; flex-direction: column; gap: 16px;">
        <div>
          <label style="font-size: 13px; font-weight: 700; color: #475569; display: block; margin-bottom: 6px;">이메일</label>
          <input v-model="email" type="email" required class="form-input" placeholder="name@example.com" />
        </div>

        <div>
          <label style="font-size: 13px; font-weight: 700; color: #475569; display: block; margin-bottom: 6px;">닉네임</label>
          <input v-model="nickname" type="text" required class="form-input" placeholder="동네에서 사용할 닉네임 (2자 이상)" />
        </div>

        <div>
          <label style="font-size: 13px; font-weight: 700; color: #475569; display: block; margin-bottom: 6px;">비밀번호</label>
          <input v-model="password" type="password" required minlength="8" class="form-input" placeholder="8자 이상 입력해주세요" />
        </div>

        <button type="submit" class="btn-primary" style="justify-content: center; padding: 12px; font-size: 15px; margin-top: 8px;">
          가입 완료하기
        </button>
      </form>

      <div style="text-align: center; margin-top: 20px; font-size: 13px; color: #64748B;">
        이미 계정이 있으신가요?
        <NuxtLink to="/login" style="color: #0F766E; font-weight: 700;">로그인하기</NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '~/stores/auth'
import { useApi } from '~/composables/useApi'

const auth = useAuthStore()
const { fetchAPI } = useApi()
const router = useRouter()

const email = ref('')
const nickname = ref('')
const password = ref('')

const handleRegister = async () => {
  const res = await fetchAPI('/auth/register', {
    method: 'POST',
    body: {
      email: email.value,
      nickname: nickname.value,
      password: password.value,
      locale: 'ko-KR',
      country_code: 'KR'
    }
  })

  if (res.data) {
    auth.setAuth(res.data.user, res.data.access_token)
    alert('회원가입이 완료되었습니다!')
    router.push('/')
  } else if (res.error) {
    alert(res.error.message || '가입 실패')
  }
}
</script>

