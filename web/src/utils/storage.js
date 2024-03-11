/**
 * LocalStorage utility functions for managing persisted data
 */

const USER_STORAGE_KEY = 'mut_user'

/**
 * Load user data from localStorage
 * @returns {Object|null} The stored user object or null if not found/invalid
 */
export const loadUserFromStorage = () => {
  try {
    const stored = localStorage.getItem(USER_STORAGE_KEY)
    return stored ? JSON.parse(stored) : null
  } catch (err) {
    console.error('Failed to load user from localStorage:', err)
    return null
  }
}

/**
 * Save user data to localStorage
 * @param {Object|null} userData - The user object to store, or null to clear
 */
export const saveUserToStorage = (userData) => {
  try {
    if (userData) {
      localStorage.setItem(USER_STORAGE_KEY, JSON.stringify(userData))
    } else {
      localStorage.removeItem(USER_STORAGE_KEY)
    }
  } catch (err) {
    console.error('Failed to save user to localStorage:', err)
  }
}

/**
 * Clear all user-related data from localStorage
 */
export const clearUserStorage = () => {
  try {
    localStorage.removeItem(USER_STORAGE_KEY)
  } catch (err) {
    console.error('Failed to clear user from localStorage:', err)
  }
}

