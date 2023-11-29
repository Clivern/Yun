<template>
  <div class="min-h-screen flex items-center justify-center bg-notion-bg px-4">
    <div class="max-w-md w-full">
      <!-- Logo and Header -->
      <div class="text-center mb-10">
        <div class="flex justify-center mb-8">
          <img src="/logo.png" alt="Yun Logo" class="h-24 w-auto">
        </div>
      </div>

      <!-- Login Form -->
      <div class="bg-white rounded-lg border border-notion-border p-8 shadow-sm">
        <h2 class="text-2xl font-semibold text-notion-text mb-6 text-center">Welcome back</h2>
        <form class="space-y-5" @submit.prevent="handleLogin">
          <!-- Email Field -->
          <div>
            <label for="email" class="block text-sm font-medium text-notion-text mb-2">
              Email
            </label>
            <input
              id="email"
              v-model="form.email"
              type="email"
              required
              class="input-field"
              placeholder="Enter your email"
              :disabled="loading"
            >
          </div>

          <!-- Password Field -->
          <div>
            <label for="password" class="block text-sm font-medium text-notion-text mb-2">
              Password
            </label>
            <input
              id="password"
              v-model="form.password"
              type="password"
              required
              class="input-field"
              placeholder="Enter your password"
              :disabled="loading"
            >
          </div>

          <!-- Remember Me -->
          <div class="flex items-center justify-between">
            <div class="flex items-center">
              <input
                id="remember-me"
                v-model="form.rememberMe"
                type="checkbox"
                class="h-4 w-4 rounded text-notion-text focus:ring-notion-text border-notion-border"
              >
              <label for="remember-me" class="ml-2 block text-sm text-notion-textLight">
                Remember me
              </label>
            </div>
          </div>

          <!-- Error Message -->
          <div v-if="error" class="rounded-md border border-red-200 bg-red-50 p-3">
            <p class="text-sm text-red-800">
              {{ error }}
            </p>
          </div>

          <!-- Submit Button -->
          <div>
            <button
              type="submit"
              class="w-full btn-primary py-2.5 disabled:opacity-50 disabled:cursor-not-allowed"
              :disabled="loading"
            >
              <span v-if="!loading">Sign in</span>
              <span v-else class="flex items-center justify-center">
                <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                </svg>
                Signing in...
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
import { reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const form = reactive({
  email: '',
  password: '',
  rememberMe: false
})

const loading = computed(() => authStore.loading)
const error = computed(() => authStore.error)

const handleLogin = async () => {
  const result = await authStore.login(form.email, form.password, form.rememberMe)
  if (result.success) {
    router.push('/dashboard')
  }
}
</script>

