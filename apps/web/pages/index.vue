<template>
  <div>
    <!-- Hero Section -->
    <section style="background: linear-gradient(135deg, #0F766E 0%, #115E59 100%); color: white; padding: 50px 0;">
      <div class="container" style="display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 30px;">
        <div style="max-width: 580px;">
          <div style="display: inline-block; background: rgba(255,255,255,0.2); padding: 4px 12px; border-radius: 9999px; font-size: 13px; font-weight: 600; margin-bottom: 12px;">
            우리 동네 하이퍼로컬 라이프스타일
          </div>
          <h1 style="font-size: 36px; font-weight: 800; line-height: 1.25; margin-bottom: 16px; letter-spacing: -0.5px;">
            가까운 이웃과 함께하는<br />따뜻하고 안전한 동네 생활
          </h1>
          <p style="font-size: 16px; opacity: 0.9; line-height: 1.6; margin-bottom: 24px;">
            걸어서 만날 수 있는 이웃과 믿을 수 있는 중고거래, 동네 유용한 정보 나눔, 그리고 우리 동네 진짜 숨은 가게들을 BN NEAR에서 만나보세요.
          </p>
          <div style="display: flex; gap: 12px;">
            <NuxtLink to="/marketplace" class="btn-primary" style="background: white; color: #0F766E; font-size: 15px; padding: 10px 22px;">
              동네 매물 둘러보기
            </NuxtLink>
            <NuxtLink to="/community" class="btn-secondary" style="background: rgba(255,255,255,0.15); color: white; border: none; font-size: 15px; padding: 10px 20px;">
              동네생활 이야기
            </NuxtLink>
          </div>
        </div>

        <div style="background: white; padding: 24px; border-radius: 16px; color: #0F172A; box-shadow: 0 20px 25px -5px rgba(0,0,0,0.2); width: 320px;">
          <div style="font-size: 13px; font-weight: 700; color: #0F766E; margin-bottom: 4px;">📍 현재 탐색 지역</div>
          <div style="font-size: 18px; font-weight: 800; margin-bottom: 12px;">{{ auth.currentNeighborhood.name }}</div>
          <div style="font-size: 13px; color: #64748B; margin-bottom: 16px;">
            내 실제 위치를 인증하면 동네 주민들과 실시간 채팅 및 직거래가 가능합니다.
          </div>
          <NuxtLink to="/me" class="btn-primary" style="width: 100%; justify-content: center;">
            내 동네 GPS 인증하기
          </NuxtLink>
        </div>
      </div>
    </section>

    <!-- Categories Row -->
    <section class="container" style="margin-top: 36px;">
      <h2 style="font-size: 20px; font-weight: 800; margin-bottom: 16px;">인기 카테고리</h2>
      <div style="display: flex; gap: 10px; overflow-x: auto; padding-bottom: 8px;">
        <button
          v-for="cat in categories"
          :key="cat.id"
          @click="selectCategory(cat.id)"
          :style="{
            padding: '8px 16px',
            borderRadius: '9999px',
            fontSize: '13px',
            fontWeight: '600',
            whiteSpace: 'nowrap',
            backgroundColor: selectedCat === cat.id ? '#0F766E' : '#FFFFFF',
            color: selectedCat === cat.id ? '#FFFFFF' : '#475569',
            border: '1px solid #E2E8F0',
            cursor: 'pointer'
          }"
        >
          {{ cat.name }}
        </button>
      </div>
    </section>

    <!-- Nearby Listings Feed -->
    <section class="container" style="margin-top: 36px;">
      <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px;">
        <div>
          <h2 style="font-size: 22px; font-weight: 800;">방금 올라온 우리 동네 물품</h2>
          <p style="font-size: 13px; color: #64748B;">이웃들이 직접 등록한 실시간 매물입니다</p>
        </div>
        <NuxtLink to="/marketplace" style="font-size: 14px; font-weight: 600; color: #0F766E;">
          전체보기 →
        </NuxtLink>
      </div>

      <div style="display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)); gap: 20px;">
        <ListingCard
          v-for="listing in listings"
          :key="listing.id"
          :listing="listing"
        />
      </div>
    </section>

    <!-- Community & Local Business Highlights -->
    <section class="container" style="margin-top: 48px; margin-bottom: 40px;">
      <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(320px, 1fr)); gap: 24px;">
        <!-- Community Box -->
        <div style="background: white; border: 1px solid #E2E8F0; border-radius: 16px; padding: 24px;">
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;">
            <h3 style="font-size: 18px; font-weight: 800; color: #0F172A;">동네생활 최신 소식</h3>
            <NuxtLink to="/community" style="font-size: 13px; color: #0F766E; font-weight: 600;">더보기 →</NuxtLink>
          </div>
          <div v-for="post in communityPosts" :key="post.id" style="border-bottom: 1px solid #F1F5F9; padding: 12px 0;">
            <div style="font-size: 11px; font-weight: 700; color: #0F766E; margin-bottom: 4px;">{{ post.category }}</div>
            <NuxtLink :to="`/community/${post.id}`" style="font-size: 14px; font-weight: 600; color: #1E293B; display: block; margin-bottom: 4px;">
              {{ post.title }}
            </NuxtLink>
            <div style="font-size: 12px; color: #94A3B8; display: flex; justify-content: space-between;">
              <span>{{ post.author_nickname }}</span>
              <span>💬 {{ post.comment_count }} · ❤️ {{ post.like_count }}</span>
            </div>
          </div>
        </div>

        <!-- Business Spotlight Box -->
        <div style="background: white; border: 1px solid #E2E8F0; border-radius: 16px; padding: 24px;">
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;">
            <h3 style="font-size: 18px; font-weight: 800; color: #0F172A;">우리 동네 착한 가게</h3>
            <NuxtLink to="/businesses" style="font-size: 13px; color: #0F766E; font-weight: 600;">더보기 →</NuxtLink>
          </div>
          <div v-for="biz in businesses" :key="biz.id" style="border-bottom: 1px solid #F1F5F9; padding: 12px 0;">
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 4px;">
              <span style="font-size: 14px; font-weight: 700; color: #0F172A;">{{ biz.name }}</span>
              <span class="badge-verified" v-if="biz.verification_status === 'VERIFIED'">인증가게</span>
            </div>
            <p style="font-size: 12px; color: #64748B; margin-bottom: 6px;">{{ biz.description }}</p>
            <div style="font-size: 11px; color: #94A3B8;">📍 {{ biz.public_address }}</div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '~/stores/auth'
import { useApi } from '~/composables/useApi'

const auth = useAuthStore()
const { fetchAPI } = useApi()
const router = useRouter()

const categories = ref<any[]>([])
const listings = ref<any[]>([])
const communityPosts = ref<any[]>([])
const businesses = ref<any[]>([])
const selectedCat = ref<string>('')

const loadData = async () => {
  const catRes = await fetchAPI('/categories')
  if (catRes.data) {
    categories.value = [{ id: '', name: '전체' }, ...catRes.data]
  }

  const listRes = await fetchAPI('/listings')
  if (listRes.data) {
    listings.value = listRes.data
  }

  const commRes = await fetchAPI('/community/posts?neighborhood_id=b0000000-0000-0000-0000-000000000001')
  if (commRes.data) {
    communityPosts.value = commRes.data.slice(0, 3)
  }

  const bizRes = await fetchAPI('/businesses?neighborhood_id=b0000000-0000-0000-0000-000000000001')
  if (bizRes.data) {
    businesses.value = bizRes.data.slice(0, 3)
  }
}

const selectCategory = (id: string) => {
  selectedCat.value = id
  if (id) {
    router.push({ path: '/marketplace', query: { category_id: id } })
  } else {
    router.push('/marketplace')
  }
}

onMounted(() => {
  loadData()
})
</script>
