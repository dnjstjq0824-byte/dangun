# BN NEAR — 31. Cloudflare Pages Automated Continuous Deployment Manual

**Order ID:** BN-NEAR-MASTER-001  
**Target:** Cloudflare Pages (Git Webhook & Wrangler CI/CD)  
**Repository:** `https://github.com/dnjstjq0824-byte/dangun.git`  
**Branch:** `main`

---

## 1. Overview & Cloudflare Integration Models

BN NEAR supports automated rebuilding and deployment on Cloudflare via two production mechanisms:

1. **Option A: Cloudflare Pages Native Git Integration (Zero-Config, Recommended)**
   - Cloudflare connects directly to `dnjstjq0824-byte/dangun` via GitHub App.
   - Every `git push` to `main` automatically triggers an edge build and global CDN rollout across Cloudflare's 330+ edge locations.
   - Built-in preview deployments for pull requests and branches.

2. **Option B: GitHub Actions + Wrangler Action (`ci/cloudflare-pages.yml`)**
   - Builds through GitHub Actions workflow and pushes compiled bundles directly to Cloudflare Pages via Wrangler.

---

## 2. Setting Up Cloudflare Pages Native Git Integration (Step-by-Step)

### A. Deploy Public Web App (Nuxt 3)

1. Log in to the [Cloudflare Dashboard](https://dash.cloudflare.com/).
2. Navigate to **Workers & Pages** -> **Create application** -> **Pages** tab.
3. Click **Connect to Git** and authorize your GitHub account (`dnjstjq0824-byte`).
4. Select the repository: `dnjstjq0824-byte/dangun`.
5. Configure the Build Settings:
   - **Project Name:** `bn-near-web`
   - **Production Branch:** `main`
   - **Framework Preset:** `Nuxt.js`
   - **Root Directory:** `apps/web`
   - **Build Command:** `pnpm build`
   - **Build Output Directory:** `dist`
6. Under **Environment variables (Advanced)**, add:
   - `NODE_VERSION`: `20`
   - `PNPM_VERSION`: `9`
   - `NITRO_PRESET`: `cloudflare-pages`
   - `API_BASE_URL`: `https://api.yourdomain.com/api/v1` (or your Go API URL)
   - `WS_BASE_URL`: `wss://api.yourdomain.com/api/v1/ws`
7. Click **Save and Deploy**.
8. **Result:** Any future `git push origin main` triggers an immediate automatic build and instant edge deployment!

---

### B. Deploy Admin Operations Cockpit (Vue 3 + Vite)

1. In Cloudflare Dashboard, go to **Workers & Pages** -> **Create application** -> **Pages** -> **Connect to Git**.
2. Select the repository: `dnjstjq0824-byte/dangun`.
3. Configure the Build Settings:
   - **Project Name:** `bn-near-admin`
   - **Production Branch:** `main`
   - **Framework Preset:** `Vite`
   - **Root Directory:** `apps/admin`
   - **Build Command:** `pnpm build`
   - **Build Output Directory:** `dist`
4. Under **Environment variables**, add:
   - `NODE_VERSION`: `20`
   - `PNPM_VERSION`: `9`
   - `VITE_API_BASE`: `https://api.yourdomain.com/api/v1`
5. Click **Save and Deploy**.
6. **Result:** Single Page Application (SPA) routing is preserved via `public/_redirects` (`/* /index.html 200`), and every commit triggers an automatic rebuild!

---

## 3. Cloudflare Pages Optimization Artifacts Included in Repo

| File | Purpose | Location |
|---|---|---|
| `apps/web/wrangler.toml` | Nuxt 3 Cloudflare Pages configuration & Node.js compat flags | [`apps/web/wrangler.toml`](file:///c:/Users/dnjst/Desktop/project/danggun/apps/web/wrangler.toml) |
| `apps/web/public/_headers` | Edge security headers & immutable static asset caching | [`apps/web/public/_headers`](file:///c:/Users/dnjst/Desktop/project/danggun/apps/web/public/_headers) |
| `apps/admin/wrangler.toml` | Vue 3 Admin Cloudflare Pages configuration | [`apps/admin/wrangler.toml`](file:///c:/Users/dnjst/Desktop/project/danggun/apps/admin/wrangler.toml) |
| `apps/admin/public/_redirects` | 200 SPA fallback for Vue Router deep links | [`apps/admin/public/_redirects`](file:///c:/Users/dnjst/Desktop/project/danggun/apps/admin/public/_redirects) |
| `apps/admin/public/_headers` | Security headers for Admin Cockpit | [`apps/admin/public/_headers`](file:///c:/Users/dnjst/Desktop/project/danggun/apps/admin/public/_headers) |
| `ci/cloudflare-pages.yml` | GitHub Actions workflow for Wrangler API deployment | [`ci/cloudflare-pages.yml`](file:///c:/Users/dnjst/Desktop/project/danggun/ci/cloudflare-pages.yml) |

