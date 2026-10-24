# BN NEAR — 04. Information Architecture & Navigation

**Order ID:** BN-NEAR-MASTER-001  

---

## 1. Information Architecture Tree

```text
BN NEAR Platform
│
├── 1. Onboarding & Account
│   ├── Language & Country Selection
│   ├── Authentication (Login / Register / Verify Email)
│   ├── Neighborhood Selector & PostGIS GPS Verification
│   └── User Profile & Preferences
│
├── 2. Marketplace (Core Tab 1: Home)
│   ├── Neighborhood Selector Bar & Verification Pill
│   ├── Category Carousel & Search Bar
│   ├── Dynamic Feed (Nearby Listings, Filtered by Radius & Condition)
│   ├── Listing Detail
│   │   ├── Media Gallery (Cover + Gallery)
│   │   ├── Title, Condition, Price, Time Ago
│   │   ├── Seller Profile Card (Manner Temp, Badges, Join Date)
│   │   ├── Description & Meeting Area (Approximate)
│   │   └── Action Bar (Favorite, Make Offer, Chat Now)
│   └── Sell Form (Core Tab 3: Sell)
│       ├── Photo Picker (up to 10 images)
│       ├── Title, Category, Condition, Price, Negotiable Toggle
│       ├── Neighborhood Assignment
│       └── Publish & Review
│
├── 3. Real-Time Chat (Core Tab 4: Chats)
│   ├── Conversation List (Unread Badges, Last Message, Listing Thumbnail)
│   └── Conversation Room
│       ├── Top Context Bar (Listing Preview, Deal State, Status Action Buttons)
│       ├── Message Feed (Text, Image, Offer Card, Meeting Card)
│       ├── In-Chat Deal Actions (Make Offer, Accept/Decline, Reserve, Mark Sold)
│       └── Chat Composer (Text Input, Photo Attachment, Quick Replies)
│
├── 4. Community (Core Tab 2: Community)
│   ├── Neighborhood Switcher & Category Tabs (All, Questions, Life, Recommendations, Lost & Found)
│   ├── Post Feed & Post Detail
│   ├── Comment Section (Threaded Replies, Like Counter)
│   └── Write Post Modal
│
├── 5. Local Businesses (Directory accessible via Home & Community)
│   ├── Category Filter (Dining, Beauty, Services, Retail)
│   ├── Store Card & Store Profile Detail
│   ├── Operating Hours, Contact Info, Location Map
│   └── Store Announcement Posts
│
├── 6. My Page (Core Tab 5: My)
│   ├── Profile Header (Avatar, Nickname, Reputation Score)
│   ├── Activity Hub (Selling, Purchased, Watchlist / Favorites)
│   ├── Verification Center (Manage primary & secondary neighborhoods, GPS verify)
│   ├── Notifications Center
│   ├── Settings (Language, Currency, Push toggles, Blocked Users)
│   └── Account Security & Deletion
│
└── 7. Admin Operations (Dedicated Web Cockpit)
    ├── Operational Overview & Live KPI Analytics
    ├── Geography Management (Countries, Cities, Neighborhood Polygons)
    ├── User & Moderation Center (Reports, Appeals, Suspensions)
    ├── Content Triage (Listings, Business verification, Community posts)
    └── System Governance (Audit logs, Feature flags, Background jobs)
```

