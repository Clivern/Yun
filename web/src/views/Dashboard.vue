<template>
  <div class="min-h-screen bg-notion-bg">
    <!-- Navigation -->
    <nav class="bg-white border-b border-notion-border">
      <div class="max-w-7xl mx-auto px-6 lg:px-8">
        <div class="flex justify-between h-14">
          <div class="flex items-center">
            <div class="flex-shrink-0 flex items-center">
              <img src="/logo.png" alt="Yun Logo" class="h-8 w-auto">
            </div>
            <div class="hidden md:ml-8 md:flex md:space-x-1">
              <router-link
                to="/dashboard"
                class="bg-notion-hover text-notion-text inline-flex items-center px-3 py-1.5 rounded-md text-sm font-medium"
              >
                Dashboard
              </router-link>
              <router-link
                to="/gateways"
                class="text-notion-textLight hover:bg-notion-hover hover:text-notion-text inline-flex items-center px-3 py-1.5 rounded-md text-sm font-medium transition-colors"
              >
                Gateways
              </router-link>
              <router-link
                to="/mcps"
                class="text-notion-textLight hover:bg-notion-hover hover:text-notion-text inline-flex items-center px-3 py-1.5 rounded-md text-sm font-medium transition-colors"
              >
                MCPs
              </router-link>
              <router-link
                to="/servers"
                class="text-notion-textLight hover:bg-notion-hover hover:text-notion-text inline-flex items-center px-3 py-1.5 rounded-md text-sm font-medium transition-colors"
              >
                Servers
              </router-link>
              <router-link
                to="/tools"
                class="text-notion-textLight hover:bg-notion-hover hover:text-notion-text inline-flex items-center px-3 py-1.5 rounded-md text-sm font-medium transition-colors"
              >
                Tools
              </router-link>
              <router-link
                to="/users"
                class="text-notion-textLight hover:bg-notion-hover hover:text-notion-text inline-flex items-center px-3 py-1.5 rounded-md text-sm font-medium transition-colors"
              >
                Users
              </router-link>
              <router-link
                to="/settings"
                class="text-notion-textLight hover:bg-notion-hover hover:text-notion-text inline-flex items-center px-3 py-1.5 rounded-md text-sm font-medium transition-colors"
              >
                Settings
              </router-link>
            </div>
          </div>
          <div class="flex items-center">
            <div class="flex items-center space-x-3">
              <button
                @click="handleLogout"
                class="btn-secondary text-sm"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </div>
    </nav>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto py-8 px-6 lg:px-8">
      <!-- Page Header -->
      <div class="mb-8">
        <h1 class="text-3xl font-semibold text-notion-text">Dashboard</h1>
      </div>

      <!-- Stats Grid -->
      <div class="mb-6">
        <div class="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-4">
          <!-- Stat Card 1 -->
          <div class="bg-white rounded-lg border border-notion-border p-5 hover:shadow-sm transition-shadow">
            <p class="text-xs font-medium text-notion-textLight mb-2 uppercase tracking-wide">Status</p>
            <p class="text-2xl font-semibold text-notion-text">{{ stats.status }}</p>
          </div>

          <!-- Stat Card 2 -->
          <div class="bg-white rounded-lg border border-notion-border p-5 hover:shadow-sm transition-shadow">
            <p class="text-xs font-medium text-notion-textLight mb-2 uppercase tracking-wide">Total Requests</p>
            <p class="text-2xl font-semibold text-notion-text">{{ stats.totalRequests.toLocaleString() }}</p>
          </div>

          <!-- Stat Card 3 -->
          <div class="bg-white rounded-lg border border-notion-border p-5 hover:shadow-sm transition-shadow">
            <p class="text-xs font-medium text-notion-textLight mb-2 uppercase tracking-wide">Avg Response</p>
            <p class="text-2xl font-semibold text-notion-text">{{ stats.avgResponse }}ms</p>
          </div>

          <!-- Stat Card 4 -->
          <div class="bg-white rounded-lg border border-notion-border p-5 hover:shadow-sm transition-shadow">
            <p class="text-xs font-medium text-notion-textLight mb-2 uppercase tracking-wide">Active Gateways</p>
            <p class="text-2xl font-semibold text-notion-text">{{ stats.activeGateways }}</p>
          </div>
        </div>
      </div>

      <!-- Content Grid -->
      <div class="grid grid-cols-1 gap-3 lg:grid-cols-2">
        <!-- Recent Activity -->
        <div class="bg-white rounded-lg border border-notion-border p-6">
          <h3 class="text-sm font-semibold text-notion-text mb-5">Recent Activity</h3>
          <div class="space-y-4">
            <div
              v-for="activity in recentActivities"
              :key="activity.id"
              class="pb-4 last:pb-0 border-b border-notion-border last:border-0"
            >
              <p class="text-sm font-medium text-notion-text">{{ activity.title }}</p>
              <p class="text-sm text-notion-textLight mt-1">{{ activity.description }}</p>
              <p class="text-xs text-notion-textLight mt-1.5">{{ activity.time }}</p>
            </div>
          </div>
        </div>

        <!-- System Info -->
        <div class="bg-white rounded-lg border border-notion-border p-6">
          <h3 class="text-sm font-semibold text-notion-text mb-5">System Information</h3>
          <div class="space-y-3">
            <div class="flex justify-between py-2.5 border-b border-notion-border">
              <span class="text-sm text-notion-textLight">Version</span>
              <span class="text-sm font-medium text-notion-text">v0.4.0</span>
            </div>
            <div class="flex justify-between py-2.5 border-b border-notion-border">
              <span class="text-sm text-notion-textLight">Uptime</span>
              <span class="text-sm font-medium text-notion-text">{{ systemInfo.uptime }}</span>
            </div>
            <div class="flex justify-between py-2.5 border-b border-notion-border">
              <span class="text-sm text-notion-textLight">Memory Usage</span>
              <span class="text-sm font-medium text-notion-text">{{ systemInfo.memory }}</span>
            </div>
            <div class="flex justify-between py-2.5">
              <span class="text-sm text-notion-textLight">CPU Usage</span>
              <span class="text-sm font-medium text-notion-text">{{ systemInfo.cpu }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="mt-3">
        <div class="bg-white rounded-lg border border-notion-border p-6">
          <h3 class="text-sm font-semibold text-notion-text mb-4">Quick Actions</h3>
          <div class="flex flex-wrap gap-2">
            <button class="btn-primary">Add Gateway</button>
            <button class="btn-secondary">Configure</button>
            <button class="btn-secondary">View Analytics</button>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { healthAPI } from '@/api'

const router = useRouter()
const authStore = useAuthStore()

const stats = ref({
  status: 'Online',
  totalRequests: 42853,
  avgResponse: 145,
  activeGateways: 3
})

const recentActivities = ref([
  {
    id: 1,
    title: 'Gateway Connected',
    description: 'New MCP gateway registered successfully',
    time: '2 minutes ago',
    statusColor: 'bg-green-500'
  },
  {
    id: 2,
    title: 'Request Processed',
    description: 'API request completed in 120ms',
    time: '5 minutes ago',
    statusColor: 'bg-blue-500'
  },
  {
    id: 3,
    title: 'Configuration Updated',
    description: 'Gateway settings modified',
    time: '15 minutes ago',
    statusColor: 'bg-yellow-500'
  },
  {
    id: 4,
    title: 'Health Check',
    description: 'All systems operational',
    time: '30 minutes ago',
    statusColor: 'bg-green-500'
  }
])

const systemInfo = ref({
  uptime: '7 days, 14 hours',
  memory: '245 MB / 2 GB',
  cpu: '12%'
})

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

onMounted(async () => {
  try {
    // Check API health
    await healthAPI.check()
  } catch (error) {
    console.error('Failed to check API health:', error)
  }
})
</script>

