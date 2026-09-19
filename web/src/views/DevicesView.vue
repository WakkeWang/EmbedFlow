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
	NInputNumber,
	NSpace,
	useMessage,
} from 'naive-ui'
import { useI18n } from 'vue-i18n'
import { deviceApi, type Device } from '../api/http'
import { useAuth } from '../store/auth'

const { t } = useI18n()
const router = useRouter()
const message = useMessage()
const { isAdmin } = useAuth()

const devices = ref<Device[]>([])
const loading = ref(true)

// create/edit form (SSH fields ride along; requirement 3.2). Editing reuses
// the same modal; an empty password keeps the stored credential.
const showForm = ref(false)
const editingId = ref<number | null>(null)
const form = ref({ name: '', project: '', ssh_host: '', ssh_port: 22, ssh_user: '', ssh_password: '', note: '' })

// DS-2A: busy popup shows occupier user + session kind + elapsed duration.
const busyInfo = ref<Device | null>(null)
// Device deletion asks twice (the device may hold session history).
const deleteInfo = ref<Device | null>(null)

// ssh-test in flight per device id (requirement 3.2: manual probe button).
const testing = ref<Record<number, boolean>>({})
const testResult = ref<Record<number, { ok: boolean; detail: string }>>({})

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

function openNew() {
	editingId.value = null
	form.value = { name: '', project: '', ssh_host: '', ssh_port: 22, ssh_user: '', ssh_password: '', note: '' }
	showForm.value = true
}

function openEdit(d: Device) {
	editingId.value = d.id
	form.value = {
		name: d.name,
		project: d.project,
		ssh_host: d.ssh_host ?? '',
		ssh_port: d.ssh_port ?? 22,
		ssh_user: d.ssh_user ?? '',
		ssh_password: '', // never echoed back; empty keeps the stored one
		note: d.note ?? '',
	}
	showForm.value = true
}

async function saveDevice() {
	if (!form.value.name.trim()) return
	try {
		if (editingId.value) {
			await deviceApi.update(editingId.value, form.value)
		} else {
			await deviceApi.createV2(form.value)
		}
		showForm.value = false
		await load()
	} catch (e) {
		message.error(String(e))
	}
}

async function sshTest(d: Device) {
	testing.value = { ...testing.value, [d.id]: true }
	try {
		const res = await deviceApi.sshTest(d.id)
		testResult.value = { ...testResult.value, [d.id]: res }
	} catch (e) {
		testResult.value = { ...testResult.value, [d.id]: { ok: false, detail: String(e) } }
	} finally {
		testing.value = { ...testing.value, [d.id]: false }
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
			<NButton v-if="isAdmin" type="primary" @click="openNew">+</NButton>
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
					<div v-if="d.ssh_host" class="device-ssh mono">
						{{ d.ssh_user }}@{{ d.ssh_host }}:{{ d.ssh_port ?? 22 }}
						<NTag v-if="d.ssh_set" size="tiny" :bordered="false">{{ t('device.credSaved') }}</NTag>
					</div>
					<div v-if="testResult[d.id]" class="ssh-test-result" :class="{ ok: testResult[d.id].ok, fail: !testResult[d.id].ok }">
						{{ testResult[d.id].ok ? 'OK' : 'FAIL' }} — {{ testResult[d.id].detail }}
					</div>
					<template #footer>
						<NSpace>
							<NButton
								:type="d.busy ? 'default' : 'primary'"
								size="small"
								@click="openSession(d)"
							>
								{{ d.busy ? t('device.busy') : t('device.open') }}
							</NButton>
							<NButton v-if="d.ssh_host" size="small" :loading="!!testing[d.id]" @click="sshTest(d)">
								SSH Test
							</NButton>
							<NButton v-if="isAdmin" size="small" quaternary @click="openEdit(d)">
								{{ t('expect.edit') }}
							</NButton>
							<NButton v-if="isAdmin" size="small" quaternary type="error" @click="askDelete(d)">
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

		<!-- create/edit with SSH fields (requirement 3.2) -->
		<NModal
			:show="showForm"
			preset="dialog"
			:title="editingId ? t('device.edit') : t('device.list')"
			:show-icon="false"
			style="width: 480px"
		>
			<div class="form-grid">
				<label>{{ t('device.name') }}</label>
				<NInput v-model:value="form.name" />
				<label>{{ t('device.project') }}</label>
				<NInput v-model:value="form.project" />
				<label>{{ t('device.sshHost') }}</label>
				<NInput v-model:value="form.ssh_host" class="mono" placeholder="192.168.1.100" />
				<label>{{ t('device.sshPort') }}</label>
				<NInputNumber v-model:value="form.ssh_port" :min="1" :max="65535" style="width: 140px" />
				<label>{{ t('device.sshUser') }}</label>
				<NInput v-model:value="form.ssh_user" class="mono" placeholder="root" />
				<label>{{ t('device.sshPassword') }}</label>
				<div>
					<NInput
						v-model:value="form.ssh_password"
						type="password"
						show-password-on="click"
						class="mono"
						:placeholder="editingId ? t('device.passwordKeep') : ''"
					/>
					<div class="hint">{{ t('device.passwordHint') }}</div>
				</div>
				<label>{{ t('device.note') }}</label>
				<NInput v-model:value="form.note" />
			</div>
			<template #action>
				<NButton @click="showForm = false">{{ t('common.cancel') }}</NButton>
				<NButton type="primary" @click="saveDevice">{{ t('common.confirm') }}</NButton>
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
.header-row h2 {
	font-size: 20px;
	font-weight: 650;
}
.device-project {
	color: rgba(0, 0, 0, 0.45);
	font-size: 12px;
}
.device-ssh {
	font-size: 12px;
	color: rgba(0, 0, 0, 0.55);
	margin-top: 4px;
	display: flex;
	align-items: center;
	gap: 6px;
}
.ssh-test-result {
	margin-top: 8px;
	font-size: 12px;
	padding: 4px 8px;
	border-radius: 6px;
	background: rgba(0, 0, 0, 0.04);
	color: rgba(0, 0, 0, 0.6);
	word-break: break-all;
}
.ssh-test-result.ok {
	color: #18a058;
}
.mono {
	font-family: 'JetBrains Mono', Consolas, monospace;
}
.form-grid {
	display: grid;
	grid-template-columns: 110px 1fr;
	gap: 10px 12px;
	align-items: center;
}
.form-grid label {
	color: rgba(0, 0, 0, 0.55);
	font-size: 13px;
}
.hint {
	color: rgba(0, 0, 0, 0.4);
	font-size: 12px;
	margin-top: 2px;
}
</style>
