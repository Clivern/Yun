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
                class="bg-notion-hover text-notion-text inline-flex items-center px-3 py-1.5 rounded-md text-sm font-medium"
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
        <h1 class="text-3xl font-semibold text-notion-text">Settings</h1>
        <p class="text-sm text-notion-textLight mt-2">Manage gateway settings and preferences</p>
      </div>

      <!-- Success/Error Messages -->
      <div v-if="successMessage" class="mb-6 rounded-md border border-green-200 bg-green-50 p-4">
        <p class="text-sm text-green-800">{{ successMessage }}</p>
      </div>
      <div v-if="errorMessage" class="mb-6 rounded-md border border-red-200 bg-red-50 p-4">
        <p class="text-sm text-red-800">{{ errorMessage }}</p>
      </div>

      <!-- Settings Form -->
      <form @submit.prevent="handleSave" class="space-y-6">
        <!-- General Settings Box -->
        <div class="bg-white rounded-lg border border-notion-border p-6 shadow-sm">
          <h2 class="text-lg font-semibold text-notion-text mb-4">General</h2>
          <div class="border-b border-notion-border mb-4"></div>
          <div class="space-y-4">
            <div>
              <label for="gatewayName" class="block text-sm font-medium text-notion-text mb-2">
                Gateway Name
              </label>
              <input
                id="gatewayName"
                v-model="form.gatewayName"
                type="text"
                required
                class="input-field max-w-md"
                placeholder="Mut"
                :disabled="loading"
              >
            </div>
            <div>
              <label for="gatewayEmail" class="block text-sm font-medium text-notion-text mb-2">
                Gateway Email
              </label>
              <input
                id="gatewayEmail"
                v-model="form.gatewayEmail"
                type="email"
                required
                class="input-field max-w-md"
                placeholder="admin@mut.ai"
                :disabled="loading"
              >
            </div>
            <div>
              <label for="gatewayURL" class="block text-sm font-medium text-notion-text mb-2">
                Gateway URL
              </label>
              <input
                id="gatewayURL"
                v-model="form.gatewayURL"
                type="url"
                required
                class="input-field max-w-md"
                placeholder="http://mut.ai"
                :disabled="loading"
              >
            </div>
            <div>
              <label for="gatewayDescription" class="block text-sm font-medium text-notion-text mb-2">
                Gateway Description
              </label>
              <textarea
                id="gatewayDescription"
                v-model="form.gatewayDescription"
                rows="4"
                class="input-field max-w-md"
                placeholder="Enter a description for your gateway"
                :disabled="loading"
              ></textarea>
            </div>
          </div>
        </div>

        <!-- SMTP Configuration Box -->
        <div class="bg-white rounded-lg border border-notion-border p-6 shadow-sm">
          <h2 class="text-lg font-semibold text-notion-text mb-4">SMTP Configuration</h2>
          <div class="border-b border-notion-border mb-4"></div>
          <div class="space-y-4">
            <div>
              <label for="smtpServer" class="block text-sm font-medium text-notion-text mb-2">
                SMTP Server
              </label>
              <input
                id="smtpServer"
                v-model="form.smtpServer"
                type="text"
                class="input-field max-w-md"
                placeholder="smtp.mut.ai"
                :disabled="loading"
              >
            </div>

            <div>
              <label for="smtpPort" class="block text-sm font-medium text-notion-text mb-2">
                SMTP Port
              </label>
              <input
                id="smtpPort"
                v-model="form.smtpPort"
                type="text"
                class="input-field max-w-xs"
                placeholder="587"
                :disabled="loading"
              >
            </div>

            <div>
              <label for="smtpFromEmail" class="block text-sm font-medium text-notion-text mb-2">
                SMTP From Email
              </label>
              <input
                id="smtpFromEmail"
                v-model="form.smtpFromEmail"
                type="email"
                class="input-field max-w-md"
                placeholder="noreply@mut.ai"
                :disabled="loading"
              >
            </div>

            <div>
              <label for="smtpUsername" class="block text-sm font-medium text-notion-text mb-2">
                SMTP Username
              </label>
              <input
                id="smtpUsername"
                v-model="form.smtpUsername"
                type="text"
                class="input-field max-w-md"
                placeholder="smtpuser"
                :disabled="loading"
              >
            </div>

            <div>
              <label for="smtpPassword" class="block text-sm font-medium text-notion-text mb-2">
                SMTP Password
              </label>
              <input
                id="smtpPassword"
                v-model="form.smtpPassword"
                type="password"
                class="input-field max-w-md"
                placeholder="Enter SMTP password"
                :disabled="loading"
              >
            </div>

            <div class="flex items-start">
              <input
                id="smtpUseTLS"
                v-model="form.smtpUseTLS"
                type="checkbox"
                class="h-4 w-4 mt-1 rounded text-notion-text focus:ring-notion-text border-notion-border"
                :disabled="loading"
              >
              <label for="smtpUseTLS" class="ml-3 block">
                <span class="text-sm font-medium text-notion-text">Use TLS</span>
                <p class="text-xs text-notion-textLight mt-1">Enable TLS for SMTP connection</p>
              </label>
            </div>
          </div>
        </div>

        <!-- Maintenance Mode Box -->
        <div class="bg-white rounded-lg border border-notion-border p-6 shadow-sm">
          <h2 class="text-lg font-semibold text-notion-text mb-4">Maintenance</h2>
          <div class="border-b border-notion-border mb-4"></div>
          <div class="flex items-start">
            <input
              id="maintenanceMode"
              v-model="form.maintenanceMode"
              type="checkbox"
              class="h-4 w-4 mt-1 rounded text-notion-text focus:ring-notion-text border-notion-border"
              :disabled="loading"
            >
            <label for="maintenanceMode" class="ml-3 block">
              <span class="text-sm font-medium text-notion-text">Enable Maintenance Mode</span>
              <p class="text-xs text-notion-textLight mt-1">Put the application in maintenance mode</p>
            </label>
          </div>
        </div>

        <!-- Submit Button -->
        <div class="flex justify-end gap-3">
          <button
            type="submit"
            class="btn-primary"
            :disabled="loading"
          >
            <span v-if="!loading">Save Settings</span>
            <span v-else class="flex items-center">
              <svg class="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
              </svg>
              Saving...
            </span>
          </button>
        </div>
      </form>
    </main>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { settingsAPI } from '@/api'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const successMessage = ref(null)
const errorMessage = ref(null)

const form = reactive({
  gatewayURL: '',
  gatewayEmail: '',
  gatewayName: '',
  gatewayDescription: '',
  maintenanceMode: false,
  smtpServer: '',
  smtpPort: '',
  smtpFromEmail: '',
  smtpUsername: '',
  smtpPassword: '',
  smtpUseTLS: true
})

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

const loadSettings = async () => {
  loading.value = true
  errorMessage.value = null

  try {
    const response = await settingsAPI.getSettings()
    const settings = response.data?.settings

    if (settings) {
      // Map API response (snake_case) to form fields
      form.gatewayURL = settings.gateway_url || ''
      form.gatewayEmail = settings.gateway_email || ''
      form.gatewayName = settings.gateway_name || ''
      form.gatewayDescription = settings.gateway_description || ''
      form.maintenanceMode = settings.maintenance_mode === '1'
      form.smtpServer = settings.smtp_server || ''
      form.smtpPort = settings.smtp_port || ''
      form.smtpFromEmail = settings.smtp_from_email || ''
      form.smtpUsername = settings.smtp_username || ''
      form.smtpPassword = settings.smtp_password || ''
      form.smtpUseTLS = settings.smtp_use_tls === '1'
    }
  } catch (err) {
    console.error('Failed to load settings:', err)
    errorMessage.value = err.response?.data?.errorMessage || 'Failed to load settings'
  } finally {
    loading.value = false
  }
}

const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

// Watch for message changes and scroll to top
watch([successMessage, errorMessage], () => {
  if (successMessage.value || errorMessage.value) {
    nextTick(() => {
      scrollToTop()
    })
  }
})

const handleSave = async () => {
  loading.value = true
  successMessage.value = null
  errorMessage.value = null

  try {
    // Map form fields to API request format (camelCase)
    const payload = {
      gatewayName: form.gatewayName,
      gatewayUrl: form.gatewayURL,
      gatewayEmail: form.gatewayEmail,
      gatewayDescription: form.gatewayDescription,
      maintenanceMode: form.maintenanceMode,
      smtpUseTLS: form.smtpUseTLS
    }

    // Only include SMTP fields if they have values
    if (form.smtpServer) { payload.smtpServer = form.smtpServer }
    if (form.smtpPort) { payload.smtpPort = parseInt(form.smtpPort, 10) || 587 }
    if (form.smtpFromEmail) { payload.smtpFromEmail = form.smtpFromEmail }
    if (form.smtpUsername) { payload.smtpUsername = form.smtpUsername }
    if (form.smtpPassword) { payload.smtpPassword = form.smtpPassword }

    const response = await settingsAPI.updateSettings(payload)

    if (response.data?.successMessage) {
      successMessage.value = response.data.successMessage
      setTimeout(() => {
        successMessage.value = null
      }, 3000)
    }
  } catch (err) {
    console.error('Failed to save settings:', err)
    if (err.response?.data?.errorMessage) {
      errorMessage.value = err.response.data.errorMessage
    } else if (err.response?.data?.errors) {
      // Handle validation errors
      const errors = err.response.data.errors
      const errorList = Object.values(errors).flat().join(', ')
      errorMessage.value = `Validation errors: ${errorList}`
    } else {
      errorMessage.value = 'Failed to save settings. Please try again.'
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadSettings()
})
</script>
