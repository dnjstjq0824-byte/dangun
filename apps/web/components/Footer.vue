<template>
  <footer class="footer">
    <div class="container" style="display: flex; flex-direction: column; gap: 20px;">
      <div style="display: flex; justify-content: space-between; align-items: flex-start; flex-wrap: wrap; gap: 20px;">
        <div>
          <div style="font-weight: 800; font-size: 16px; color: #0F766E; margin-bottom: 6px;">BN NEAR</div>
          <p style="font-size: 13px; color: #64748B;">우리 동네의 모든 연결 · Booking Network Co., Ltd.</p>
          <p style="font-size: 12px; color: #94A3B8; margin-top: 4px;">글로벌 지역 기반 중고거래 및 동네 생활 슈퍼앱 플랫폼</p>
        </div>
        <div style="display: flex; gap: 24px; font-size: 13px;">
          <div>
            <div style="font-weight: 600; color: #334155; margin-bottom: 8px;">서비스</div>
            <div style="display: flex; flex-direction: column; gap: 4px; color: #64748B;">
              <NuxtLink to="/marketplace">중고거래</NuxtLink>
              <NuxtLink to="/community">동네생활</NuxtLink>
              <NuxtLink to="/businesses">동네가게</NuxtLink>
            </div>
          </div>
          <div>
            <div style="font-weight: 600; color: #334155; margin-bottom: 8px;">신뢰와 안전</div>
            <div style="display: flex; flex-direction: column; gap: 4px; color: #64748B;">
              <NuxtLink to="/safety">안전거래 가이드</NuxtLink>
              <NuxtLink to="/terms">이용약관</NuxtLink>
              <NuxtLink to="/privacy">위치정보 및 개인정보처리방침</NuxtLink>
            </div>
          </div>
        </div>
      </div>
      <div style="border-top: 1px solid #E2E8F0; padding-top: 16px; font-size: 12px; color: #94A3B8; display: flex; justify-content: space-between;">
        <span>© 2026 Booking Network Co., Ltd. All rights reserved.</span>
        <span>Order ID: BN-NEAR-MASTER-001</span>
      </div>
    </div>
  </footer>
</template>
