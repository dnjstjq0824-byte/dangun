<template>
  <NuxtLink :to="`/marketplace/${listing.id}`" class="card" style="display: flex; flex-direction: column;">
    <!-- Image Box -->
    <div style="position: relative; width: 100%; padding-top: 100%; background-color: #E2E8F0; overflow: hidden;">
      <img
        :src="listing.images?.[0] || 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?auto=format&fit=crop&w=600&q=80'"
        :alt="listing.title"
        style="position: absolute; top: 0; left: 0; width: 100%; height: 100%; object-fit: cover;"
      />
      <div v-if="listing.listing_status === 'RESERVED'" style="position: absolute; top: 10px; left: 10px; background: #EA580C; color: white; font-size: 11px; font-weight: 700; padding: 2px 8px; border-radius: 4px;">
        예약중
      </div>
      <div v-else-if="listing.listing_status === 'SOLD'" style="position: absolute; top: 10px; left: 10px; background: #64748B; color: white; font-size: 11px; font-weight: 700; padding: 2px 8px; border-radius: 4px;">
        거래완료
      </div>
    </div>

    <!-- Info Box -->
    <div style="padding: 12px; display: flex; flex-direction: column; flex: 1;">
      <div style="font-size: 14px; font-weight: 600; color: #1E293B; margin-bottom: 4px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">
        {{ listing.title }}
      </div>

      <div style="font-size: 16px; font-weight: 800; color: #0F172A; margin-bottom: 6px;">
        <span v-if="listing.is_free" style="color: #16A34A;">무료나눔</span>
        <span v-else>{{ formatPrice(listing.price_minor, listing.currency_code) }}</span>
      </div>

      <div style="font-size: 12px; color: #64748B; display: flex; align-items: center; justify-content: space-between; margin-top: auto; padding-top: 8px;">
        <span>{{ listing.neighborhood_name || '청라동' }}</span>
        <div style="display: flex; align-items: center; gap: 8px; color: #94A3B8;">
          <span v-if="listing.chat_count">💬 {{ listing.chat_count }}</span>
          <span v-if="listing.favorite_count">❤️ {{ listing.favorite_count }}</span>
        </div>
      </div>
    </div>
  </NuxtLink>
</template>

<script setup lang="ts">
defineProps<{
  listing: any
}>()

const formatPrice = (minor: number, currency: string) => {
  if (!minor) return '0원'
  if (currency === 'KRW') {
    return minor.toLocaleString('ko-KR') + '원'
  }
  return currency + ' ' + (minor / 100).toFixed(2)
}
</script>

