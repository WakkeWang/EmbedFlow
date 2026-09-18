import { reactive } from 'vue'

// DS-4A: dark theme by default, light switchable (DS-5A tokens).
export const darkMode = reactive({ isDark: true })

export function toggleDark() {
	darkMode.isDark = !darkMode.isDark
}
