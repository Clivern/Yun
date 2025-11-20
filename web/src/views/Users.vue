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
                class="bg-notion-hover text-notion-text inline-flex items-center px-3 py-1.5 rounded-md text-sm font-medium"
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
      <div class="mb-8 flex justify-between items-center">
        <div>
          <h1 class="text-3xl font-semibold text-notion-text">Users</h1>
          <p class="text-notion-textLight mt-2">Manage user accounts and permissions</p>
        </div>
        <button
          @click="openCreateModal"
          class="btn-primary"
        >
          Add User
        </button>
      </div>

      <!-- Success/Error Messages -->
      <div v-if="successMessage" class="mb-6 rounded-md border border-green-200 bg-green-50 p-4">
        <p class="text-sm text-green-800">{{ successMessage }}</p>
      </div>
      <div v-if="errorMessage" class="mb-6 rounded-md border border-red-200 bg-red-50 p-4">
        <p class="text-sm text-red-800">{{ errorMessage }}</p>
      </div>

      <!-- Users Table -->
      <div class="bg-white rounded-lg border border-notion-border overflow-hidden">
        <div v-if="loading" class="p-8 text-center">
          <svg class="animate-spin h-6 w-6 mx-auto text-notion-text" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
          <p class="text-notion-textLight mt-2">Loading users...</p>
        </div>

        <div v-else-if="users.length === 0" class="p-8 text-center">
          <p class="text-notion-textLight">No users found</p>
        </div>

        <table v-else class="min-w-full divide-y divide-notion-border">
          <thead class="bg-notion-hover">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-notion-textLight uppercase tracking-wider">
                Email
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-notion-textLight uppercase tracking-wider">
                Role
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-notion-textLight uppercase tracking-wider">
                Status
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-notion-textLight uppercase tracking-wider">
                Last Login
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-notion-textLight uppercase tracking-wider">
                Created
              </th>
              <th class="px-6 py-3 text-right text-xs font-medium text-notion-textLight uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-notion-border">
            <tr v-for="user in users" :key="user.id" class="hover:bg-notion-hover">
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm font-medium text-notion-text">{{ user.email }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="px-2 py-1 text-xs font-medium rounded-full"
                  :class="{
                    'bg-purple-100 text-purple-800': user.role === 'admin',
                    'bg-blue-100 text-blue-800': user.role === 'user',
                    'bg-gray-100 text-gray-800': user.role === 'readonly'
                  }">
                  {{ user.role }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="px-2 py-1 text-xs font-medium rounded-full"
                  :class="{
                    'bg-green-100 text-green-800': user.isActive,
                    'bg-red-100 text-red-800': !user.isActive
                  }">
                  {{ user.isActive ? 'Active' : 'Inactive' }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-notion-textLight">
                {{ formatDate(user.lastLoginAt) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-notion-textLight">
                {{ formatDate(user.createdAt) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <button
                  @click="openEditModal(user)"
                  class="text-blue-600 hover:text-blue-900 mr-4"
                >
                  Edit
                </button>
                <button
                  @click="openDeleteModal(user)"
                  class="text-red-600 hover:text-red-900"
                >
                  Delete
                </button>
              </td>
            </tr>
          </tbody>
        </table>

        <!-- Pagination -->
        <div v-if="total > 0" class="bg-white px-6 py-4 border-t border-notion-border flex items-center justify-between">
          <div class="text-sm text-notion-textLight">
            Showing {{ offset + 1 }} to {{ Math.min(offset + limit, total) }} of {{ total }} users
          </div>
          <div class="flex items-center space-x-2">
            <button
              @click="goToPage(offset - limit)"
              :disabled="offset === 0"
              class="btn-secondary text-sm disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Previous
            </button>
            <button
              @click="goToPage(offset + limit)"
              :disabled="offset + limit >= total"
              class="btn-secondary text-sm disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Next
            </button>
          </div>
        </div>
      </div>
    </main>

    <!-- Create User Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 z-50 overflow-y-auto" @click.self="closeCreateModal">
      <div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 transition-opacity bg-gray-500 bg-opacity-75" @click="closeCreateModal"></div>
        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <form @submit.prevent="handleCreate">
            <div class="bg-white px-6 pt-6 pb-4">
              <div class="flex justify-between items-center mb-4">
                <h3 class="text-lg font-semibold text-notion-text">Create User</h3>
                <button
                  type="button"
                  @click="closeCreateModal"
                  class="text-notion-textLight hover:text-notion-text"
                >
                  <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
              <div class="space-y-4">
                <div>
                  <label for="create-email" class="block text-sm font-medium text-notion-text mb-2">
                    Email *
                  </label>
                  <input
                    id="create-email"
                    v-model="createForm.email"
                    type="email"
                    required
                    class="input-field"
                    placeholder="user@example.com"
                  />
                </div>
                <div>
                  <label for="create-password" class="block text-sm font-medium text-notion-text mb-2">
                    Password *
                  </label>
                  <input
                    id="create-password"
                    v-model="createForm.password"
                    type="password"
                    required
                    class="input-field"
                    placeholder="Minimum 8 characters"
                  />
                </div>
                <div>
                  <label for="create-role" class="block text-sm font-medium text-notion-text mb-2">
                    Role *
                  </label>
                  <select
                    id="create-role"
                    v-model="createForm.role"
                    required
                    class="input-field"
                  >
                    <option value="admin">Admin</option>
                    <option value="user">User</option>
                    <option value="readonly">Readonly</option>
                  </select>
                </div>
                <div class="flex items-start">
                  <input
                    id="create-isActive"
                    v-model="createForm.isActive"
                    type="checkbox"
                    class="h-4 w-4 mt-1 rounded text-notion-text focus:ring-notion-text border-notion-border"
                  />
                  <label for="create-isActive" class="ml-3 block">
                    <span class="text-sm font-medium text-notion-text">Active</span>
                    <p class="text-xs text-notion-textLight mt-1">User account is active and can log in</p>
                  </label>
                </div>
              </div>
            </div>
            <div class="bg-notion-hover px-6 py-4 flex justify-end space-x-3">
              <button
                type="button"
                @click="closeCreateModal"
                class="btn-secondary"
              >
                Cancel
              </button>
              <button
                type="submit"
                class="btn-primary"
                :disabled="createLoading"
              >
                <span v-if="!createLoading">Create</span>
                <span v-else class="flex items-center">
                  <svg class="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                  </svg>
                  Creating...
                </span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Edit User Modal -->
    <div v-if="showEditModal" class="fixed inset-0 z-50 overflow-y-auto" @click.self="closeEditModal">
      <div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 transition-opacity bg-gray-500 bg-opacity-75" @click="closeEditModal"></div>
        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <form @submit.prevent="handleUpdate">
            <div class="bg-white px-6 pt-6 pb-4">
              <div class="flex justify-between items-center mb-4">
                <h3 class="text-lg font-semibold text-notion-text">Edit User</h3>
                <button
                  type="button"
                  @click="closeEditModal"
                  class="text-notion-textLight hover:text-notion-text"
                >
                  <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
              <div class="space-y-4">
                <div>
                  <label for="edit-email" class="block text-sm font-medium text-notion-text mb-2">
                    Email *
                  </label>
                  <input
                    id="edit-email"
                    v-model="editForm.email"
                    type="email"
                    required
                    class="input-field"
                    placeholder="user@example.com"
                  />
                </div>
                <div>
                  <label for="edit-password" class="block text-sm font-medium text-notion-text mb-2">
                    Password
                  </label>
                  <input
                    id="edit-password"
                    v-model="editForm.password"
                    type="password"
                    class="input-field"
                    placeholder="Leave empty to keep current password"
                  />
                  <p class="text-xs text-notion-textLight mt-1">Leave empty to keep current password</p>
                </div>
                <div>
                  <label for="edit-role" class="block text-sm font-medium text-notion-text mb-2">
                    Role *
                  </label>
                  <select
                    id="edit-role"
                    v-model="editForm.role"
                    required
                    class="input-field"
                  >
                    <option value="admin">Admin</option>
                    <option value="user">User</option>
                    <option value="readonly">Readonly</option>
                  </select>
                </div>
                <div class="flex items-start">
                  <input
                    id="edit-isActive"
                    v-model="editForm.isActive"
                    type="checkbox"
                    class="h-4 w-4 mt-1 rounded text-notion-text focus:ring-notion-text border-notion-border"
                  />
                  <label for="edit-isActive" class="ml-3 block">
                    <span class="text-sm font-medium text-notion-text">Active</span>
                    <p class="text-xs text-notion-textLight mt-1">User account is active and can log in</p>
                  </label>
                </div>
              </div>
            </div>
            <div class="bg-notion-hover px-6 py-4 flex justify-end space-x-3">
              <button
                type="button"
                @click="closeEditModal"
                class="btn-secondary"
              >
                Cancel
              </button>
              <button
                type="submit"
                class="btn-primary"
                :disabled="editLoading"
              >
                <span v-if="!editLoading">Update</span>
                <span v-else class="flex items-center">
                  <svg class="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                  </svg>
                  Updating...
                </span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="fixed inset-0 z-50 overflow-y-auto" @click.self="closeDeleteModal">
      <div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 transition-opacity bg-gray-500 bg-opacity-75" @click="closeDeleteModal"></div>
        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <div class="bg-white px-6 pt-6 pb-4">
            <div class="flex justify-between items-center mb-4">
              <h3 class="text-lg font-semibold text-notion-text">Delete User</h3>
              <button
                type="button"
                @click="closeDeleteModal"
                class="text-notion-textLight hover:text-notion-text"
              >
                <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            <p class="text-sm text-notion-textLight">
              Are you sure you want to delete user <strong class="text-notion-text">{{ userToDelete?.email }}</strong>?
              This action cannot be undone.
            </p>
          </div>
          <div class="bg-notion-hover px-6 py-4 flex justify-end space-x-3">
            <button
              type="button"
              @click="closeDeleteModal"
              class="btn-secondary"
            >
              Cancel
            </button>
            <button
              @click="handleDelete"
              class="px-4 py-2 bg-red-600 text-white text-sm font-medium rounded-md hover:bg-red-700 focus:outline-none transition-all duration-150"
              :disabled="deleteLoading"
            >
              <span v-if="!deleteLoading">Delete</span>
              <span v-else class="flex items-center">
                <svg class="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                </svg>
                Deleting...
              </span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { userAPI } from '@/api'

const router = useRouter()
const authStore = useAuthStore()

// State
const loading = ref(false)
const users = ref([])
const total = ref(0)
const limit = ref(50)
const offset = ref(0)
const successMessage = ref(null)
const errorMessage = ref(null)

// Modal states
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteModal = ref(false)
const createLoading = ref(false)
const editLoading = ref(false)
const deleteLoading = ref(false)
const userToDelete = ref(null)

// Forms
const createForm = reactive({
  email: '',
  password: '',
  role: 'user',
  isActive: true
})

const editForm = reactive({
  id: null,
  email: '',
  password: '',
  role: 'user',
  isActive: true
})

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

const formatDate = (dateString) => {
  if (!dateString) return 'Never'
  try {
    const date = new Date(dateString)

    // Check if date is invalid or represents a zero time (year 1 or before Unix epoch)
    if (isNaN(date.getTime()) || date.getFullYear() < 1970) {
      return 'Never'
    }

    const now = new Date()
    const diffMs = now - date
    const diffMins = Math.floor(diffMs / 60000)
    const diffHours = Math.floor(diffMs / 3600000)
    const diffDays = Math.floor(diffMs / 86400000)

    // Relative time for recent dates
    if (diffMins < 1) return 'Just now'
    if (diffMins < 60) return `${diffMins} ${diffMins === 1 ? 'minute' : 'minutes'} ago`
    if (diffHours < 24) return `${diffHours} ${diffHours === 1 ? 'hour' : 'hours'} ago`
    if (diffDays < 7) return `${diffDays} ${diffDays === 1 ? 'day' : 'days'} ago`

    // For older dates, show formatted date
    const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
    const month = months[date.getMonth()]
    const day = date.getDate()
    const year = date.getFullYear()
    const isCurrentYear = year === now.getFullYear()

    if (isCurrentYear) {
      return `${month} ${day}`
    } else {
      return `${month} ${day}, ${year}`
    }
  } catch {
    return 'Never'
  }
}

const loadUsers = async () => {
  loading.value = true
  errorMessage.value = null

  try {
    const response = await userAPI.getUsers({
      limit: limit.value,
      offset: offset.value
    })

    if (response.data) {
      users.value = response.data.users || []
      total.value = response.data._meta?.total || 0
    }
  } catch (err) {
    console.error('Failed to load users:', err)
    errorMessage.value = err.response?.data?.errorMessage || 'Failed to load users'
  } finally {
    loading.value = false
  }
}

const goToPage = (newOffset) => {
  if (newOffset < 0) return
  if (newOffset >= total.value) return
  offset.value = newOffset
  loadUsers()
}

const openCreateModal = () => {
  createForm.email = ''
  createForm.password = ''
  createForm.role = 'user'
  createForm.isActive = true
  showCreateModal.value = true
}

const closeCreateModal = () => {
  showCreateModal.value = false
  createForm.email = ''
  createForm.password = ''
  createForm.role = 'user'
  createForm.isActive = true
}

const handleCreate = async () => {
  createLoading.value = true
  errorMessage.value = null
  successMessage.value = null

  try {
    const response = await userAPI.createUser({
      email: createForm.email,
      password: createForm.password,
      role: createForm.role,
      isActive: createForm.isActive
    })

    if (response.data) {
      successMessage.value = 'User created successfully'
      closeCreateModal()
      loadUsers()
      setTimeout(() => {
        successMessage.value = null
      }, 3000)
    }
  } catch (err) {
    console.error('Failed to create user:', err)
    if (err.response?.data?.errorMessage) {
      errorMessage.value = err.response.data.errorMessage
    } else if (err.response?.data?.errors) {
      const errors = err.response.data.errors
      const errorList = Object.values(errors).flat().join(', ')
      errorMessage.value = `Validation errors: ${errorList}`
    } else {
      errorMessage.value = 'Failed to create user. Please try again.'
    }
  } finally {
    createLoading.value = false
  }
}

const openEditModal = (user) => {
  editForm.id = user.id
  editForm.email = user.email
  editForm.password = ''
  editForm.role = user.role
  editForm.isActive = user.isActive
  showEditModal.value = true
}

const closeEditModal = () => {
  showEditModal.value = false
  editForm.id = null
  editForm.email = ''
  editForm.password = ''
  editForm.role = 'user'
  editForm.isActive = true
}

const handleUpdate = async () => {
  editLoading.value = true
  errorMessage.value = null
  successMessage.value = null

  try {
    const payload = {
      email: editForm.email,
      role: editForm.role,
      isActive: editForm.isActive
    }

    // Only include password if it's provided
    if (editForm.password) {
      payload.password = editForm.password
    }

    const response = await userAPI.updateUser(editForm.id, payload)

    if (response.data) {
      successMessage.value = 'User updated successfully'
      closeEditModal()
      loadUsers()
      setTimeout(() => {
        successMessage.value = null
      }, 3000)
    }
  } catch (err) {
    console.error('Failed to update user:', err)
    if (err.response?.data?.errorMessage) {
      errorMessage.value = err.response.data.errorMessage
    } else if (err.response?.data?.errors) {
      const errors = err.response.data.errors
      const errorList = Object.values(errors).flat().join(', ')
      errorMessage.value = `Validation errors: ${errorList}`
    } else {
      errorMessage.value = 'Failed to update user. Please try again.'
    }
  } finally {
    editLoading.value = false
  }
}

const openDeleteModal = (user) => {
  userToDelete.value = user
  showDeleteModal.value = true
}

const closeDeleteModal = () => {
  showDeleteModal.value = false
  userToDelete.value = null
}

const handleDelete = async () => {
  if (!userToDelete.value) return

  deleteLoading.value = true
  errorMessage.value = null
  successMessage.value = null

  try {
    await userAPI.deleteUser(userToDelete.value.id)
    successMessage.value = 'User deleted successfully'
    closeDeleteModal()
    loadUsers()
    setTimeout(() => {
      successMessage.value = null
    }, 3000)
  } catch (err) {
    console.error('Failed to delete user:', err)
    if (err.response?.data?.errorMessage) {
      errorMessage.value = err.response.data.errorMessage
    } else {
      errorMessage.value = 'Failed to delete user. Please try again.'
    }
  } finally {
    deleteLoading.value = false
  }
}

onMounted(() => {
  loadUsers()
})
</script>
