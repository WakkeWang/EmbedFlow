<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { NLayout, NLayoutSider, NLayoutContent, NMenu, NButton } from 'naive-ui'
import { useI18n } from 'vue-i18n'
import { darkMode, toggleDark } from '../theme'
import { clearToken } from '../api/http'

const { t, locale } = useI18n()
const router = useRouter()

const menuOptions = computed(() => [
	{ label: t('nav.devices'), key: '/devices' },
	{ label: t('nav.expect'), key: '/expect' },
	{ label: t('nav.history'), key: '/history' },
	{ label: t('nav.settings'), key: '/settings', disabled: true },
])

function onMenu(key: string) {
	router.push(key)
}

function switchLang() {
	locale.value = locale.value === 'zh' ? 'en' : 'zh'
}

function logout() {
	clearToken()
	router.push('/login')
}
</script>

<template>
	<!-- DS-1D: Soybean-style left sidebar; M1 items: devices / expect / history / settings. -->
	<NLayout position="absolute" class="root">
		<NLayoutSider bordered :width="200" class="sider">
			<div class="brand">EmbedFlow</div>
			<NMenu :options="menuOptions" @update:value="onMenu" />
			<div class="sider-footer">
				<NButton quaternary size="small" @click="toggleDark">
					{{ darkMode.isDark ? 'Light' : 'Dark' }}
				</NButton>
				<NButton quaternary size="small" @click="switchLang">
					{{ t('nav.language') }}
				</NButton>
				<NButton quaternary size="small" @click="logout">
					{{ t('nav.logout') }}
				</NButton>
			</div>
		</NLayoutSider>
		<NLayoutContent content-style="padding: 16px;" class="content">
			<router-view />
		</NLayoutContent>
	</NLayout>
</template>

<style scoped>
.root {
	height: 100vh;
}
.sider {
	display: flex;
	flex-direction: column;
}
.sider :deep(.n-layout-sider-scroll-container) {
	display: flex;
	flex-direction: column;
}
.sider-footer {
	margin-top: auto;
	padding: 12px;
	display: flex;
	flex-direction: column;
	gap: 4px;
}
.content {
	height: 100vh;
}
</style>
