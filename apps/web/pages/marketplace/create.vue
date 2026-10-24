<template>
  <div class="container" style="padding-top: 30px; padding-bottom: 60px; max-width: 680px;">
    <div style="background: white; border: 1px solid #E2E8F0; border-radius: 16px; padding: 30px; box-shadow: var(--shadow-sm);">
      <h1 style="font-size: 24px; font-weight: 800; color: #0F172A; margin-bottom: 20px;">
        내 물건 판매하기
      </h1>

      <form @submit.prevent="submitListing" style="display: flex; flex-direction: column; gap: 20px;">
        <!-- Photo input -->
        <div>
          <label style="display: block; font-size: 14px; font-weight: 700; color: #334155; margin-bottom: 8px;">
            상품 사진 URL (대표 이미지)
          </label>
          <input
            v-model="imageUrl"
            type="url"
            class="form-input"
            placeholder="https://images.unsplash.com/..."
          />
        </div>

        <!-- Title -->
        <div>
          <label style="display: block; font-size: 14px; font-weight: 700; color: #334155; margin-bottom: 8px;">
            글 제목
          </label>
          <input
            v-model="title"
            type="text"
            required
            class="form-input"
            placeholder="상품명과 주요 특징을 입력해주세요"
          />
        </div>

        <!-- Category -->
        <div>
          <label style="display: block; font-size: 14px; font-weight: 700; color: #334155; margin-bottom: 8px;">
            카테고리 선택
          </label>
          <select v-model="categoryID" required class="form-input">
            <option value="">카테고리를 선택하세요</option>
            <option v-for="cat in categories" :key="cat.id" :value="cat.id">
              {{ cat.name }}
            </option>
          </select>
        </div>

        <!-- Condition -->
        <div>
          <label style="display: block; font-size: 14px; font-weight: 700; color: #334155; margin-bottom: 8px;">
            상품 상태
          </label>
          <div style="display: flex; gap: 8px; flex-wrap: wrap;">
            <button
              type="button"
              v-for="c in conditions"
              :key="c.value"
              @click="condition = c.value"
              :style="{
                padding: '8px 14px',
                borderRadius: '8px',
                fontSize: '13px',
                fontWeight: '600',
                border: '1px solid #CBD5E1',
                backgroundColor: condition === c.value ? '#0F766E' : '#FFFFFF',
                color: condition === c.value ? '#FFFFFF' : '#475569'
              }"
            >
              {{ c.label }}
            </button>
          </div>
        </div>

        <!-- Price -->
        <div>
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px;">
            <label style="font-size: 14px; font-weight: 700; color: #334155;">판매 가격 (원)</label>
            <label style="font-size: 13px; font-weight: 600; color: #16A34A; display: flex; align-items: center; gap: 4px; cursor: pointer;">
              <input type="checkbox" v-model="isFree" />
              무료나눔
            </label>
          </div>
          <input
            v-if="!isFree"
            v-model.number="price"
            type="number"
            min="0"
            required
            class="form-input"
            placeholder="₩ 가격을 입력해주세요"
          />
          <div v-else style="padding: 10px 14px; background: #ECFDF5; border-radius: 8px; color: #16A34A; font-weight: 700; font-size: 14px;">
            이웃에게 무료로 따뜻하게 나눔합니다 ✨
          </div>
        </div>

        <!-- Negotiable toggle -->
        <label v-if="!isFree" style="display: flex; align-items: center; gap: 8px; font-size: 14px; font-weight: 600; cursor: pointer;">
          <input type="checkbox" v-model="isNegotiable" />
          가격 제안 받기 (이웃과 가격을 협상할 수 있습니다)
        </label>

        <!-- Description -->
        <div>
          <label style="display: block; font-size: 14px; font-weight: 700; color: #334155; margin-bottom: 8px;">
            자세한 설명
          </label>
          <textarea
            v-model="description"
            rows="6"
            required
            class="form-input"
            placeholder="구매 시기, 사용감, 하자 유무, 거래 선호 장소 등을 솔직하고 자세히 적어주세요."
          ></textarea>
        </div>

        <!-- Submit Button -->
        <button type="submit" class="btn-primary" style="justify-content: center; padding: 14px; font-size: 16px; margin-top: 10px;">
          등록 완료하기
        </button>
      </form>
    </div>
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

const title = ref('')
const categoryID = ref('')
const condition = ref('GOOD')
const price = ref<number>(0)
const isFree = ref(false)
const isNegotiable = ref(true)
const description = ref('')
const imageUrl = ref('https://images.unsplash.com/photo-1546868871-7041f2a55e12?auto=format&fit=crop&w=800&q=80')
const categories = ref<any[]>([])

const conditions = [
  { value: 'NEW', label: '미개봉 새상품' },
  { value: 'LIKE_NEW', label: '거의 새것' },
  { value: 'GOOD', label: '사용감 적음' },
  { value: 'FAIR', label: '사용감 있음' },
  { value: 'FOR_PARTS', label: '고장/부품용' }
]

const loadCategories = async () => {
  const res = await fetchAPI('/categories')
  if (res.data) {
    categories.value = res.data
  }
}

const submitListing = async () => {
  if (!auth.isAuthenticated) {
    router.push('/login')
    return
  }

  const payload = {
    title: title.value,
    description: description.value,
    category_id: categoryID.value,
    condition: condition.value,
    price_minor: isFree.value ? 0 : price.value,
    currency_code: 'KRW',
    is_free: isFree.value,
    is_negotiable: isNegotiable.value,
    neighborhood_id: auth.currentNeighborhood.id,
    images: [imageUrl.value]
  }

  const res = await fetchAPI('/listings', {
    method: 'POST',
    body: payload
  })

  if (res.data) {
    router.push(`/marketplace/${res.data.id}`)
  } else if (res.error) {
    alert(res.error.message || '등록 실패')
  }
}

onMounted(() => {
  if (!auth.isAuthenticated) {
    router.push('/login')
    return
  }
  loadCategories()
})
</script>

