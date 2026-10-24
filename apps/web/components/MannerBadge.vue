<template>
  <div style="display: inline-flex; align-items: center; gap: 8px;">
    <div style="text-align: right;">
      <div style="font-size: 15px; font-weight: 800; color: #0F766E;">
        {{ temperature.toFixed(1) }}°C
      </div>
      <div style="font-size: 10px; color: #64748B;">매너온도</div>
    </div>
    <div style="width: 36px; height: 36px; border-radius: 9999px; background: #ECFDF5; border: 2px solid #0F766E; display: flex; align-items: center; justify-content: center; font-size: 16px;">
      <span v-if="temperature >= 40">🔥</span>
      <span v-else-if="temperature >= 37.5">😊</span>
      <span v-else>🙂</span>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  temperature: number
}>()
</script>

