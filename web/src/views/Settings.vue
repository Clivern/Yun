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
            <!-- Application Name -->
            <div>
              <label for="app-name" class="block text-sm font-medium text-notion-text mb-2">
                Gateway Name
              </label>
              <input
                id="app-name"
                v-model="settings.appName"
                type="text"
                class="input-field max-w-md"
                placeholder="My Gateway"
              >
              <p class="text-xs text-notion-textLight mt-1.5">The name displayed throughout the application</p>
            </div>

            <!-- Application URL -->
            <div>
              <label for="app-url" class="block text-sm font-medium text-notion-text mb-2">
                Gateway URL
              </label>
              <input
                id="app-url"
                v-model="settings.appUrl"
                type="url"
                class="input-field max-w-md"
                placeholder="https://mut.com"
              >
              <p class="text-xs text-notion-textLight mt-1.5">The public URL where this application is accessible</p>
            </div>

            <!-- Description -->
            <div>
              <label for="description" class="block text-sm font-medium text-notion-text mb-2">
                Description
              </label>
              <textarea
                id="description"
                v-model="settings.description"
                rows="3"
                class="input-field max-w-2xl"
                placeholder="A brief description of your gateway..."
              ></textarea>
              <p class="text-xs text-notion-textLight mt-1.5">Optional description for your gateway</p>
            </div>
          </div>
        </div>

        <!-- Security Settings -->
        <div class="bg-white rounded-lg border border-notion-border p-6 shadow-sm">
          <h2 class="text-lg font-semibold text-notion-text mb-5">Security</h2>
          <div class="space-y-5">
            <!-- Session Timeout -->
            <div>
              <label for="session-timeout" class="block text-sm font-medium text-notion-text mb-2">
                Session Timeout (minutes)
              </label>
              <input
                id="session-timeout"
                v-model="settings.sessionTimeout"
                type="number"
                min="5"
                max="1440"
                class="input-field max-w-xs"
                placeholder="30"
              >
              <p class="text-xs text-notion-textLight mt-1.5">How long users can stay logged in without activity</p>
            </div>

            <!-- Require Strong Passwords -->
            <div class="flex items-start">
              <input
                id="strong-passwords"
                v-model="settings.requireStrongPasswords"
                type="checkbox"
                class="h-4 w-4 mt-1 rounded text-notion-text focus:ring-notion-text border-notion-border"
              >
              <label for="strong-passwords" class="ml-3 block">
                <span class="text-sm font-medium text-notion-text">Require Strong Passwords</span>
                <p class="text-xs text-notion-textLight mt-1">Enforce password complexity requirements for all users</p>
              </label>
            </div>

            <!-- Enable Two-Factor Authentication -->
            <div class="flex items-start">
              <input
                id="two-factor"
                v-model="settings.enableTwoFactor"
                type="checkbox"
                class="h-4 w-4 mt-1 rounded text-notion-text focus:ring-notion-text border-notion-border"
              >
              <label for="two-factor" class="ml-3 block">
                <span class="text-sm font-medium text-notion-text">Enable Two-Factor Authentication</span>
                <p class="text-xs text-notion-textLight mt-1">Allow users to set up 2FA for their accounts</p>
              </label>
            </div>
          </div>
        </div>

        <!-- API Settings -->
        <div class="bg-white rounded-lg border border-notion-border p-6 shadow-sm">
          <h2 class="text-lg font-semibold text-notion-text mb-5">API Configuration</h2>
          <div class="space-y-5">
            <!-- API Rate Limit -->
            <div>
              <label for="rate-limit" class="block text-sm font-medium text-notion-text mb-2">
                Rate Limit (requests per minute)
              </label>
              <input
                id="rate-limit"
                v-model="settings.rateLimit"
                type="number"
                min="10"
                max="10000"
                class="input-field max-w-xs"
                placeholder="100"
              >
              <p class="text-xs text-notion-textLight mt-1.5">Maximum number of API requests per minute per IP</p>
            </div>

            <!-- API Timeout -->
            <div>
              <label for="api-timeout" class="block text-sm font-medium text-notion-text mb-2">
                API Timeout (seconds)
              </label>
              <input
                id="api-timeout"
                v-model="settings.apiTimeout"
                type="number"
                min="5"
                max="300"
                class="input-field max-w-xs"
                placeholder="30"
              >
              <p class="text-xs text-notion-textLight mt-1.5">Maximum time to wait for API responses</p>
            </div>

            <!-- Enable API Logging -->
            <div class="flex items-start">
              <input
                id="api-logging"
                v-model="settings.enableApiLogging"
                type="checkbox"
                class="h-4 w-4 mt-1 rounded text-notion-text focus:ring-notion-text border-notion-border"
              >
              <label for="api-logging" class="ml-3 block">
                <span class="text-sm font-medium text-notion-text">Enable API Request Logging</span>
                <p class="text-xs text-notion-textLight mt-1">Log all API requests for debugging and monitoring</p>
              </label>
            </div>
          </div>
        </div>

        <!-- Email Settings -->
        <div class="bg-white rounded-lg border border-notion-border p-6 shadow-sm">
          <h2 class="text-lg font-semibold text-notion-text mb-5">Email Notifications</h2>
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
              <label for="from-email" class="block text-sm font-medium text-notion-text mb-2">
                From Email
              </label>
              <input
                id="from-email"
                v-model="settings.fromEmail"
                type="email"
                class="input-field max-w-md"
                placeholder="noreply@example.com"
              >
              <p class="text-xs text-notion-textLight mt-1.5">Email address used for sending notifications</p>
            </div>

            <!-- Enable Email Notifications -->
            <div class="flex items-start">
              <input
                id="email-notifications"
                v-model="settings.enableEmailNotifications"
                type="checkbox"
                class="h-4 w-4 mt-1 rounded text-notion-text focus:ring-notion-text border-notion-border"
              >
              <label for="email-notifications" class="ml-3 block">
                <span class="text-sm font-medium text-notion-text">Enable Email Notifications</span>
                <p class="text-xs text-notion-textLight mt-1">Send email notifications for important events</p>
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

            <!-- Maintenance Message -->
            <div v-if="settings.maintenanceMode">
              <label for="maintenance-message" class="block text-sm font-medium text-notion-text mb-2">
                Maintenance Message
              </label>
              <textarea
                id="maintenance-message"
                v-model="settings.maintenanceMessage"
                rows="3"
                class="input-field max-w-2xl"
                placeholder="We're currently performing scheduled maintenance. We'll be back shortly."
              ></textarea>
              <p class="text-xs text-notion-textLight mt-1.5">Message displayed to users during maintenance</p>
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
  appName: 'Mut Gateway',
  appUrl: 'https://gateway.example.com',
  description: '',

  // Security
  sessionTimeout: 30,
  requireStrongPasswords: true,
  enableTwoFactor: false,

  // API
  rateLimit: 100,
  apiTimeout: 30,
  enableApiLogging: true,

  // Email
  smtpServer: '',
  smtpPort: 587,
  fromEmail: '',
  enableEmailNotifications: false,

  // Maintenance
  maintenanceMode: false,
  maintenanceMessage: "We're currently performing scheduled maintenance. We'll be back shortly."
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
    settings.appName = 'Mut Gateway'
    settings.appUrl = 'https://gateway.example.com'
    settings.description = ''
    settings.sessionTimeout = 30
    settings.requireStrongPasswords = true
    settings.enableTwoFactor = false
    settings.rateLimit = 100
    settings.apiTimeout = 30
    settings.enableApiLogging = true
    settings.smtpServer = ''
    settings.smtpPort = 587
    settings.fromEmail = ''
    settings.enableEmailNotifications = false
    settings.maintenanceMode = false
    settings.maintenanceMessage = "We're currently performing scheduled maintenance. We'll be back shortly."

    successMessage.value = 'Settings reset to defaults'
    setTimeout(() => {
      successMessage.value = null
    }, 3000)
  }
}
</script>

