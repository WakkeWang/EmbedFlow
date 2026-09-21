<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import {
	NCard,
	NButton,
	NInput,
	NInputNumber,
	NSelect,
	NTag,
	NEmpty,
	NSpace,
	NModal,
	useMessage,
} from 'naive-ui'
import { useI18n } from 'vue-i18n'
import { deployRuleApi, deviceApi, projectApi, type DeployRuleRecord, type DeployPayload, type Device, type ProjectRecord } from '../api/http'
import { useProject } from '../store/project'
import { useAuth } from '../store/auth'

// Deploy rule editor (requirement 3.1): three forms -- manual (markdown
// steps), ssh (command templates + params), flash (expect steps).
const { t } = useI18n()
const message = useMessage()
const { currentId } = useProject()
const { isAdmin } = useAuth()

interface RuleVM extends DeployRuleRecord {
	payload: DeployPayload
}

const rules = ref<RuleVM[]>([])
const devices = ref<Device[]>([])
const loading = ref(true)
const editor = ref<(Partial<RuleVM> & { payload: DeployPayload }) | null>(null)
// flash mode edits raw expect-step JSON in v1 (the dedicated card editor
// stays on the M1 flash-rules page; flash rules there are the same engine).
const flashStepsText = ref('[]')

const modeOptions = computed(() => [
	{ label: t('deploy.modeManual'), value: 'manual' },
	{ label: t('deploy.modeSSH'), value: 'ssh' },
	{ label: t('deploy.modeFlash'), value: 'flash' },
])

const deviceOptions = computed(() =>
	devices.value.map((d) => ({ label: d.name, value: d.id as number })),
)

onMounted(load)
// Project switch re-pull: rules are project-scoped; drop a stale editor.
watch(currentId, () => {
	editor.value = null
	load()
})

async function load() {
	if (!currentId.value) {
		rules.value = []
		loading.value = false
		return
	}
	loading.value = true
	try {
		rules.value = await deployRuleApi.list(currentId.value)
	} catch {
		rules.value = []
	}
	try {
		devices.value = await deviceApi.list()
	} catch {
		devices.value = []
	} finally {
		loading.value = false
	}
}

function newRule() {
	editor.value = {
		name: '',
		payload: { mode: 'manual', steps_md: '', ssh_commands: [], ssh_params: [] },
	}
	flashStepsText.value = '[]'
}

// Ensure the optional usb_disk section exists for editing (undefined in
// stored payloads that predate it, and in new rules).
function usbOf(p: DeployPayload): NonNullable<DeployPayload['usb_disk']> {
	if (!p.usb_disk) {
		p.usb_disk = { root_name: '', layout: [], checksums: [] }
	}
	return p.usb_disk
}

function editRule(r: RuleVM) {
	editor.value = { ...r, payload: { ...r.payload } }
	flashStepsText.value = r.payload.flash_steps_json || '[]'
}

async function save() {
	if (!editor.value || !currentId.value) return
	if (!editor.value.name?.trim()) {
		message.error(t('deploy.nameRequired'))
		return
	}
	const p = editor.value.payload
	if (p.mode === 'flash') {
		// Validate the step JSON server-side anyway; fail fast client-side.
		try {
			JSON.parse(flashStepsText.value)
		} catch {
			message.error(t('deploy.badFlashJSON'))
			return
		}
		p.flash_steps_json = flashStepsText.value
	}
	const body = {
		name: editor.value.name?.trim(),
		mode: p.mode,
		steps_md: p.steps_md,
		ssh_device_id: p.ssh_device_id,
		ssh_commands: p.ssh_commands,
		ssh_params: p.ssh_params,
		ssh_timeout_sec: p.ssh_timeout_sec,
		flash_device_id: p.flash_device_id,
		flash_steps: p.flash_steps_json,
		flash_timeout_sec: p.flash_timeout_sec,
		usb_disk: p.usb_disk && (p.usb_disk.root_name || p.usb_disk.layout?.length || p.usb_disk.checksums?.length)
			? p.usb_disk
			: null,
	}
	try {
		if (editor.value.id) {
			await deployRuleApi.update(editor.value.id, body)
		} else {
			await deployRuleApi.create(currentId.value, body)
		}
		message.success(t('common.confirm'))
		editor.value = null
		await load()
	} catch (e) {
		message.error(String(e))
	}
}

async function removeRule(r: RuleVM) {
	try {
		await deployRuleApi.remove(r.id)
		await load()
	} catch (e) {
		message.error(String(e))
	}
}

function addParam() {
	editor.value?.payload.ssh_params?.push({ name: '', value: '' })
}

function removeParam(i: number) {
	editor.value?.payload.ssh_params?.splice(i, 1)
}

function addCommand() {
	editor.value?.payload.ssh_commands?.push('')
}

function removeCommand(i: number) {
	editor.value?.payload.ssh_commands?.splice(i, 1)
}

function addUsbLayout(p: DeployPayload) {
	usbOf(p).layout?.push({ glob: '', dir: '.' })
}

function removeUsbLayout(p: DeployPayload, i: number) {
	usbOf(p).layout?.splice(i, 1)
}

function addUsbChecksum(p: DeployPayload) {
	usbOf(p).checksums?.push({ name: '', glob: '' })
}

function removeUsbChecksum(p: DeployPayload, i: number) {
	usbOf(p).checksums?.splice(i, 1)
}

// Cross-project copy (requirement 1.6): pick a target project, the rule
// lands there as a new row with the same payload.
const copyTarget = ref<DeployRuleRecord | null>(null)
const copyProjects = ref<ProjectRecord[]>([])
const copyTo = ref<number | null>(null)

async function askCopy(r: DeployRuleRecord) {
	copyTarget.value = r
	if (copyProjects.value.length === 0) {
		try {
			copyProjects.value = await projectApi.list()
		} catch {
			copyProjects.value = []
		}
	}
	const others = copyProjects.value.filter((p) => p.id !== currentId.value)
	copyTo.value = others[0]?.id ?? null
}

async function confirmCopy() {
	if (!copyTarget.value || !copyTo.value) return
	try {
		await projectApi.copyConfigObject(copyTarget.value.id, copyTo.value)
		message.success(t('project.copied'))
	} catch (e) {
		message.error(String(e))
	} finally {
		copyTarget.value = null
	}
}

const modeTag = (m: string) =>
	m === 'ssh' ? 'info' : m === 'flash' ? 'warning' : 'default'
</script>

<template>
	<div>
		<div class="header-row">
			<h2>{{ t('deploy.rulesTitle') }}</h2>
			<NButton v-if="isAdmin" type="primary" @click="newRule">{{ t('deploy.newRule') }}</NButton>
		</div>

		<NEmpty v-if="!loading && rules.length === 0 && !editor" :description="t('deploy.rulesEmpty')" />

		<div v-if="rules.length > 0 && !editor" class="rule-list">
			<NCard v-for="r in rules" :key="r.id" size="small" :title="r.name">
				<template #header-extra>
					<NTag size="small" :type="modeTag(r.mode)">{{ t('deploy.mode' + r.mode[0].toUpperCase() + r.mode.slice(1)) }}</NTag>
				</template>
				<div class="muted">
					{{
						r.mode === 'manual'
							? t('deploy.modeManualDesc')
							: r.mode === 'ssh'
								? t('deploy.modeSSHDesc', { n: r.payload.ssh_commands?.length ?? 0 })
								: t('deploy.modeFlashDesc')
					}}
				</div>
				<NSpace>
					<NButton v-if="isAdmin" size="small" @click="editRule(r)">{{ t('expect.edit') }}</NButton>
					<NButton v-if="isAdmin" size="small" @click="askCopy(r)">{{ t('project.copyTo') }}</NButton>
					<NButton v-if="isAdmin" size="small" type="error" quaternary @click="removeRule(r)">
						{{ t('expect.delete') }}
					</NButton>
				</NSpace>
			</NCard>
		</div>

		<div v-if="editor" class="editor">
			<div class="header-row">
				<NInput v-model:value="editor.name" :placeholder="t('deploy.ruleName')" style="width: 240px" />
				<NButton type="primary" @click="save">{{ t('expect.save') }}</NButton>
				<NButton quaternary size="small" @click="editor = null">{{ t('expect.backToList') }}</NButton>
			</div>

			<NCard size="small" :title="t('deploy.mode')" class="form-card">
				<NSelect v-model:value="editor.payload.mode" :options="modeOptions" style="width: 200px" />
			</NCard>

			<!-- USB stick packaging (requirement 3.1.3 second half): optional
			     for every mode -- the zip is generated from a build record's
			     artifacts on demand. -->
			<NCard size="small" :title="t('deploy.usbTitle')" class="form-card">
				<div class="hint">{{ t('deploy.usbHint') }}</div>
				<div class="form-grid">
					<label>{{ t('deploy.usbRoot') }}</label>
					<NInput v-model:value="usbOf(editor.payload).root_name" :placeholder="t('deploy.usbRootHint')" style="width: 320px" />
				</div>
				<div class="hint">{{ t('deploy.usbLayoutHint') }}</div>
				<div v-for="(l, i) in usbOf(editor.payload).layout" :key="'l' + i" class="param-row">
					<NInput v-model:value="l.glob" :placeholder="t('deploy.usbGlobPh')" style="width: 260px" />
					<span class="muted">-></span>
					<NInput v-model:value="l.dir" placeholder="." style="width: 180px" />
					<NButton size="tiny" quaternary type="error" @click="removeUsbLayout(editor.payload, i)">{{ t('expect.delete') }}</NButton>
				</div>
				<NButton size="small" @click="addUsbLayout(editor.payload)">{{ t('deploy.usbAddLayout') }}</NButton>
				<div class="hint" style="margin-top: 10px">{{ t('deploy.usbCsumHint') }}</div>
				<div v-for="(c, i) in usbOf(editor.payload).checksums" :key="'c' + i" class="param-row">
					<NInput v-model:value="c.name" :placeholder="t('deploy.usbFieldPh')" style="width: 160px" />
					<span class="muted">=</span>
					<NInput v-model:value="c.glob" :placeholder="t('deploy.usbGlobPh')" style="width: 260px" />
					<NButton size="tiny" quaternary type="error" @click="removeUsbChecksum(editor.payload, i)">{{ t('expect.delete') }}</NButton>
				</div>
				<NButton size="small" @click="addUsbChecksum(editor.payload)">{{ t('deploy.usbAddCsum') }}</NButton>
			</NCard>

			<!-- manual: markdown steps (requirement 3.1.1) -->
			<NCard v-if="editor.payload.mode === 'manual'" size="small" :title="t('deploy.stepsMD')" class="form-card">
				<textarea
					v-model="editor.payload.steps_md"
					class="n-input mono"
					rows="8"
					:placeholder="'1. Download artifact\n2. Copy to /tmp\n3. Run installer'"
					spellcheck="false"
				/>
				<div class="hint">{{ t('deploy.stepsMDHint') }}</div>
			</NCard>

			<!-- ssh: command templates + params (requirement 3.1.2) -->
			<template v-if="editor.payload.mode === 'ssh'">
				<NCard size="small" :title="t('deploy.sshDevice')" class="form-card">
					<div class="form-grid">
						<label>{{ t('deploy.targetDevice') }}</label>
						<NSelect v-model:value="editor.payload.ssh_device_id" :options="deviceOptions" style="width: 240px" clearable />
						<label>{{ t('deploy.sshTimeout') }}</label>
						<NInputNumber v-model:value="editor.payload.ssh_timeout_sec" :min="0" :step="30" style="width: 180px">
							<template #suffix>s</template>
						</NInputNumber>
					</div>
				</NCard>
				<NCard size="small" :title="t('deploy.commands')" class="form-card">
					<div class="hint">{{ t('deploy.commandsHint') }}</div>
					<div v-for="(_, i) in editor.payload.ssh_commands" :key="i" class="cmd-row">
						<input v-model="editor.payload.ssh_commands![i]" class="n-input mono" spellcheck="false" />
						<NButton size="tiny" quaternary type="error" @click="removeCommand(i)">{{ t('expect.delete') }}</NButton>
					</div>
					<NButton size="small" @click="addCommand">{{ t('deploy.addCommand') }}</NButton>
				</NCard>
				<NCard size="small" :title="t('deploy.params')" class="form-card">
					<div class="hint">{{ t('deploy.paramsHint') }}</div>
					<div v-for="(p, i) in editor.payload.ssh_params" :key="i" class="param-row">
						<NInput v-model:value="p.name" :placeholder="t('deploy.paramName')" style="width: 160px" />
						<NInput v-model:value="p.value" :placeholder="t('deploy.paramDefault')" style="width: 220px" />
						<NButton size="tiny" quaternary type="error" @click="removeParam(i)">{{ t('expect.delete') }}</NButton>
					</div>
					<NButton size="small" @click="addParam">{{ t('deploy.addParam') }}</NButton>
				</NCard>
			</template>

			<!-- flash: expect steps JSON (requirement 3.1.3; the M1 card editor
			     remains the friendly path) -->
			<template v-if="editor.payload.mode === 'flash'">
				<NCard size="small" :title="t('deploy.flashConfig')" class="form-card">
					<div class="form-grid">
						<label>{{ t('deploy.targetDevice') }}</label>
						<NSelect v-model:value="editor.payload.flash_device_id" :options="deviceOptions" style="width: 240px" clearable />
						<label>{{ t('deploy.flashTimeout') }}</label>
						<NInputNumber v-model:value="editor.payload.flash_timeout_sec" :min="0" :step="30" style="width: 180px">
							<template #suffix>s</template>
						</NInputNumber>
					</div>
					<div class="hint">{{ t('deploy.flashJSONHint') }}</div>
					<textarea v-model="flashStepsText" class="n-input mono" rows="10" spellcheck="false" />
				</NCard>
			</template>
		</div>

		<!-- copy-to-project modal (requirement 1.6) -->
		<NModal
			:show="copyTarget !== null"
			preset="dialog"
			:title="t('project.copyTo')"
			:show-icon="false"
			style="width: 420px"
			@positive-click="confirmCopy"
			@negative-click="copyTarget = null"
		>
			<div class="form-grid">
				<label>{{ t('project.copyTarget') }}</label>
				<NSelect
					v-model:value="copyTo"
					:options="copyProjects.filter((p) => p.id !== currentId).map((p) => ({ label: p.name, value: p.id }))"
					:placeholder="t('project.copyNoTarget')"
				/>
			</div>
			<template #action>
				<NButton @click="copyTarget = null">{{ t('common.cancel') }}</NButton>
				<NButton type="primary" :disabled="!copyTo" @click="confirmCopy">{{ t('common.confirm') }}</NButton>
			</template>
		</NModal>
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
	margin-top: 8px;
}
.hint {
	color: rgba(0, 0, 0, 0.45);
	font-size: 12.5px;
}
.muted {
	color: rgba(0, 0, 0, 0.45);
	font-size: 13px;
	margin-bottom: 8px;
}
.cmd-row,
.param-row {
	display: flex;
	align-items: center;
	gap: 8px;
	margin-bottom: 6px;
}
.cmd-row input {
	flex: 1;
}
</style>
