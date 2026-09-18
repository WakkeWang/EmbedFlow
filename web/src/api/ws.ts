import {
	Frame,
	ProtocolVersion,
	parseControl,
	encode,
	type FrameType,
	type ParsedFrame,
} from './frames'
import { getToken } from './http'

export type ConnState = 'connecting' | 'authing' | 'open' | 'closed'

interface WSOpts {
	// onBinary receives raw serial bytes (data frames).
	onBinary?: (data: Uint8Array) => void
	// onControl receives every known control frame after auth.
	onControl?: (f: ParsedFrame) => void
	// onState reports connection transitions for the disconnect overlay (DS-2B).
	onState?: (s: ConnState, attempt: number) => void
}

// wsURL derives the WebSocket endpoint from the page location.
function wsURL(): string {
	const proto = location.protocol === 'https:' ? 'wss' : 'ws'
	return `${proto}://${location.host}/ws/client`
}

// WSClient is the browser side of the serial tunnel: dual-frame protocol,
// heartbeat, reconnect with exponential backoff (decision 18A: native WS,
// self-rolled; matches server protocol 4A).
export class WSClient {
	private ws: WebSocket | null = null
	private state: ConnState = 'closed'
	private attempt = 0
	private closedByUser = false
	private heartbeatTimer: ReturnType<typeof setInterval> | null = null
	private seq = 0
	private opts: WSOpts

	constructor(opts: WSOpts) {
		this.opts = opts
	}

	connect() {
		this.closedByUser = false
		this.setState('connecting')
		const ws = new WebSocket(wsURL())
		this.ws = ws
		ws.binaryType = 'arraybuffer'

		ws.onopen = () => {
			this.setState('authing')
			// First frame: auth with protocol version (CEO-11A).
			this.sendControl(Frame.Auth, {
				token: getToken(),
				protocol_version: ProtocolVersion,
				client_version: 'web-dev',
			})
			this.startHeartbeat()
		}

		ws.onmessage = (ev) => {
			if (typeof ev.data === 'string') {
				this.onText(ev.data)
			} else {
				// Binary data frame: raw serial bytes (decision 4A).
				this.opts.onBinary?.(new Uint8Array(ev.data as ArrayBuffer))
			}
		}

		ws.onclose = () => {
			this.stopHeartbeat()
			if (this.closedByUser) {
				this.setState('closed')
				return
			}
			this.scheduleReconnect()
		}

		ws.onerror = () => {
			// onclose follows; nothing to do here.
		}
	}

	close() {
		this.closedByUser = true
		this.stopHeartbeat()
		this.ws?.close()
		this.ws = null
		this.setState('closed')
	}

	sendControl(type: number, body: unknown) {
		if (this.ws?.readyState === WebSocket.OPEN) {
			this.ws.send(encode(type as FrameType, body))
		}
	}

	sendBinary(data: Uint8Array) {
		if (this.ws?.readyState === WebSocket.OPEN) {
			this.ws.send(data.buffer.slice(data.byteOffset, data.byteOffset + data.byteLength) as ArrayBuffer)
		}
	}

	private onText(wire: string) {
		let f: ParsedFrame
		try {
			f = parseControl(wire)
		} catch {
			return // malformed control frame: ignore (server tolerates too)
		}
		switch (f.type) {
			case Frame.AuthOK:
				this.attempt = 0
				this.setState('open')
				break
			case Frame.AuthFail:
				// Version mismatch or bad token: do not retry the same way.
				this.closedByUser = true
				this.stopHeartbeat()
				this.setState('closed')
				this.opts.onControl?.(f)
				break
			default:
				this.opts.onControl?.(f)
		}
	}

	private startHeartbeat() {
		this.stopHeartbeat()
		this.heartbeatTimer = setInterval(() => {
			this.seq++
			this.sendControl(Frame.Heartbeat, { seq: this.seq })
		}, 2000)
	}

	private stopHeartbeat() {
		if (this.heartbeatTimer) {
			clearInterval(this.heartbeatTimer)
			this.heartbeatTimer = null
		}
	}

	private scheduleReconnect() {
		this.attempt++
		this.setState('connecting')
		// Exponential backoff, 1s..30s cap (DS-2B: the overlay shows attempt N).
		const delay = Math.min(1000 * 2 ** Math.min(this.attempt - 1, 5), 30000)
		setTimeout(() => {
			if (!this.closedByUser) this.connect()
		}, delay)
	}

	private setState(s: ConnState) {
		this.state = s
		this.opts.onState?.(s, this.attempt)
	}

	get connState(): ConnState {
		return this.state
	}
}
