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

      <!-- Settings Sections -->
      <div class="space-y-6">
        <!-- General Settings -->
        <div class="bg-white rounded-lg border border-notion-border p-6 shadow-sm">
          <h2 class="text-lg font-semibold text-notion-text mb-5">General</h2>
          <div class="space-y-5">
            <!-- Gateway Name -->
            <div>
              <label for="gateway-name" class="block text-sm font-medium text-notion-text mb-2">
                Gateway Name
              </label>
              <input
                id="gateway-name"
                v-model="settings.gatewayName"
                type="text"
                class="input-field max-w-md"
                placeholder="My Gateway"
              >
              <p class="text-xs text-notion-textLight mt-1.5">The name displayed throughout the application</p>
            </div>

            <!-- Gateway URL -->
            <div>
              <label for="gateway-url" class="block text-sm font-medium text-notion-text mb-2">
                Gateway URL
              </label>
              <input
                id="gateway-url"
                v-model="settings.gatewayUrl"
                type="url"
                class="input-field max-w-md"
                placeholder="https://mut.com"
              >
              <p class="text-xs text-notion-textLight mt-1.5">The public URL where this application is accessible</p>
            </div>

            <!-- Gateway Email -->
            <div>
              <label for="gateway-email" class="block text-sm font-medium text-notion-text mb-2">
                Gateway Email
              </label>
              <input
                id="gateway-email"
                v-model="settings.gatewayEmail"
                type="email"
                class="input-field max-w-md"
                placeholder="admin@example.com"
              >
              <p class="text-xs text-notion-textLight mt-1.5">Primary contact email displayed in notifications</p>
            </div>

            <!-- Description -->
            <div>
              <label for="gateway-description" class="block text-sm font-medium text-notion-text mb-2">
                Description
              </label>
              <textarea
                id="gateway-description"
                v-model="settings.gatewayDescription"
                rows="3"
                class="input-field max-w-2xl"
                placeholder="A brief description of your gateway..."
              ></textarea>
              <p class="text-xs text-notion-textLight mt-1.5">Optional description for your gateway</p>
            </div>
          </div>
        </div>

        <!-- SMTP Settings -->
        <div class="bg-white rounded-lg border border-notion-border p-6 shadow-sm">
          <h2 class="text-lg font-semibold text-notion-text mb-5">SMTP</h2>
          <div class="space-y-5">
            <!-- SMTP Server -->
            <div>
              <label for="smtp-server" class="block text-sm font-medium text-notion-text mb-2">
                SMTP Server
              </label>
              <input
                id="smtp-server"
                v-model="settings.smtpServer"
                type="text"
                class="input-field max-w-md"
                placeholder="smtp.example.com"
              >
            </div>

            <!-- SMTP Port -->
            <div>
              <label for="smtp-port" class="block text-sm font-medium text-notion-text mb-2">
                SMTP Port
              </label>
              <input
                id="smtp-port"
                v-model="settings.smtpPort"
                type="number"
                min="1"
                max="65535"
                class="input-field max-w-xs"
                placeholder="587"
              >
            </div>

            <!-- From Email -->
            <div>
              <label for="smtp-from-email" class="block text-sm font-medium text-notion-text mb-2">
                From Email
              </label>
              <input
                id="smtp-from-email"
                v-model="settings.smtpFromEmail"
                type="email"
                class="input-field max-w-md"
                placeholder="noreply@example.com"
              >
              <p class="text-xs text-notion-textLight mt-1.5">Email address used as the sender for notifications</p>
            </div>

            <!-- SMTP Username -->
            <div>
              <label for="smtp-username" class="block text-sm font-medium text-notion-text mb-2">
                SMTP Username
              </label>
              <input
                id="smtp-username"
                v-model="settings.smtpUsername"
                type="text"
                class="input-field max-w-md"
                placeholder="user@example.com"
              >
            </div>

            <!-- SMTP Password -->
            <div>
              <label for="smtp-password" class="block text-sm font-medium text-notion-text mb-2">
                SMTP Password
              </label>
              <input
                id="smtp-password"
                v-model="settings.smtpPassword"
                type="password"
                class="input-field max-w-md"
                placeholder="••••••••"
              >
            </div>

            <!-- Use TLS -->
            <div class="flex items-start">
              <input
                id="smtp-use-tls"
                v-model="settings.smtpUseTLS"
                type="checkbox"
                class="h-4 w-4 mt-1 rounded text-notion-text focus:ring-notion-text border-notion-border"
              >
              <label for="smtp-use-tls" class="ml-3 block">
                <span class="text-sm font-medium text-notion-text">Use TLS</span>
                <p class="text-xs text-notion-textLight mt-1">Enable TLS when connecting to the SMTP server</p>
              </label>
            </div>
          </div>
        </div>

        <!-- Maintenance Mode -->
        <div class="bg-white rounded-lg border border-notion-border p-6 shadow-sm">
          <h2 class="text-lg font-semibold text-notion-text mb-5">Maintenance</h2>
          <div class="space-y-5">
            <!-- Maintenance Mode Toggle -->
            <div class="flex items-start">
              <input
                id="maintenance-mode"
                v-model="settings.maintenanceMode"
                type="checkbox"
                class="h-4 w-4 mt-1 rounded text-notion-text focus:ring-notion-text border-notion-border"
              >
              <label for="maintenance-mode" class="ml-3 block">
                <span class="text-sm font-medium text-notion-text">Enable Maintenance Mode</span>
                <p class="text-xs text-notion-textLight mt-1">Put the application in maintenance mode (only admins can access)</p>
              </label>
            </div>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex justify-end gap-3 pt-2">
          <button
            @click="handleReset"
            class="btn-secondary"
            :disabled="saving"
          >
            Reset to Defaults
          </button>
          <button
            @click="handleSave"
            class="btn-primary"
            :disabled="saving"
          >
            <span v-if="!saving">Save Changes</span>
            <span v-else class="flex items-center">
              <svg class="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
              </svg>
              Saving...
            </span>
          </button>
        </div>

        <!-- Success/Error Messages -->
        <div v-if="successMessage" class="rounded-md border border-green-200 bg-green-50 p-4">
          <p class="text-sm text-green-800">{{ successMessage }}</p>
        </div>
        <div v-if="errorMessage" class="rounded-md border border-red-200 bg-red-50 p-4">
          <p class="text-sm text-red-800">{{ errorMessage }}</p>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const saving = ref(false)
const successMessage = ref(null)
const errorMessage = ref(null)

const settings = reactive({
  // General
  gatewayName: 'Mut Gateway',
  gatewayUrl: 'https://gateway.example.com',
  gatewayEmail: 'admin@example.com',
  gatewayDescription: '',

  // SMTP
  smtpServer: '',
  smtpPort: 587,
  smtpFromEmail: '',
  smtpUsername: '',
  smtpPassword: '',
  smtpUseTLS: false,

  // Maintenance
  maintenanceMode: false
})

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

const handleSave = async () => {
  saving.value = true
  successMessage.value = null
  errorMessage.value = null

  try {
    // TODO: Add API call to save settings
    // await settingsAPI.update(settings)

    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000))

    successMessage.value = 'Settings saved successfully!'

    // Clear success message after 3 seconds
    setTimeout(() => {
      successMessage.value = null
    }, 3000)
  } catch (err) {
    console.error('Failed to save settings:', err)
    errorMessage.value = 'Failed to save settings. Please try again.'
  } finally {
    saving.value = false
  }
}

const handleReset = () => {
  if (confirm('Are you sure you want to reset all settings to their default values?')) {
    // Reset to default values
    settings.gatewayName = 'Mut Gateway'
    settings.gatewayUrl = 'https://gateway.example.com'
    settings.gatewayEmail = 'admin@example.com'
    settings.gatewayDescription = ''
    settings.smtpServer = ''
    settings.smtpPort = 587
    settings.smtpFromEmail = ''
    settings.smtpUsername = ''
    settings.smtpPassword = ''
    settings.smtpUseTLS = false
    settings.maintenanceMode = false

    successMessage.value = 'Settings reset to defaults'
    setTimeout(() => {
      successMessage.value = null
    }, 3000)
  }
}
</script>

