# BN NEAR — 01. Product Requirements Document

**Order ID:** BN-NEAR-MASTER-001  
**Project:** BN NEAR — Global Hyperlocal Marketplace & Community Platform  
**Owner:** Booking Network co., Ltd — Chief Product Owner  
**Status:** Approved for Implementation  

---

## 1. Executive Summary & Vision
BN NEAR is an independent, global hyperlocal marketplace and community platform engineered to connect neighbors worldwide. Users can safely discover, purchase, and sell second-hand goods, chat in real-time, arrange transactions, write mutual reviews, share local community news, and discover verified local businesses.

### Core Value Pillars
1. **Hyperlocal Proximity:** Prioritize items and stories within walking and short transit distance.
2. **Speed & Clarity:** Fast listing creation, instant messaging, and frictionless deal reservations.
3. **Mutual Trust & Safety:** Cryptographically signed sessions, location boundary verifications, privacy-preserving geo approximations, and verified transaction reviews.
4. **Community Cohesion:** Local neighborhood boards for questions, recommendations, and lost & found.
5. **Global Extensibility:** Multi-country, multi-currency, and multi-language support with region-specific launch states.

---

## 2. Competitive Differentiation
Unlike legacy classified boards or single-market apps:
- **No Clone Policy:** Original visual language (`#0F766E` deep teal, modern slate surface, distinct card hierarchy).
- **Location Privacy First:** Exact GPS coordinates are never exposed to other users or leaked in public APIs. Neighborhoods are represented by validated administrative boundaries and approximate public centroids.
- **Two-Phase Deal Lifecycle:** Seller mark-sold and buyer confirmation state machine prevents ghost reviews and reputation gaming.
- **Enterprise-Grade Admin & Safety:** Audited moderation triage, appeals workflows, and zero unmoderated surveillance.

---

## 3. High-Level Feature Modules
1. **Identity & Auth:** Email verification, secure sessions, device tracking, account deletion.
2. **Geographical Engine:** Hierarchical Geo (Country, Region, City, Neighborhood), PostGIS polygon checks.
3. **Location Verification:** Accuracy-weighted GPS boundary validation, retry budgets, anti-tamper heuristics.
4. **Marketplace Engine:** Multi-image uploads, category classification, condition grading, price negotiation flags, geospatial radius search, favorites.
5. **Real-time Chat:** WebSocket duplex communication, sequence number ordering, client message deduplication, delivery receipts, offline push queues.
6. **Deals & Negotiations:** Price offers, reservation locks, safe meeting place coordinator, mutual completion.
7. **Reputation Engine:** Verified deal reviews, multi-criteria manner ratings, anti-manipulation calculation.
8. **Community Boards:** Neighborhood-scoped discussions, threaded comments, upvotes, reporting.
9. **Local Business Directory:** Verified merchant profiles, operating hours, local announcements.
10. **Trust & Moderation:** Report triage, case escalation, enforcement actions, user appeals.
11. **Admin Operations Console:** Web-based operational cockpit with role-based access control and tamper-proof audit trails.
