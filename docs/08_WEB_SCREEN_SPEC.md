# BN NEAR — 08. Public Web Screen Specifications (Nuxt 3)

**Order ID:** BN-NEAR-MASTER-001  
**Framework:** Nuxt 3 / Vue 3 / TypeScript / Pinia / TailwindCSS  
**Rendering Strategy:** Hybrid SSR (Server-Side Rendering for SEO pages, Client-Side for Chat & Account)  

---

## 1. Public Pages Hierarchy & SEO Structure

```text
/ [locale]                          --> Landing page with localized popular listings and neighborhood spotlight
/ [locale]/marketplace              --> Marketplace grid with faceted search, geo radius, categories, price filters
/ [locale]/marketplace/[slug]       --> SSR Product Detail with OpenGraph & Schema.org Product structured data
/ [locale]/categories/[slug]        --> Category hub with localized description and active listings
/ [locale]/neighborhoods/[slug]     --> Neighborhood hub with local statistics, boundaries, active items and community highlights
/ [locale]/community                --> Neighborhood community discussion board
/ [locale]/community/[id]           --> Discussion post detail with threaded comments
/ [locale]/businesses               --> Local merchant directory
/ [locale]/businesses/[slug]        --> Merchant profile, store hours, contact, and store posts
/ [locale]/login                    --> Authentication entry point
/ [locale]/register                 --> Registration flow with email verification step
/ [locale]/me/*                     --> Authenticated user dashboard (Listings, Favorites, Deals, Verification)
/ [locale]/chats                    --> Full-featured Web real-time chat interface (duplex WebSocket)
/ [locale]/safety                   --> Safety tips, prohibited items policy, anti-fraud guides
/ [locale]/terms & /privacy         --> Privacy policy, location data processing terms
```

---

## 2. Desktop vs. Mobile Web Responsiveness
- **Desktop (>= 1024px):**
  - Persistent sticky top navbar: Logo, Neighborhood switcher dropdown, Omnibox search with auto-suggestions, "Sell" CTA button, Chat indicator with unread pill, Notifications popover, Profile menu.
  - Multi-column layout: Left sidebar for categories and distance filters; main area 4-column product grid; right sidebar for local business spotlights and recent community topics.
- **Mobile Web (< 1024px):**
  - Adapts to app-like experience with bottom navigation bar matching the Flutter layout.
  - Mobile touch-friendly swipe gallery, sliding drawers for filter selection, and fixed bottom CTA on product detail.
