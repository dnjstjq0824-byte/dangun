# BN NEAR — 29. Known Limitations & Technical Constraints

**Order ID:** BN-NEAR-MASTER-001  

---

## 1. Regulatory & Scope Boundaries
1. **No In-App Payment Escrow:** In v1.0, peer-to-peer exchanges occur in person via mutual agreement. The platform does not hold user funds in escrow or operate payment custody.
2. **Third-Party Identity Verification:** While email verification and location boundary verification are implemented, national resident identity card verification (e.g. PASS/KMC in Korea) requires third-party enterprise vendor contracts and is marked as deferred.
3. **Apple iOS TestFlight Distribution:** Building signed `.ipa` archives requires active Apple Developer Team credentials. The Flutter code is 100% iOS-compliant and ready for Xcode archive when credentials are provided.

---

## 2. Technical Boundaries
1. **GPS Jitter at Edges:** Users standing within 10-20 meters of an administrative boundary line may experience occasional verification failures due to GPS multi-path reflection; an accuracy-weighted buffer tolerance is provided to mitigate this.
2. **Offline Image Sync:** Full offline photo synchronization is scheduled for subsequent mobile app iterations.

