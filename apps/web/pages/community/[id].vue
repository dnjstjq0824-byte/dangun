<template>
  <div class="container" style="padding-top: 24px; padding-bottom: 60px; max-width: 760px;">
    <NuxtLink to="/community" style="font-size: 14px; color: #0F766E; font-weight: 700; margin-bottom: 16px; display: inline-block;">
      ← 목록으로 돌아가기
    </NuxtLink>

    <div v-if="post" style="background: white; border: 1px solid #E2E8F0; border-radius: 16px; padding: 28px; box-shadow: var(--shadow-sm); margin-bottom: 24px;">
      <div style="font-size: 13px; font-weight: 700; color: #0F766E; margin-bottom: 8px;">
        {{ post.category }}
      </div>

      <h1 style="font-size: 24px; font-weight: 800; color: #0F172A; margin-bottom: 14px;">
        {{ post.title }}
      </h1>

      <div style="display: flex; align-items: center; gap: 10px; margin-bottom: 20px; font-size: 13px; color: #64748B; border-bottom: 1px solid #F1F5F9; padding-bottom: 14px;">
        <span style="font-weight: 700; color: #1E293B;">{{ post.author_nickname }}</span>
        <span>·</span>
        <span>{{ post.neighborhood_name }}</span>
      </div>

      <div style="font-size: 16px; line-height: 1.8; color: #334155; white-space: pre-wrap; margin-bottom: 24px;">
        {{ post.body }}
      </div>

      <div style="display: flex; gap: 16px; border-top: 1px solid #F1F5F9; padding-top: 16px;">
        <button @click="toggleLike" style="font-weight: 700; color: #0F766E; display: flex; align-items: center; gap: 6px;">
          <span>{{ post.is_liked ? '❤️' : '🤍' }}</span>
          <span>공감 {{ post.like_count }}</span>
        </button>
        <span style="color: #64748B; font-weight: 600;">💬 댓글 {{ comments.length }}</span>
      </div>
    </div>

    <!-- Comments Section -->
    <div style="background: white; border: 1px solid #E2E8F0; border-radius: 16px; padding: 24px;">
      <h3 style="font-size: 18px; font-weight: 800; margin-bottom: 18px;">댓글 {{ comments.length }}개</h3>

      <!-- Comments List -->
      <div style="display: flex; flex-direction: column; gap: 16px; margin-bottom: 24px;">
        <div v-for="c in comments" :key="c.id" style="border-bottom: 1px solid #F1F5F9; padding-bottom: 12px;">
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 4px;">
            <span style="font-weight: 700; font-size: 14px; color: #1E293B;">{{ c.author_nickname }}</span>
            <span style="font-size: 11px; color: #94A3B8;">{{ formatTime(c.created_at) }}</span>
          </div>
          <p style="font-size: 14px; color: #475569; line-height: 1.5;">{{ c.body }}</p>
        </div>
      </div>

      <!-- Add Comment Form -->
      <form @submit.prevent="submitComment" style="display: flex; gap: 10px;">
        <input
          v-model="commentBody"
          type="text"
          class="form-input"
          placeholder="이웃에게 따뜻한 댓글을 남겨보세요..."
          required
          style="flex: 1;"
        />
        <button type="submit" class="btn-primary" style="padding: 10px 20px;">
          등록
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '~/stores/auth'
import { useApi } from '~/composables/useApi'

const route = useRoute()
const auth = useAuthStore()
const { fetchAPI } = useApi()

const post = ref<any>(null)
const comments = ref<any[]>([])
const commentBody = ref('')

const formatTime = (ts: string) => {
  if (!ts) return ''
  return new Date(ts).toLocaleString([], { dateStyle: 'short', timeStyle: 'short' })
}

const loadPost = async () => {
  const id = route.params.id as string
  const pRes = await fetchAPI(`/community/posts/${id}`)
  if (pRes.data) {
    post.value = pRes.data
  }

  const cRes = await fetchAPI(`/community/posts/${id}/comments`)
  if (cRes.data) {
    comments.value = cRes.data
  }
}

const toggleLike = async () => {
  if (!auth.isAuthenticated || !post.value) return
  const res = await fetchAPI(`/community/posts/${post.value.id}/likes`, { method: 'POST' })
  if (res.data) {
    post.value.is_liked = res.data.liked
    post.value.like_count += post.value.is_liked ? 1 : -1
  }
}

const submitComment = async () => {
  if (!auth.isAuthenticated || !commentBody.value.trim() || !post.value) return
  const res = await fetchAPI(`/community/posts/${post.value.id}/comments`, {
    method: 'POST',
    body: { body: commentBody.value.trim() }
  })
  if (res.data) {
    comments.value.push(res.data)
    commentBody.value = ''
  }
}

onMounted(() => {
  loadPost()
})
</script>
