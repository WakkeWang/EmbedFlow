<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NSpace, NButton, NTag, NProgress, NModal, NSelect, NInput } from 'naive-ui'
import { useI18n } from 'vue-i18n'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import '@xterm/xterm/css/xterm.css'
import { WSClient } from '../api/ws'
import { Frame, type SessionStateFrame, type ExpectProgressFrame, type ConfirmFrame } from '../api/frames'
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
const progress = ref<{ cur: number; total: number; phase: string; desc?: string } | null>(null)
const showAbortConfirm = ref(false)
const rules = ref<{ id: number; name: string }[]>([])
const selectedRule = ref<number | null>(null)

// Human confirmation cards (issue #9): pending list with PASS/FAIL+note.
interface ConfirmCard {
	id: number
	prompt: string
	state: string
	result?: string
	note?: string
}
const confirmCards = ref<ConfirmCard[]>([])
const confirmPrompt = ref('')
const showConfirmDialog = ref(false)

let term: Terminal | null = null
let fit: FitAddon | null = null
let ws: WSClient | null = null
let everOpened = false

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
				} else if (st.state === 'conn_warning') {
					// Serial-client connection lost; task session survives the
					// grace window (CEO-16A). Show the countdown detail.
					term?.write(`\r\n[${t('terminal.serialLinkDown', { n: st.detail ?? '' })}]\r\n`)
				} else if (st.state === 'conn_down') {
					term?.write(`\r\n[${t('terminal.serialLinkDown', { n: st.detail ?? '' })}]\r\n`)
				} else if (st.state === 'conn_recovered') {
					term?.write(`\r\n[${t('terminal.serialLinkRecovered')}]\r\n`)
				} else if (st.state === 'closed' || st.state === 'failed') {
					term?.write(`\r\n[session ${st.state}: ${st.detail ?? ''}]\r\n`)
				}
			} else if (f.type === Frame.ExpectProgress) {
				const p = f.body as ExpectProgressFrame
				progress.value = { cur: p.step_index + 1, total: p.step_total, phase: p.phase, desc: p.step_desc }
			} else if (f.type === Frame.Confirm) {
				// Confirmation card state changes from the server.
				const cf = f.body as ConfirmFrame
				const i = confirmCards.value.findIndex((c) => c.id === cf.confirm_id)
				if (i >= 0) {
					confirmCards.value[i] = {
						id: cf.confirm_id,
						prompt: cf.prompt,
						state: cf.state,
						result: cf.result,
						note: cf.note,
					}
				} else {
					confirmCards.value.push({
						id: cf.confirm_id,
						prompt: cf.prompt,
						state: cf.state,
						result: cf.result,
						note: cf.note,
					})
				}
			}
		},
		onState: (s, a) => {
			const prev = connState.value
			connState.value = s
			attempt.value = a
			// CEO-17A: after a reconnect (closed -> connecting -> open), re-open
			// the session so a network blip drops us back into our own session
			// instead of a dead socket. The initial mount sends open via the
			// interval below; this covers every later reconnect.
			if (prev !== 'open' && s === 'open' && everOpened) {
				tryOpen()
			}
			if (s === 'open') {
				everOpened = true
			}
		},
		onAuthFail: (reason) => {
			term?.write(`\r\n[auth failed: ${reason}]\r\n`)
		},
	})
	ws.connect()

	// Open (or idempotently rejoin) our session on this device; a busy
	// device owned by someone else puts us in read-only follow (CEO-15A).
	const followMode = route.query.follow === '1'
	const command = followMode ? 'follow' : 'open'
	const tryOpen = () => {
		ws?.sendControl(Frame.SessionCtrl, { command, device_id: deviceID })
	}
	// Wait for the WS to be open before sending the open command.
	const iv = setInterval(() => {
		if (ws && ws.connState === 'open') {
			clearInterval(iv)
			tryOpen()
		}
	}, 200)

	// Restore confirmation cards after a page refresh (they persist in the
	// store; without this, pending cards vanish until someone re-resolves).
	// session id arrives async on the active frame; retry until known.
	const loadCards = setInterval(async () => {
		if (!sessionID.value) return
		clearInterval(loadCards)
		try {
			const list = await sessionApi.confirmations(sessionID.value)
			for (const c of list as { id: number; prompt: string; state: string; result?: string; note?: string }[]) {
				if (!confirmCards.value.find((x) => x.id === c.id)) {
					confirmCards.value.push(c)
				}
			}
		} catch {
			// Non-fatal: the server may not have any cards.
		}
	}, 300)

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
	if (!sessionID.value || !confirmPrompt.value.trim()) return
	try {
		const res = await sessionApi.insertConfirm(sessionID.value, confirmPrompt.value.trim())
		confirmCards.value.push({
			id: res.id,
			prompt: confirmPrompt.value.trim(),
			state: 'pending',
		})
		confirmPrompt.value = ''
		showConfirmDialog.value = false
	} catch (e) {
		term?.write(`\r\n[confirm insert failed: ${e}]\r\n`)
	}
}

// Resolve a card (issue #9: explicit PASS/FAIL + note, never auto-dismiss).
function resolveCard(card: ConfirmCard, result: 'pass' | 'fail') {
	ws?.sendControl(Frame.Confirm, {
		session_id: sessionID.value,
		confirm_id: card.id,
		prompt: card.prompt,
		state: 'resolved',
		result,
		note: '',
	})
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
				<template v-else-if="progress && progress.phase === 'aborted'">
					<NTag size="small" type="warning">
						{{ t('expect.aborted', { n: progress.cur }) }}
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
				<NButton size="small" @click="showConfirmDialog = true">+ {{ t('terminal.newConfirmCard') }}</NButton>
				<NButton size="small" @click="closeSession">{{ t('terminal.closeSession') }}</NButton>
			</NSpace>
		</div>

		<!-- Issue #9: pending confirmation cards beside the terminal. -->
		<div v-if="confirmCards.length > 0" class="confirm-bar">
			<div v-for="card in confirmCards" :key="card.id" class="confirm-card">
				<span class="confirm-prompt">{{ card.prompt }}</span>
				<NTag v-if="card.state === 'resolved'" size="small" :type="card.result === 'pass' ? 'success' : 'error'">
					{{ card.result?.toUpperCase() }}
				</NTag>
				<template v-else>
					<NButton size="tiny" type="success" @click="resolveCard(card, 'pass')">PASS</NButton>
					<NButton size="tiny" type="error" @click="resolveCard(card, 'fail')">FAIL</NButton>
				</template>
			</div>
		</div>

		<NModal :show="showConfirmDialog" preset="dialog" :title="t('terminal.newConfirmCard')" :show-icon="false">
			<NInput v-model:value="confirmPrompt" placeholder="e.g. LED on? / LED 是否亮" @keyup.enter="insertConfirm" />
			<template #action>
				<NButton @click="showConfirmDialog = false">{{ t('common.cancel') }}</NButton>
				<NButton type="primary" @click="insertConfirm">{{ t('common.confirm') }}</NButton>
			</template>
		</NModal>

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
.confirm-bar {
	display: flex;
	gap: 8px;
	flex-wrap: wrap;
	padding: 6px 0;
	flex: 0 0 auto;
}
.confirm-card {
	display: flex;
	align-items: center;
	gap: 8px;
	padding: 4px 10px;
	border: 1px solid rgba(255, 255, 255, 0.2);
	border-radius: 4px;
	font-size: 13px;
}
.confirm-prompt {
	max-width: 280px;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}
</style>
