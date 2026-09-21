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
	// The role rides the token lifecycle (logout / 401 both clear it).
	localStorage.removeItem('embedflow.role')
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
	ssh_host?: string
	ssh_port?: number
	ssh_user?: string
	ssh_set?: boolean
	note?: string
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
	// Cross-project copy + JSON export/import (requirement 1.6).
	copyConfigObject: (id: number, targetProjectId: number) =>
		api<{ id: number }>(`/api/config-objects/${id}/copy`, { method: 'POST', body: { target_project_id: targetProjectId } }),
	exportURL: (projectId: number, kinds?: string) =>
		`/api/projects/${projectId}/export${kinds ? `?kinds=${kinds}` : ''}`,
	importBundle: (projectId: number, bundle: Record<string, any>) =>
		api<{ imported: number; skipped: number }>(`/api/projects/${projectId}/import`, { method: 'POST', body: bundle }),
}

export const deviceApi = {
	list: () => api<Device[]>('/api/devices'),
	create: (name: string, project: string) =>
		api<{ id: number }>('/api/devices', { method: 'POST', body: { name, project } }),
	remove: (id: number) => api(`/api/devices/${id}`, { method: 'DELETE' }),
	// SSH-aware variants (requirement 3.2). Password is write-only.
	createV2: (d: {
		name: string
		project: string
		ssh_host?: string
		ssh_port?: number
		ssh_user?: string
		ssh_password?: string
		note?: string
	}) => api<{ id: number }>('/api/devices/v2', { method: 'POST', body: d }),
	update: (
		id: number,
		d: {
			name: string
			project: string
			ssh_host?: string
			ssh_port?: number
			ssh_user?: string
			ssh_password?: string
			note?: string
		},
	) => api(`/api/devices/${id}`, { method: 'PUT', body: d }),
	sshTest: (id: number) => api<{ ok: boolean; detail: string }>(`/api/devices/${id}/ssh-test`, { method: 'POST' }),
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

// Build module (M2): build items, batches, records, artifacts, settings.

export interface BuildItemRecord {
	id: number
	project_id: number
	name: string
	source_type: 'git' | 'local'
	git_url?: string
	git_branch?: string
	git_commit?: string
	check_latest?: boolean
	local_path?: string
	command: string
	artifacts: string[]
	timeout_sec?: number
	version_cmd?: string
	prereq_json: string
}

export interface BatchRecord {
	id: number
	project_id: number
	status: 'queued' | 'running' | 'completed' | 'failed' | 'canceled'
	items_json: string
	created_by: string
	created_at: string
	started_at?: string
	ended_at?: string
}

export interface BuildRecordRecord {
	id: number
	batch_id: number
	item_id: number
	project_id: number
	status: 'pending' | 'building' | 'succeeded' | 'failed' | 'canceled' | 'skipped'
	commit_sha: string
	branch?: string
	version_info: string
	exit_code?: number | null
	executor: string
	started_at: string
	ended_at: string
}

export interface ArtifactRecord {
	id: number
	build_record_id: number
	name: string
	size: number
	checksum: string
}

export interface UserRecord {
	id: number
	username: string
	role: 'admin' | 'member'
}

export const buildItemApi = {
	list: (projectId: number) => api<BuildItemRecord[]>(`/api/projects/${projectId}/build-items`),
	listAll: () => api<BuildItemRecord[]>('/api/build-items'),
	create: (projectId: number, item: Partial<BuildItemRecord>) =>
		api<{ id: number }>(`/api/projects/${projectId}/build-items`, { method: 'POST', body: item }),
	update: (id: number, item: Partial<BuildItemRecord>) =>
		api(`/api/build-items/${id}`, { method: 'PUT', body: item }),
	remove: (id: number) => api(`/api/build-items/${id}`, { method: 'DELETE' }),
}

export const batchApi = {
	create: (projectId: number, itemIds: number[]) =>
		api<{ id: number }>('/api/batches', { method: 'POST', body: { project_id: projectId, item_ids: itemIds } }),
	list: (projectId: number) => api<BatchRecord[]>(`/api/batches?project_id=${projectId}`),
	get: (id: number) => api<{ batch: BatchRecord; records: BuildRecordRecord[] }>(`/api/batches/${id}`),
	cancel: (id: number) => api(`/api/batches/${id}/cancel`, { method: 'POST' }),
}

export const buildRecordApi = {
	list: (projectId: number) => api<BuildRecordRecord[]>(`/api/build-records?project_id=${projectId}`),
	get: (id: number) => api<BuildRecordRecord>(`/api/build-records/${id}`),
	artifacts: (id: number) => api<ArtifactRecord[]>(`/api/build-records/${id}/artifacts`),
	logTail: (id: number, n = 2000) => api<{ tail: string }>(`/api/build-records/${id}/log/tail?n=${n}`),
	logDownloadURL: (id: number) => `/api/build-records/${id}/log/download?token=${encodeURIComponent(getToken())}`,
	remove: (id: number, mode: 'record' | 'artifacts') =>
		api(`/api/build-records/${id}?mode=${mode}`, { method: 'DELETE' }),
	removeAll: (projectId: number, mode: 'record' | 'artifacts') =>
		api('/api/build-records/delete-all', { method: 'POST', body: { project_id: projectId, mode } }),
	artifactDownloadURL: (id: number) => `/api/artifacts/${id}/download?token=${encodeURIComponent(getToken())}`,
}

export const settingsApi = {
	get: () => api<Record<string, string>>('/api/settings'),
	put: (values: Record<string, string>) => api('/api/settings', { method: 'PUT', body: values }),
}

// Deploy module (M3): deploy rules, deploy records.

export interface DeployParamDef {
	name: string
	value?: string
}

export interface USBLayoutDef {
	glob: string
	dir?: string
}

export interface USBChecksumDef {
	name: string
	glob: string
}

export interface USBDiskConfig {
	root_name: string
	layout?: USBLayoutDef[]
	checksums?: USBChecksumDef[]
}

export interface DeployPayload {
	mode: 'manual' | 'ssh' | 'flash'
	steps_md?: string
	ssh_device_id?: number
	ssh_commands?: string[]
	ssh_params?: DeployParamDef[]
	ssh_timeout_sec?: number
	flash_device_id?: number
	flash_steps_json?: string
	flash_timeout_sec?: number
	usb_disk?: USBDiskConfig | null
}

export interface DeployRuleRecord {
	id: number
	project_id: number
	name: string
	mode: 'manual' | 'ssh' | 'flash'
	created_at: string
	updated_at: string
	payload: DeployPayload
}

export interface DeployRecordRecord {
	id: number
	project_id: number
	rule_id: number
	device_id: number
	build_record_id: number
	executor: string
	status: 'running' | 'succeeded' | 'failed' | 'canceled'
	detail: string
	exec_record_id?: number
	session_id?: number
	started_at: string
	ended_at: string
}

export const deployRuleApi = {
	list: (projectId: number) => api<DeployRuleRecord[]>(`/api/projects/${projectId}/deploy-rules`),
	get: (id: number) => api<{ rule: DeployRuleRecord; payload: DeployPayload }>(`/api/deploy-rules/${id}`),
	create: (projectId: number, body: Partial<DeployRuleRecord> & { payload?: DeployPayload }) =>
		api<{ id: number }>(`/api/projects/${projectId}/deploy-rules`, { method: 'POST', body }),
	update: (id: number, body: Partial<DeployRuleRecord> & { payload?: DeployPayload }) =>
		api(`/api/deploy-rules/${id}`, { method: 'PUT', body }),
	remove: (id: number) => api(`/api/deploy-rules/${id}`, { method: 'DELETE' }),
}

export const deployApi = {
	trigger: (ruleId: number, buildRecordId: number, params: Record<string, string>) =>
		api<{ id: number }>('/api/deployments', { method: 'POST', body: { rule_id: ruleId, build_record_id: buildRecordId, params } }),
	list: (projectId: number) => api<DeployRecordRecord[]>(`/api/projects/${projectId}/deployments`),
	get: (id: number) => api<DeployRecordRecord>(`/api/deployments/${id}`),
	cancel: (id: number) => api(`/api/deployments/${id}/cancel`, { method: 'POST' }),
	remove: (id: number) => api(`/api/deployments/${id}`, { method: 'DELETE' }),
	logTail: (id: number, n = 1000) => api<{ tail: string }>(`/api/deployments/${id}/log/tail?n=${n}`),
}

export const userApi = {
	list: () => api<UserRecord[]>('/api/users'),
	create: (username: string, password: string, role: string) =>
		api('/api/users', { method: 'POST', body: { username, password, role } }),
	changeOwnPassword: (oldpw: string, newpw: string) =>
		api('/api/users/self/password', { method: 'POST', body: { old: oldpw, new: newpw } }),
	resetPassword: (id: number, password: string) =>
		api(`/api/users/${id}/password`, { method: 'POST', body: { password } }),
}
