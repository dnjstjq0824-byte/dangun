<template>
  <div class="container" style="padding-top: 24px; padding-bottom: 50px; max-width: 860px;">
    <!-- Header -->
    <div style="display: flex; justify-content: space-between; align-items: flex-end; margin-bottom: 24px;">
      <div>
        <h1 style="font-size: 26px; font-weight: 800; color: #0F172A;">동네생활</h1>
        <p style="font-size: 14px; color: #64748B;">{{ auth.currentNeighborhood.name }} 주민들과 나누는 소소한 이야기와 정보</p>
      </div>

      <button v-if="auth.isAuthenticated" @click="showWriteModal = true" class="btn-primary">
        ✏️ 동네 글쓰기
      </button>
    </div>

    <!-- Category Tabs -->
    <div style="display: flex; gap: 8px; margin-bottom: 20px; overflow-x: auto; padding-bottom: 6px;">
      <button
        v-for="cat in categories"
        :key="cat"
        @click="selectedCategory = cat; loadPosts()"
        :style="{
          padding: '6px 14px',
          borderRadius: '9999px',
          fontSize: '13px',
          fontWeight: '600',
          border: '1px solid #CBD5E1',
          backgroundColor: selectedCategory === cat ? '#0F766E' : '#FFFFFF',
          color: selectedCategory === cat ? '#FFFFFF' : '#475569'
        }"
      >
        {{ cat }}
      </button>
    </div>

    <!-- Posts Feed -->
    <div style="display: flex; flex-direction: column; gap: 14px;">
      <div
        v-for="post in posts"
        :key="post.id"
        class="card"
        style="padding: 20px;"
      >
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px;">
          <span style="font-size: 12px; font-weight: 700; color: #0F766E; background: #F0FDFA; padding: 2px 8px; border-radius: 4px;">
            {{ post.category }}
          </span>
          <span style="font-size: 12px; color: #94A3B8;">{{ post.author_nickname }} · {{ post.neighborhood_name }}</span>
        </div>

        <NuxtLink :to="`/community/${post.id}`" style="font-size: 17px; font-weight: 800; color: #0F172A; display: block; margin-bottom: 8px;">
          {{ post.title }}
        </NuxtLink>

        <p style="font-size: 14px; color: #475569; line-height: 1.6; margin-bottom: 14px;">
          {{ post.body }}
        </p>

        <div style="display: flex; justify-content: space-between; align-items: center; font-size: 13px; color: #94A3B8;">
          <button @click="toggleLike(post)" style="display: flex; align-items: center; gap: 4px; color: #64748B; font-weight: 600;">
            <span>{{ post.is_liked ? '❤️' : '🤍' }}</span>
            <span>공감 {{ post.like_count }}</span>
          </button>
          <NuxtLink :to="`/community/${post.id}`" style="color: #64748B; font-weight: 600;">
            💬 댓글 {{ post.comment_count }}
          </NuxtLink>
        </div>
      </div>
    </div>

    <!-- Write Modal -->
    <div v-if="showWriteModal" style="position: fixed; top: 0; left: 0; width: 100vw; height: 100vh; background: rgba(0,0,0,0.5); z-index: 100; display: flex; align-items: center; justify-content: center;">
      <div style="background: white; border-radius: 16px; padding: 24px; width: 500px; max-width: 90%;">
        <h3 style="font-size: 18px; font-weight: 800; margin-bottom: 16px;">동네생활 글쓰기</h3>
        <form @submit.prevent="submitPost" style="display: flex; flex-direction: column; gap: 14px;">
          <div>
            <label style="font-size: 13px; font-weight: 700; color: #475569; display: block; margin-bottom: 4px;">카테고리</label>
            <select v-model="newPostCategory" class="form-input" required>
              <option value="동네질문">동네질문</option>
              <option value="생활정보">생활정보</option>
              <option value="동네맛집/추천">동네맛집/추천</option>
              <option value="분실/실종">분실/실종</option>
              <option value="일상">일상</option>
            </select>
          </div>
          <div>
            <label style="font-size: 13px; font-weight: 700; color: #475569; display: block; margin-bottom: 4px;">제목</label>
            <input v-model="newPostTitle" type="text" class="form-input" required placeholder="제목을 입력하세요" />
          </div>
          <div>
            <label style="font-size: 13px; font-weight: 700; color: #475569; display: block; margin-bottom: 4px;">내용</label>
            <textarea v-model="newPostBody" class="form-input" rows="5" required placeholder="이웃들과 나누고 싶은 이야기를 적어보세요."></textarea>
          </div>
          <div style="display: flex; gap: 10px; justify-content: flex-end; margin-top: 10px;">
            <button type="button" @click="showWriteModal = false" class="btn-secondary">취소</button>
            <button type="submit" class="btn-primary">등록하기</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '~/stores/auth'
import { useApi } from '~/composables/useApi'

const auth = useAuthStore()
const { fetchAPI } = useApi()

const categories = ['전체', '동네질문', '생활정보', '동네맛집/추천', '분실/실종', '일상']
const selectedCategory = ref('전체')
const posts = ref<any[]>([])

const showWriteModal = ref(false)
const newPostCategory = ref('동네질문')
const newPostTitle = ref('')
const newPostBody = ref('')

const loadPosts = async () => {
  let url = `/community/posts?neighborhood_id=${auth.currentNeighborhood.id}`
  if (selectedCategory.value !== '전체') {
    url += `&category=${encodeURIComponent(selectedCategory.value)}`
  }
  const res = await fetchAPI(url)
  if (res.data) {
    posts.value = res.data
  }
}

const toggleLike = async (post: any) => {
  if (!auth.isAuthenticated) return
  const res = await fetchAPI(`/community/posts/${post.id}/likes`, { method: 'POST' })
  if (res.data) {
    post.is_liked = res.data.liked
    post.like_count += post.is_liked ? 1 : -1
  }
}

const submitPost = async () => {
  const res = await fetchAPI('/community/posts', {
    method: 'POST',
    body: {
      neighborhood_id: auth.currentNeighborhood.id,
      category: newPostCategory.value,
      title: newPostTitle.value,
      body: newPostBody.value
    }
  })
  if (res.data) {
    posts.value.unshift(res.data)
    showWriteModal.value = false
    newPostTitle.value = ''
    newPostBody.value = ''
  }
}

onMounted(() => {
  loadPosts()
})
</script>
