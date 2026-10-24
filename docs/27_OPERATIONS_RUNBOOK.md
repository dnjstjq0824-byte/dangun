# BN NEAR — 27. Operations & Incident Runbook

**Order ID:** BN-NEAR-MASTER-001  

---

## 1. Operational Health Checks
- **Liveness Probe:** `GET /health` -> Returns `200 OK` with uptime metadata.
- **Readiness Probe:** `GET /ready` -> Validates active connection pool to PostgreSQL and Redis.

---

## 2. Incident Response Playbooks

### Incident 1: WebSocket Connection Storm / Reconnection Floods
- **Symptom:** High spike in `/api/v1/ws` handshake requests and CPU load.
- **Action:**
  1. Inspect network gateway logs for client flapping.
  2. Ensure exponential backoff with jitter is respected by client applications.
  3. Verify Redis Pub/Sub memory utilization.

### Incident 2: Spam / Mass Content Influx
- **Symptom:** Rapid listing creation or message flood from single IP or user ring.
- **Action:**
  1. Check rate limit counters in Redis.
  2. Use Admin Console (`/admin/users`) to enact temporary account restriction.
  3. Batch-takedown affected listings via moderation queue (`/admin/listings`).

### Incident 3: Boundary Polygon Data Corruption
- **Symptom:** High failure rate of location verifications in specific city.
- **Action:**
  1. Inspect `location_verification_attempts` for error reason codes.
  2. Verify polygon geometry using GeoJSON validator in Admin Console.
  3. Re-import official municipal boundary GeoJSON.
