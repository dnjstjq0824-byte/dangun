# BN NEAR — 10. System Architecture & Topology

**Order ID:** BN-NEAR-MASTER-001  
**Architecture Style:** Modular Monolith (Services decomposed by domain boundaries, ready for horizontal microservice extraction)  

---

## 1. High-Level Architecture Topology

```mermaid
flowchart TD
    subgraph Clients["Clients Layer"]
        Mobile["Flutter Mobile App\n(Android / iOS / Web)"]
        PublicWeb["Nuxt 3 Public Web\n(SSR + Pinia)"]
        AdminWeb["Vue 3 Admin Cockpit\n(Vite + TypeScript)"]
    end

    subgraph Gateway["Edge & Gateway Layer"]
        Traefik["Reverse Proxy / TLS / CDN"]
        APIGateway["Go API Gateway\n(/api/v1)"]
        WSGateway["Go WebSocket Hub\n(/api/v1/ws)"]
    end

    subgraph CoreServices["Go Domain Engine"]
        IdentitySvc["Identity & Auth Service"]
        GeoSvc["Geography & Location Engine"]
        VerifySvc["Location Verification Engine"]
        ListingSvc["Marketplace & Search Service"]
        ChatSvc["Chat & Delivery Hub"]
        DealSvc["Deals & Reservation FSM"]
        RepSvc["Reputation & Reviews Engine"]
        CommSvc["Community & Groups Service"]
        BizSvc["Local Business Service"]
        ModSvc["Trust, Safety & Moderation"]
        NotifSvc["Event & Notification Dispatcher"]
        AdminSvc["Admin & Audit Logging"]
    end

    subgraph DataStorage["Data & State Layer"]
        PG[("PostgreSQL 16 + PostGIS 3.4\n(ACID, Spatial Indexes, FTS)")]
        RedisCache[("Redis 7\n(Pub/Sub, Rate Limits, Cache)")]
        ObjectStore[("S3 Compatible Object Storage\n(Media Assets, Thumbnails)")]
    end

    subgraph Workers["Asynchronous Workers"]
        MediaWorker["Media Sanitizer & EXIF Stripper"]
        SearchIndexer["Geo & Text Index Reindexer"]
        CleanupJob["Verification Expiration & Temp Sweeper"]
    end

    Clients --> Traefik
    Traefik --> APIGateway
    Traefik --> WSGateway

    APIGateway --> CoreServices
    WSGateway --> ChatSvc
    ChatSvc <--> RedisCache

    CoreServices --> PG
    CoreServices --> RedisCache
    CoreServices --> ObjectStore

    CoreServices -. Outbox Events .-> Workers
    Workers --> PG
    Workers --> ObjectStore
```

---

## 2. Key Architecture Decisions
1. **Unified Storage Contract:** The core data layer leverages PostgreSQL 16 with PostGIS spatial extensions for exact boundary queries (`ST_Contains`, `ST_DWithin`). In testing environments or autonomous single-binary runs, the repository layer provides an interchangeable high-speed geospatial in-memory storage engine adhering to identical interfaces.
2. **WebSocket Sequence Guarantee:** Chat messages are assigned sequential monotonically increasing sequence numbers per conversation within a database transaction. Clients confirm receipt with ACKs, and reconnection protocols synchronize missed sequence ranges seamlessly.
3. **Reservation Row-Level Locking:** To eliminate double-booking or race conditions when multiple buyers attempt to reserve an item, the Deal Engine uses optimistic versioning (`version BIGINT`) and strict database transactions.
4. **Zero-Knowledge GPS Policy:** The system strictly separates internal raw coordinates (processed transiently for verification calculations and immediately discarded) from public catalog representations (which only store neighborhood IDs or randomized centroid offsets).
