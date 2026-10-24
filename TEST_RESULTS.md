# BN NEAR — Test Results Ledger

**Order ID:** BN-NEAR-MASTER-001  
**Last Updated:** 2026-10-24T22:48:00+09:00  

---

## Test Execution Summary

| Test Category | Total Tests | Passed | Failed | Skipped | Status |
|---|---|---|---|---|---|
| **Go Domain & Unit Lifecycle** | 1 (13 sub-assertions) | 1 | 0 | 0 | **PASS** |
| **Go Integration & REST APIs** | 35 Scenarios | 35 | 0 | 0 | **PASS** |
| **WebSocket Real-Time Chat Suite** | 8 Protocol Checks | 8 | 0 | 0 | **PASS** |
| **Deals Concurrency & Lock Tests** | 2 Scenarios | 2 | 0 | 0 | **PASS** |
| **Location Boundary Verification** | 3 Geometries | 3 | 0 | 0 | **PASS** |
| **Vue Web (Nuxt 3) Production Build** | Full SSR Client + Server | Pass | 0 | 0 | **PASS** |
| **Vue 3 Admin Cockpit Build & Typecheck** | Vite + vue-tsc | Pass | 0 | 0 | **PASS** |
| **Flutter Mobile Analyze & Tests** | flutter analyze & test | Pass (0 issues) | 0 | 0 | **PASS** |

---

### 1. Backend Integration & Domain Test Suite (`services/api/api_test.go` & `services/api/ws_test.go`)
- **TC-01, TC-02:** User Registration, Password Bcrypt Hashing, JWT Auth generation.
- **TC-03:** Session validation, token expiration, revocation.
- **TC-04:** Geography Hierarchy (Country -> City -> Neighborhood).
- **TC-05:** Ray-Casting Polygon Boundary Verification inside Cheongna polygon (37.5385, 126.6558 -> Cheongna-dong).
- **TC-06:** Out-of-bounds location rejection (returns HTTP 422 Unprocessable Entity).
- **TC-07:** Low GPS accuracy rejection (>150m accuracy threshold returns HTTP 422).
- **TC-09, TC-11:** Marketplace Listing creation & Negative price validation (rejected with 400 Bad Request).
- **TC-13:** Spatial Proximity Haversine Radius query (<5000m).
- **TC-15:** Concurrent Listing Favorites toggle & counter increment.
- **TC-16, TC-18:** 1:1 Conversation creation & Monotonic Sequence No generation (`seq_no: 1, 2, 3...`).
- **TC-19:** Client Message ID deduplication (retransmissions return existing sequence without creating duplicates).
- **TC-20:** Real-time read receipt updates (`last_read_seq_no`).
- **TC-22:** Price offer creation, negotiation, acceptance.
- **TC-23, TC-24:** Reservation row-level collision locking (rejecting concurrent second buyer with HTTP 409 Conflict).
- **TC-25:** Two-phase deal lifecycle: Seller marks sold (`COMPLETED_CLAIMED`) -> Buyer confirms receipt (`COMPLETED_CONFIRMED`).
- **TC-26:** Premature review gating (attempting review before buyer confirmation blocked with HTTP 412 Precondition Failed).
- **TC-27, TC-28:** Mutual review submission, Manner Temperature update (38.2°C -> higher), duplicate review rejection (HTTP 409 Conflict).
- **TC-29, TC-30:** Community board posts, comments, likes.
- **TC-32:** Safety abuse reports and auto-triage queue.
- **TC-34, TC-35:** Admin RBAC permissions (HTTP 403 Forbidden for non-admin) & tamper-proof audit logging with mandatory reason enforcement.
- **WebSocket Protocol Suite:** Connection handshake (`connection.ready`), `chat.message.send`, server `chat.message.ack`, broadcast `chat.message.receive`, ping/pong heartbeat.

---

### 2. Frontend Applications Verification
- **Nuxt 3 Public Web App (`apps/web`):**
  - Command: `pnpm build`
  - Output: Nitro 2.13.5 node-server preset built, client manifest (216 modules), server bundles (155 modules), static assets generated into `.output/server` and `.output/public`. 0 errors.
- **Vue 3 Admin Cockpit (`apps/admin`):**
  - Command: `vue-tsc -b && vite build`
  - Output: 45 modules transformed, typecheck clean, production bundle output in `dist/`. 0 errors.
- **Flutter Mobile App (`apps/mobile`):**
  - Command: `flutter analyze` -> `No issues found! (ran in 7.1s)`
  - Command: `flutter test` -> `All tests passed! (00:00 +1)`
