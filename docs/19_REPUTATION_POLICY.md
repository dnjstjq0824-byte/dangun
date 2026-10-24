# BN NEAR — 19. Reputation & Manner Score Policy

**Order ID:** BN-NEAR-MASTER-001  
**Module:** `services/api/internal/reviews`  

---

## 1. Manner Temperature Formula
To build an intuitive, empathetic representation of reliability, each user possesses a "Manner Temperature" initialized at standard human body temperature: **36.5°C**.

### Calculation Model:
$$ T_{\text{current}} = 36.5 + \sum_{i=1}^{N} \Delta T_i $$

Where for each confirmed transaction review $i$:
- **5-Star Rating (Excellent):** $+0.5^\circ\text{C}$ (plus bonus $+0.1^\circ\text{C}$ for positive manner tags such as "Punctual", "Kind and responsive").
- **4-Star Rating (Good):** $+0.2^\circ\text{C}$.
- **3-Star Rating (Neutral):** $0.0^\circ\text{C}$.
- **2-Star Rating (Disappointing):** $-0.5^\circ\text{C}$.
- **1-Star Rating (Bad / Misleading):** $-1.5^\circ\text{C}$.
- **Substantiated Moderation Sanction:** $-3.0^\circ\text{C}$ to $-10.0^\circ\text{C}$.

### Boundaries:
- Minimum Score: **0.0°C**
- Maximum Score: **99.9°C**

---

## 2. Anti-Gaming Constraints
1. **Verified Deal Requirement:** Reviews can NEVER be posted arbitrarily. A review requires an explicit `deal_id` in status `COMPLETED_CONFIRMED` in which the reviewer is either `buyer_user_id` or `seller_user_id`.
2. **One Review per Party per Deal:** A strict unique constraint `UNIQUE (deal_id, reviewer_user_id)` prevents repetitive rating stuffing.
3. **Collusion Dampening:** If the same pair of users completes multiple transactions within 30 days, subsequent reviews have diminishing marginal temperature weight (100% -> 50% -> 10% -> 0%).
4. **No Artificial Character Defamation:** The temperature score is explicitly labeled as an in-platform community manner metric and not a legal or personal character judgment.
