<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import {
	NCard,
	NButton,
	NInput,
	NInputNumber,
	NSelect,
	NTag,
	NEmpty,
	NCheckbox,
	NModal,
	NSpace,
	useMessage,
} from 'naive-ui'
import { useI18n } from 'vue-i18n'
import { buildItemApi, type BuildItemRecord } from '../api/http'
import { useProject } from '../store/project'
import { useAuth } from '../store/auth'

// Build item editor (requirement 2.1): source (git/local), command,
// artifact globs, timeout, version command, prerequisite groups
// (group-internal OR, group-to-group AND -- requirement 2.3).
const { t } = useI18n()
const message = useMessage()
const { currentId } = useProject()
const { isAdmin } = useAuth()

interface VM extends Partial<BuildItemRecord> {
	id?: number
	name: string
	source_type: 'git' | 'local'
	command: string
	artifactsText: string
	prereqGroups: number[][]
}

const items = ref<BuildItemRecord[]>([])
// All items across projects: prerequisite references may cross projects
// (requirement 2.1), so the picker lists everything with project labels.
const allItems = ref<BuildItemRecord[]>([])
const editor = ref<VM | null>(null)
const loading = ref(true)

// Prerequisite editor state: groups of item ids (from any project via raw
// number entry for cross-project refs; the picker lists this project's).
const prereqPickerOpen = ref(false)
const prereqGroupIndex = ref(0)

onMounted(load)

async function load() {
	if (!currentId.value) {
		items.value = []
		loading.value = false
		return
	}
	loading.value = true
	try {
		items.value = await buildItemApi.list(currentId.value)
	} catch {
		items.value = []
	}
	// Cross-project picker source (2.1); non-fatal if it fails.
	try {
		allItems.value = await buildItemApi.listAll()
	} catch {
		allItems.value = items.value
	} finally {
		loading.value = false
	}
}

function emptyVM(): VM {
	return { name: '', source_type: 'local', command: '', artifactsText: '', prereqGroups: [] }
}

function toVM(it: BuildItemRecord): VM {
	let groups: number[][] = []
	try {
		groups = JSON.parse(it.prereq_json || '[]') as number[][]
	} catch {
		groups = []
	}
	return {
		id: it.id,
		name: it.name,
		source_type: it.source_type,
		git_url: it.git_url,
		git_branch: it.git_branch,
		git_commit: it.git_commit,
		check_latest: it.check_latest,
		local_path: it.local_path,
		command: it.command,
		artifactsText: (it.artifacts ?? []).join('\n'),
		timeout_sec: it.timeout_sec,
		version_cmd: it.version_cmd,
		prereqGroups: groups,
	}
}

function newRule() {
	editor.value = emptyVM()
}

function editRule(it: BuildItemRecord) {
	editor.value = toVM(it)
}

function closeEditor() {
	editor.value = null
}

async function save() {
	if (!editor.value || !currentId.value) return
	if (!editor.value.name.trim() || !editor.value.command.trim()) {
		message.error(t('build.nameCommandRequired'))
		return
	}
	const artifacts = editor.value.artifactsText
		.split('\n')
		.map((s) => s.trim())
		.filter(Boolean)
	const body = {
		name: editor.value.name.trim(),
		source_type: editor.value.source_type,
		git_url: editor.value.git_url ?? '',
		git_branch: editor.value.git_branch ?? '',
		git_commit: editor.value.git_commit ?? '',
		check_latest: editor.value.check_latest ?? false,
		local_path: editor.value.local_path ?? '',
		command: editor.value.command,
		artifacts,
		timeout_sec: editor.value.timeout_sec ?? 0,
		version_cmd: editor.value.version_cmd ?? '',
		prereq_json: JSON.stringify(editor.value.prereqGroups ?? []),
	}
	try {
		if (editor.value.id) {
			await buildItemApi.update(editor.value.id, body)
		} else {
			await buildItemApi.create(currentId.value, body)
		}
		message.success(t('common.confirm'))
		editor.value = null
		await load()
	} catch (e) {
		message.error(String(e))
	}
}

async function removeItem(it: BuildItemRecord) {
	try {
		await buildItemApi.remove(it.id)
		await load()
	} catch (e) {
		message.error(String(e))
	}
}

// --- prerequisite group editing ---

function addGroup() {
	editor.value?.prereqGroups.push([])
}

function removeGroup(gi: number) {
	editor.value?.prereqGroups.splice(gi, 1)
}

const pickerOptions = computed(() =>
	allItems.value
		.filter((it) => it.id !== editor.value?.id)
		.map((it) => ({
			label: it.project_id === currentId.value ? it.name : `${it.name} (project ${it.project_id})`,
			value: it.id as number,
		})),
)

function openPicker(gi: number) {
	prereqGroupIndex.value = gi
	prereqPickerOpen.value = true
}

function pickPrereq(id: number) {
	const g = editor.value?.prereqGroups[prereqGroupIndex.value]
	if (g && !g.includes(id)) {
		g.push(id)
	}
	prereqPickerOpen.value = false
}

function removePrereq(gi: number, pi: number) {
	editor.value?.prereqGroups[gi].splice(pi, 1)
}

function prereqName(id: number): string {
	const it = allItems.value.find((x) => x.id === id)
	if (!it) return `#${id}`
	return it.project_id === currentId.value ? it.name : `${it.name} (project ${it.project_id})`
}

const sourceOptions = [
	{ label: 'git', value: 'git' },
	{ label: 'local', value: 'local' },
]
</script>

<template>
	<div>
		<div class="header-row">
			<h2>{{ t('build.itemsTitle') }}</h2>
			<NButton v-if="isAdmin" type="primary" @click="newRule">{{ t('build.newItem') }}</NButton>
		</div>

		<NEmpty v-if="!loading && items.length === 0 && !editor" :description="t('build.itemsEmpty')" />

		<div v-if="items.length > 0 && !editor" class="rule-list">
			<NCard v-for="it in items" :key="it.id" size="small" :title="it.name">
				<template #header-extra>
					<NTag size="small">{{ it.source_type }}</NTag>
				</template>
				<div class="mono item-cmd">{{ it.command }}</div>
				<NSpace>
					<NButton v-if="isAdmin" size="small" @click="editRule(it)">{{ t('expect.edit') }}</NButton>
					<NButton v-if="isAdmin" size="small" type="error" quaternary @click="removeItem(it)">
						{{ t('expect.delete') }}
					</NButton>
				</NSpace>
			</NCard>
		</div>

		<div v-if="editor" class="editor">
			<div class="header-row">
				<NInput v-model:value="editor.name" :placeholder="t('build.name')" style="width: 240px" />
				<NButton type="primary" @click="save">{{ t('expect.save') }}</NButton>
				<NButton quaternary size="small" @click="closeEditor">{{ t('expect.backToList') }}</NButton>
			</div>

			<NCard size="small" :title="t('build.source')" class="form-card">
				<div class="form-grid">
					<label>{{ t('build.sourceType') }}</label>
					<NSelect v-model:value="editor.source_type" :options="sourceOptions" style="width: 160px" />

					<template v-if="editor.source_type === 'git'">
						<label>git URL</label>
						<input v-model="editor.git_url" class="n-input mono" spellcheck="false" />
					</template>
					<template v-else>
						<label>{{ t('build.localPath') }}</label>
						<input v-model="editor.local_path" class="n-input mono" spellcheck="false" />
					</template>

					<!-- Branch/commit apply to both source kinds (local clone
					     takes --branch; a pinned commit builds the commit object,
					     skipping the worktree dirty gate). -->
					<label>{{ t('build.branch') }}</label>
					<input v-model="editor.git_branch" class="n-input mono" spellcheck="false" />
					<label>{{ t('build.commit') }}</label>
					<div>
						<input v-model="editor.git_commit" class="n-input mono" spellcheck="false" />
						<div v-if="editor.source_type === 'local'" class="hint">{{ t('build.localCommitHint') }}</div>
					</div>

					<template v-if="editor.source_type === 'git'">
						<label>{{ t('build.checkLatest') }}</label>
						<NCheckbox v-model:checked="editor.check_latest">{{ t('build.checkLatestHint') }}</NCheckbox>
					</template>
				</div>
			</NCard>

			<NCard size="small" :title="t('build.build')" class="form-card">
				<div class="form-grid">
					<label>{{ t('build.command') }}</label>
					<input v-model="editor.command" class="n-input mono" spellcheck="false" />
					<label>{{ t('build.artifacts') }}</label>
					<div>
						<textarea
							v-model="editor.artifactsText"
							class="n-input mono"
							rows="3"
							:placeholder="'bootfw*.bin\nu-boot.bin'"
							spellcheck="false"
						/>
						<div class="hint">{{ t('build.artifactsHint') }}</div>
					</div>
					<label>{{ t('build.timeout') }}</label>
					<NInputNumber v-model:value="editor.timeout_sec" :min="0" :step="60" style="width: 180px">
						<template #suffix>s</template>
					</NInputNumber>
					<label>{{ t('build.versionCmd') }}</label>
					<input v-model="editor.version_cmd" class="n-input mono" spellcheck="false" />
					<div></div>
					<div class="hint">{{ t('build.versionHint') }}</div>
				</div>
			</NCard>

			<NCard size="small" :title="t('build.prereqs')" class="form-card">
				<div class="hint prereq-hint">{{ t('build.prereqHint') }}</div>
				<div v-for="(g, gi) in editor.prereqGroups" :key="gi" class="prereq-group">
					<div class="prereq-head">
						<NTag size="small" type="info">{{ t('build.groupN', { n: gi + 1 }) }}</NTag>
						<span class="hint">{{ t('build.groupOr') }}</span>
						<NButton size="tiny" quaternary @click="openPicker(gi)">+</NButton>
						<NButton size="tiny" quaternary type="error" @click="removeGroup(gi)">
							{{ t('expect.delete') }}
						</NButton>
					</div>
					<div class="prereq-items">
						<NTag v-for="(pid, pi) in g" :key="pid" size="small" closable @close="removePrereq(gi, pi)">
							{{ prereqName(pid) }}
						</NTag>
						<span v-if="g.length === 0" class="hint">{{ t('build.groupEmpty') }}</span>
					</div>
				</div>
				<NButton size="small" @click="addGroup">{{ t('build.addGroup') }}</NButton>
			</NCard>

			<NModal :show="prereqPickerOpen" preset="dialog" :title="t('build.pickPrereq')" :show-icon="false">
				<NSelect :options="pickerOptions" @update:value="(v: number) => pickPrereq(v)" />
			</NModal>
		</div>
	</div>
</template>

<style scoped>
.header-row {
	display: flex;
	align-items: center;
	gap: 12px;
	margin-bottom: 12px;
}
.header-row h2 {
	font-size: 20px;
	font-weight: 650;
}
.rule-list {
	display: flex;
	flex-direction: column;
	gap: 8px;
	max-width: 720px;
}
.item-cmd {
	font-size: 13px;
	color: rgba(0, 0, 0, 0.55);
	margin-bottom: 8px;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}
.form-card {
	margin-bottom: 12px;
	max-width: 860px;
}
.form-grid {
	display: grid;
	grid-template-columns: 130px 1fr;
	gap: 10px 14px;
	align-items: center;
}
.form-grid label {
	color: rgba(0, 0, 0, 0.55);
	font-size: 13px;
}
.mono {
	font-family: 'JetBrains Mono', Consolas, monospace;
}
input.n-input,
textarea.n-input {
	width: 100%;
	padding: 6px 10px;
	border-radius: 6px;
	border: 1px solid rgba(0, 0, 0, 0.18);
	background: transparent;
	color: inherit;
	font-size: 13px;
}
textarea.n-input {
	resize: vertical;
}
.hint {
	color: rgba(0, 0, 0, 0.45);
	font-size: 12.5px;
}
.prereq-hint {
	margin-bottom: 10px;
}
.prereq-group {
	border: 1px dashed rgba(0, 0, 0, 0.15);
	border-radius: 8px;
	padding: 8px 12px;
	margin-bottom: 8px;
}
.prereq-head {
	display: flex;
	align-items: center;
	gap: 8px;
	margin-bottom: 6px;
}
.prereq-items {
	display: flex;
	gap: 6px;
	flex-wrap: wrap;
}
</style>
