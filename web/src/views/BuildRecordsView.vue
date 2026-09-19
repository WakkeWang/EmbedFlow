<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import {
	NButton,
	NTag,
	NEmpty,
	NDataTable,
	NSpace,
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
	batchApi,
	buildRecordApi,
	buildItemApi,
	type BuildRecordRecord,
	type ArtifactRecord,
} from '../api/http'
import { useProject } from '../store/project'

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

onMounted(async () => {
	await load()
	// Deep link: ?open=<id> (from the batch drawer's detail button).
	const open = Number(route.query.open)
	if (open) await openRecordById(open)
})

async function load() {
	if (!currentId.value) {
		records.value = []
		return
	}
	try {
		const items = await buildItemApi.list(currentId.value)
		itemNames.value = Object.fromEntries(items.map((it) => [it.id, it.name]))
		records.value = await fetchRecords()
	} catch {
		records.value = []
	}
}

async function fetchRecords(): Promise<BuildRecordRecord[]> {
	// Records are listed per batch; walk the project's batches.
	try {
		const batches = await batchApi.list(currentId.value)
		const out: BuildRecordRecord[] = []
		for (const b of batches) {
			const res = await batchApi.get(b.id)
			out.push(...res.records)
		}
		return out
	} catch {
		return []
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

const logSegments = computed(() => {
	// Group log lines by phase tag: "<ts> [phase] content".
	const segs: { phase: string; lines: string[] }[] = []
	for (const line of (logText.value || '').split('\n')) {
		const m = line.match(/^\S+ \S+ \[([a-z]+)\] (.*)$/)
		const phase = m?.[1] ?? 'other'
		const text = m?.[2] ?? line
		const last = segs[segs.length - 1]
		if (last && last.phase === phase) {
			last.lines.push(text)
		} else {
			segs.push({ phase, lines: [text] })
		}
	}
	return segs
})

function statusType(s: string): 'default' | 'info' | 'success' | 'error' | 'warning' {
	switch (s) {
		case 'pending':
			return 'default'
		case 'building':
			return 'info'
		case 'succeeded':
			return 'success'
		case 'failed':
			return 'error'
		case 'canceled':
			return 'warning'
		default:
			return 'default'
	}
}

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
			h('a', { href: buildRecordApi.artifactDownloadURL(a.id), target: '_blank' }, t('history.download')),
	},
])

const columns = computed<DataTableColumns<BuildRecordRecord>>(() => [
	{ title: '#', key: 'id', width: 60 },
	{ title: t('build.batch'), key: 'batch_id', width: 80, render: (r) => `#${r.batch_id}` },
	{ title: t('build.item'), key: 'item', render: (r) => itemName(r.item_id) },
	{
		title: t('history.state'),
		key: 'status',
		width: 130,
		render: (r) => h(NTag, { size: 'small', type: statusType(r.status) }, { default: () => t('build.rec_' + r.status) }),
	},
	{ title: t('build.commit'), key: 'commit_sha', width: 120, render: (r) => (r.commit_sha ? r.commit_sha.slice(0, 8) : '-') },
	{ title: t('history.start'), key: 'started_at', render: (r) => fmtTime(r.started_at) },
	{
		title: '',
		key: 'actions',
		width: 170,
		render: (r) =>
			h(NSpace, { size: 4 }, {
				default: () => [
					h(NButton, { size: 'tiny', quaternary: true, onClick: () => openRecord(r) }, { default: () => t('build.detail') }),
					h(
						NButton,
						{
							size: 'tiny',
							quaternary: true,
							type: 'error',
							disabled: r.status === 'pending' || r.status === 'building',
							onClick: () => (deleteTarget.value = r),
						},
						{ default: () => t('expect.delete') },
					),
				],
			}),
	},
])
</script>

<template>
	<div>
		<div class="header-row">
			<h2>{{ t('build.recordsTitle') }}</h2>
		</div>

		<NEmpty v-if="records.length === 0" :description="t('build.recordsEmpty')" />

		<NDataTable v-else :columns="columns" :data="records" size="small" class="records-table" />

		<NDrawer v-model:show="showDetail" :width="760">
			<NDrawerContent :title="`${t('build.record')} #${detail?.id}`" closable>
				<template v-if="detail">
					<div class="meta">
						<NTag size="small" :type="statusType(detail.status)">{{ t('build.rec_' + detail.status) }}</NTag>
						<span v-if="detail.commit_sha" class="mono meta-item">commit {{ detail.commit_sha.slice(0, 12) }}</span>
						<span class="meta-item">{{ detail.executor }}</span>
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
	</div>
</template>

<style scoped>
.header-row h2 {
	font-size: 20px;
	font-weight: 650;
	margin-bottom: 12px;
}
.records-table {
	max-width: 1080px;
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
</style>
