<template>
  <div>
    <h2 style="font-size: 22px; font-weight: 800; margin-bottom: 20px;">글로벌 지역 및 동네 경계 (PostGIS) 거버넌스</h2>

    <div style="display: grid; grid-template-columns: 1fr 2fr; gap: 24px;">
      <!-- Countries & Cities -->
      <div>
        <h3 style="font-size: 16px; font-weight: 800; margin-bottom: 12px;">지원 국가 및 도시</h3>
        <div style="background: white; border: 1px solid var(--admin-border); border-radius: 12px; padding: 16px; display: flex; flex-direction: column; gap: 10px;">
          <div v-for="c in countries" :key="c.code" style="padding: 10px; border-bottom: 1px solid #F1F5F9; display: flex; justify-content: space-between; align-items: center;">
            <div>
              <span style="font-weight: 700; font-size: 14px;">{{ c.name }} ({{ c.code }})</span>
              <div style="font-size: 11px; color: #64748B;">통화: {{ c.default_currency }} · 언어: {{ c.default_locale }}</div>
            </div>
            <span style="font-size: 11px; font-weight: 700; background: #ECFDF5; color: #16A34A; padding: 2px 6px; border-radius: 4px;">
              {{ c.status }}
            </span>
          </div>
        </div>
      </div>

      <!-- Neighborhoods Boundary List -->
      <div>
        <h3 style="font-size: 16px; font-weight: 800; margin-bottom: 12px;">등록된 동네 폴리곤 경계</h3>
        <table class="data-table">
          <thead>
            <tr>
              <th>동네 명칭</th>
              <th>중심 좌표 (Centroid)</th>
              <th>폴리곤 꼭짓점 수</th>
              <th>상태</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="n in neighborhoods" :key="n.id">
              <td style="font-weight: 700;">{{ n.name }}</td>
              <td style="font-family: monospace; font-size: 12px;">
                {{ n.public_center.latitude.toFixed(4) }}, {{ n.public_center.longitude.toFixed(4) }}
              </td>
              <td>
                <span style="font-size: 12px; font-weight: 600; color: #0F766E;">
                  {{ n.boundary_vertices?.length || 4 }}개 좌표
                </span>
              </td>
              <td>
                <span style="font-size: 11px; font-weight: 700; background: #ECFDF5; color: #16A34A; padding: 2px 6px; border-radius: 4px;">
                  {{ n.status }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const countries = ref<any[]>([])
const neighborhoods = ref<any[]>([])

const loadGeo = async () => {
  try {
    const cRes = await fetch('http://localhost:8080/api/v1/countries')
    const cData = await cRes.json()
    if (cData.data) countries.value = cData.data

    const nRes = await fetch('http://localhost:8080/api/v1/neighborhoods/search')
    const nData = await nRes.json()
    if (nData.data) neighborhoods.value = nData.data
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadGeo()
})
</script>
