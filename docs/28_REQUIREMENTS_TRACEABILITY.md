# BN NEAR — 28. Complete Requirements Traceability Ledger

**Order ID:** BN-NEAR-MASTER-001  

---

## 1. Traceability Mapping: Requirements -> Architecture -> Implementation -> Tests

| Requirement | Design Document | Code Module | Verification Test | Status |
|---|---|---|---|:---:|
| **REQ-IDN-01** Auth & Registration | `docs/01_PRODUCT_REQUIREMENTS.md` | `services/api/internal/identity` | `TC-01`, `TC-02` | VERIFIED |
| **REQ-IDN-02** Sessions & JWT | `docs/10_SYSTEM_ARCHITECTURE.md` | `services/api/internal/identity` | `TC-02`, `TC-03` | VERIFIED |
| **REQ-LOC-01** Geo Hierarchy | `docs/11_DATABASE_SCHEMA.md` | `services/api/internal/geography` | `TC-04` | VERIFIED |
| **REQ-LOC-02** PostGIS Boundary | `docs/15_LOCATION_VERIFICATION.md` | `database/migrations` | `TC-05` | VERIFIED |
| **REQ-LOC-03** Verification Engine | `docs/15_LOCATION_VERIFICATION.md` | `services/api/internal/verification` | `TC-05`, `TC-06`, `TC-07` | VERIFIED |
| **REQ-LOC-04** Location Privacy | `docs/23_PRIVACY_AND_LOCATION_POLICY.md` | `services/api/internal/verification` | `TC-05` | VERIFIED |
| **REQ-MKT-01** Listings CRUD | `docs/17_LISTING_LIFECYCLE.md` | `services/api/internal/listings` | `TC-09`, `TC-10` | VERIFIED |
| **REQ-MKT-05** Geo Radius Search | `docs/16_SEARCH_AND_GEO_INDEXING.md` | `services/api/internal/search` | `TC-13` | VERIFIED |
| **REQ-MKT-07** Listing Favorites | `docs/11_DATABASE_SCHEMA.md` | `services/api/internal/favorites` | `TC-15` | VERIFIED |
| **REQ-CHT-01** 1:1 Conversation | `docs/14_WEBSOCKET_PROTOCOL.md` | `services/api/internal/chat` | `TC-16` | VERIFIED |
| **REQ-CHT-02** WebSocket Gateway | `contracts/websocket/protocol.md` | `services/api/internal/chat` | `TC-17` | VERIFIED |
| **REQ-CHT-03** Message Deduplication | `contracts/websocket/protocol.md` | `services/api/internal/chat` | `TC-18`, `TC-19` | VERIFIED |
| **REQ-CHT-04** Read Receipts | `contracts/websocket/protocol.md` | `services/api/internal/chat` | `TC-20` | VERIFIED |
| **REQ-CHT-06** Blocked User Guard | `docs/22_SECURITY_ARCHITECTURE.md` | `services/api/internal/chat` | `TC-21` | VERIFIED |
| **REQ-DEA-01** Price Offers | `docs/18_DEAL_STATE_MACHINE.md` | `services/api/internal/offers` | `TC-22` | VERIFIED |
| **REQ-DEA-03** Reservation Lock | `docs/18_DEAL_STATE_MACHINE.md` | `services/api/internal/deals` | `TC-23`, `TC-24` | VERIFIED |
| **REQ-DEA-04** Two-Phase Completion| `docs/18_DEAL_STATE_MACHINE.md` | `services/api/internal/deals` | `TC-25` | VERIFIED |
| **REQ-REP-01** Verified Deal Review| `docs/19_REPUTATION_POLICY.md` | `services/api/internal/reviews` | `TC-26`, `TC-28` | VERIFIED |
| **REQ-REP-03** Manner Temperature | `docs/19_REPUTATION_POLICY.md` | `services/api/internal/reviews` | `TC-27` | VERIFIED |
| **REQ-COM-01** Community Board | `docs/04_INFORMATION_ARCHITECTURE.md` | `services/api/internal/community` | `TC-29`, `TC-30` | VERIFIED |
| **REQ-BIZ-01** Local Business | `docs/21_BUSINESS_VERIFICATION.md` | `services/api/internal/business` | `TC-31` | VERIFIED |
| **REQ-MOD-01** Universal Report | `docs/20_COMMUNITY_MODERATION.md` | `services/api/internal/moderation` | `TC-32`, `TC-33` | VERIFIED |
| **REQ-ADM-01** Admin Cockpit | `docs/09_ADMIN_SCREEN_SPEC.md` | `apps/admin` | `TC-34` | VERIFIED |
| **REQ-ADM-06** Admin RBAC | `docs/24_RBAC_MATRIX.md` | `services/api/internal/admin` | `TC-34` | VERIFIED |
| **REQ-ADM-07** Audit Logging | `docs/11_DATABASE_SCHEMA.md` | `services/api/internal/admin` | `TC-35` | VERIFIED |
| **REQ-APP-01** Flutter Mobile App | `docs/07_FLUTTER_SCREEN_SPEC.md` | `apps/mobile` | `Flutter Analyze` | VERIFIED |
| **REQ-SEO-01** Nuxt 3 Public Web | `docs/08_WEB_SCREEN_SPEC.md` | `apps/web` | `Nuxt Build` | VERIFIED |

