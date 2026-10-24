# BN NEAR — 25. Master Test Matrix & Verification Protocols

**Order ID:** BN-NEAR-MASTER-001  

---

## 1. Test Scenarios Matrix (35 Mandatory Functional Verifications)

| Test ID | Test Scenario | Expected Outcome | Verification Tier |
|---|---|---|---|
| **TC-01** | User Registration & Email Verification | User account created, password hashed with Bcrypt, activation flag set | Go Integration |
| **TC-02** | User Login & JWT Session Issuance | Valid access & refresh tokens issued; invalid password returns 401 | Go Unit / API |
| **TC-03** | Token Refresh & Session Invalidation | Refresh token rotation succeeds; revoked session returns 401 | Go API |
| **TC-04** | Country, City, Neighborhood Hierarchy Retrieval | Correct hierarchical relation returned; polygon boundaries present | Go API |
| **TC-05** | Location Verification Inside Neighborhood Boundary | Point inside polygon with accuracy < 100m returns VERIFIED (30d expiry) | Go Verification |
| **TC-06** | Location Verification Outside Neighborhood Boundary | Point outside polygon returns OUT_OF_BOUNDS error | Go Verification |
| **TC-07** | Location Verification Low GPS Accuracy | Accuracy > 150m returns LOW_ACCURACY error | Go Verification |
| **TC-08** | Multi-neighborhood Registration Quota | 1st and 2nd neighborhood succeed; 3rd registration rejected | Go API |
| **TC-09** | Create Listing in Verified Neighborhood | Listing stored in ACTIVE status with sanitized metadata | Go API |
| **TC-10** | Create Listing in Unverified Neighborhood | Creation blocked if neighborhood not verified | Go API |
| **TC-11** | Negative Price or Invalid Currency Validation | Rejection with 422 Unprocessable Entity | Go Unit |
| **TC-12** | Upload Multiple Photos & Cover Selection | Media assets associated, EXIF stripped, cover marked | Go API |
| **TC-13** | Spatial Radius Search (e.g. 5km) | Only listings within 5km radius returned; sorted by distance | Go Geo Search |
| **TC-14** | Category & Condition Search Filters | Results filtered accurately matching criteria | Go Search |
| **TC-15** | Favorite Listing Toggle & Counter Concurrency | Favorite record toggled; listing favorite_count accurately updated | Go Concurrency |
| **TC-16** | Initialize Chat Conversation from Listing | Unique conversation record retrieved or created | Go Chat API |
| **TC-17** | WebSocket Token Handshake & Connection | Authenticated connection established; connection.ready event received | Go WS Integration |
| **TC-18** | Real-time Message Send & Sequence Monotonicity | Message saved; sequence_no increments; client receives ACK | Go WS Integration |
| **TC-19** | Client Message Deduplication | Resending identical client_message_id returns original sequence without duplicate | Go WS Integration |
| **TC-20** | Real-time Read Receipt Synchronization | last_read_sequence updated; broadcasted to sender | Go WS Integration |
| **TC-21** | Blocked User Communication Guard | Blocked user cannot initiate chat or send messages (HTTP 403 / WS error) | Go Chat Security |
| **TC-22** | Price Offer Negotiation Lifecycle | Buyer proposes offer; seller accepts; status shifts to OFFER_ACCEPTED | Go Deal API |
| **TC-23** | Reservation State Transition with Row Lock | Seller reserves item for buyer; listing shifts to RESERVED | Go Deal API |
| **TC-24** | Reservation Collision Prevention | Concurrent reservation attempt on same listing fails with 409 Conflict | Go Concurrency |
| **TC-25** | Two-Phase Deal Completion | Seller marks sold (CLAIMED) -> Buyer confirms (CONFIRMED) | Go Deal State |
| **TC-26** | Review Eligibility Guard | Review rejected if deal is NOT in CONFIRMED state | Go Review API |
| **TC-27** | Mutual Reviews & Manner Score Calculation | Both parties review; manner temperature recalculated based on formula | Go Review API |
| **TC-28** | Duplicate Review Submission Prevention | Second review attempt on same deal by same user rejected (409) | Go Review API |
| **TC-29** | Community Post Creation (Verified Resident) | Verified member successfully publishes community post | Go Community API |
| **TC-30** | Community Threaded Comments & Upvotes | Comment threads and like counters update atomically | Go Community API |
| **TC-31** | Local Business Profile & Verification Flow | Store registered; admin approves; Verified badge displayed | Go Business API |
| **TC-32** | Abuse Report Creation & Moderation Triage | Report filed; case generated in moderation queue | Go Moderation API |
| **TC-33** | Moderator Enforcement Action & User Appeal | Moderator sanctions user; user files appeal; senior admin reviews | Go Moderation API |
| **TC-34** | Admin RBAC Permission Enforcement | Non-admin forbidden (403); Moderator restricted by role permissions | Go Admin RBAC |
| **TC-35** | Tamper-Proof Audit Logging Integrity | All administrative mutations generate immutable audit records | Go Audit Log |
