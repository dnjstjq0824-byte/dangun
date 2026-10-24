# BN NEAR — 06. Global Design System Specification

**Order ID:** BN-NEAR-MASTER-001  
**Brand Philosophy:** LOCAL · FRIENDLY · MODERN · FAST · TRUSTWORTHY  

---

## 1. Color Palette

```text
/* Primary Brand Identity */
--color-brand-primary:      #0F766E; /* Deep Hyperlocal Teal */
--color-brand-primary-hover:#115E59; /* Forest Teal */
--color-brand-accent:       #2563EB; /* Trust Royal Blue */

/* Neutral Surfaces */
--color-bg-base:            #F8FAFC; /* Slate 50 */
--color-surface-card:       #FFFFFF; /* Pure White */
--color-surface-subtle:     #F1F5F9; /* Slate 100 */
--color-border-default:     #E2E8F0; /* Slate 200 */
--color-border-strong:      #CBD5E1; /* Slate 300 */

/* Typography Colors */
--color-text-primary:       #0F172A; /* Slate 900 */
--color-text-secondary:     #475569; /* Slate 600 */
--color-text-muted:         #94A3B8; /* Slate 400 */
--color-text-on-primary:    #FFFFFF;

/* Status Indicators */
--color-status-success:     #16A34A; /* Green 600 */
--color-status-warning:     #D97706; /* Amber 600 */
--color-status-danger:      #DC2626; /* Red 600 */
--color-status-reserved:    #EA580C; /* Orange 600 */
--color-status-sold:        #64748B; /* Slate 500 */
```

---

## 2. Typography Scale
- **Display 1:** 32px / 40px Bold (Page titles, hero sections)
- **Title 1:** 24px / 32px SemiBold (Section headers, product titles)
- **Title 2:** 20px / 28px SemiBold (Modal headers, card titles)
- **Body Large:** 16px / 24px Regular/Medium (Primary body, listings, chat bubbles)
- **Body Regular:** 14px / 20px Regular/Medium (Secondary labels, metadata, buttons)
- **Caption:** 12px / 16px Regular (Timestamps, counts, badges)

---

## 3. Elevation & Radii
- **Border Radius Small (`rounded-sm`):** 4px (Badges, tags)
- **Border Radius Medium (`rounded-md`):** 8px (Input fields, buttons)
- **Border Radius Large (`rounded-lg`):** 12px (Cards, modal dialogs)
- **Border Radius Full (`rounded-full`):** 9999px (Pills, avatar, floating buttons)
- **Shadows:** Subtly elevated using `box-shadow: 0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1)`.

---

## 4. Key Atomic Components
1. **Neighborhood Verification Pill:** Displays current verified neighborhood with a green shield icon; if unverified or expired, shows an amber prompt to verify.
2. **Product Feed Card:** 1:1 aspect ratio rounded image, title, condition badge, price in bold, neighborhood name, time-ago, favorite heart icon, and chat count counter.
3. **Manner Temperature Badge:** Visual gauge displaying the user's reputation score (e.g. 36.5°C base standard, color-coded from slate to energetic emerald).
4. **Chat Bubble Suite:** Alternating right-aligned (sender, primary teal background with white text) and left-aligned (receiver, subtle slate background with dark text) bubbles with message sequence and read status markers.
5. **Deal Action Banner:** Sticky bottom sheet in chat facilitating "Make Offer", "Reserve", "Mark Sold", and "Write Review".

