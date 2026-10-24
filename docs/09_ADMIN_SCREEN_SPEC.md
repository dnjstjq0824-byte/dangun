# BN NEAR — 09. Vue Admin Screen Specifications

**Order ID:** BN-NEAR-MASTER-001  
**Framework:** Vue 3 / Vite / TypeScript / Pinia / Tailwind  
**Role:** Operational Cockpit for Platform Governance, Moderation & Geo Management  

---

## 1. Master Layout & Navigation Shell
Left navigation sidebar organized by functional operational domain:
1. **Overview & Analytics:**
   - Real-time KPI counters: Active Users, Verified Neighborhoods, Active Listings, Completed Deals, Live WebSocket Sessions.
   - Geographical adoption chart (growth by country and city).
2. **Geographical Management:**
   - Country & Currency Config: Launch status (`ACTIVE`, `BETA`, `READ_ONLY`, `NOT_LAUNCHED`), default currency, timezone.
   - City & Boundary Editor: GeoJSON import/export for neighborhood polygon boundaries and centroid management.
3. **Location Verification Operations:**
   - Verification attempt audit log, failure rate diagnostics, manual review queue for edge-case boundary appeals.
4. **User & Identity Governance:**
   - User search by ID/email/nickname, account status flags, session force-revocation, restriction history.
5. **Marketplace & Listings Moderation:**
   - Listing queue: Pending review, flagged items, prohibited goods classifier results.
   - One-click take-down with canned explanation notice to seller.
6. **Chat Safety Evidence Inspector:**
   - Privacy-gated inspection tool: Accessible ONLY upon user report or legal escalation with mandatory audit log reason.
   - Shows reported message snippet with sender, receiver, timestamp, and report reason.
7. **Trust & Moderation Cases:**
   - Triage workbench for reported users, fraudulent deal attempts, and review manipulation.
   - Case resolution workflow: Assign moderator, record decision (Dismiss, Warn, Suspend, Ban), record reasoning.
8. **Appeals Management:**
   - Secondary review queue for sanctioned users to submit evidence and appeal moderation decisions.
9. **Business Profile Verification:**
   - Review merchant proof documents, verify physical store location, approve "Verified Merchant" badge.
10. **System Governance & Audit Logs:**
    - Tamper-proof append-only audit log: Actor, Action, Resource, Timestamp, IP, Redacted Metadata.
    - Background job monitor and Feature flag toggle matrix.
