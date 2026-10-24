<template>
  <div class="container" style="padding-top: 30px; padding-bottom: 60px; max-width: 900px;">
    <div v-if="listing" style="background: white; border: 1px solid #E2E8F0; border-radius: 16px; overflow: hidden; box-shadow: var(--shadow-sm);">
      <!-- Main Photo -->
      <div style="width: 100%; height: 450px; background: #000; position: relative;">
        <img
          :src="listing.images?.[0] || 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?auto=format&fit=crop&w=1200&q=80'"
          :alt="listing.title"
          style="width: 100%; height: 100%; object-fit: contain;"
        />
        <div v-if="listing.listing_status === 'RESERVED'" style="position: absolute; top: 20px; left: 20px; background: #EA580C; color: white; padding: 4px 12px; border-radius: 6px; font-weight: 700; font-size: 14px;">
          예약중인 상품입니다
        </div>
      </div>

      <!-- Content Container -->
      <div style="padding: 30px;">
        <!-- Seller Info Ribbon -->
        <div style="display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid #F1F5F9; padding-bottom: 20px; margin-bottom: 24px;">
          <div style="display: flex; align-items: center; gap: 14px;">
            <div style="width: 48px; height: 48px; border-radius: 9999px; background: #E2E8F0; display: flex; align-items: center; justify-content: center; font-size: 20px;">
              👤
            </div>
            <div>
              <div style="font-weight: 700; font-size: 16px; color: #0F172A;">{{ listing.seller_nickname || '판매자' }}</div>
              <div style="font-size: 13px; color: #64748B;">{{ listing.neighborhood_name }} · 직거래 선호</div>
            </div>
          </div>
          <MannerBadge :temperature="listing.seller_manner_temp || 36.5" />
        </div>

        <!-- Title & Metadata -->
        <div style="margin-bottom: 20px;">
          <div style="display: flex; align-items: center; gap: 8px; margin-bottom: 8px;">
            <span style="font-size: 13px; color: #0F766E; font-weight: 700;">{{ listing.category_name || '디지털기기' }}</span>
            <span style="color: #CBD5E1;">·</span>
            <span class="badge-condition">{{ listing.condition }}</span>
            <span v-if="listing.is_negotiable" style="font-size: 12px; color: #2563EB; font-weight: 600;">가격제안 가능</span>
          </div>

          <h1 style="font-size: 26px; font-weight: 800; color: #0F172A; line-height: 1.3; margin-bottom: 12px;">
            {{ listing.title }}
          </h1>

          <div style="font-size: 28px; font-weight: 900; color: #0F766E;">
            <span v-if="listing.is_free" style="color: #16A34A;">무료나눔</span>
            <span v-else>{{ listing.price_minor?.toLocaleString() }}원</span>
          </div>
        </div>

        <!-- Description -->
        <div style="font-size: 16px; line-height: 1.8; color: #334155; white-space: pre-wrap; margin-bottom: 30px; border-bottom: 1px solid #F1F5F9; padding-bottom: 30px;">
          {{ listing.description }}
        </div>

        <!-- Safe Public Meeting Indicator -->
        <div style="background: #F8FAFC; border: 1px solid #E2E8F0; border-radius: 12px; padding: 16px; margin-bottom: 30px; display: flex; align-items: center; gap: 14px;">
          <div style="font-size: 24px;">🛡️</div>
          <div>
            <div style="font-size: 14px; font-weight: 700; color: #0F172A;">BN NEAR 안전 직거래 구역</div>
            <div style="font-size: 13px; color: #64748B;">판매자의 상세 거주지 주소는 공개되지 않습니다. {{ listing.neighborhood_name }} 공공장소(지하철역, 주민센터 등)에서 거래하세요.</div>
          </div>
        </div>

        <!-- Action Buttons -->
        <div style="display: flex; gap: 14px; align-items: center;">
          <button @click="toggleFavorite" :style="{
            padding: '12px 20px',
            borderRadius: '10px',
            border: '1px solid #E2E8F0',
            backgroundColor: isFav ? '#FEF2F2' : '#FFFFFF',
            color: isFav ? '#DC2626' : '#475569',
            fontWeight: '700',
            display: 'flex',
            alignItems: 'center',
            gap: '6px'
          }">
            <span>{{ isFav ? '❤️' : '🤍' }}</span>
            <span>관심 {{ favCount }}</span>
          </button>

          <button v-if="listing.is_negotiable && !isMyListing" @click="showOfferModal = true" class="btn-secondary" style="padding: 12px 24px; font-size: 15px;">
            가격 제안하기
          </button>

          <button v-if="!isMyListing" @click="startChat" class="btn-primary" style="padding: 12px 32px; font-size: 15px; margin-left: auto;">
            💬 채팅으로 거래하기
          </button>
          <div v-else style="margin-left: auto; font-size: 14px; font-weight: 600; color: #0F766E;">
            내 등록 물품입니다
          </div>
        </div>
      </div>
    </div>

    <!-- Offer Modal -->
    <div v-if="showOfferModal" style="position: fixed; top: 0; left: 0; width: 100vw; height: 100vh; background: rgba(0,0,0,0.5); z-index: 100; display: flex; align-items: center; justify-content: center;">
      <div style="background: white; border-radius: 16px; padding: 24px; width: 380px; max-width: 90%;">
        <h3 style="font-size: 18px; font-weight: 800; margin-bottom: 12px;">희망 가격 제안</h3>
        <p style="font-size: 13px; color: #64748B; margin-bottom: 16px;">판매자가 수락하면 거래가 성사됩니다. 신중하게 제안해주세요.</p>
        <div style="margin-bottom: 20px;">
          <label style="font-size: 12px; font-weight: 700; color: #475569; display: block; margin-bottom: 6px;">제안 금액 (원)</label>
          <input v-model.number="offerAmount" type="number" class="form-input" placeholder="예: 350000" />
        </div>
        <div style="display: flex; gap: 10px; justify-content: flex-end;">
          <button @click="showOfferModal = false" class="btn-secondary">취소</button>
          <button @click="submitOffer" class="btn-primary">제안 전송</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '~/stores/auth'
import { useApi } from '~/composables/useApi'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const { fetchAPI } = useApi()

const listing = ref<any>(null)
const isFav = ref(false)
const favCount = ref(0)
const showOfferModal = ref(false)
const offerAmount = ref<number>(0)

const isMyListing = computed(() => {
  return auth.user?.id && listing.value?.seller_user_id === auth.user.id
})

const loadListing = async () => {
  const id = route.params.id as string
  const res = await fetchAPI(`/listings/${id}`)
  if (res.data) {
    listing.value = res.data
    isFav.value = res.data.is_favorited || false
    favCount.value = res.data.favorite_count || 0
    if (res.data.price_minor) {
      offerAmount.value = Math.floor(res.data.price_minor * 0.9)
    }
  }
}

const toggleFavorite = async () => {
  if (!auth.isAuthenticated) {
    router.push('/login')
    return
  }
  const res = await fetchAPI('/me/favorites', {
    method: 'POST',
    body: { listing_id: listing.value.id }
  })
  if (res.data) {
    isFav.value = res.data.is_favorited
    favCount.value += isFav.value ? 1 : -1
  }
}

const startChat = async () => {
  if (!auth.isAuthenticated) {
    router.push('/login')
    return
  }
  const res = await fetchAPI('/conversations', {
    method: 'POST',
    body: { listing_id: listing.value.id }
  })
  if (res.data) {
    router.push(`/chats?id=${res.data.id}`)
  }
}

const submitOffer = async () => {
  if (!auth.isAuthenticated) {
    router.push('/login')
    return
  }
  const res = await fetchAPI(`/listings/${listing.value.id}/price-offers`, {
    method: 'POST',
    body: { amount_minor: offerAmount.value }
  })
  if (res.data) {
    alert('가격 제안이 전송되었습니다!')
    showOfferModal.value = false
  } else if (res.error) {
    alert(res.error.message || '제안 실패')
  }
}

onMounted(() => {
  loadListing()
})
</script>

