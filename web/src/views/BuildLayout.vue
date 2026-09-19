<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

// Build section shell: sub-navigation tabs (build items / batches /
// records). Mirrors DeployLayout's tab pattern.
const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const tabs = computed(() => [
	{ key: 'items', label: t('nav.buildItems') },
	{ key: 'batches', label: t('nav.buildBatches') },
	{ key: 'records', label: t('nav.buildRecords') },
])

const activeKey = computed(() => route.path.split('/')[2] ?? 'items')
</script>

<template>
	<div class="deploy-shell">
		<div class="tabs">
			<button
				v-for="tab in tabs"
				:key="tab.key"
				class="tab"
				:class="{ active: activeKey === tab.key }"
				@click="router.push('/build/' + tab.key)"
			>
				{{ tab.label }}
			</button>
		</div>
		<router-view />
	</div>
</template>

<style scoped>
.tabs {
	display: flex;
	gap: 4px;
	margin-bottom: 18px;
	border-bottom: 1px solid rgba(0, 0, 0, 0.08);
}
.tab {
	appearance: none;
	background: none;
	border: none;
	border-bottom: 2px solid transparent;
	color: rgba(0, 0, 0, 0.5);
	font-size: 15px;
	font-weight: 500;
	font-family: inherit;
	padding: 10px 18px;
	cursor: pointer;
	transition: color 0.15s, border-color 0.15s;
}
.tab:hover {
	color: rgba(0, 0, 0, 0.8);
}
.tab.active {
	color: #2f6bff;
	border-bottom-color: #2f6bff;
}
</style>
