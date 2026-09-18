<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
	NGrid,
	NGridItem,
	NCard,
	NButton,
	NTag,
	NModal,
	NEmpty,
	NInput,
	NSpace,
	useMessage,
} from 'naive-ui'
import { useI18n } from 'vue-i18n'
import { deviceApi, type Device } from '../api/http'

const { t } = useI18n()
const router = useRouter()
const message = useMessage()

const devices = ref<Device[]>([])
const loading = ref(true)
const showNew = ref(false)
const newName = ref('')
const newProject = ref('')

// DS-2A: busy popup shows occupier user + session kind + elapsed duration.
const busyInfo = ref<Device | null>(null)
// Device deletion asks twice (the device may hold session history).
const deleteInfo = ref<Device | null>(null)

onMounted(load)

async function load() {
	loading.value = true
	try {
		devices.value = await deviceApi.list()
	} catch {
		devices.value = []
	} finally {
		loading.value = false
	}
}

async function createDevice() {
	if (!newName.value.trim()) return
	try {
		await deviceApi.create(newName.value.trim(), newProject.value.trim())
		showNew.value = false
		newName.value = ''
		newProject.value = ''
		await load()
	} catch (e) {
		message.error(String(e))
	}
}

function openSession(d: Device) {
	if (d.busy && d.owner) {
		// DS-2A: full occupier info in the popup; positive action = follow.
		busyInfo.value = d
		return
	}
	router.push(`/terminal/device-${d.id}`)
}

// CEO-15A: the second viewer follows read-only.
function followSession() {
	if (busyInfo.value) {
		router.push(`/terminal/device-${busyInfo.value.id}?follow=1`)
	}
	busyInfo.value = null
}

function askDelete(d: Device) {
	deleteInfo.value = d
}

async function confirmDelete() {
	if (!deleteInfo.value) return
	try {
		await deviceApi.remove(deleteInfo.value.id)
		deleteInfo.value = null
		await load()
	} catch (e) {
		message.error(String(e))
	}
}

function elapsed(d: Device) {
	if (!d.since) return ''
	const mins = Math.floor((Date.now() - new Date(d.since).getTime()) / 60000)
	return mins >= 1 ? `${mins}m` : '<1m'
}
</script>

<template>
	<div>
		<div class="header-row">
			<h2>{{ t('device.list') }}</h2>
			<NButton type="primary" @click="showNew = true">+</NButton>
		</div>

		<NEmpty v-if="!loading && devices.length === 0" :description="t('device.emptyHint')">
			<template #icon></template>
		</NEmpty>

		<NGrid v-else :cols="4" :x-gap="12" :y-gap="12" responsive="screen" item-responsive>
			<NGridItem v-for="d in devices" :key="d.id" span="1 s:1 m:2 l:1">
				<!-- DS-1A: card per device; name + status light + busy badge first. -->
				<NCard :title="d.name" size="small">
					<template #header-extra>
						<NTag v-if="d.busy" size="small" type="warning">
							{{ t('device.sessionOf', { user: d.owner, kind: d.kind === 'task' ? t('device.task') : t('device.manual') }) }}
							{{ elapsed(d) }}
						</NTag>
						<NTag v-else-if="d.shared" size="small" type="info">{{ t('device.online') }}</NTag>
						<NTag v-else size="small" type="default">{{ t('device.offline') }}</NTag>
					</template>
					<div class="device-project">{{ d.project }}</div>
					<template #footer>
						<NSpace>
							<NButton
								:type="d.busy ? 'default' : 'primary'"
								size="small"
								@click="openSession(d)"
							>
								{{ d.busy ? t('device.busy') : t('device.open') }}
							</NButton>
							<NButton size="small" quaternary type="error" @click="askDelete(d)">
								{{ t('device.delete') }}
							</NButton>
						</NSpace>
					</template>
				</NCard>
			</NGridItem>
		</NGrid>

		<!-- DS-2A: busy modal with occupier identity + read-only follow entry
		     (CEO-15A: a second viewer can join instead of walking away). -->
		<NModal
			:show="busyInfo !== null"
			preset="dialog"
			:title="t('device.busy')"
			:positive-text="t('terminal.readOnly')"
			:negative-text="t('common.cancel')"
			@positive-click="followSession"
			@negative-click="busyInfo = null"
		>
			<span v-if="busyInfo">
				{{ t('device.sessionOf', { user: busyInfo.owner, kind: busyInfo.kind === 'task' ? t('device.task') : t('device.manual') }) }}
				<span v-if="busyInfo.since"> · {{ t('device.duration', { n: elapsed(busyInfo) }) }}</span>
				<span v-if="busyInfo.kind === 'task'"> · {{ t('device.taskRunning') }}</span>
			</span>
		</NModal>

		<NModal :show="showNew" preset="dialog" :title="t('device.list')" :show-icon="false">
			<NInput v-model:value="newName" placeholder="name" style="margin-bottom: 8px" />
			<NInput v-model:value="newProject" placeholder="project" />
			<template #action>
				<NButton @click="showNew = false">{{ t('common.cancel') }}</NButton>
				<NButton type="primary" @click="createDevice">{{ t('common.confirm') }}</NButton>
			</template>
		</NModal>

		<!-- Device deletion: second confirmation (sessions history goes with it). -->
		<NModal
			:show="deleteInfo !== null"
			preset="dialog"
			:title="t('device.delete')"
			:content="t('device.deleteConfirm')"
			:positive-text="t('common.confirm')"
			:negative-text="t('common.cancel')"
			@positive-click="confirmDelete"
			@negative-click="deleteInfo = null"
		/>
	</div>
</template>

<style scoped>
.header-row {
	display: flex;
	align-items: center;
	justify-content: space-between;
	margin-bottom: 12px;
}
.device-project {
	color: rgba(255, 255, 255, 0.52);
	font-size: 12px;
}
</style>
