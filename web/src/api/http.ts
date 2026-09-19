import { ofetch } from 'ofetch'

// HTTP client (decision 18A: ofetch). Token rides the Authorization header.
export const TOKEN_KEY = 'embedflow.token'

export function getToken(): string {
	return localStorage.getItem(TOKEN_KEY) ?? ''
}

export function setToken(t: string) {
	localStorage.setItem(TOKEN_KEY, t)
}

export function clearToken() {
	localStorage.removeItem(TOKEN_KEY)
}

export const api = ofetch.create({
	retry: false,
	onRequest({ options }) {
		const t = getToken()
		if (t) {
			options.headers = new Headers(options.headers)
			options.headers.set('Authorization', `Bearer ${t}`)
		}
	},
	onResponseError({ response }) {
		if (response.status === 401) {
			clearToken()
			window.location.href = '/login'
		}
	},
})

export interface ProjectRecord {
	id: number
	name: string
	note: string
	created_at: string
}

export interface Device {
	id: number
	name: string
	project: string
	busy: boolean
	owner?: string
	kind?: string
	session_id?: number
	shared: boolean
	since?: string
}

export interface SessionRecord {
	id: number
	device_id: number
	kind: string
	owner: string
	state: string
	started_at: string
	ended_at: string
	end_reason: string
	log_incomplete: boolean
}

export interface ExpectRuleRecord {
	id: number
	name: string
	steps_json: string
	updated_at: string
}

export const projectApi = {
	list: () => api<ProjectRecord[]>('/api/projects'),
	create: (name: string, note: string) =>
		api<{ id: number }>('/api/projects', { method: 'POST', body: { name, note } }),
	remove: (id: number) => api(`/api/projects/${id}`, { method: 'DELETE' }),
}

export const deviceApi = {
	list: () => api<Device[]>('/api/devices'),
	create: (name: string, project: string) =>
		api<{ id: number }>('/api/devices', { method: 'POST', body: { name, project } }),
	remove: (id: number) => api(`/api/devices/${id}`, { method: 'DELETE' }),
}

export const sessionApi = {
	history: (deviceId: number) => api<SessionRecord[]>(`/api/devices/${deviceId}/sessions`),
	close: (sessionId: number) => api(`/api/sessions/${sessionId}/close`, { method: 'POST' }),
	logDownloadURL: (sessionId: number) =>
		`/api/sessions/${sessionId}/log/download?token=${encodeURIComponent(getToken())}`,
	logTail: (sessionId: number, n = 1000) =>
		api<{ tail: string }>(`/api/sessions/${sessionId}/log/tail?n=${n}`),
	insertConfirm: (sessionId: number, prompt: string) =>
		api<{ id: number }>(`/api/sessions/${sessionId}/confirm`, { method: 'POST', body: { prompt } }),
	confirmations: (sessionId: number) =>
		api<unknown[]>(`/api/sessions/${sessionId}/confirmations`),
}

export const ruleApi = {
	list: () => api<ExpectRuleRecord[]>('/api/expect-rules'),
	create: (name: string, steps: unknown) =>
		api<{ id: number }>('/api/expect-rules', { method: 'POST', body: { name, steps } }),
	get: (id: number) => api<ExpectRuleRecord>(`/api/expect-rules/${id}`),
	update: (id: number, name: string, steps: unknown) =>
		api(`/api/expect-rules/${id}`, { method: 'PUT', body: { name, steps } }),
	remove: (id: number) => api(`/api/expect-rules/${id}`, { method: 'DELETE' }),
}
