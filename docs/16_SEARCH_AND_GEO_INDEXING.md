# BN NEAR — 16. Search & Geospatial Indexing Specification

**Order ID:** BN-NEAR-MASTER-001  
**Module:** `services/api/internal/search`  

---

## 1. Geospatial Distance Indexing
Listings in BN NEAR are discovered based on physical proximity rather than arbitrary postal codes.

### PostGIS Query Implementation:
```sql
SELECT 
    l.id, l.title, l.price_minor, l.currency_code, l.condition, l.listing_status,
    l.neighborhood_id, n.name AS neighborhood_name,
    ST_Distance(l.public_location_geom::geography, ST_SetSRID(ST_MakePoint(:target_lng, :target_lat), 4326)::geography) AS distance_meters
FROM listings l
JOIN neighborhoods n ON l.neighborhood_id = n.id
WHERE 
    l.listing_status IN ('ACTIVE', 'RESERVED')
    AND ST_DWithin(l.public_location_geom::geography, ST_SetSRID(ST_MakePoint(:target_lng, :target_lat), 4326)::geography, :radius_meters)
ORDER BY distance_meters ASC, l.published_at DESC
LIMIT :limit OFFSET :offset;
```

### In-Memory Equivalence (Haversine & Bounding Box):
In test suites and lightweight deployments:
```go
func HaversineDistanceMeters(lat1, lon1, lat2, lon2 float64) float64 {
    const R = 6371000.0 // Earth radius in meters
    dLat := (lat2 - lat1) * (math.Pi / 180.0)
    dLon := (lon2 - lon1) * (math.Pi / 180.0)
    a := math.Sin(dLat/2)*math.Sin(dLat/2) +
         math.Cos(lat1*(math.Pi/180.0))*math.Cos(lat2*(math.Pi/180.0))*
         math.Sin(dLon/2)*math.Sin(dLon/2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    return R * c
}
```

---

## 2. Text Search & Ranking
1. **Title & Description Tokenization:** Multi-language tokenization with trigram matching (`pg_trgm`) for fuzzy spelling matching.
2. **Scoring Formula:**
   $$ \text{Score} = w_1 \cdot \text{TextMatch} + w_2 \cdot \frac{1}{1 + \text{Distance}} + w_3 \cdot \frac{1}{1 + \text{AgeHours}} $$
3. **Status Filter:** Only items with status `ACTIVE` or `RESERVED` are returned in public searches. `SOLD` or `ARCHIVED` items are excluded unless queried directly by participants.
