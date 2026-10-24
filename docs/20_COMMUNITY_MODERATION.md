# BN NEAR — 20. Community Moderation & Safety Policy

**Order ID:** BN-NEAR-MASTER-001  
**Module:** `services/api/internal/moderation`  

---

## 1. Universal Reporting System
Users can flag inappropriate content or behavior from any platform surface:
- **Listing Report:** Counterfeit items, prohibited weapons/drugs, fraud, price abuse.
- **User Profile Report:** Impersonation, harassment, scam behavior.
- **Chat Conversation Report:** Toxic language, stalking, fraudulent off-platform payment links.
- **Community Post/Comment Report:** Hate speech, doxxing / personal privacy leaks, spam advertising.

---

## 2. Moderation Workflow
```mermaid
flowchart TD
    ReportReceived["User Submits Report\n(Target, Reason, Evidence)"] --> Triage["Automated Triage Engine\n(Severity & Rate-check)"]
    Triage --> CaseCreated["Moderation Case Created\n(Priority: P0-P3)"]
    CaseCreated --> Assign["Assigned to Operations Moderator"]
    Assign --> EvidenceReview["Review Evidence\n(Privacy-redacted view)"]
    EvidenceReview --> Decision{"Moderation Decision"}
    Decision -->|No Violation| Dismiss["Case Dismissed\n(Reason recorded in audit)"]
    Decision -->|Minor Issue| Warning["Official In-App Warning"]
    Decision -->|Content Violation| TakeDown["Content Removed\n(Notify Author)"]
    Decision -->|Severe Malice| Suspension["Temporary / Permanent Account Suspension"]
    Suspension --> Notice["User Notice with Sanction Reason"]
    Notice --> Appeal{"User Submits Appeal?"}
    Appeal -->|Yes| SecondaryReview["Senior Admin Review"]
    Appeal -->|No| Closed["Sanction Enforced"]
```

---

## 3. Privacy-Preserving Chat Evidence Review
Moderators CANNOT browse private chats freely. When a user reports a chat, the client submits a snapshot containing up to 10 surrounding messages as cryptographic evidence. This snapshot is sealed into the moderation case record and accessible only to authorized moderators with explicit audit logging of every view action.
