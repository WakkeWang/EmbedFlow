import { ofetch } from 'ofetch'

// HTTP client (decision 18A: ofetch). Base URL empty = same origin; token
// rides the Authorization header until the real login slice lands.
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

export interface Device {
	id: number
	name: string
	project: string
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

// M1 interim: device list and history endpoints will land server-side with
// the session open flow (T5 server part). These helpers target them.
export const deviceApi = {
	list: () => api<Device[]>('/api/devices'),
}

export const sessionApi = {
	history: (deviceId: number) => api<SessionRecord[]>(`/api/devices/${deviceId}/sessions`),
	close: (sessionId: number) =>
		api(`/api/sessions/${sessionId}/close`, { method: 'POST' }),
	logDownloadURL: (sessionId: number) =>
		`/api/sessions/${sessionId}/log/download?token=${encodeURIComponent(getToken())}`,
}
