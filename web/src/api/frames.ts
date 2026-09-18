// Client-side mirror of pkg/protocol (decision 4A): control frames are JSON
// text frames, data frames binary. Frame type numbers must match the server.
export const ProtocolVersion = 1

export const Frame = {
	Auth: 1,
	AuthOK: 2,
	AuthFail: 3,
	Heartbeat: 4,
	ShareRequest: 5,
	DeviceStatus: 6,
	SessionCtrl: 7,
	SessionState: 8,
	ExpectProgress: 9,
	Confirm: 10,
} as const

export type FrameType = (typeof Frame)[keyof typeof Frame]

export interface AuthFrame {
	token: string
	protocol_version: number
	client_version?: string
}
export interface AuthOKFrame {
	server_version?: string
}
export interface AuthFailFrame {
	reason: string
}
export interface HeartbeatFrame {
	seq: number
}
export interface SessionCtrlFrame {
	command: string
	session_id?: number
	device_id?: number
}
export interface SessionStateFrame {
	session_id: number
	device_id: number
	state: string
	detail?: string
	occupier_user?: string
	occupier_kind?: string
	occupier_since?: string
	owner_is_you?: boolean
}

export interface ExpectProgressFrame {
	session_id: number
	step_index: number
	step_total: number
	step_desc?: string
	phase: string
	detail?: string
}

export interface ConfirmFrame {
	session_id: number
	confirm_id: number
	prompt: string
	state: string
	result?: string
	note?: string
}

export type ControlBody =
	| AuthFrame
	| AuthOKFrame
	| AuthFailFrame
	| HeartbeatFrame
	| SessionCtrlFrame
	| SessionStateFrame

export interface Envelope {
	type: number
	body?: unknown
}

export interface ParsedFrame {
	type: number
	body: any
}

// ParseControl decodes a text frame; unknown types come back with body as
// raw object so the caller can ignore them (CEO-3A).
export function parseControl(wire: string): ParsedFrame {
	const env = JSON.parse(wire) as Envelope
	return { type: env.type, body: env.body }
}

// Encode serializes a body under its frame type number.
export function encode(type: FrameType, body: unknown): string {
	return JSON.stringify({ type, body })
}
