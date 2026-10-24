# BN NEAR — 15. Location Verification Engine Specification

**Order ID:** BN-NEAR-MASTER-001  
**Module:** `services/api/internal/verification`  

---

## 1. Objectives & Guarantees
The Location Verification Engine confirms that a user is physically present within a declared neighborhood boundary before granting neighborhood-exclusive privileges (selling items, writing community posts, claiming verified resident status).

### Privacy Mandate:
- **Zero Raw GPS Exposure:** The user's exact latitude and longitude are transmitted only over TLS to the verification endpoint.
- **Ephemeral Calculation:** Coordinates are used exclusively in-memory for boundary containment queries and distance-to-edge calculations. They are NEVER stored in public tables or user profiles.
- **Audit Trails:** Only the verification outcome (`VERIFIED`, `OUT_OF_BOUNDS`, `EXCESSIVE_ERROR`), accuracy meters, and timestamp are preserved in `location_verification_attempts`.

---

## 2. Mathematical Boundary Verification Algorithm
```text
Inputs:
  target_neighborhood_id (UUID)
  client_lat (float64)
  client_lng (float64)
  accuracy_meters (float64)
  timestamp (time.Time)

Process:
  1. Validate client timestamp freshness (|now - timestamp| <= 5 minutes).
  2. Validate reported GPS accuracy (accuracy_meters <= 150.0m). If accuracy > 150m, reject with ERROR_LOW_ACCURACY.
  3. Retrieve neighborhood boundary polygon:
     PostGIS: ST_Contains(boundary_geom, ST_SetSRID(ST_MakePoint(lng, lat), 4326))
     In-Memory Fallback: Ray-Casting algorithm on polygon coordinates with bounding-box pre-filtering.
  4. If Point is inside Polygon:
     -> Verification Status: VERIFIED
     -> Set expires_at: now + 30 days
     -> Update user_neighborhoods record
  5. If Point is outside Polygon:
     -> Compute minimum distance to boundary: ST_Distance(boundary_geom, ST_SetSRID(ST_MakePoint(lng, lat), 4326))
     -> If distance <= accuracy_meters: Edge tolerance grant or mark for review
     -> Otherwise: Reject with ERROR_OUT_OF_BOUNDS
```

---

## 3. Quota & Fraud Prevention Rules
- A user may register a maximum of **2 active neighborhoods** (e.g. Home and Workplace).
- Consecutive failed attempts are limited to **5 attempts per hour** per user IP/account.
- Re-verification is required every **30 days**.

