<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import {
	NDataTable,
	type DataTableColumns,
	NButton,
	NDrawer,
	NDrawerContent,
	NEmpty,
	NTag,
} from 'naive-ui'
import { useI18n } from 'vue-i18n'
import { h } from 'vue'
import { sessionApi, deviceApi, type SessionRecord, type Device } from '../api/http'

const { t } = useI18n()

const devices = ref<Device[]>([])
const records = ref<SessionRecord[]>([])
const showDrawer = ref(false)

onMounted(async () => {
	try {
		devices.value = await deviceApi.list()
	} catch {
		devices.value = []
	}
})

async function openHistory(d: Device) {
	try {
		records.value = await sessionApi.history(d.id)
	} catch {
		records.value = []
	}
	showDrawer.value = true
}

function fmtTime(v: string) {
	return v ? new Date(v).toLocaleString() : '-'
}

function fmtDuration(r: SessionRecord) {
	if (!r.ended_at) return '-'
	const ms = new Date(r.ended_at).getTime() - new Date(r.started_at).getTime()
	if (ms < 0) return '-'
	const s = Math.floor(ms / 1000)
	const m = Math.floor(s / 60)
	return m > 0 ? `${m}m ${s % 60}s` : `${s}s`
}

const columns = computed<DataTableColumns<SessionRecord>>(() => [
	{ title: '#', key: 'id', width: 60 },
	{ title: t('history.state'), key: 'state', width: 100, render: renderState },
	{ title: t('history.start'), key: 'started_at', render: (r) => fmtTime(r.started_at) },
	{ title: t('history.end'), key: 'ended_at', render: (r) => fmtTime(r.ended_at) },
	{ title: t('history.duration'), key: 'dur', render: fmtDuration },
	{ title: t('history.log'), key: 'log', render: renderLog },
])

function renderState(r: SessionRecord) {
	const label =
		r.state === 'active'
			? t('common.active')
			: r.state === 'failed'
				? t('common.failed')
				: t('common.closed')
	const type = r.state === 'active' ? 'info' : r.state === 'failed' ? 'error' : 'default'
	// Log-incomplete gets its own warning badge (design doc disk-full gap).
	if (r.log_incomplete) {
		return h('span', [h(NTag, { size: 'small', type }, { default: () => label }), ' !'])
	}
	return h(NTag, { size: 'small', type }, { default: () => label })
}

function renderLog(r: SessionRecord) {
	return h('a', { href: sessionApi.logDownloadURL(r.id), target: '_blank' }, t('history.download'))
}
</script>

<template>
	<!-- DS-7C: per-device history drawer. -->
	<div>
		<h2>{{ t('history.title') }}</h2>
		<div class="device-chips">
			<NButton v-for="d in devices" :key="d.id" quaternary @click="openHistory(d)">
				{{ d.name }}
			</NButton>
		</div>

		<NDrawer v-model:show="showDrawer" :width="680">
			<NDrawerContent :title="t('history.title')" closable>
				<NEmpty v-if="records.length === 0" :description="t('history.empty')" />
				<NDataTable v-else :columns="columns" :data="records" size="small" />
			</NDrawerContent>
		</NDrawer>
	</div>
</template>

<style scoped>
.device-chips {
	display: flex;
	gap: 8px;
	flex-wrap: wrap;
	margin-bottom: 12px;
}
</style>
