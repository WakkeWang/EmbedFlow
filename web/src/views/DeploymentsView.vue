<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import {
	NButton,
	NSelect,
	NTag,
	NEmpty,
	NModal,
	NSpace,
	NInput,
	useMessage,
} from 'naive-ui'
import { useI18n } from 'vue-i18n'
import {
	deployApi,
	deployRuleApi,
	buildRecordApi,
	type DeployRecordRecord,
	type DeployRuleRecord,
	type BuildRecordRecord,
	type ArtifactRecord,
} from '../api/http'
import { useProject } from '../store/project'
import { Frame, encode, parseControl } from '../api/frames'

// Deployments (requirement 3.6): trigger = rule + build record + params;
// the record list shows live status via DeployEvent WS frames.
const { t } = useI18n()
const message = useMessage()
const { currentId } = useProject()

const records = ref<DeployRecordRecord[]>([])
const rules = ref<DeployRuleRecord[]>([])
const buildRecords = ref<BuildRecordRecord[]>([])
const loading = ref(true)

// trigger form state
const showTrigger = ref(false)
const trigRule = ref<number | null>(null)
const trigRecord = ref<number | null>(null)
const trigParams = ref<Record<string, string>>({})
const trigArtifacts = ref<ArtifactRecord[]>([])

// live detail state
const detail = ref<DeployRecordRecord | null>(null)
const detailLog = ref('')
const detailOpen = ref(false)
// flash cancel needs the second confirmation (requirement 3.6: 鍗婂埛鎻愮ず)
const cancelTarget = ref<DeployRecordRecord | null>(null)
// terminal records are deletable; the confirm dialog guards the click
const deleteTarget = ref<DeployRecordRecord | null>(null)

let ws: WebSocket | null = null
let subbed = new Set<number>()

const ruleOptions = computed(() =>
	rules.value.map((r) => ({ label: `${r.name} (${t('deploy.mode' + r.mode[0].toUpperCase() + r.mode.slice(1))})`, value: r.id as number })),
)
const recordOptions = computed(() =>
	buildRecords.value
		.filter((r) => r.status === 'succeeded')
		.map((r) => ({ label: `#${r.id} build ${r.item_id} ${r.status}`, value: r.id as number })),
)

onMounted(() => {
	load()
	connect()
})

onUnmounted(() => {
	if (ws) {
		ws.close()
		ws = null
	}
})

function wsURL(): string {
	const proto = location.protocol === 'https:' ? 'wss' : 'ws'
	const token = localStorage.getItem('embedflow.token') ?? ''
	return `${proto}://${location.host}/ws/client?token=${encodeURIComponent(token)}`
}


function connect() {
	ws = new WebSocket(wsURL())
	ws.onopen = () => {
		// The WS path requires an auth frame first (same as the terminal).
		ws?.send(encode(Frame.Auth, { token: localStorage.getItem('embedflow.token') ?? '', protocol_version: 1 }))
		resubscribe()
	}
	ws.onmessage = (ev) => {
		if (typeof ev.data !== 'string') return
		const f = parseControl(ev.data)
		if (f.type === Frame.AuthOK) {
			resubscribe()
			return
		}
		if (f.type === Frame.DeployEvent) {
			const body = f.body as { deploy_id: number; phase: string; detail?: string }
			onDeployEvent(body.deploy_id, body.phase, body.detail)
		}
	}
	ws.onclose = () => {
		// Reconnect with backoff; the sub set persists across reconnects.
		setTimeout(connect, 3000)
	}
}

function resubscribe() {
	for (const id of subbed) {
		ws?.send(encode(Frame.DeployCtrl, { command: 'subscribe', deploy_id: id }))
	}
}

function onDeployEvent(id: number, phase: string, detailText?: string) {
	const rec = records.value.find((r) => r.id === id)
	if (rec && phase !== 'log') {
		rec.status = phase as DeployRecordRecord['status']
		if (detailText) rec.detail = detailText
	}
	if (detail.value?.id === id) {
		if (phase === 'log') {
			detailLog.value += (detailText ?? '') + '\n'
		} else if (detailText) {
			detail.value.status = phase as DeployRecordRecord['status']
			detail.value.detail = detailText
		}
	}
}

async function load() {
	if (!currentId.value) {
		records.value = []
		loading.value = false
		return
	}
	loading.value = true
	try {
		records.value = await deployApi.list(currentId.value)
	} catch {
		records.value = []
	}
	try {
		rules.value = await deployRuleApi.list(currentId.value)
	} catch {
		rules.value = []
	}
	try {
		buildRecords.value = await buildRecordApi.list(currentId.value)
	} catch {
		buildRecords.value = []
	} finally {
		loading.value = false
	}
	// Follow running deployments.
	const running = records.value.filter((r) => r.status === 'running').map((r) => r.id)
	for (const id of running) {
		if (!subbed.has(id)) {
			subbed.add(id)
			ws?.send(encode(Frame.DeployCtrl, { command: 'subscribe', deploy_id: id }))
		}
	}
}

function openTrigger() {
	trigRule.value = null
	trigRecord.value = null
	trigParams.value = {}
	trigArtifacts.value = []
	showTrigger.value = true
}

async function onPickRecord(id: number | null) {
	trigArtifacts.value = []
	if (!id) return
	try {
		trigArtifacts.value = await buildRecordApi.artifacts(id)
	} catch {
		/* artifacts are optional */
	}
}

async function trigger() {
	if (!trigRule.value) return
	try {
		const res = await deployApi.trigger(trigRule.value, trigRecord.value ?? 0, trigParams.value)
		showTrigger.value = false
		message.success(t('deploy.triggered'))
		await load()
		const rec = records.value.find((x) => x.id === res.id)
		if (rec) openDetail(rec)
	} catch (e) {
		message.error(String(e))
	}
}

async function openDetail(r: DeployRecordRecord) {
	detail.value = r
	detailOpen.value = true
	detailLog.value = ''
	try {
		const tail = await deployApi.logTail(r.id)
		detailLog.value = tail.tail
	} catch {
		/* no log yet */
	}
	if (!subbed.has(r.id)) {
		subbed.add(r.id)
		ws?.send(encode(Frame.DeployCtrl, { command: 'subscribe', deploy_id: r.id }))
	}
}

function askCancel(r: DeployRecordRecord) {
	cancelTarget.value = r
}

async function confirmCancel() {
	if (!cancelTarget.value) return
	try {
		await deployApi.cancel(cancelTarget.value.id)
		message.info(t('deploy.cancelSent'))
	} catch (e) {
		message.error(String(e))
	} finally {
		cancelTarget.value = null
	}
}

function askDelete(r: DeployRecordRecord) {
	deleteTarget.value = r
}

async function confirmDelete() {
	if (!deleteTarget.value) return
	try {
		await deployApi.remove(deleteTarget.value.id)
		message.success(t('deploy.deleted'))
		if (detail.value?.id === deleteTarget.value.id) {
			detailOpen.value = false
		}
		await load()
	} catch (e) {
		message.error(String(e))
	} finally {
		deleteTarget.value = null
	}
}

// USB stick zip for the detail's rule + build record (requirement 3.1.3
// second half). A plain browser navigation carries the auth cookie? No --
// the API is Bearer-token gated, so fetch the blob with the header and
// save it client-side.
const usbBusy = ref(false)
async function downloadUSBZip() {
	const d = detail.value
	if (!d || !d.rule_id || !d.build_record_id) {
		message.warning(t('deploy.usbPickBuildFirst'))
		return
	}
	usbBusy.value = true
	try {
		const res = await fetch(`/api/deploy-rules/${d.rule_id}/usb-zip?build_record_id=${d.build_record_id}`, {
			headers: { Authorization: `Bearer ${localStorage.getItem('embedflow.token') ?? ''}` },
		})
		if (!res.ok) {
			let detailMsg = String(res.status)
			try {
				const body = await res.json()
				detailMsg = body.error ?? detailMsg
			} catch {
				/* non-JSON error body */
			}
			message.error(detailMsg)
			return
		}
		const blob = await res.blob()
		const url = URL.createObjectURL(blob)
		const a = document.createElement('a')
		a.href = url
		a.download = `usb-deploy-rule${d.rule_id}-build${d.build_record_id}.zip`
		a.click()
		URL.revokeObjectURL(url)
	} catch (e) {
		message.error(String(e))
	} finally {
		usbBusy.value = false
	}
}

const statusType = (s: string) =>
	s === 'succeeded' ? 'success' : s === 'failed' ? 'error' : s === 'running' ? 'info' : 'warning'

// The detail's rule carries the usb_disk section when stick packaging is
// configured (the button only shows then).
const detailHasUSB = computed(() => {
	const d = detail.value
	if (!d || !d.rule_id || !d.build_record_id) return false
	const rule = rules.value.find((x) => x.id === d.rule_id)
	return !!rule?.payload?.usb_disk
})
</script>

<template>
	<div>
		<div class="header-row">
			<h2>{{ t('deploy.recordsTitle') }}</h2>
			<NButton type="primary" :disabled="rules.length === 0" @click="openTrigger">{{ t('deploy.trigger') }}</NButton>
		</div>

		<NEmpty
			v-if="!loading && records.length === 0"
			:description="rules.length === 0 ? t('deploy.noRulesHint') : t('deploy.recordsEmpty')"
		/>

		<div v-else class="record-list">
			<div v-for="r in records" :key="r.id" class="record-row" @click="openDetail(r)">
				<span class="rec-id">#{{ r.id }}</span>
				<NTag size="small" :type="statusType(r.status)">{{ t('deploy.status_' + r.status) }}</NTag>
				<span class="muted">{{ r.rule_id ? t('deploy.ruleN', { n: r.rule_id }) : '' }}</span>
				<span class="muted">{{ r.device_id ? t('deploy.deviceN', { n: r.device_id }) : '' }}</span>
				<span class="muted">{{ r.build_record_id ? t('deploy.buildN', { n: r.build_record_id }) : '' }}</span>
				<span class="muted exec">{{ r.executor }}</span>
				<span class="rec-detail">{{ r.detail }}</span>
				<NButton
					v-if="r.status === 'running'"
					size="tiny"
					type="error"
					@click.stop="askCancel(r)"
				>
					{{ t('deploy.cancel') }}
				</NButton>
				<NButton
					v-if="r.status !== 'running'"
					size="tiny"
					quaternary
					@click.stop="askDelete(r)"
				>
					{{ t('deploy.delete') }}
				</NButton>
			</div>
		</div>

		<!-- trigger modal -->
		<NModal :show="showTrigger" preset="dialog" :title="t('deploy.trigger')" :show-icon="false" style="width: 560px" @update:show="showTrigger = $event">
			<div class="form-grid">
				<label>{{ t('deploy.pickRule') }}</label>
				<NSelect v-model:value="trigRule" :options="ruleOptions" />
				<label>{{ t('deploy.pickBuild') }}</label>
				<NSelect v-model:value="trigRecord" :options="recordOptions" clearable @update:value="onPickRecord" />
			</div>
			<div v-if="trigArtifacts.length" class="artifact-box">
				<div class="hint">{{ t('deploy.artifacts') }}</div>
				<div v-for="a in trigArtifacts" :key="a.id" class="artifact-line mono">
					{{ a.name }} <span class="muted">({{ a.size }} B, {{ a.checksum.slice(0, 12) }}鈥?</span>
				</div>
			</div>
			<template v-for="r in [rules.find((x) => x.id === trigRule)]" :key="r?.id">
				<div v-if="r?.payload.ssh_params?.length" class="param-box">
					<div class="hint">{{ t('deploy.deployParams') }}</div>
					<div v-for="p in r.payload.ssh_params" :key="p.name" class="param-row">
						<span class="mono">{{ p.name }}</span>
						<NInput v-model:value="trigParams[p.name]" :placeholder="p.value || ''" size="small" />
					</div>
				</div>
			</template>
			<template #action>
				<NButton @click="showTrigger = false">{{ t('common.cancel') }}</NButton>
				<NButton type="primary" :disabled="!trigRule" @click="trigger">
					{{ t('deploy.go') }}
				</NButton>
			</template>
		</NModal>

		<!-- detail drawer: status + log tail -->
		<NModal :show="detailOpen" preset="card" :title="t('deploy.detailN', { n: detail?.id })" style="width: 720px" @update:show="detailOpen = $event">
			<NSpace v-if="detail">
				<NTag :type="statusType(detail.status)">{{ t('deploy.status_' + detail.status) }}</NTag>
				<span class="muted">{{ detail.detail }}</span>
				<NButton v-if="detailHasUSB" size="small" type="primary" :loading="usbBusy" @click="downloadUSBZip">
					{{ t('deploy.usbDownload') }}
				</NButton>
			</NSpace>
			<pre class="log-box mono">{{ detailLog || t('deploy.noLog') }}</pre>
		</NModal>

		<!-- flash cancel: double confirmation (requirement 3.6, device may be
		     half-flashed; DS-7B) -->
		<NModal
			:show="cancelTarget !== null"
			preset="dialog"
			type="error"
			:title="t('deploy.cancelConfirmTitle')"
			:content="t('deploy.cancelConfirmBody')"
			:positive-text="t('deploy.cancelYes')"
			:negative-text="t('common.cancel')"
			@positive-click="confirmCancel"
			@negative-click="cancelTarget = null"
		/>

		<!-- record delete: terminal records only, guarded -->
		<NModal
			:show="deleteTarget !== null"
			preset="dialog"
			type="warning"
			:title="t('deploy.deleteConfirmTitle')"
			:content="t('deploy.deleteConfirmBody')"
			:positive-text="t('deploy.delete')"
			:negative-text="t('common.cancel')"
			@positive-click="confirmDelete"
			@negative-click="deleteTarget = null"
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
.record-list {
	display: flex;
	flex-direction: column;
	max-width: 900px;
}
.record-row {
	display: flex;
	align-items: center;
	gap: 10px;
	padding: 9px 10px;
	border-radius: 8px;
	cursor: pointer;
	font-size: 13.5px;
}
.record-row:hover {
	background: rgba(0, 0, 0, 0.04);
}
.rec-id {
	font-weight: 600;
	min-width: 42px;
}
.muted {
	color: rgba(0, 0, 0, 0.45);
	font-size: 12.5px;
}
.exec {
	min-width: 70px;
}
.rec-detail {
	flex: 1;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
	color: rgba(0, 0, 0, 0.6);
}
.form-grid {
	display: grid;
	grid-template-columns: 110px 1fr;
	gap: 10px 12px;
	align-items: center;
	margin-bottom: 10px;
}
.form-grid label {
	color: rgba(0, 0, 0, 0.55);
	font-size: 13px;
}
.artifact-box,
.param-box {
	border: 1px dashed rgba(0, 0, 0, 0.15);
	border-radius: 8px;
	padding: 8px 12px;
	margin-bottom: 8px;
}
.artifact-line {
	font-size: 12.5px;
	padding: 2px 0;
}
.param-row {
	display: flex;
	align-items: center;
	gap: 8px;
	margin: 4px 0;
}
.param-row .mono {
	font-family: 'JetBrains Mono', Consolas, monospace;
	font-size: 13px;
	min-width: 100px;
}
.log-box {
	background: rgba(0, 0, 0, 0.04);
	border-radius: 8px;
	padding: 12px;
	max-height: 420px;
	overflow: auto;
	font-family: 'JetBrains Mono', Consolas, monospace;
	font-size: 12.5px;
	white-space: pre-wrap;
	margin-top: 10px;
}
.hint {
	color: rgba(0, 0, 0, 0.45);
	font-size: 12.5px;
	margin-bottom: 4px;
}
.mono {
	font-family: 'JetBrains Mono', Consolas, monospace;
}
</style>
