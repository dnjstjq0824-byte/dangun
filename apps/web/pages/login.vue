<template>
  <div class="container" style="padding-top: 50px; padding-bottom: 70px; max-width: 440px;">
    <div style="background: white; border: 1px solid #E2E8F0; border-radius: 16px; padding: 32px; box-shadow: var(--shadow-sm);">
      <div style="text-align: center; margin-bottom: 24px;">
        <span class="logo-badge" style="font-size: 16px; padding: 6px 12px;">NEAR</span>
        <h1 style="font-size: 24px; font-weight: 800; color: #0F172A; margin-top: 12px;">
          BN NEAR 로그인
        </h1>
        <p style="font-size: 13px; color: #64748B; margin-top: 4px;">우리 동네 이웃들과의 따뜻한 거래</p>
      </div>

      <form @submit.prevent="handleLogin" style="display: flex; flex-direction: column; gap: 16px;">
        <div>
          <label style="font-size: 13px; font-weight: 700; color: #475569; display: block; margin-bottom: 6px;">이메일</label>
          <input v-model="email" type="email" required class="form-input" placeholder="name@example.com" />
        </div>

        <div>
          <label style="font-size: 13px; font-weight: 700; color: #475569; display: block; margin-bottom: 6px;">비밀번호</label>
          <input v-model="password" type="password" required class="form-input" placeholder="비밀번호를 입력하세요" />
        </div>

        <button type="submit" class="btn-primary" style="justify-content: center; padding: 12px; font-size: 15px; margin-top: 8px;">
          로그인
        </button>
      </form>

      <!-- Quick Demo Logins -->
      <div style="margin-top: 24px; padding-top: 20px; border-top: 1px solid #F1F5F9;">
        <div style="font-size: 12px; font-weight: 700; color: #94A3B8; text-align: center; margin-bottom: 12px;">
          빠른 테스트 계정 바로 로그인
        </div>
        <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 8px;">
          <button @click="fillAccount('seller@bn-near.com', 'Password123!')" class="btn-secondary" style="font-size: 12px; padding: 6px;">
            🧑‍💼 판매자 (David)
          </button>
          <button @click="fillAccount('buyer@bn-near.com', 'Password123!')" class="btn-secondary" style="font-size: 12px; padding: 6px;">
            🙋‍♀️ 구매자 (Minji)
          </button>
          <button @click="fillAccount('merchant@bn-near.com', 'Password123!')" class="btn-secondary" style="font-size: 12px; padding: 6px;">
            🥖 동네가게 (MasterKim)
          </button>
          <button @click="fillAccount('admin@bn-near.com', 'Password123!')" class="btn-secondary" style="font-size: 12px; padding: 6px;">
            🛡️ 관리자 (Admin)
          </button>
        </div>
      </div>

      <div style="text-align: center; margin-top: 20px; font-size: 13px; color: #64748B;">
        아직 계정이 없으신가요?
        <NuxtLink to="/register" style="color: #0F766E; font-weight: 700;">회원가입하기</NuxtLink>
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
const password = ref('')

const fillAccount = (e: string, p: string) => {
  email.value = e
  password.value = p
}

const handleLogin = async () => {
  const res = await fetchAPI('/auth/login', {
    method: 'POST',
    body: {
      email: email.value,
      password: password.value,
      device_id: 'web-browser'
    }
  })

  if (res.data) {
    auth.setAuth(res.data.user, res.data.access_token)
    router.push('/')
  } else if (res.error) {
    alert(res.error.message || '로그인 실패')
  }
}
</script>

