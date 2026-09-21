<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import {
	NButton,
	NTag,
	NEmpty,
	NDataTable,
	NModal,
	NCheckbox,
	NDrawer,
	NDrawerContent,
	type DataTableColumns,
	useMessage,
} from 'naive-ui'
import { h } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import {
	buildRecordApi,
	buildItemApi,
	type BuildRecordRecord,
	type ArtifactRecord,
} from '../api/http'
import { useProject } from '../store/project'
import { statusType } from '../api/status'

// Build records across batches (requirement 2.4): history, log with phase
// segments, artifacts with checksums, delete with record/artifacts choice.
const { t } = useI18n()
const message = useMessage()
const route = useRoute()
const { currentId } = useProject()

const records = ref<BuildRecordRecord[]>([])
const itemNames = ref<Record<number, string>>({})

// Record detail drawer: log tail (phase-grouped) + artifacts.
const detail = ref<BuildRecordRecord | null>(null)
const logText = ref('')
const artifacts = ref<ArtifactRecord[]>([])
const showDetail = ref(false)
const deleteTarget = ref<BuildRecordRecord | null>(null)
const deleteMode = ref<'record' | 'artifacts'>('record')
const withArtifacts = computed({
	get: () => deleteMode.value === 'artifacts',
	set: (v: boolean) => (deleteMode.value = v ? 'artifacts' : 'record'),
})

// Delete-all (requirement 2.4: 删除/全部删除) with the same mode choice.
const showDeleteAll = ref(false)

// Row selection (0.80 feedback: select-all + delete selected, mirroring the
// batches page). Only terminal records are checkable: live ones are not
// deletable (the server skips them either way).
const checkedRecords = ref<number[]>([])
const deletableRecords = computed(() => records.value.filter((r) => r.status !== 'pending' && r.status !== 'building'))
const allRecordsChecked = computed(
	() => deletableRecords.value.length > 0 && deletableRecords.value.every((r) => checkedRecords.value.includes(r.id)),
)

function toggleRecord(id: number, v: boolean) {
	if (v) {
		if (!checkedRecords.value.includes(id)) checkedRecords.value = [...checkedRecords.value, id]
	} else {
		checkedRecords.value = checkedRecords.value.filter((x) => x !== id)
	}
}

function toggleAllRecords(v: boolean) {
	checkedRecords.value = v ? deletableRecords.value.map((r) => r.id) : []
}

// Selected-records delete reuses the delete-all endpoint's ids form, with
// the same record-vs-artifacts mode choice.
const deleteTargetRecords = ref<number[] | null>(null)

async function confirmDeleteRecords() {
	if (!deleteTargetRecords.value) return
	try {
		await buildRecordApi.removeMany(deleteTargetRecords.value, deleteMode.value)
		message.success(t('common.confirm'))
		checkedRecords.value = []
		await load()
	} catch (e) {
		message.error(String(e))
	} finally {
		deleteTargetRecords.value = null
	}
}

async function confirmDeleteAll() {
	if (!currentId.value) return
	try {
		await buildRecordApi.removeAll(currentId.value, deleteMode.value)
		showDeleteAll.value = false
		await load()
		message.success(t('common.confirm'))
	} catch (e) {
		message.error(String(e))
	}
}

onMounted(async () => {
	await load()
	// Deep link: ?open=<id> (from the batch drawer's detail button).
	const open = Number(route.query.open)
	if (open) await openRecordById(open)
})

// Project switch re-pull: records are project-scoped; close a live detail
// drawer and clear a stale selection.
watch(currentId, () => {
	if (showDetail.value) showDetail.value = false
	checkedRecords.value = []
	load()
})

async function load() {
	if (!currentId.value) {
		records.value = []
		return
	}
	try {
		const items = await buildItemApi.list(currentId.value)
		itemNames.value = Object.fromEntries(items.map((it) => [it.id, it.name]))
		records.value = await buildRecordApi.list(currentId.value)
	} catch {
		records.value = []
	}
}

async function openRecordById(id: number) {
	try {
		detail.value = await buildRecordApi.get(id)
		logText.value = ''
		artifacts.value = []
		showDetail.value = true
		try {
			logText.value = (await buildRecordApi.logTail(id)).tail
		} catch {
			logText.value = ''
		}
		try {
			artifacts.value = await buildRecordApi.artifacts(id)
		} catch {
			artifacts.value = []
		}
	} catch (e) {
		message.error(String(e))
	}
}

function openRecord(r: BuildRecordRecord) {
	openRecordById(r.id)
}

// Artifacts quick-download modal (0.80 feedback: 下载列).
const artifactsOpen = ref(false)
const artifactsFor = ref<BuildRecordRecord | null>(null)
const quickArts = ref<ArtifactRecord[]>([])

async function openArtifacts(r: BuildRecordRecord) {
	artifactsFor.value = r
	quickArts.value = []
	artifactsOpen.value = true
	try {
		quickArts.value = await buildRecordApi.artifacts(r.id)
	} catch {
		quickArts.value = []
	}
}

async function confirmDelete() {
	if (!deleteTarget.value) return
	try {
		await buildRecordApi.remove(deleteTarget.value.id, deleteMode.value)
		deleteTarget.value = null
		showDetail.value = false
		await load()
	} catch (e) {
		message.error(String(e))
	}
}

// Group log lines by phase tag: "<ts> [phase] content".
function segmentLog(text: string): { phase: string; lines: string[] }[] {
	const segs: { phase: string; lines: string[] }[] = []
	for (const line of (text || '').split('\n')) {
		const m = line.match(/^\S+ \S+ \[([a-z]+)\] (.*)$/)
		const phase = m?.[1] ?? 'other'
		const segText = m?.[2] ?? line
		const last = segs[segs.length - 1]
		if (last && last.phase === phase) {
			last.lines.push(segText)
		} else {
			segs.push({ phase, lines: [segText] })
		}
	}
	return segs
}

const logSegments = computed(() => segmentLog(logText.value))

function itemName(id: number): string {
	return itemNames.value[id] ?? `#${id}`
}

function fmtSize(n: number): string {
	if (n < 1024) return `${n} B`
	if (n < 1024 * 1024) return `${(n / 1024).toFixed(1)} KB`
	return `${(n / 1024 / 1024).toFixed(1)} MB`
}

function fmtTime(v: string): string {
	return v ? new Date(v).toLocaleString() : '-'
}

const artifactColumns = computed<DataTableColumns<ArtifactRecord>>(() => [
	{ title: t('build.name'), key: 'name' },
	{ title: t('build.size'), key: 'size', width: 100, render: (a) => fmtSize(a.size) },
	{
		title: 'checksum',
		key: 'checksum',
		render: (a) => h('span', { class: 'mono checksum' }, a.checksum),
	},
	{
		title: '',
		key: 'dl',
		width: 110,
		render: (a) =>
			h('a', { href: buildRecordApi.artifactDownloadURL(a.id), target: '_blank' }, t('build.downloadArtifacts')),
	},
])

// Zero-padded id display: #0001 (4 digits, grows past 9999 gracefully).
function fmtNo(n: number): string {
	return '#' + String(n).padStart(4, '0')
}

const columns = computed<DataTableColumns<BuildRecordRecord>>(() => [
	{
		// Selection: header checkbox = select all DELETABLE records (live
		// ones are never checkable), row checkboxes for the rest.
		title: () =>
			h(NCheckbox, {
				checked: allRecordsChecked.value,
				indeterminate: checkedRecords.value.length > 0 && !allRecordsChecked.value,
				disabled: deletableRecords.value.length === 0,
				onUpdateChecked: toggleAllRecords,
			}),
		key: 'check',
		width: 44,
		render: (r) =>
			h(NCheckbox, {
				checked: checkedRecords.value.includes(r.id),
				disabled: r.status === 'pending' || r.status === 'building',
				onUpdateChecked: (v: boolean) => toggleRecord(r.id, v),
			}),
	},
	{
		// The record id IS the entry point: click it to open the detail
		// drawer (log + artifacts + download), no separate actions column.
		title: t('build.recordNo'),
		key: 'id',
		width: 84,
		render: (r) =>
			h('a', {
				class: 'mono rec-link',
				href: 'javascript:;',
				onClick: () => openRecord(r),
			}, fmtNo(r.id)),
	},
	{ title: t('build.batch'), key: 'batch_id', width: 84, render: (r) => fmtNo(r.batch_id) },
	{ title: t('build.item'), key: 'item', width: 150, ellipsis: { tooltip: true }, render: (r) => itemName(r.item_id) },
	{
		title: t('history.state'),
		key: 'status',
		width: 82,
		render: (r) => h(NTag, { size: 'small', type: statusType(r.status) }, { default: () => t('build.rec_' + r.status) }),
	},
	{ title: t('build.branch'), key: 'branch', width: 110, render: (r) =>
		h('span', { class: 'mono', style: 'font-size:12px' }, r.branch || '-') },
	{ title: t('build.commitShort'), key: 'commit_sha', width: 86, render: (r) =>
		h('span', { class: 'mono', style: 'font-size:12px' }, r.commit_sha ? r.commit_sha.slice(0, 8) : '-') },
	{ title: t('history.start'), key: 'started_at', width: 140, render: (r) => fmtTime(r.started_at) },
	{
		// Artifacts quick download (0.80 feedback): one click lists the
		// record's artifacts for direct download, no drawer detour.
		title: t('build.artifactsShort'),
		key: 'artifacts',
		width: 90,
		render: (r) =>
			h(
				NButton,
				{
					size: 'tiny',
					quaternary: true,
					disabled: r.status !== 'succeeded',
					onClick: () => openArtifacts(r),
				},
				{ default: () => t('build.downloadArtifacts') },
			),
	},
])
</script>

<template>
	<div>
		<div class="header-row">
			<h2>{{ t('build.recordsTitle') }}</h2>
			<NButton v-if="records.length > 0" size="small" quaternary type="error" @click="showDeleteAll = true">
				{{ t('build.deleteAll') }}
			</NButton>
		</div>

		<NEmpty v-if="records.length === 0" :description="t('build.recordsEmpty')" />

		<template v-else>
			<div class="records-toolbar">
				<NCheckbox
					:checked="allRecordsChecked"
					:indeterminate="checkedRecords.length > 0 && !allRecordsChecked"
					@update:checked="toggleAllRecords"
				>
					{{ t('build.selectAll') }}
				</NCheckbox>
				<NButton
					size="small"
					type="error"
					quaternary
					:disabled="checkedRecords.length === 0"
					@click="deleteTargetRecords = [...checkedRecords]"
				>
					{{ t('build.deleteSelected', { n: checkedRecords.length }) }}
				</NButton>
			</div>
			<NDataTable
				:columns="columns"
				:data="records"
				size="small"
				class="records-table"
				:scroll-x="1028"
			/>
		</template>

		<NDrawer v-model:show="showDetail" :width="760">
			<NDrawerContent :title="`${t('build.record')} ${fmtNo(detail?.id ?? 0)}`" closable>
				<template v-if="detail">
					<div class="meta">
						<NTag size="small" :type="statusType(detail.status)">{{ t('build.rec_' + detail.status) }}</NTag>
						<span v-if="detail.commit_sha" class="mono meta-item">commit {{ detail.commit_sha.slice(0, 12) }}</span>
						<span class="meta-item">{{ detail.executor }}</span>
						<span class="flex1"></span>
						<a
							class="log-dl"
							:href="buildRecordApi.logDownloadURL(detail.id)"
							target="_blank"
						>{{ t('history.downloadLog') }}</a>
						<NButton
							size="tiny"
							type="error"
							quaternary
							:disabled="detail.status === 'pending' || detail.status === 'building'"
							@click="deleteTarget = detail"
						>
							{{ t('expect.delete') }}
						</NButton>
					</div>
					<pre v-if="detail.version_info" class="version-pre">{{ detail.version_info }}</pre>

					<h3 class="sec-title">{{ t('history.log') }}</h3>
					<div v-if="logSegments.length === 0" class="hint">{{ t('history.noLog') }}</div>
					<div v-for="(seg, si) in logSegments" :key="si" class="log-seg">
						<NTag size="tiny" :bordered="false">{{ seg.phase }}</NTag>
						<pre class="log-pre">{{ seg.lines.join('\n') }}</pre>
					</div>

					<h3 class="sec-title">{{ t('build.artifacts') }}</h3>
					<NDataTable
						v-if="artifacts.length > 0"
						:columns="artifactColumns"
						:data="artifacts"
						size="small"
					/>
					<div v-else class="hint">{{ t('build.noArtifacts') }}</div>
				</template>
			</NDrawerContent>
		</NDrawer>

		<NModal
			:show="deleteTarget !== null"
			preset="dialog"
			:title="t('expect.delete')"
			:positive-text="t('common.confirm')"
			:negative-text="t('common.cancel')"
			@positive-click="confirmDelete"
			@negative-click="deleteTarget = null"
		>
			<div v-if="deleteTarget">
				<NCheckbox v-model:checked="withArtifacts" size="small">{{ t('build.deleteArtifacts') }}</NCheckbox>
				<div class="hint" style="margin-top: 6px">{{ t('build.deleteHint') }}</div>
			</div>
		</NModal>

		<NModal
			:show="showDeleteAll"
			preset="dialog"
			:title="t('build.deleteAll')"
			:positive-text="t('common.confirm')"
			:negative-text="t('common.cancel')"
			@positive-click="confirmDeleteAll"
			@negative-click="showDeleteAll = false"
		>
			<div>
				<NCheckbox v-model:checked="withArtifacts" size="small">{{ t('build.deleteArtifacts') }}</NCheckbox>
				<div class="hint" style="margin-top: 6px">{{ t('build.deleteAllHint') }}</div>
			</div>
		</NModal>

		<!-- selected-records delete: same record-vs-artifacts choice -->
		<NModal
			:show="deleteTargetRecords !== null"
			preset="dialog"
			type="warning"
			:title="t('build.deleteSelectedTitle')"
			:positive-text="t('common.confirm')"
			:negative-text="t('common.cancel')"
			@positive-click="confirmDeleteRecords"
			@negative-click="deleteTargetRecords = null"
		>
			<div>
				<NCheckbox v-model:checked="withArtifacts" size="small">{{ t('build.deleteArtifacts') }}</NCheckbox>
				<div class="hint" style="margin-top: 6px">{{ t('build.deleteSelectedHint', { n: deleteTargetRecords?.length ?? 0 }) }}</div>
			</div>
		</NModal>

		<!-- artifacts quick download (the records table's 下载 column) -->
		<NModal
			:show="artifactsOpen"
			preset="card"
			:title="`${t('build.artifacts')} - ${t('build.record')} ${fmtNo(artifactsFor?.id ?? 0)}`"
			style="width: 620px"
			@update:show="artifactsOpen = $event"
		>
			<div v-if="quickArts.length > 0" class="quick-arts">
				<div v-for="a in quickArts" :key="a.id" class="quick-art-row">
					<span class="mono quick-name">{{ a.name }}</span>
					<span class="muted">{{ fmtSize(a.size) }}</span>
					<a :href="buildRecordApi.artifactDownloadURL(a.id)" target="_blank">{{ t('build.downloadArtifacts') }}</a>
				</div>
			</div>
			<div v-else class="hint">{{ t('build.noArtifacts') }}</div>
		</NModal>

	</div>
</template>

<style scoped>
.header-row h2 {
	font-size: 20px;
	font-weight: 650;
	margin-bottom: 12px;
}
.records-table {
	/* The column set (check/record/batch/item/status/branch/commit/start/
	   artifacts) sums past 1200px: no cap, the outer scroll owns
	   narrow viewports. */
	max-width: none;
}
.records-toolbar {
	display: flex;
	align-items: center;
	justify-content: space-between;
	margin-bottom: 8px;
}
.meta {
	display: flex;
	align-items: center;
	gap: 10px;
	margin-bottom: 10px;
}
.meta-item {
	font-size: 13px;
	color: rgba(0, 0, 0, 0.55);
}
.flex1 {
	flex: 1;
}
.mono {
	font-family: 'JetBrains Mono', Consolas, monospace;
}
.version-pre {
	background: rgba(0, 0, 0, 0.04);
	border-radius: 8px;
	padding: 10px 12px;
	font-size: 13px;
	white-space: pre-wrap;
	margin: 0 0 14px;
}
.sec-title {
	font-size: 15px;
	font-weight: 600;
	margin: 14px 0 8px;
}
.log-dl {
	font-size: 12.5px;
	font-weight: 400;
	margin-left: 10px;
}
.log-seg {
	margin-bottom: 8px;
}
.log-pre {
	background: rgba(0, 0, 0, 0.04);
	border-radius: 8px;
	padding: 8px 12px;
	font-family: 'JetBrains Mono', Consolas, monospace;
	font-size: 12px;
	white-space: pre-wrap;
	word-break: break-all;
	margin: 4px 0 0;
	max-height: 300px;
	overflow: auto;
}
.hint {
	color: rgba(0, 0, 0, 0.45);
	font-size: 13px;
}
.checksum {
	font-size: 12px;
	word-break: break-all;
}
.quick-arts {
	display: flex;
	flex-direction: column;
	gap: 8px;
}
.quick-art-row {
	display: flex;
	align-items: center;
	gap: 12px;
	padding: 8px 10px;
	border-radius: 8px;
	background: rgba(0, 0, 0, 0.03);
}
.quick-name {
	flex: 1;
	font-size: 13px;
	word-break: break-all;
}
.quick-art-row a {
	font-weight: 550;
}
</style>
