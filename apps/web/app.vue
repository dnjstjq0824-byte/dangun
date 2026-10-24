<template>
  <div style="min-height: 100vh; display: flex; flex-direction: column;">
    <Navbar />
    <main style="flex: 1;">
      <NuxtPage />
    </main>
    <Footer />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from '~/stores/auth'
import '~/assets/css/main.css'

const auth = useAuthStore()

onMounted(() => {
  auth.initAuth()
})
</script>

