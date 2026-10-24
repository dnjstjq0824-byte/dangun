<template>
  <div class="container" style="padding-top: 20px; padding-bottom: 40px; height: calc(100vh - 140px); display: flex;">
    <div style="background: white; border: 1px solid #E2E8F0; border-radius: 16px; overflow: hidden; width: 100%; display: flex; box-shadow: var(--shadow-sm);">
      <!-- Left: Conversations List -->
      <div style="width: 320px; border-right: 1px solid #E2E8F0; display: flex; flex-direction: column;">
        <div style="padding: 16px; border-bottom: 1px solid #E2E8F0; font-size: 18px; font-weight: 800; color: #0F172A;">
          채팅 대화목록
        </div>

        <div style="flex: 1; overflow-y: auto;">
          <div
            v-for="conv in conversations"
            :key="conv.id"
            @click="selectConversation(conv)"
            :style="{
              padding: '14px 16px',
              borderBottom: '1px solid #F1F5F9',
              cursor: 'pointer',
              backgroundColor: activeConv?.id === conv.id ? '#F0FDFA' : '#FFFFFF',
              display: 'flex',
              gap: '12px',
              alignItems: 'center'
            }"
          >
            <div style="width: 44px; height: 44px; border-radius: 9999px; background: #E2E8F0; display: flex; align-items: center; justify-content: center; font-size: 18px;">
              👤
            </div>
            <div style="flex: 1; overflow: hidden;">
              <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 2px;">
                <span style="font-size: 14px; font-weight: 700; color: #0F172A;">
                  {{ conv.buyer_user_id === auth.user?.id ? conv.seller_nickname : conv.buyer_nickname }}
                </span>
                <span v-if="conv.unread_count > 0" style="background: #0F766E; color: white; border-radius: 9999px; font-size: 11px; font-weight: 700; padding: 1px 6px;">
                  {{ conv.unread_count }}
                </span>
              </div>
              <div style="font-size: 12px; color: #64748B; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
                {{ conv.last_message || '새로운 대화가 시작되었습니다.' }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Right: Active Chat Room -->
      <div v-if="activeConv" style="flex: 1; display: flex; flex-direction: column;">
        <!-- Top Context Header -->
        <div style="padding: 14px 20px; border-bottom: 1px solid #E2E8F0; background: #F8FAFC; display: flex; justify-content: space-between; align-items: center;">
          <div style="display: flex; align-items: center; gap: 12px;">
            <img
              :src="activeConv.listing_image || 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?auto=format&fit=crop&w=100&q=80'"
              style="width: 44px; height: 44px; border-radius: 6px; object-fit: cover;"
            />
            <div>
              <div style="font-size: 14px; font-weight: 700; color: #0F172A;">{{ activeConv.listing_title }}</div>
              <div style="font-size: 13px; font-weight: 800; color: #0F766E;">
                {{ activeConv.listing_price?.toLocaleString() }}원
              </div>
            </div>
          </div>

          <!-- Deal Lifecycle Actions -->
          <div style="display: flex; gap: 8px;">
            <!-- If Seller -->
            <template v-if="isSeller">
              <button
                v-if="dealStatus === 'INITIATED' || dealStatus === 'ACTIVE'"
                @click="reserveItem"
                class="btn-secondary"
                style="font-size: 12px; padding: 6px 12px;"
              >
                거래 예약하기
              </button>
              <button
                v-if="dealStatus === 'RESERVED'"
                @click="markSold"
                class="btn-primary"
                style="font-size: 12px; padding: 6px 12px;"
              >
                판매 완료하기
              </button>
            </template>

            <!-- If Buyer -->
            <template v-else>
              <button
                v-if="dealStatus === 'COMPLETED_CLAIMED'"
                @click="confirmPurchase"
                class="btn-primary"
                style="font-size: 12px; padding: 6px 12px; background: #16A34A;"
              >
                구매 확정하기
              </button>
            </template>

            <button
              v-if="dealStatus === 'COMPLETED_CONFIRMED'"
              @click="showReviewModal = true"
              class="btn-primary"
              style="font-size: 12px; padding: 6px 12px;"
            >
              거래 후기 작성 ⭐
            </button>
          </div>
        </div>

        <!-- Messages Stream -->
        <div ref="messagesContainer" style="flex: 1; overflow-y: auto; padding: 20px; display: flex; flex-direction: column; gap: 12px;">
          <div
            v-for="msg in messages"
            :key="msg.id || msg.sequence_no"
            :style="{
              alignSelf: msg.sender_user_id === auth.user?.id ? 'flex-end' : 'flex-start',
              maxWidth: '65%'
            }"
          >
            <div
              :style="{
                padding: '10px 14px',
                borderRadius: '14px',
                fontSize: '14px',
                backgroundColor: msg.sender_user_id === auth.user?.id ? '#0F766E' : '#F1F5F9',
                color: msg.sender_user_id === auth.user?.id ? '#FFFFFF' : '#0F172A',
                boxShadow: 'var(--shadow-sm)'
              }"
            >
              {{ msg.body }}
            </div>
            <div style="font-size: 10px; color: #94A3B8; margin-top: 2px; text-align: right;">
              {{ formatTime(msg.created_at) }}
            </div>
          </div>
        </div>

        <!-- Composer -->
        <div style="padding: 14px 20px; border-top: 1px solid #E2E8F0; display: flex; gap: 10px; align-items: center;">
          <input
            v-model="inputMessage"
            @keyup.enter="sendMessage"
            type="text"
            placeholder="메시지를 입력하세요..."
            class="form-input"
            style="flex: 1;"
          />
          <button @click="sendMessage" class="btn-primary" style="padding: 10px 18px;">
            전송
          </button>
        </div>
      </div>

      <!-- No Active Conversation Selected -->
      <div v-else style="flex: 1; display: flex; align-items: center; justify-content: center; color: #94A3B8;">
        채팅할 대화방을 선택하세요.
      </div>
    </div>

    <!-- Review Modal -->
    <div v-if="showReviewModal" style="position: fixed; top: 0; left: 0; width: 100vw; height: 100vh; background: rgba(0,0,0,0.5); z-index: 100; display: flex; align-items: center; justify-content: center;">
      <div style="background: white; border-radius: 16px; padding: 24px; width: 360px;">
        <h3 style="font-size: 18px; font-weight: 800; margin-bottom: 12px;">거래 후기 작성</h3>
        <p style="font-size: 13px; color: #64748B; margin-bottom: 16px;">이웃과의 따뜻한 거래 경험을 공유해주세요.</p>
        <div style="margin-bottom: 16px;">
          <label style="font-size: 12px; font-weight: 700; color: #475569; display: block; margin-bottom: 6px;">만족도 별점</label>
          <select v-model.number="reviewRating" class="form-input">
            <option :value="5">⭐⭐⭐⭐⭐ 5점 (최고예요)</option>
            <option :value="4">⭐⭐⭐⭐ 4점 (좋아요)</option>
            <option :value="3">⭐⭐⭐ 3점 (보통이에요)</option>
          </select>
        </div>
        <div style="margin-bottom: 20px;">
          <label style="font-size: 12px; font-weight: 700; color: #475569; display: block; margin-bottom: 6px;">한줄 후기</label>
          <input v-model="reviewComment" type="text" class="form-input" placeholder="친절하고 매너가 좋아요!" />
        </div>
        <div style="display: flex; gap: 10px; justify-content: flex-end;">
          <button @click="showReviewModal = false" class="btn-secondary">닫기</button>
          <button @click="submitReview" class="btn-primary">후기 등록</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '~/stores/auth'
import { useApi } from '~/composables/useApi'

const auth = useAuthStore()
const { fetchAPI } = useApi()
const route = useRoute()
const router = useRouter()

const conversations = ref<any[]>([])
const activeConv = ref<any>(null)
const messages = ref<any[]>([])
const inputMessage = ref('')
const messagesContainer = ref<HTMLElement | null>(null)

const dealStatus = ref('ACTIVE')
const currentDealID = ref('')
const showReviewModal = ref(false)
const reviewRating = ref(5)
const reviewComment = ref('')

let ws: WebSocket | null = null

const isSeller = computed(() => {
  return auth.user?.id && activeConv.value?.seller_user_id === auth.user.id
})

const formatTime = (ts: string) => {
  if (!ts) return ''
  const d = new Date(ts)
  return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

const loadConversations = async () => {
  const res = await fetchAPI('/conversations')
  if (res.data) {
    conversations.value = res.data
    const targetID = route.query.id as string
    if (targetID) {
      const target = conversations.value.find(c => c.id === targetID)
      if (target) {
        selectConversation(target)
      }
    } else if (conversations.value.length > 0) {
      selectConversation(conversations.value[0])
    }
  }
}

const selectConversation = async (conv: any) => {
  activeConv.value = conv
  const res = await fetchAPI(`/conversations/${conv.id}/messages`)
  if (res.data) {
    messages.value = res.data
    scrollToBottom()
  }

  // Load deal state
  const dealsRes = await fetchAPI('/me/deals')
  if (dealsRes.data) {
    const matchingDeal = dealsRes.data.find((d: any) => d.listing_id === conv.listing_id)
    if (matchingDeal) {
      currentDealID.value = matchingDeal.id
      dealStatus.value = matchingDeal.status
    } else {
      dealStatus.value = conv.listing_status || 'ACTIVE'
    }
  }
}

const initWebSocket = () => {
  if (!auth.token || !process.client) return
  const wsURL = `ws://localhost:8080/api/v1/ws?token=${auth.token}`
  ws = new WebSocket(wsURL)

  ws.onmessage = (event) => {
    const data = JSON.parse(event.data)
    if (data.type === 'chat.message.receive' || data.type === 'chat.message.ack') {
      if (activeConv.value && data.conversation_id === activeConv.value.id) {
        // Append if not duplicate
        if (!messages.value.some(m => m.sequence_no === data.sequence_no)) {
          messages.value.push({
            sequence_no: data.sequence_no,
            sender_user_id: data.sender_user_id || auth.user?.id,
            body: data.payload?.text || inputMessage.value,
            created_at: data.created_at || new Date().toISOString()
          })
          scrollToBottom()
        }
      }
    }
  }
}

const sendMessage = async () => {
  if (!inputMessage.value.trim() || !activeConv.value) return
  const text = inputMessage.value.trim()
  inputMessage.value = ''

  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({
      type: 'chat.message.send',
      conversation_id: activeConv.value.id,
      client_message_id: crypto.randomUUID(),
      payload: {
        message_type: 'TEXT',
        text: text
      }
    }))
  } else {
    // REST Fallback
    const res = await fetchAPI(`/conversations/${activeConv.value.id}/messages`, {
      method: 'POST',
      body: {
        client_message_id: crypto.randomUUID(),
        message_type: 'TEXT',
        body: text
      }
    })
    if (res.data) {
      messages.value.push(res.data)
      scrollToBottom()
    }
  }
}

const reserveItem = async () => {
  const res = await fetchAPI(`/deals/${activeConv.value.listing_id}/reserve`, {
    method: 'POST',
    body: { buyer_user_id: activeConv.value.buyer_user_id }
  })
  if (res.data) {
    dealStatus.value = 'RESERVED'
    currentDealID.value = res.data.id
    alert('상품이 예약되었습니다!')
  }
}

const markSold = async () => {
  if (!currentDealID.value) return
  const res = await fetchAPI(`/deals/${currentDealID.value}/seller-complete`, {
    method: 'POST'
  })
  if (res.data) {
    dealStatus.value = 'COMPLETED_CLAIMED'
    alert('판매 완료 처리되었습니다. 구매자 확인 대기 중입니다.')
  }
}

const confirmPurchase = async () => {
  if (!currentDealID.value) return
  const res = await fetchAPI(`/deals/${currentDealID.value}/buyer-confirm`, {
    method: 'POST'
  })
  if (res.data) {
    dealStatus.value = 'COMPLETED_CONFIRMED'
    alert('구매가 확정되었습니다! 거래 후기를 남겨보세요.')
  }
}

const submitReview = async () => {
  if (!currentDealID.value) return
  const res = await fetchAPI(`/deals/${currentDealID.value}/reviews`, {
    method: 'POST',
    body: {
      rating: reviewRating.value,
      tags: ['친절하고 매너가 좋아요', '시간 약속을 잘 지켜요'],
      comment: reviewComment.value
    }
  })
  if (res.data) {
    alert('후기가 정상적으로 등록되었습니다!')
    showReviewModal.value = false
  } else if (res.error) {
    alert(res.error.message || '후기 등록 실패')
  }
}

onMounted(() => {
  if (!auth.isAuthenticated) {
    router.push('/login')
    return
  }
  loadConversations()
  initWebSocket()
})
</script>
