<script setup lang="ts">
import { ref, onMounted, computed, onBeforeUnmount } from 'vue'
import {
	NCard,
	NButton,
	NTag,
	NEmpty,
	NDataTable,
	NModal,
	type DataTableColumns,
	useMessage,
} from 'naive-ui'
import { h } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import {
	batchApi,
	buildItemApi,
	type BatchRecord,
	type BuildItemRecord,
	type BuildRecordRecord,
} from '../api/http'
import { useProject } from '../store/project'
import { useAuth } from '../store/auth'
import { WSClient } from '../api/ws'
import { Frame, type BuildEventFrame } from '../api/frames'

// Batch trigger + list (requirement 2.3/2.4). Selection closes over
// prerequisites client-side for the "N items will build" preview; the
// server's kernel is authoritative.
const { t } = useI18n()
const message = useMessage()
const router = useRouter()
const { currentId } = useProject()
const { isAdmin } = useAuth()

const items = ref<BuildItemRecord[]>([])
const batches = ref<BatchRecord[]>([])
const checked = ref<number[]>([])
const triggering = ref(false)

// Closure preview: which items will actually build (selected + prereqs).
function closureOf(selectedIds: number[]): number[] {
	const byId = new Map(items.value.map((it) => [it.id, it]))
	const seen = new Set<number>()
	const walk = (id: number) => {
		if (seen.has(id)) return
		seen.add(id)
		const it = byId.get(id)
		if (!it) return
		let groups: number[][] = []
		try {
			groups = JSON.parse(it.prereq_json || '[]') as number[][]
		} catch {
			groups = []
		}
		for (const g of groups) for (const p of g) walk(p)
	}
	for (const id of selectedIds) walk(id)
	return [...seen]
}

const closurePreview = computed(() => closureOf(checked.value))
const pulledIn = computed(
	() => closurePreview.value.filter((id) => !checked.value.includes(id)),
)

// Live batch detail (WS-subscribed while the drawer is open).
const detailBatch = ref<BatchRecord | null>(null)
const detailRecords = ref<BuildRecordRecord[]>([])
const showDetail = ref(false)
let ws: WSClient | null = null

onMounted(async () => {
	await load()
	ws = new WSClient({ onControl: onFrame })
	ws.connect()
})

onBeforeUnmount(() => {
	if (detailBatch.value) {
		ws?.sendControl(Frame.BuildCtrl, { command: 'unsubscribe', batch_id: detailBatch.value.id })
	}
	ws?.close()
})

async function load() {
	if (!currentId.value) {
		items.value = []
		batches.value = []
		return
	}
	try {
		items.value = await buildItemApi.list(currentId.value)
	} catch {
		items.value = []
	}
	try {
		batches.value = await batchApi.list(currentId.value)
	} catch {
		batches.value = []
	}
}

async function trigger() {
	if (!currentId.value || checked.value.length === 0) return
	triggering.value = true
	try {
		const res = await batchApi.create(currentId.value, checked.value)
		checked.value = []
		message.success(t('build.batchCreated'))
		await load()
		await openDetail(res.id)
	} catch (e) {
		message.error(String(e))
	} finally {
		triggering.value = false
	}
}

async function openDetail(id: number) {
	try {
		const res = await batchApi.get(id)
		detailBatch.value = res.batch
		detailRecords.value = res.records
		showDetail.value = true
		ws?.sendControl(Frame.BuildCtrl, { command: 'subscribe', batch_id: id })
	} catch (e) {
		message.error(String(e))
	}
}

function closeDetail() {
	if (detailBatch.value) {
		ws?.sendControl(Frame.BuildCtrl, { command: 'unsubscribe', batch_id: detailBatch.value.id })
	}
	showDetail.value = false
}

function onFrame(f: { type: number; body: any }) {
	if (f.type !== Frame.BuildEvent) return
	const ev = f.body as BuildEventFrame
	if (!detailBatch.value || ev.batch_id !== detailBatch.value.id) return
	applyEvent(ev)
}

function applyEvent(ev: BuildEventFrame) {
	if (ev.phase === 'batch_done' || ev.phase === 'batch_status') {
		if (detailBatch.value && (ev.phase === 'batch_done' || !detailBatch.value.status || detailBatch.value.status === 'queued' || detailBatch.value.status === 'running')) {
			detailBatch.value.status = (ev.detail as BatchRecord['status']) ?? detailBatch.value.status
		}
		if (ev.phase === 'batch_done') {
			load() // refresh list badges
		}
		return
	}
	if (!ev.record_id) return
	const i = detailRecords.value.findIndex((r) => r.id === ev.record_id)
	if (i < 0) return
	switch (ev.phase) {
		case 'started':
			detailRecords.value[i].status = 'building'
			break
		case 'succeeded':
			detailRecords.value[i].status = 'succeeded'
			break
		case 'failed':
			detailRecords.value[i].status = 'failed'
			break
		case 'skipped':
			detailRecords.value[i].status = 'skipped'
			break
		case 'canceled':
			detailRecords.value[i].status = 'canceled'
			break
	}
}

async function cancelBatch(b: BatchRecord) {
	try {
		await batchApi.cancel(b.id)
		message.success(t('common.confirm'))
		await load()
		if (detailBatch.value?.id === b.id) await openDetail(b.id)
	} catch (e) {
		message.error(String(e))
	}
}

function statusType(s: string): 'default' | 'info' | 'success' | 'error' | 'warning' {
	switch (s) {
		case 'running':
		case 'building':
			return 'info'
		case 'completed':
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

function itemNameOf(record: BuildRecordRecord): string {
	const found = items.value.find((it) => it.id === record.item_id)
	return found?.name ?? `#${record.item_id}`
}

function gotoRecord(r: BuildRecordRecord) {
	router.push(`/build/records?open=${r.id}`)
}

const detailColumns = computed<DataTableColumns<BuildRecordRecord>>(() => [
	{ title: '#', key: 'id', width: 60 },
	{ title: t('build.item'), key: 'item', render: (r) => itemNameOf(r) },
	{
		title: t('history.state'),
		key: 'status',
		width: 110,
		render: (r) => h(NTag, { size: 'small', type: statusType(r.status) }, { default: () => t('build.rec_' + r.status) }),
	},
	{ title: '', key: 'go', width: 90, render: (r) => h(NButton, { size: 'tiny', quaternary: true, onClick: () => gotoRecord(r) }, { default: () => t('build.detail') }) },
])

const itemOptions = computed(() =>
	items.value.map((it) => ({ label: it.name, value: it.id as number })),
)
</script>

<template>
	<div>
		<div class="header-row">
			<h2>{{ t('build.batchesTitle') }}</h2>
		</div>

		<NEmpty v-if="items.length === 0" :description="t('build.itemsEmpty')" />

		<template v-else>
			<NCard size="small" :title="t('build.trigger')" class="trigger-card">
				<NSelect
					v-model:value="checked"
					multiple
					:options="itemOptions"
					:placeholder="t('build.pickItems')"
					style="margin-bottom: 10px"
				/>
				<div class="trigger-row">
					<span class="hint">
						{{
							t('build.willBuild', {
								n: closurePreview.length,
								m: pulledIn.length,
							})
						}}
					</span>
					<NButton type="primary" :loading="triggering" :disabled="checked.length === 0" @click="trigger">
						{{ t('build.trigger') }}
					</NButton>
				</div>
			</NCard>

			<NEmpty v-if="batches.length === 0" :description="t('build.batchesEmpty')" style="margin-top: 20px" />
			<NCard v-for="b in batches" :key="b.id" size="small" class="batch-card" @click="openDetail(b.id)">
				<div class="batch-row">
					<span class="batch-id mono">#{{ b.id }}</span>
					<NTag size="small" :type="statusType(b.status)">{{ t('build.batch_' + b.status) }}</NTag>
					<span class="hint">{{ b.created_by }}</span>
					<span class="flex1"></span>
					<NButton
						v-if="(b.status === 'queued' || b.status === 'running') && isAdmin"
						size="tiny"
						type="error"
						quaternary
						@click.stop="cancelBatch(b)"
					>
						{{ t('expect.abort') }}
					</NButton>
				</div>
			</NCard>
		</template>

		<NModal :show="showDetail" preset="card" :title="`${t('build.batch')} #${detailBatch?.id}`" style="width: 760px" @after-leave="closeDetail">
			<NDataTable :columns="detailColumns" :data="detailRecords" size="small" />
		</NModal>
	</div>
</template>

<style scoped>
.header-row h2 {
	font-size: 20px;
	font-weight: 650;
	margin-bottom: 12px;
}
.trigger-card {
	max-width: 860px;
	margin-bottom: 16px;
}
.trigger-row {
	display: flex;
	align-items: center;
	justify-content: space-between;
	gap: 12px;
}
.hint {
	color: rgba(0, 0, 0, 0.45);
	font-size: 13px;
}
.batch-card {
	max-width: 860px;
	margin-bottom: 8px;
	cursor: pointer;
}
.batch-row {
	display: flex;
	align-items: center;
	gap: 10px;
}
.batch-id {
	font-weight: 600;
}
.flex1 {
	flex: 1;
}
.mono {
	font-family: 'JetBrains Mono', Consolas, monospace;
}
</style>
