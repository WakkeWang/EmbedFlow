<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NLayout, NLayoutSider, NLayoutContent, NMenu, NButton } from 'naive-ui'
import { useI18n } from 'vue-i18n'
import { darkMode, toggleDark } from '../theme'
import { clearToken } from '../api/http'
import ProjectSwitcher from '../components/ProjectSwitcher.vue'
import { useProject } from '../store/project'

const { t, locale } = useI18n()
const router = useRouter()
const { load } = useProject()

onMounted(load)

// DS-1D evolved: the sidebar mirrors the requirement-1.6 organization —
// one container (project) above, the four sections below. Settings is the
// admin page (requirement 1.6); non-admins are redirected by the view.
const menuOptions = computed(() => [
	{ label: t('nav.dashboard'), key: '/' },
	{ label: t('nav.build'), key: '/build/items' },
	{ label: t('nav.deploy'), key: '/deploy/devices' },
	{ label: t('nav.test'), key: '/test' },
	{ label: t('nav.release'), key: '/release' },
	{ label: t('nav.settings'), key: '/settings' },
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
	<!-- Full-height sidebar rail (user feedback: the rail owns one column). -->
	<NLayout position="absolute" class="root" :class="{ dark: darkMode.isDark }">
		<NLayoutSider bordered class="sider">
			<div class="brand">
				<div class="brand-mark">EF</div>
				<div class="brand-text">
					<div class="brand-name">EmbedFlow</div>
					<div class="brand-sub">{{ t('project.title') }}</div>
				</div>
			</div>
			<ProjectSwitcher />
			<NMenu :options="menuOptions" @update:value="onMenu" class="nav-menu" />
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
		<NLayoutContent class="content-outer">
			<div class="content-scroll">
				<router-view />
			</div>
		</NLayoutContent>
	</NLayout>
</template>

<style scoped>
.root {
	height: 100vh;
}
.root:not(.dark) {
	background: #f6f7f9;
}
.root.dark {
	background: #101014;
}
.sider {
	display: flex;
	flex-direction: column;
	height: 100vh;
	background: #ffffff;
}
.root.dark .sider {
	background: #17171c;
}
.sider :deep(.n-layout-sider-scroll-container) {
	display: flex;
	flex-direction: column;
	height: 100%;
}
.brand {
	display: flex;
	align-items: center;
	gap: 11px;
	padding: 20px 18px 16px;
}
.brand-mark {
	width: 38px;
	height: 38px;
	border-radius: 10px;
	background: linear-gradient(135deg, #2f6bff, #1f52d6);
	color: #fff;
	font-weight: 700;
	font-size: 15px;
	display: flex;
	align-items: center;
	justify-content: center;
	letter-spacing: 0.5px;
	flex: 0 0 auto;
}
.brand-name {
	font-size: 17px;
	font-weight: 700;
	letter-spacing: 0.2px;
	line-height: 1.2;
	color: #1f52d6;
}
.root.dark .brand-name {
	color: #7ea2ff;
}
.brand-sub {
	font-size: 11px;
	color: rgba(0, 0, 0, 0.38);
	letter-spacing: 1px;
	text-transform: uppercase;
}
.root.dark .brand-sub {
	color: rgba(255, 255, 255, 0.4);
}
.nav-menu {
	flex: 1;
	padding: 0 8px;
}
.nav-menu :deep(.n-menu-item-content) {
	font-weight: 550;
	letter-spacing: 0.2px;
}
.sider-footer {
	padding: 12px;
	display: flex;
	flex-direction: column;
	gap: 4px;
	border-top: 1px solid rgba(0, 0, 0, 0.07);
}
.root.dark .sider-footer {
	border-top-color: rgba(255, 255, 255, 0.08);
}
.sider-footer .n-button {
	justify-content: flex-start;
}
.content-outer {
	height: 100vh;
	background: #f6f7f9;
}
.root.dark .content-outer {
	background: #101014;
}
.content-scroll {
	height: 100vh;
	overflow-y: auto;
	padding: 26px 32px 48px;
	max-width: 1360px;
}
</style>
