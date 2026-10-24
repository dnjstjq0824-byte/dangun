# BN NEAR — 07. Flutter Mobile Screen Specifications

**Order ID:** BN-NEAR-MASTER-001  
**Target Runtimes:** Android, iOS, Flutter Web  
**Primary Framework:** Flutter 3.47+ / Dart 3.13+ with Riverpod & GoRouter  

---

## 1. Bottom Navigation Shell (5 Core Tabs)
The bottom navigation bar persists across all top-level destinations:
- **Tab 1: HOME (`/home`)**
  - Header: Neighborhood selector dropdown with verification badge (`Verified / Unverified`), notification bell, search icon.
  - Category horizontal chips (`All`, `Digital`, `Appliances`, `Furniture`, `Fashion`, `Baby`, `Sports`, `Books`, etc.).
  - Refreshable ListView with product cards (cover thumbnail, title, price, neighborhood, elapsed time, favorites, chat count).
  - Floating Action Button (FAB): "Sell Item" leading to `/listings/create`.
- **Tab 2: COMMUNITY (`/community`)**
  - Header: Neighborhood name, topic filter chips (`All`, `Questions`, `Recommendations`, `Daily`, `Lost & Found`).
  - Feed list of posts with author avatar, verified badge, title, excerpt, image thumbnail, like counter, and comment counter.
  - FAB: "Write Post" leading to `/community/create`.
- **Tab 3: SELL (`/listings/create`)**
  - Direct listing creation flow: Photo picker (up to 10 photos, drag to reorder, set cover), title, category picker modal, item condition selector, price input with currency formatting and "Free" checkbox, description markdown-ready editor, neighborhood selector.
- **Tab 4: CHATS (`/chats`)**
  - Conversation list: Counterpart avatar, counterpart nickname, listing miniature preview, last message snippet, timestamp, unread message count badge.
  - Swipe to leave / mute.
- **Tab 5: MY PAGE (`/my`)**
  - User profile card: Avatar, display name, Manner Temperature bar (e.g. 38.2°C), neighborhood verification status.
  - Quick action grids: "My Listings (Selling/Reserved/Sold)", "Purchase History", "Saved / Favorites", "Community Posts".
  - Management items: "Neighborhood Verification Settings", "Notification Preferences", "Blocked Users", "Help & Safety Guidelines", "Account Settings & Logout".

---

## 2. Secondary & Modal Screens
- **Listing Detail Screen (`/listings/:id`)**:
  - Image carousel with page indicators.
  - Seller profile ribbon with manner temperature and link to seller's other items.
  - Item details: Title, category, condition tag, posted date, full description.
  - Approximate public meeting area map preview (without private coordinates).
  - Sticky bottom action bar: Favorite toggle, price display, "Make Offer" button, "Chat" button.
- **Chat Room Screen (`/chats/:id`)**:
  - Top context card: Miniature listing photo, price, deal status badge (`ACTIVE`, `RESERVED`, `SOLD`), and contextual action buttons ("Reserve", "Mark Sold", "Confirm Purchase").
  - Reverse ListView for real-time messages with date dividers.
  - Message bubble: Text, time, read status checkmarks, attachment preview.
  - Deal action cards inside chat: Offer proposal card (with Accept/Decline actions), Meeting place recommendation card.
  - Bottom composer: Attachment plus button, text input, send icon.
- **Location Verification Screen (`/verification`)**:
  - Interactive map illustrating the target neighborhood boundary polygon.
  - GPS accuracy indicator and "Verify My Location Now" button.
  - Success feedback dialog with granted validity duration (30 days).
- **Review Submission Modal (`/deals/:id/review`)**:
  - Star rating picker (1 to 5 stars).
  - Quick positive / constructive tags ("Punctual", "Friendly", "Great condition", etc.).
  - Optional written comment.
