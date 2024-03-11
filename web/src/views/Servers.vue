<template>
  <div class="min-h-screen bg-notion-bg">
    <!-- Navigation -->
    <nav class="bg-white border-b border-notion-border">
      <div class="max-w-7xl mx-auto px-6 lg:px-8">
        <div class="flex justify-between h-14">
          <div class="flex items-center">
            <div class="flex-shrink-0 flex items-center">
              <img src="/logo.png" alt="Mut Logo" class="h-8 w-auto">
            </div>
            <div class="hidden md:ml-8 md:flex md:space-x-1">
              <router-link
                to="/dashboard"
                class="text-notion-textLight hover:bg-notion-hover hover:text-notion-text inline-flex items-center px-3 py-1.5 rounded-md text-sm font-medium transition-colors"
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
                class="bg-notion-hover text-notion-text inline-flex items-center px-3 py-1.5 rounded-md text-sm font-medium"
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
        <h1 class="text-3xl font-semibold text-notion-text">Servers</h1>
        <p class="text-notion-textLight mt-2">Manage your server configurations</p>
      </div>

      <!-- Content Area -->
      <div class="bg-white rounded-lg border border-notion-border p-6">
        <p class="text-notion-textLight text-center py-8">Server content will be added here</p>
      </div>
    </main>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}
</script>

