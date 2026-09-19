import { ref, computed } from 'vue'

// Auth state (requirement 1.5: two-level roles). The login response's role
// rides localStorage so a refresh keeps the admin/member distinction.
const ROLE_KEY = 'embedflow.role'

const role = ref<string>(localStorage.getItem(ROLE_KEY) ?? '')

export function setRole(r: string) {
	role.value = r
	localStorage.setItem(ROLE_KEY, r)
}

export function clearRole() {
	role.value = ''
	localStorage.removeItem(ROLE_KEY)
}

export function useAuth() {
	const isAdmin = computed(() => role.value === 'admin')
	return { role, isAdmin }
}
