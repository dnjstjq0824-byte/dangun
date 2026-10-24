# BN NEAR — 24. Role-Based Access Control (RBAC) Matrix

**Order ID:** BN-NEAR-MASTER-001  
**Module:** `services/api/internal/admin`  

---

## 1. Operational Roles Definition
- **SUPER_ADMIN:** Full unconstrained system administration and role assignment.
- **OPERATIONS_ADMIN:** Platform configuration, feature flags, background jobs.
- **GEO_ADMIN:** Country, city, and neighborhood boundary polygon management.
- **TRUST_SAFETY_ADMIN:** High-severity user sanctions, suspension appeals, chat safety evidence inspection.
- **CONTENT_MODERATOR:** Review flagged listings, community posts, comments, take-down actions.
- **BUSINESS_ADMIN:** Verification of merchant credentials and local business profiles.
- **SUPPORT_AGENT:** Help desk tickets and customer inquiries.
- **AUDITOR:** Read-only access to audit logs and system telemetry.

---

## 2. Granular Permissions Matrix

| Permission Code | SUPER_ADMIN | OPERATIONS_ADMIN | GEO_ADMIN | TRUST_SAFETY | CONTENT_MOD | BUSINESS_ADMIN | AUDITOR |
|---|:---:|:---:|:---:|:---:|:---:|:---:|:---:|
| `users.read` | [x] | [x] | [ ] | [x] | [x] | [ ] | [x] |
| `users.suspend` | [x] | [ ] | [ ] | [x] | [ ] | [ ] | [ ] |
| `geo.manage` | [x] | [x] | [x] | [ ] | [ ] | [ ] | [ ] |
| `listings.moderate` | [x] | [ ] | [ ] | [x] | [x] | [ ] | [ ] |
| `chat.inspect_evidence` | [x] | [ ] | [ ] | [x] | [ ] | [ ] | [ ] |
| `business.verify` | [x] | [ ] | [ ] | [ ] | [ ] | [x] | [ ] |
| `audit.read` | [x] | [x] | [ ] | [x] | [ ] | [ ] | [x] |
| `system.flags` | [x] | [x] | [ ] | [ ] | [ ] | [ ] | [ ] |
