# BN NEAR — 22. Security Architecture & Threat Model

**Order ID:** BN-NEAR-MASTER-001  
**Standards:** OWASP Top 10, ASVS Level 2, MASVS  

---

## 1. Threat Mitigation Strategies

| Threat Vector | Mitigation Strategy | Implementation Point |
|---|---|---|
| **IDOR (Insecure Direct Object Reference)** | Enforce resource ownership checks in repository layer; never trust client-supplied owner IDs | `services/api/internal/middleware` |
| **SQL / Spatial Injection** | Parameterized queries with prepared statements via `database/sql` / `pgx` | All database repositories |
| **Cross-Site Scripting (XSS)** | HTML escaping & input sanitization; strict Content Security Policy (CSP) headers | `apps/web`, `apps/admin` |
| **Broken Authentication** | Bcrypt / Argon2 password hashing; cryptographically random JWT refresh tokens | `services/api/internal/identity` |
| **WebSocket Hijacking** | JWT bearer validation during upgrade; origin verification; per-message conversation ACL | `services/api/internal/chat` |
| **Brute Force & Flooding** | IP and user-based token bucket rate limiting on auth, messaging, and search | `services/api/internal/middleware` |
| **Malicious Media Uploads** | File magic byte validation (JPEG, PNG, WebP only); EXIF metadata stripping; size limits | `services/worker/internal/media` |
| **Tampered Location Fixes** | GPS accuracy thresholds, timestamp validity windows, velocity/teleportation heuristics | `services/api/internal/verification` |
| **Admin Privilege Escalation** | Fine-grained RBAC matrix; immutable append-only audit logging of all admin mutations | `services/api/internal/admin` |

