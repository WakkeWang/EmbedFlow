<script setup lang="ts">
import { computed } from 'vue'
import { NConfigProvider, NMessageProvider, darkTheme, zhCN, dateZhCN } from 'naive-ui'
import { darkMode } from './theme'

// M1.5 visual pass v2 (user feedback 2026-09-19): light default, 16px base
// type, Linear/Vercel-grade neutral palette with a blue brand voice.
// Backgrounds: sidebar #ffffff, content #f6f7f9 (set in LayoutView CSS).
const themeOverrides = computed(() => ({
	common: {
		primaryColor: '#2f6bff',
		primaryColorHover: '#5083ff',
		primaryColorPressed: '#1f52d6',
		primaryColorSuppl: '#5083ff',

		fontFamily:
			"-apple-system, 'Segoe UI Variable', 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', 'Noto Sans CJK SC', sans-serif",
		fontFamilyMono: "'JetBrains Mono', 'Cascadia Code', Consolas, 'Courier New', monospace",
		fontSize: '16px',
		fontSizeMedium: '16px',
		fontSizeSmall: '15px',
		fontSizeLarge: '17px',
		fontSizeHuge: '18px',

		borderRadius: '8px',
		borderRadiusSmall: '6px',

		heightMedium: '38px',
		heightSmall: '32px',
		heightLarge: '44px',

		bodyColor: '#f6f7f9',
		cardColor: '#ffffff',
		modalColor: '#ffffff',
		popoverColor: '#ffffff',
	},
	Card: {
		paddingMedium: '20px 22px',
		paddingSmall: '16px 18px',
		titleFontSizeMedium: '17px',
		titleFontWeight: '600',
		borderRadius: '14px',
	},
	Menu: {
		fontSizeMedium: '20px',
		fontSizeLarge: '20px',
		itemHeightMedium: '56px',
		borderRadius: '10px',
	},
	DataTable: {
		fontSizeMedium: '15px',
		thPaddingMedium: '12px 14px',
		tdPaddingMedium: '12px 14px',
	},
	Button: {
		fontSizeMedium: '16px',
		fontSizeSmall: '15px',
		fontWeight: '500',
	},
	Input: {
		heightMedium: '38px',
	},
	Tag: {
		fontSizeSmall: '14px',
		fontSizeMedium: '15px',
	},
	...(darkMode.isDark ? {} : {}),
}))

const theme = computed(() => (darkMode.isDark ? darkTheme : null))
</script>

<template>
  <NConfigProvider
    :theme="theme"
    :locale="zhCN"
    :date-locale="dateZhCN"
    :theme-overrides="themeOverrides"
  >
    <NMessageProvider>
      <router-view />
    </NMessageProvider>
  </NConfigProvider>
</template>
