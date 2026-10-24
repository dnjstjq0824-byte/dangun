# BN NEAR — 02. Product Scope & Release Boundaries

**Order ID:** BN-NEAR-MASTER-001  
**Project:** BN NEAR — Global Hyperlocal Marketplace & Community Platform  

---

## 1. Core Release Scope (v1.0 - Must Deliver)

### In-Scope (P0 / Mandatory):
- **User Lifecycle:** Registration, email verification, login/logout, profile management, block user, account deletion.
- **Location & Neighborhoods:** Country, City, Neighborhood hierarchy; PostGIS boundary checks; GPS location verification with validity window; primary/secondary neighborhood setup.
- **Marketplace Listings:** Create listing with up to 10 photos, category selection, item condition (`NEW` through `FOR_PARTS`), price in minor currency units, negotiable flag, neighborhood scoping.
- **Search & Discovery:** Text search, category filtering, distance radius (1km, 3km, 5km, 10km), price range, condition filter, sort by newest/distance/price.
- **Real-Time 1:1 Chat:** WebSocket chat gateway, listing context header, text & image messages, sequence IDs, read status synchronization, offline push event dispatch.
- **Deals & Negotiations:** In-chat price offers (accept/decline/counter), meeting place & time coordinator, reservation status lock, two-phase transaction completion (`COMPLETED_CLAIMED` -> `COMPLETED_CONFIRMED`).
- **Reviews & Reputation:** Deal-gated reviews, 5-star rating, manner tags, calculated manner temperature/reputation score.
- **Community:** Neighborhood-restricted posts (questions, daily life, recommendations, lost & found), threaded comments, post likes.
- **Local Business:** Store profile (hours, location, description, category), owner verification request, promotional business posts.
- **Safety & Moderation:** Report content/user dialogs, admin moderation queue, case decisions (warn, suspend, remove), user appeal submission and review.
- **Admin Console:** Vue 3 web admin for geographical data, user moderation, listing audits, verification overrides, and system audit logs.
- **Client Platforms:** 
  - Flutter Mobile (Android & iOS native binaries + responsive web layout)
  - Nuxt 3 Public Web (SSR for public listings, category pages, community)
  - Vue 3 Admin Cockpit

---

## 2. Deferred & Future Expansion Scope (Post-v1.0)
The following modules are strictly governed by legal and regulatory frameworks and must remain behind feature flags until approved:
- **Direct In-App Financial Escrow & Payments:** Requires payment service provider licenses.
- **Real Estate & Automotive Brokerage:** Requires state brokerage licenses.
- **Jobs & Employment Matching:** Subject to labor intermediary regulations.
- **Paid Advertising / Sponsored Listings:** Post-launch commercial feature.
