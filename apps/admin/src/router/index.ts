import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from '../pages/Dashboard.vue'
import Users from '../pages/Users.vue'
import Listings from '../pages/Listings.vue'
import Geography from '../pages/Geography.vue'
import ChatSafety from '../pages/ChatSafety.vue'
import Moderation from '../pages/Moderation.vue'
import Businesses from '../pages/Businesses.vue'
import AuditLogs from '../pages/AuditLogs.vue'

const routes = [
  { path: '/', name: 'Dashboard', component: Dashboard },
  { path: '/users', name: 'Users', component: Users },
  { path: '/listings', name: 'Listings', component: Listings },
  { path: '/geography', name: 'Geography', component: Geography },
  { path: '/chat-safety', name: 'ChatSafety', component: ChatSafety },
  { path: '/moderation', name: 'Moderation', component: Moderation },
  { path: '/businesses', name: 'Businesses', component: Businesses },
  { path: '/audit-logs', name: 'AuditLogs', component: AuditLogs }
]

export const router = createRouter({
  history: createWebHistory(),
  routes
})

