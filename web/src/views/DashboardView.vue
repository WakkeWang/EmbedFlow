<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NGrid, NGridItem, NCard, NButton, NEmpty, NTag, NModal, NInput, NSpace, useMessage } from 'naive-ui'
import { useI18n } from 'vue-i18n'
import { useProject } from '../store/project'
import { projectApi, getToken } from '../api/http'

const { t } = useI18n()
const router = useRouter()
const message = useMessage()
const { projects, current, load, create, remove } = useProject()

const showNew = ref(false)
const newName = ref('')
const newNote = ref('')
const deleteTarget = ref<{ id: number; name: string } | null>(null)

// Import: a JSON bundle file loads into the currently selected project.
const importInput = ref<HTMLInputElement | null>(null)

onMounted(load)

const sections = [
	{ key: 'build', path: '/build/items', glyph: 'hammer', state: 'live' },
	{ key: 'deploy', path: '/deploy/devices', glyph: 'chip', state: 'live' },
	{ key: 'test', path: '/test', glyph: 'flask', state: 'M4' },
	{ key: 'release', path: '/release', glyph: 'rocket', state: 'M5' },
] as const

async function createProject() {
	if (!newName.value.trim()) return
	try {
		await create(newName.value.trim(), newNote.value.trim())
		showNew.value = false
		newName.value = ''
		newNote.value = ''
		message.success(t('project.create'))
	} catch (e) {
		message.error(String(e))
	}
}

async function confirmDelete() {
	if (!deleteTarget.value) return
	try {
		await remove(deleteTarget.value.id)
		deleteTarget.value = null
	} catch (e) {
		message.error(String(e))
	}
}

// Export the project's rules (+ devices) as one JSON download. The export
// endpoint is Bearer-gated like every API route, so fetch with the header
// instead of a plain navigation.
async function exportProject(id: number, name: string) {
	try {
		const res = await fetch(projectApi.exportURL(id), {
			headers: { Authorization: `Bearer ${getToken()}` },
		})
		if (!res.ok) {
			message.error(String(res.status))
			return
		}
		const blob = await res.blob()
		const url = URL.createObjectURL(blob)
		const a = document.createElement('a')
		a.href = url
		a.download = `embedflow-export-${name}.json`
		a.click()
		URL.revokeObjectURL(url)
	} catch (e) {
		message.error(String(e))
	}
}

async function onImportFile(ev: Event) {
	const file = (ev.target as HTMLInputElement).files?.[0]
	;(ev.target as HTMLInputElement).value = ''
	if (!file) return
	if (!current.value) {
		message.warning(t('project.importFileFirst'))
		return
	}
	let bundle: Record<string, any>
	try {
		bundle = JSON.parse(await file.text())
	} catch {
		message.error(t('project.importBad'))
		return
	}
	try {
		const res = await projectApi.importBundle(current.value.id, bundle)
		message.success(t('project.importedN', { n: res.imported }))
	} catch (e) {
		message.error(t('project.importBad'))
	}
}

function go(path: string) {
	if (!current.value) return
	router.push(path)
}
</script>

<template>
	<div>
		<div class="header-row">
			<div>
				<h2>{{ t('dashboard.title') }}</h2>
				<div v-if="current" class="subtitle">
					{{ current.name }}<span v-if="current.note" class="note"> · {{ current.note }}</span>
				</div>
			</div>
			<NButton type="primary" @click="showNew = true">{{ t('project.new') }}</NButton>
		</div>

		<NEmpty v-if="projects.length === 0" :description="t('dashboard.empty')" class="empty-block">
			<template #extra>
				<div class="empty-extra">
					<span class="hint">{{ t('dashboard.hint') }}</span>
					<NButton type="primary" size="small" @click="showNew = true">{{ t('project.new') }}</NButton>
				</div>
			</template>
		</NEmpty>

		<template v-else>
			<NGrid :cols="4" :x-gap="14" :y-gap="14" responsive="screen" item-responsive>
				<NGridItem v-for="s in sections" :key="s.key" span="4 m:2 l:1">
					<NCard
						size="small"
						class="section-card"
						:class="{ 'section-live': s.state === 'live' }"
						hoverable
						@click="go(s.path)"
					>
						<div class="section-head">
							<span class="glyph" :data-g="s.glyph"></span>
							<NTag
								size="tiny"
								:type="s.state === 'live' ? 'success' : 'default'"
								:bordered="false"
							>
								{{ s.state === 'live' ? t('dashboard.live') : s.state }}
							</NTag>
						</div>
						<div class="section-title">{{ t('dashboard.sections.' + s.key) }}</div>
						<div class="section-desc muted">
							{{ s.state === 'live' ? t('dashboard.deployHint') : t('dashboard.' + s.key + 'Empty') }}
						</div>
					</NCard>
				</NGridItem>
			</NGrid>

			<div class="project-list">
				<div class="list-head">
					<h3>{{ t('project.list') }}</h3>
					<NSpace>
						<NButton size="tiny" @click="importInput?.click()">{{ t('project.import') }}</NButton>
						<input ref="importInput" type="file" accept=".json,application/json" style="display: none" @change="onImportFile" />
					</NSpace>
				</div>
				<div v-for="p in projects" :key="p.id" class="project-row" :class="{ active: p.id === current?.id }">
					<span class="project-dot"></span>
					<span class="project-name">{{ p.name }}</span>
					<span class="muted note">{{ p.note }}</span>
					<NButton size="tiny" quaternary class="row-del" @click="exportProject(p.id, p.name)">
						{{ t('project.export') }}
					</NButton>
					<NButton
						size="tiny"
						quaternary
						type="error"
						class="row-del"
						@click="deleteTarget = { id: p.id, name: p.name }"
					>
						{{ t('project.delete') }}
					</NButton>
				</div>
			</div>
		</template>

		<NModal :show="showNew" preset="dialog" :title="t('project.new')" :show-icon="false">
			<NInput v-model:value="newName" :placeholder="t('project.name')" style="margin-bottom: 8px" />
			<NInput v-model:value="newNote" :placeholder="t('project.desc')" />
			<template #action>
				<NButton @click="showNew = false">{{ t('common.cancel') }}</NButton>
				<NButton type="primary" @click="createProject">{{ t('project.create') }}</NButton>
			</template>
		</NModal>

		<NModal
			:show="deleteTarget !== null"
			preset="dialog"
			:title="t('project.delete')"
			:content="t('project.deleteConfirm')"
			:positive-text="t('common.confirm')"
			:negative-text="t('common.cancel')"
			@positive-click="confirmDelete"
			@negative-click="deleteTarget = null"
		/>
	</div>
</template>

<style scoped>
.header-row {
	display: flex;
	align-items: flex-start;
	justify-content: space-between;
	margin-bottom: 18px;
}
.header-row h2 {
	font-size: 22px;
	font-weight: 650;
	margin-bottom: 2px;
}
.subtitle {
	color: rgba(0, 0, 0, 0.55);
	font-size: 14px;
}
.subtitle .note {
	opacity: 0.7;
}
.empty-block {
	margin-top: 80px;
}
.empty-extra {
	display: flex;
	flex-direction: column;
	align-items: center;
	gap: 12px;
}
.hint {
	color: rgba(0, 0, 0, 0.45);
	font-size: 13px;
	max-width: 420px;
	line-height: 1.6;
}
.section-card {
	cursor: pointer;
	transition: transform 0.15s ease, border-color 0.15s ease;
}
.section-card:hover {
	transform: translateY(-2px);
}
.section-card.section-live {
	border-color: rgba(47, 107, 255, 0.4);
}
.section-head {
	display: flex;
	align-items: center;
	justify-content: space-between;
	margin-bottom: 10px;
}
.glyph {
	width: 34px;
	height: 34px;
	border-radius: 9px;
	background: rgba(0, 0, 0, 0.05);
	position: relative;
	flex: 0 0 auto;
}
.glyph::after {
	content: '';
	position: absolute;
	inset: 0;
	margin: auto;
	width: 14px;
	height: 14px;
	background: rgba(0, 0, 0, 0.55);
	border-radius: 3px;
}
.glyph[data-g='chip']::after {
	background: #2f6bff;
	border-radius: 2px;
}
.section-title {
	font-size: 16px;
	font-weight: 600;
	margin-bottom: 4px;
}
.section-desc {
	line-height: 1.5;
	min-height: 38px;
}
.muted {
	color: rgba(0, 0, 0, 0.45);
	font-size: 13px;
}
.project-list {
	margin-top: 26px;
}
.list-head {
	display: flex;
	align-items: center;
	justify-content: space-between;
}
.list-head h3 {
	font-size: 15px;
	font-weight: 600;
	margin-bottom: 8px;
}
.project-row {
	display: flex;
	align-items: center;
	gap: 12px;
	padding: 10px 12px;
	border-radius: 8px;
}
.project-row:hover {
	background: rgba(0, 0, 0, 0.04);
}
.project-row.active {
	background: rgba(47, 107, 255, 0.08);
}
.project-dot {
	width: 8px;
	height: 8px;
	border-radius: 50%;
	background: rgba(0, 0, 0, 0.2);
	flex: 0 0 auto;
}
.project-row.active .project-dot {
	background: #2f6bff;
}
.project-name {
	font-size: 14px;
	font-weight: 550;
	min-width: 140px;
}
.note {
	flex: 1;
}
.row-del {
	opacity: 0;
	transition: opacity 0.15s;
}
.project-row:hover .row-del {
	opacity: 1;
}
</style>
