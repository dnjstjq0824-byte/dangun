# BN NEAR — Master Requirements Traceability Matrix

**Project:** BN NEAR — Global Hyperlocal Marketplace & Community Platform  
**Document ID:** BN-NEAR-RTM-001  
**Total Requirement Items:** 60 Core Specifications  

| Req ID | Domain | Specification Summary | Source Article | Implementation Target | Verification Method |
|---|---|---|---|---|---|
| **REQ-IDN-01** | Identity | Email/Password Registration, validation, Argon2/Bcrypt hash | Art. 8, 20.2, 22 | `services/api/internal/identity` | Unit + E2E Auth Test |
| **REQ-IDN-02** | Identity | JWT Session management, Refresh tokens, Device identification | Art. 8, 20.2, 22 | `services/api/internal/identity` | Token Rotation Test |
| **REQ-IDN-03** | Identity | User Profile, bio, preferences, currency, locale, timezone | Art. 8, 20.2, 22 | `services/api/internal/identity` | Profile API Test |
| **REQ-IDN-04** | Identity | Account Deletion, Soft delete, session revocation, privacy sanitization | Art. 8, 35 | `services/api/internal/identity` | Deletion Test |
| **REQ-LOC-01** | Location | Global Geo Hierarchy: Country -> Region -> City -> Neighborhood | Art. 9, 20.3, 23 | `services/api/internal/geography` | Hierarchy Query Test |
| **REQ-LOC-02** | Location | PostGIS Polygonal Boundary & Center Coordinates storage | Art. 9, 20.3, 23 | `database/migrations` | Spatial Query Test |
| **REQ-LOC-03** | Location | Real-time Location Verification (Lat/Lng, Accuracy, Expiry, Distance) | Art. 9, 23 | `services/api/internal/verification` | Verification Engine Test |
| **REQ-LOC-04** | Location | Privacy preservation: No raw GPS exposure, approximate public coordinates | Art. 9.4, 35 | `services/api/internal/verification` | Privacy Audit Test |
| **REQ-LOC-05** | Location | Multi-neighborhood registration (Primary + Secondary), quota limit | Art. 9.5, 23 | `services/api/internal/verification` | Quota Guard Test |
| **REQ-MKT-01** | Marketplace | Listing creation with category, condition, price_minor, currency, attributes | Art. 11, 20.4, 24 | `services/api/internal/listings` | CRUD Test |
| **REQ-MKT-02** | Marketplace | Multiple image uploads, cover image selection, EXIF scrub | Art. 11.4, 20.4 | `services/api/internal/listings` | Media Security Test |
| **REQ-MKT-03** | Marketplace | Item Condition (`NEW`, `LIKE_NEW`, `GOOD`, `FAIR`, `FOR_PARTS`) | Art. 11.2 | `services/api/internal/listings` | Enum Validation Test |
| **REQ-MKT-04** | Marketplace | Listing Status Transition (`DRAFT`, `ACTIVE`, `RESERVED`, `SOLD`, `ARCHIVED`) | Art. 11.3 | `services/api/internal/listings` | State Transition Test |
| **REQ-MKT-05** | Marketplace | Hyperlocal Geospatial Discovery with PostGIS Radius & Bounding Box | Art. 12, 20.14, 24 | `services/api/internal/search` | Geo Distance Query Test |
| **REQ-MKT-06** | Marketplace | Multi-criteria Filtering (Keyword, Category, Price Range, Condition, Free) | Art. 12.1, 24 | `services/api/internal/search` | Search Filter Test |
| **REQ-MKT-07** | Marketplace | Listing Favorites & Counter concurrency tracking | Art. 20.4, 24 | `services/api/internal/favorites` | Concurrent Favorites Test |
| **REQ-CHT-01** | Chat | 1:1 Conversation initialization tied to Listing (Buyer + Seller) | Art. 14, 20.5, 25 | `services/api/internal/chat` | Conversation Init Test |
| **REQ-CHT-02** | Chat | WebSocket Gateway (`/api/v1/ws`) with token authentication | Art. 14.3, 25 | `services/api/internal/chat` | WS Connection Test |
| **REQ-CHT-03** | Chat | Guaranteed Message Storage: Client Message ID, Sequence No, Server ACK | Art. 14.4, 14.5 | `services/api/internal/chat` | Message ACK & Ordering Test |
| **REQ-CHT-04** | Chat | Real-time Read Receipts (`last_read_sequence`) & Unread counts | Art. 14.1, 25 | `services/api/internal/chat` | Read Receipt Test |
| **REQ-CHT-05** | Chat | Session Reconnection & Missing message catch-up synchronization | Art. 14.5, 25 | `services/api/internal/chat` | Catch-up Sync Test |
| **REQ-CHT-06** | Chat | User blocking enforcement: Blocked users cannot message | Art. 14.6, 20.2 | `services/api/internal/chat` | Block Guard Test |
| **REQ-DEA-01** | Deals | Price Offer creation, negotiation, acceptance, decline, expiration | Art. 15, 20.6, 26 | `services/api/internal/offers` | Offer Lifecycle Test |
| **REQ-DEA-02** | Deals | Deal Record creation & Meeting schedule (time, safe public place) | Art. 15.1, 20.6, 26 | `services/api/internal/deals` | Deal Meeting Test |
| **REQ-DEA-03** | Deals | Reservation Collision Lock: Only 1 active reservation per listing | Art. 15.5, 26 | `services/api/internal/deals` | Concurrency Lock Test |
| **REQ-DEA-04** | Deals | Two-phase Completion: Seller marks sold, Buyer confirms | Art. 15.2, 26 | `services/api/internal/deals` | 2-Phase State Machine Test |
| **REQ-REP-01** | Reputation | Deal-verified Reviews: Review permitted only upon confirmed Deal record | Art. 16, 20.7, 27 | `services/api/internal/reviews` | Review Eligibility Test |
| **REQ-REP-02** | Reputation | Mutual reviews: Both buyer and seller can review independently | Art. 16.1, 27 | `services/api/internal/reviews` | Mutual Review Test |
| **REQ-REP-03** | Reputation | User Reputation Snapshot score calculation with manner metrics | Art. 16.3, 20.7 | `services/api/internal/reviews` | Reputation Calculation Test |
| **REQ-COM-01** | Community | Neighborhood Board: Question, Info, Recommendation, Lost & Found | Art. 17, 20.8, 28 | `services/api/internal/community` | Community Feed Test |
| **REQ-COM-02** | Community | Localized Permissions: Only verified neighborhood members can post | Art. 17.3, 28 | `services/api/internal/community` | Verification Guard Test |
| **REQ-COM-03** | Community | Threaded Comments, Replies & Likes | Art. 17.2, 20.8, 28 | `services/api/internal/community` | Thread & Like Test |
| **REQ-COM-04** | Community | Neighborhood Interest Groups & Membership | Art. 17.1, 20.8, 28 | `services/api/internal/community` | Group Test |
| **REQ-BIZ-01** | Business | Local Store Profile (Name, Category, Hours, Address, Media) | Art. 18, 20.9, 29 | `services/api/internal/business` | Store Profile Test |
| **REQ-BIZ-02** | Business | Business Verification workflow with administrative review | Art. 18.2, 29 | `services/api/internal/business` | Business Verify Test |
| **REQ-BIZ-03** | Business | Store Posts & Community Announcements | Art. 18.3, 20.9 | `services/api/internal/business` | Store Posts Test |
| **REQ-MOD-01** | Moderation | Universal Reporting: Listing, User, Chat, Community, Review | Art. 20.10, 31.14 | `services/api/internal/moderation`| Report Dispatch Test |
| **REQ-MOD-02** | Moderation | Automated Triage & Moderator Case Assignment | Art. 31.14 | `services/api/internal/moderation`| Case Workflow Test |
| **REQ-MOD-03** | Moderation | Enforcement Actions: Warning, Content Removal, Suspension | Art. 20.10, 31.14 | `services/api/internal/moderation`| Enforcement Test |
| **REQ-MOD-04** | Moderation | User Appeals System & Admin Review | Art. 20.10, 31.14 | `services/api/internal/moderation`| Appeals Test |
| **REQ-ADM-01** | Admin | Admin Dashboard Overview metrics (Users, Listings, Deals, Chats) | Art. 31.2, 32 | `apps/admin/src/pages/Dashboard` | Metric Aggregation Test |
| **REQ-ADM-02** | Admin | Country, City & Neighborhood Boundary Geo Manager | Art. 31.3-31.5, 32 | `apps/admin/src/pages/Geography` | Geo Admin Test |
| **REQ-ADM-03** | Admin | User & Session Management (Search, Suspension, Revocation) | Art. 31.7, 32 | `apps/admin/src/pages/Users` | User Ops Test |
| **REQ-ADM-04** | Admin | Listing Content Moderation (Approve, Take-down, Restore) | Art. 31.8, 32 | `apps/admin/src/pages/Listings` | Content Ops Test |
| **REQ-ADM-05** | Admin | Privacy-respecting Chat Safety Evidence Inspector | Art. 31.9 | `apps/admin/src/pages/ChatSafety`| Privacy-gated Review Test |
| **REQ-ADM-06** | Admin | Role-Based Access Control (RBAC) with granular permissions | Art. 33 | `services/api/internal/admin` | RBAC Auth Guard Test |
| **REQ-ADM-07** | Admin | Tamper-proof Audit Logging for all administrative actions | Art. 20.13, 31.16 | `services/api/internal/admin` | Audit Log Integrity Test |
| **REQ-NOT-01** | Notifications| Event-driven User Notifications (Chat, Offer, Deal, Review) | Art. 20.11, 30 | `services/api/internal/notifications`| Notification Dispatch Test |
| **REQ-NOT-02** | Notifications| Push Device token registration & user preference controls | Art. 20.11, 30 | `services/api/internal/notifications`| Device & Prefs Test |
| **REQ-GLO-01** | Globalization| Multi-language localization (KR, EN, JA, ZH-CN, ES, FR) | Art. 19.1 | `apps/web`, `apps/mobile` | i18n Render Test |
| **REQ-GLO-02** | Globalization| Multi-currency support (KRW, USD, JPY, EUR) & minor unit integers | Art. 19.3, 20.1 | `services/api/internal/listings` | Currency Format Test |
| **REQ-SEC-01** | Security | IDOR Protection: Strict resource ownership verification | Art. 38 | `services/api/internal/middleware`| IDOR Security Test |
| **REQ-SEC-02** | Security | Rate Limiting on Auth, Messaging, Search, and Verification | Art. 21, 38 | `services/api/internal/middleware`| Rate Limit Test |
| **REQ-SEC-03** | Security | Input sanitization, XSS, CSRF & SQL injection prevention | Art. 38 | `services/api/internal/middleware`| Sec Scan Test |
| **REQ-SEO-01** | SEO | Nuxt SSR for Public Listings, Neighborhoods & Categories | Art. 4.1, 37 | `apps/web` | SSR Meta & Schema Test |
| **REQ-SEO-02** | SEO | Sitemap generation, Robots.txt, OpenGraph & Structured JSON-LD | Art. 37 | `apps/web` | SEO Validation Test |
| **REQ-APP-01** | Mobile | Flutter Cross-Platform App (Android / iOS / Web) | Art. 4.2, 6.1 | `apps/mobile` | Flutter Analyze & Test |
| **REQ-APP-02** | Mobile | 5-Tab Navigation: HOME | COMMUNITY | SELL | CHATS | MY | Art. 6.1 | `apps/mobile` | Tab Navigation Test |
| **REQ-APP-03** | Mobile | Native Location Permissions, GPS Fix & Verification UI | Art. 8.1, 9.3 | `apps/mobile` | Location Service Test |
| **REQ-OPS-01** | Operations | Docker Compose stack for PostGIS, Redis, API, Worker, Web, Admin | Art. 4.6, 41 | `infrastructure/docker` | Container Health Test |

