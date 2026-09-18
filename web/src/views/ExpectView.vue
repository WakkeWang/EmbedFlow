<script setup lang="ts">
import { ref, computed } from 'vue'
import {
	NCard,
	NButton,
	NInput,
	NInputNumber,
	NSelect,
	NTag,
	NEmpty,
	NModal,
	useMessage,
} from 'naive-ui'
import { useI18n } from 'vue-i18n'
import {
	emptyStep,
	validateSteps,
	type ExpectRule,
} from '../api/expect'

const { t } = useI18n()
const message = useMessage()

// M1 interim: rules persist via /api/expect when the server part lands;
// until then they live in localStorage so the editor is fully usable.
const STORAGE_KEY = 'embedflow.expect.rules'

function loadRules(): ExpectRule[] {
	try {
		return JSON.parse(localStorage.getItem(STORAGE_KEY) ?? '[]')
	} catch {
		return []
	}
}

const rules = ref<ExpectRule[]>(loadRules())
const current = ref<ExpectRule | null>(null)
const errors = ref<Map<number, string>>(new Map())

const stepCount = computed(() => current.value?.steps.length ?? 0)

function jumpTo(n: number) {
	document.getElementById(`step-${n}`)?.scrollIntoView({ behavior: 'smooth', block: 'center' })
}

function newRule() {
	current.value = { name: '', steps: [emptyStep()] }
}

function persist() {
	localStorage.setItem(STORAGE_KEY, JSON.stringify(rules.value))
}

function save() {
	if (!current.value) return
	errors.value = validateSteps(current.value.steps)
	if (errors.value.size > 0) {
		// DS-7A: jump to the first offending step.
		const first = Math.min(...errors.value.keys())
		jumpTo(first)
		return
	}
	if (!rules.value.includes(current.value)) {
		rules.value.push(current.value)
	}
	persist()
	message.success(t('expect.save'))
}

function editRule(r: ExpectRule) {
	current.value = r
}

function addStep() {
	current.value?.steps.push(emptyStep())
}

function insertAfter(i: number) {
	current.value?.steps.splice(i + 1, 0, emptyStep())
}

function copyStep(i: number) {
	if (!current.value) return
	current.value.steps.splice(i + 1, 0, { ...current.value.steps[i] })
}

function moveUp(i: number) {
	if (!current.value || i === 0) return
	const [s] = current.value.steps.splice(i, 1)
	current.value.steps.splice(i - 1, 0, s)
}

function moveDown(i: number) {
	if (!current.value || !current.value.steps || i >= current.value.steps.length - 1) return
	const [s] = current.value.steps.splice(i, 1)
	current.value.steps.splice(i + 1, 0, s)
}

function deleteStep(i: number) {
	current.value?.steps.splice(i, 1)
}

const onFailOptions = computed(() => [
	{ label: t('expect.abort'), value: 'abort' },
	{ label: t('expect.continue'), value: 'continue' },
	{ label: t('expect.retry'), value: 'retry' },
])

// Control-character insertion (issue #7): the operator never hand-types
// escapes; buttons append to the focused send input via the model.
function appendSend(i: number, esc: string) {
	if (!current.value) return
	current.value.steps[i].send += esc
}

const showAbortConfirm = ref(false)
</script>

<template>
	<div>
		<div class="header-row">
			<h2>{{ t('expect.list') }}</h2>
			<NButton type="primary" @click="newRule">{{ t('expect.new') }}</NButton>
		</div>

		<NEmpty
			v-if="rules.length === 0 && !current"
			:description="t('expect.emptyHint')"
		/>

		<div v-if="rules.length > 0 && !current" class="rule-list">
			<NCard v-for="r in rules" :key="r.name" size="small" :title="r.name">
				<template #header-extra>
					<NTag size="small">{{ r.steps.length }} {{ t('expect.steps') }}</NTag>
				</template>
				<NButton size="small" @click="editRule(r)">{{ t('expect.save') }}/{{ t('common.close') }}</NButton>
			</NCard>
		</div>

		<div v-if="current" class="editor">
			<div class="header-row">
				<NInput v-model:value="current.name" :placeholder="'rule name'" style="width: 240px" />
				<span class="step-total">{{ stepCount }} {{ t('expect.steps') }}</span>
				<!-- DS-7A: step count + jump-to-step at page top, no grouping. -->
				<NSelect
					v-if="stepCount > 3"
					:options="Array.from({ length: stepCount }, (_, i) => ({ label: t('expect.stepNo', { n: i + 1 }), value: i }))"
					:placeholder="t('expect.stepNo', { n: 1 })"
					style="width: 140px"
					@update:value="(v: number) => jumpTo(v)"
				/>
				<NButton type="primary" @click="save">{{ t('expect.save') }}</NButton>
				<NButton @click="addStep">{{ t('expect.addStep') }}</NButton>
			</div>

			<!-- DS-1C: vertical card list; per-step card with inline elements. -->
			<NCard
				v-for="(s, i) in current.steps"
				:id="`step-${i}`"
				:key="i"
				size="small"
				class="step-card"
				:class="{ 'step-error': errors.has(i) }"
			>
				<template #header>
					<span class="mono">{{ t('expect.stepNo', { n: i + 1 }) }}</span>
					<NTag v-if="errors.has(i)" type="error" size="small" style="margin-left: 8px">
						{{ errors.get(i) }}
					</NTag>
				</template>
				<template #header-extra>
					<NButton quaternary size="tiny" @click="insertAfter(i)">{{ t('expect.insert') }}</NButton>
					<NButton quaternary size="tiny" @click="copyStep(i)">{{ t('expect.copy') }}</NButton>
					<NButton quaternary size="tiny" @click="moveUp(i)">{{ t('expect.up') }}</NButton>
					<NButton quaternary size="tiny" @click="moveDown(i)">{{ t('expect.down') }}</NButton>
					<NButton quaternary size="tiny" type="error" @click="deleteStep(i)">
						{{ t('expect.delete') }}
					</NButton>
				</template>

				<div class="step-grid">
					<label>{{ t('expect.awaits') }}</label>
					<input v-model="s.await" class="mono n-input" spellcheck="false" />

					<label>{{ t('expect.send') }}</label>
					<div>
						<input v-model="s.send" class="mono n-input" spellcheck="false" />
						<div class="ctrl-row">
							<NButton size="tiny" @click="appendSend(i, '\\r')">{{ t('expect.crlf') }}</NButton>
							<NButton size="tiny" @click="appendSend(i, '\\C')">{{ t('expect.ctrlC') }}</NButton>
							<NButton size="tiny" @click="appendSend(i, '\\x')">{{ t('expect.hex') }}</NButton>
						</div>
					</div>

					<label>{{ t('expect.timeout') }}</label>
					<NInputNumber v-model:value="s.timeout_ms" :min="0" :step="1000" style="width: 160px" />

					<label>{{ t('expect.delay') }}</label>
					<NInputNumber v-model:value="s.delay_ms" :min="0" :step="100" style="width: 160px" />

					<label>{{ t('expect.onFail') }}</label>
					<NSelect v-model:value="s.on_fail" :options="onFailOptions" style="width: 160px" clearable />
				</div>
			</NCard>
		</div>

		<NModal
			:show="showAbortConfirm"
			preset="dialog"
			:title="t('expect.running')"
			:content="t('expect.abortConfirm')"
			:positive-text="t('common.confirm')"
			:negative-text="t('common.cancel')"
			@positive-click="showAbortConfirm = false"
			@negative-click="showAbortConfirm = false"
		/>
	</div>
</template>

<style scoped>
.header-row {
	display: flex;
	align-items: center;
	gap: 12px;
	margin-bottom: 12px;
}
.rule-list {
	display: flex;
	flex-direction: column;
	gap: 8px;
	max-width: 640px;
}
.step-card {
	margin-bottom: 10px;
}
.step-card.step-error {
	border: 1px solid #e88080;
}
.step-grid {
	display: grid;
	grid-template-columns: 110px 1fr;
	gap: 8px 12px;
	align-items: center;
}
.step-grid label {
	color: rgba(255, 255, 255, 0.65);
	font-size: 12px;
}
.mono {
	font-family: 'JetBrains Mono', Consolas, monospace;
}
input.n-input,
:deep(input.n-input) {
	width: 100%;
	padding: 5px 9px;
	border-radius: 3px;
	border: 1px solid rgba(255, 255, 255, 0.24);
	background: transparent;
	color: inherit;
}
.ctrl-row {
	display: flex;
	gap: 6px;
	margin-top: 4px;
}
.step-total {
	color: rgba(255, 255, 255, 0.65);
	font-size: 13px;
}
</style>
