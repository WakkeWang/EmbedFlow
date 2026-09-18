<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
	NSpace,
	NButton,
	NTag,
	NProgress,
	NModal,
} from 'naive-ui'
import { useI18n } from 'vue-i18n'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import '@xterm/xterm/css/xterm.css'
import { WSClient } from '../api/ws'
import { Frame, type SessionStateFrame } from '../api/frames'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const sessionId = route.params.sessionId as string
const termEl = ref<HTMLDivElement | null>(null)
const connState = ref<'connecting' | 'authing' | 'open' | 'closed'>('connecting')
const attempt = ref(0)

// Expect run progress (DS-3A): current step / total; red freeze on failure
// (DS-7B abort with confirm).
const progress = ref<{ cur: number; total: number } | null>(null)
const showAbortConfirm = ref(false)

let term: Terminal | null = null
let fit: FitAddon | null = null
let ws: WSClient | null = null

function mountTerminal() {
	if (!termEl.value) return
	term = new Terminal({
		fontFamily: "'JetBrains Mono', Consolas, monospace",
		fontSize: 14,
		theme: {
			background: '#101014',
			foreground: '#e0e0e0',
		},
		cursorBlink: true,
	})
	fit = new FitAddon()
	term.loadAddon(fit)
	term.open(termEl.value)
	fit.fit()

	// Keyboard input -> device (TX binary frame, decision 4A).
	term.onData((data) => {
		ws?.sendBinary(new TextEncoder().encode(data))
	})
}

onMounted(() => {
	mountTerminal()

	ws = new WSClient({
		onBinary: (data) => {
			// Device -> browser (RX).
			term?.write(data)
		},
		onControl: (f) => {
			if (f.type === Frame.SessionState) {
				const st = f.body as SessionStateFrame
				if (st.state === 'rejected') return
			}
		},
		onState: (s, a) => {
			connState.value = s
			attempt.value = a
		},
	})
	ws.connect()

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
	ws?.sendControl(Frame.SessionCtrl, {
		command: 'close',
		session_id: Number(sessionId),
	})
	router.push('/devices')
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
				<template v-if="progress">
					<span class="mono">{{ t('expect.stepOf', { cur: progress.cur, total: progress.total }) }}</span>
					<NProgress
						type="line"
						:percentage="(progress.cur / Math.max(progress.total, 1)) * 100"
						:show-indicator="false"
						style="width: 200px"
					/>
					<NButton size="small" type="error" @click="showAbortConfirm = true">
						{{ t('expect.abort') }}
					</NButton>
				</template>
			</NSpace>
			<NButton size="small" @click="closeSession">{{ t('terminal.closeSession') }}</NButton>
		</div>

		<!-- DS-2B: translucent overlay keeps received data visible underneath. -->
		<div class="term-wrap">
			<div ref="termEl" class="term"></div>
			<div v-if="showOverlay" class="overlay">
				<span>{{ overlayText }}</span>
			</div>
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
	min-height: 0;
}
.term {
	position: absolute;
	inset: 0;
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
