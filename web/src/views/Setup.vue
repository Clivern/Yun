<template>
  <div class="min-h-screen flex items-center justify-center bg-notion-bg px-4 py-12">
    <div class="max-w-md w-full">
      <!-- Logo and Header -->
      <div class="text-center mb-10">
        <div class="flex justify-center mb-6">
          <img src="/logo.png" alt="Yun Logo" class="h-24 w-auto">
        </div>
        <h1 class="text-2xl font-semibold text-notion-text mb-2">Yun Setup</h1>
        <p class="text-sm text-notion-textLight">Configure your gateway and create an admin account</p>
      </div>

      <!-- Setup Form -->
      <div class="bg-white rounded-lg border border-notion-border p-8 shadow-sm">
        <form class="space-y-5" @submit.prevent="handleSetup">
          <!-- Gateway Name Field -->
          <div>
            <label for="gateway-name" class="block text-sm font-medium text-notion-text mb-2">
              Gateway Name
            </label>
            <input
              id="gateway-name"
              v-model="form.gatewayName"
              type="text"
              required
              minlength="2"
              maxlength="50"
              class="input-field"
              placeholder="My Gateway"
              :disabled="loading"
            >
            <p class="text-xs text-notion-textLight mt-1.5">A friendly name for this gateway</p>
          </div>

          <!-- Gateway URL Field -->
          <div>
            <label for="gateway-url" class="block text-sm font-medium text-notion-text mb-2">
              Gateway URL
            </label>
            <input
              id="gateway-url"
              v-model="form.gatewayURL"
              type="url"
              required
              class="input-field"
              placeholder="https://gateway.example.com"
              :disabled="loading"
            >
            <p class="text-xs text-notion-textLight mt-1.5">The public URL where this gateway will be accessible</p>
          </div>

          <!-- Gateway Email Field -->
          <div>
            <label for="gateway-email" class="block text-sm font-medium text-notion-text mb-2">
              Gateway Email
            </label>
            <input
              id="gateway-email"
              v-model="form.gatewayEmail"
              type="email"
              required
              class="input-field"
              placeholder="gateway@yun.com"
              :disabled="loading"
            >
            <p class="text-xs text-notion-textLight mt-1.5">The contact email for this gateway</p>
          </div>

          <!-- Divider -->
          <div class="relative py-3">
            <div class="absolute inset-0 flex items-center">
              <div class="w-full border-t border-notion-border"></div>
            </div>
            <div class="relative flex justify-center text-xs">
              <span class="px-3 bg-white text-notion-textLight font-medium">Admin Account</span>
            </div>
          </div>

          <!-- Admin Email Field -->
          <div>
            <label for="admin-email" class="block text-sm font-medium text-notion-text mb-2">
              Admin Email
            </label>
            <input
              id="admin-email"
              v-model="form.adminEmail"
              type="email"
              required
              class="input-field"
              placeholder="admin@yun.com"
              :disabled="loading"
            >
            <p class="text-xs text-notion-textLight mt-1.5">The email address for the admin account</p>
          </div>

          <!-- Admin Password Field -->
          <div>
            <label for="admin-password" class="block text-sm font-medium text-notion-text mb-2">
              Admin Password
            </label>
            <input
              id="admin-password"
              v-model="form.adminPassword"
              type="password"
              required
              minlength="8"
              class="input-field"
              placeholder="Enter a secure password (min. 8 characters)"
              :disabled="loading"
            >
            <p class="text-xs text-notion-textLight mt-1.5">Use a strong password with at least 8 characters</p>
          </div>

          <!-- Confirm Password Field -->
          <div>
            <label for="confirm-password" class="block text-sm font-medium text-notion-text mb-2">
              Confirm Password
            </label>
            <input
              id="confirm-password"
              v-model="form.confirmPassword"
              type="password"
              required
              minlength="8"
              class="input-field"
              placeholder="Re-enter your password"
              :disabled="loading"
            >
          </div>

          <!-- Error Message -->
          <div v-if="error" class="rounded-md border border-red-200 bg-red-50 p-3">
            <p class="text-sm text-red-800">
              {{ error }}
            </p>
          </div>

          <!-- Success Message -->
          <div v-if="success" class="rounded-md border border-green-200 bg-green-50 p-3">
            <p class="text-sm text-green-800">
              {{ success }}
            </p>
          </div>

          <!-- Submit Button -->
          <div>
            <button
              type="submit"
              class="w-full btn-primary py-2.5 disabled:opacity-50 disabled:cursor-not-allowed"
              :disabled="loading"
            >
              <span v-if="!loading">Complete Setup</span>
              <span v-else class="flex items-center justify-center">
                <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                </svg>
                Setting up...
              </span>
            </button>
          </div>
        </form>
      </div>

      <!-- Footer -->
      <p class="text-center text-xs text-notion-textLight mt-8">
        Copyright © 2025 Yun. All rights reserved.
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { setupAPI } from '@/api'

const router = useRouter()

const form = reactive({
  gatewayURL: '',
  gatewayEmail: '',
  gatewayName: '',
  adminEmail: '',
  adminPassword: '',
  confirmPassword: ''
})

const loading = ref(false)
const error = ref(null)
const success = ref(null)

// Check if system is already set up on component mount
onMounted(async () => {
  try {
    const response = await setupAPI.checkInstalled()
    if (response.data.installed) {
      // System already set up, redirect to login
      router.push('/login')
    }
  } catch (err) {
    // If check fails, allow access to setup page
    console.error('Failed to check setup status:', err)
  }
})

const validateForm = () => {
  // Check if passwords match
  if (form.adminPassword !== form.confirmPassword) {
    error.value = 'Passwords do not match'
    return false
  }

  // Check password length
  if (form.adminPassword.length < 8) {
    error.value = 'Password must be at least 8 characters long'
    return false
  }

  // Basic URL validation
  try {
    new URL(form.gatewayURL)
  } catch {
    error.value = 'Please enter a valid gateway URL'
    return false
  }

  return true
}

const handleSetup = async () => {
  loading.value = true
  error.value = null
  success.value = null

  if (!validateForm()) {
    loading.value = false
    return
  }

  try {
    const response = await setupAPI.install({
      gatewayURL: form.gatewayURL,
      gatewayEmail: form.gatewayEmail,
      gatewayName: form.gatewayName,
      adminEmail: form.adminEmail,
      adminPassword: form.adminPassword
    })

    success.value = 'Gateway setup completed successfully! Redirecting to login...'

    // Redirect to login after 2 seconds
    setTimeout(() => {
      router.push('/login')
    }, 2000)
  } catch (err) {
    if (err.response?.data?.errorMessage) {
      error.value = err.response.data.errorMessage
    } else if (err.message) {
      error.value = err.message
    } else {
      error.value = 'Setup failed. Please try again.'
    }
  } finally {
    loading.value = false
  }
}
</script>

