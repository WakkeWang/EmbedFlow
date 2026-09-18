<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NSpace, NButton, NTag, NProgress, NModal, NSelect } from 'naive-ui'
import { useI18n } from 'vue-i18n'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import '@xterm/xterm/css/xterm.css'
import { WSClient } from '../api/ws'
import { Frame, type SessionStateFrame, type ExpectProgressFrame } from '../api/frames'
import { sessionApi, ruleApi } from '../api/http'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const deviceParam = route.params.sessionId as string // device-<id>
const deviceID = Number(deviceParam.replace('device-', ''))
const termEl = ref<HTMLDivElement | null>(null)
const connState = ref<'connecting' | 'authing' | 'open' | 'closed'>('connecting')
const attempt = ref(0)
const sessionID = ref<number>(0)
const readonly = ref(false)

// Expect run progress (DS-3A): current step / total; the run entry selects
// a rule (CEO-17A run button).
const progress = ref<{ cur: number; total: number; phase: string } | null>(null)
const showAbortConfirm = ref(false)
const rules = ref<{ id: number; name: string }[]>([])
const selectedRule = ref<number | null>(null)

let term: Terminal | null = null
let fit: FitAddon | null = null
let ws: WSClient | null = null

function mountTerminal() {
	if (!termEl.value) return
	term = new Terminal({
		fontFamily: "'JetBrains Mono', Consolas, monospace",
		fontSize: 14,
		theme: { background: '#101014', foreground: '#e0e0e0' },
		cursorBlink: true,
	})
	fit = new FitAddon()
	term.loadAddon(fit)
	term.open(termEl.value)
	fit.fit()

	// Keyboard input -> device (TX binary frame, decision 4A). Blocked while
	// an expect run is live (CEO-17A: no interleaved keystrokes).
	term.onData((data) => {
		if (readonly.value || (progress.value && progress.value.phase === 'running')) {
			term?.write(`\r\n[read-only or run in progress: input not sent]\r\n`)
			return
		}
		ws?.sendBinary(new TextEncoder().encode(data))
	})
}

onMounted(async () => {
	mountTerminal()

	ws = new WSClient({
		onBinary: (data) => term?.write(data),
		onControl: (f) => {
			if (f.type === Frame.SessionState) {
				const st = f.body as SessionStateFrame
				if (st.state === 'active' && st.owner_is_you) {
					sessionID.value = st.session_id
					readonly.value = false
				} else if (st.state === 'readonly') {
					sessionID.value = st.session_id
					readonly.value = true
				} else if (st.state === 'busy') {
					// Owner-is-you idempotent hit: adopt the live session.
					if (st.owner_is_you && st.session_id) {
						sessionID.value = st.session_id
					}
				} else if (st.state === 'closed' || st.state === 'failed') {
					term?.write(`\r\n[session ${st.state}: ${st.detail ?? ''}]\r\n`)
				}
			} else if (f.type === Frame.ExpectProgress) {
				const p = f.body as ExpectProgressFrame
				progress.value = { cur: p.step_index + 1, total: p.step_total, phase: p.phase }
			}
		},
		onState: (s, a) => {
			connState.value = s
			attempt.value = a
		},
		onAuthFail: (reason) => {
			term?.write(`\r\n[auth failed: ${reason}]\r\n`)
		},
	})
	ws.connect()

	// Open (or idempotently rejoin) our session on this device.
	const tryOpen = () => {
		ws?.sendControl(Frame.SessionCtrl, { command: 'open', device_id: deviceID })
	}
	// Wait for the WS to be open before sending the open command.
	const iv = setInterval(() => {
		if (ws && ws.connState === 'open') {
			clearInterval(iv)
			tryOpen()
		}
	}, 200)

	// Rule list for the run entry.
	try {
		const list = await ruleApi.list()
		rules.value = list.map((r) => ({ id: r.id, name: r.name }))
	} catch (e) {
		// Visible diagnostics for demo-mode visitors; rules stay empty on a
		// transient failure and a page reload retries.
		term?.write(`\r\n[rules load failed: ${e}]\r\n`)
		rules.value = []
	}
	;(window as unknown as Record<string, unknown>).__efRules = rules.value

	window.addEventListener('resize', onResize)
})

onBeforeUnmount(() => {
	window.removeEventListener('resize', onResize)
	ws?.close()
	term?.dispose()
})

function onResize() {
	fit?.fit()
}

const overlayText = computed(() => {
	switch (connState.value) {
		case 'connecting':
			return attempt.value > 0 ? t('terminal.reconnecting', { n: attempt.value }) : '...'
		case 'authing':
			return '...'
		case 'closed':
			return t('terminal.disconnected')
		default:
			return ''
	}
})

const showOverlay = computed(() => connState.value !== 'open')

function closeSession() {
	if (sessionID.value) {
		ws?.sendControl(Frame.SessionCtrl, { command: 'close', session_id: sessionID.value })
	}
	router.push('/devices')
}

function runRule() {
	if (!selectedRule.value) return
	// CEO-17A + device mutex (1A): a manual session occupies the device; the
	// run replaces it. Close ours first, then create the task session.
	const doRun = () =>
		ws?.sendControl(Frame.SessionCtrl, {
			command: 'run',
			device_id: deviceID,
			rule_id: selectedRule.value,
		})
	if (sessionID.value && !readonly.value) {
		ws?.sendControl(Frame.SessionCtrl, { command: 'close', session_id: sessionID.value })
		setTimeout(doRun, 150)
	} else {
		doRun()
	}
}

function abortRun() {
	if (sessionID.value) {
		ws?.sendControl(Frame.SessionCtrl, { command: 'abort', session_id: sessionID.value })
	}
	showAbortConfirm.value = false
}

// Insert a human confirmation card (issue #9) via the REST path.
async function insertConfirm() {
	if (!sessionID.value) return
	const prompt = window.prompt('Prompt (e.g. "LED on?")')
	if (!prompt) return
	try {
		await sessionApi.insertConfirm(sessionID.value, prompt)
		term?.write(`\r\n[confirmation card created: ${prompt}]\r\n`)
	} catch (e) {
		term?.write(`\r\n[confirm insert failed: ${e}]\r\n`)
	}
}
</script>

<template>
	<!-- DS-1B: terminal fills the viewport; top bar + status fixed height. -->
	<div class="terminal-page">
		<div class="top-bar">
			<NSpace align="center">
				<NTag :type="connState === 'open' ? 'success' : 'warning'" size="small">
					{{ connState === 'open' ? t('terminal.connected') : overlayText }}
				</NTag>
				<NTag v-if="readonly" size="small" type="warning">{{ t('terminal.readOnly') }}</NTag>
				<template v-if="progress && progress.phase === 'running'">
					<span class="mono">{{ t('expect.stepOf', { cur: progress.cur, total: progress.total }) }}</span>
					<NProgress
						type="line"
						:percentage="(progress.cur / Math.max(progress.total, 1)) * 100"
						:show-indicator="false"
						style="width: 180px"
						status="success"
					/>
					<NButton size="small" type="error" @click="showAbortConfirm = true">
						{{ t('expect.abort') }}
					</NButton>
				</template>
				<template v-else-if="progress && progress.phase === 'failed'">
					<NTag size="small" type="error">
						{{ t('expect.failedAt', { n: progress.cur }) }}
					</NTag>
				</template>
				<template v-else-if="progress && progress.phase === 'completed'">
					<NTag size="small" type="success">{{ t('expect.completed') }}</NTag>
				</template>
			</NSpace>
			<NSpace>
				<!-- CEO-17A: run entry on the terminal page. -->
				<NSelect
					v-model:value="selectedRule"
					:options="rules.map((r) => ({ label: r.name, value: r.id }))"
					:placeholder="t('expect.run')"
					size="small"
					style="width: 180px"
				/>
				<NButton size="small" type="primary" :disabled="!selectedRule" @click="runRule">
					{{ t('expect.run') }}
				</NButton>
				<NButton size="small" @click="insertConfirm">+ Confirm</NButton>
				<NButton size="small" @click="closeSession">{{ t('terminal.closeSession') }}</NButton>
			</NSpace>
		</div>

		<!-- DS-2B: translucent overlay keeps received data visible underneath. -->
		<div class="term-wrap">
			<div ref="termEl" class="term"></div>
			<div v-if="showOverlay" class="overlay">
				<span>{{ overlayText }}</span>
			</div>
		</div>

		<!-- DS-7B: abort asks a second confirmation, warns half-flash. -->
		<NModal
			:show="showAbortConfirm"
			preset="dialog"
			:title="t('expect.running')"
			:content="t('expect.abortConfirm')"
			:positive-text="t('common.confirm')"
			:negative-text="t('common.cancel')"
			@positive-click="abortRun"
			@negative-click="showAbortConfirm = false"
		/>
	</div>
</template>

<style scoped>
.terminal-page {
	display: flex;
	flex-direction: column;
	height: calc(100vh - 32px);
}
.top-bar {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 6px 0;
	flex: 0 0 auto;
}
.term-wrap {
	position: relative;
	flex: 1 1 auto;
	min-height: 320px;
	height: calc(100vh - 100px);
	overflow: hidden;
}
.term {
	position: absolute;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
}
.overlay {
	position: absolute;
	inset: 0;
	background: rgba(16, 16, 20, 0.55);
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: 15px;
	pointer-events: none;
}
.mono {
	font-family: 'JetBrains Mono', Consolas, monospace;
}
</style>
