import { reactive } from 'vue'

// DS-4A revised (user feedback 2026-09-19): light theme by default; the
// dark option stays switchable. The light surface is #f6f7f9 (content) on
// white (sidebar), not clinical white-on-white.
export const darkMode = reactive({ isDark: false })

export function toggleDark() {
	darkMode.isDark = !darkMode.isDark
}
