# BN NEAR — 26. Deployment & Infrastructure Manual

**Order ID:** BN-NEAR-MASTER-001  
**Stack:** Docker Compose / Kubernetes / Linux AMD64  

---

## 1. Containerized Service Topology

```text
infrastructure/docker/
  ├── docker-compose.yml       (Orchestration for API, Worker, Web, Admin, Postgres, Redis)
  ├── Dockerfile.api           (Go Multi-stage binary builder)
  ├── Dockerfile.worker        (Go Worker binary builder)
  ├── Dockerfile.web           (Node 24 Nuxt SSR production container)
  └── Dockerfile.admin         (Nginx container serving Vue 3 Vite SPA)
```

---

## 2. Environment Configuration
Key environment variables:
- `PORT`: HTTP listener port (default `8080`)
- `DATABASE_URL`: PostgreSQL connection string (`postgres://bn_near:secret@postgres:5432/bn_near_db?sslmode=disable`)
- `REDIS_URL`: Redis connection string (`redis://redis:6379/0`)
- `JWT_SECRET`: 256-bit cryptographically random signing key
- `CORS_ALLOWED_ORIGINS`: Comma-separated trusted origins
- `STORAGE_DRIVER`: `local` or `s3`
- `GEO_VERIFICATION_ACCURACY_THRESHOLD`: Maximum allowed GPS uncertainty (default `150` meters)

---

## 3. Quick Start (Development & Production)
```bash
# 1. Start all container services
docker-compose -f infrastructure/docker/docker-compose.yml up -d --build

# 2. Verify container health
docker-compose -f infrastructure/docker/docker-compose.yml ps

# 3. Check API readiness
curl http://localhost:8080/health
```

