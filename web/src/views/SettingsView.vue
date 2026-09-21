<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import {
	NCard,
	NButton,
	NInput,
	NSelect,
	NTag,
	NDataTable,
	NSpace,
	NModal,
	type DataTableColumns,
	useMessage,
} from 'naive-ui'
import { h } from 'vue'
import { useI18n } from 'vue-i18n'
import { settingsApi, userApi, type UserRecord } from '../api/http'
import { useAuth } from '../store/auth'
import { useRouter } from 'vue-router'

// System settings (requirement 1.6: four admin-managed items) + account
// management (requirement 1.5: admins create accounts, everyone changes
// their own password).
const { t } = useI18n()
const message = useMessage()
const router = useRouter()
const { isAdmin } = useAuth()

const settings = ref<Record<string, string>>({
	max_parallel_batches: '4',
	checksum: 'sha256',
	publish_whitelist: '',
	build_tmp_dir: '',
})
const saving = ref(false)

const users = ref<UserRecord[]>([])
const newUserOpen = ref(false)
const newUsername = ref('')
const newPassword = ref('')
const newRole = ref('member')
const resetTarget = ref<UserRecord | null>(null)
const resetPassword = ref('')
const selfPwOpen = ref(false)
const selfOld = ref('')
const selfNew = ref('')

onMounted(async () => {
	if (!isAdmin.value) {
		router.replace('/')
		return
	}
	try {
		settings.value = { ...settings.value, ...(await settingsApi.get()) }
	} catch {
		// defaults stand
	}
	await loadUsers()
})

async function loadUsers() {
	try {
		users.value = await userApi.list()
	} catch {
		users.value = []
	}
}

async function save() {
	saving.value = true
	try {
		await settingsApi.put({
			max_parallel_batches: settings.value.max_parallel_batches ?? '4',
			checksum: settings.value.checksum ?? 'sha256',
			publish_whitelist: settings.value.publish_whitelist ?? '',
			build_tmp_dir: settings.value.build_tmp_dir ?? '',
		})
		message.success(t('common.confirm'))
	} catch (e) {
		message.error(String(e))
	} finally {
		saving.value = false
	}
}

async function createUser() {
	if (!newUsername.value.trim() || !newPassword.value) return
	try {
		await userApi.create(newUsername.value.trim(), newPassword.value, newRole.value)
		newUserOpen.value = false
		newUsername.value = ''
		newPassword.value = ''
		await loadUsers()
		message.success(t('common.confirm'))
	} catch (e) {
		message.error(String(e))
	}
}

async function doReset() {
	if (!resetTarget.value || !resetPassword.value) return
	try {
		await userApi.resetPassword(resetTarget.value.id, resetPassword.value)
		resetTarget.value = null
		resetPassword.value = ''
		message.success(t('common.confirm'))
	} catch (e) {
		message.error(String(e))
	}
}

// The modal's X / mask click only fires update:show(false) -- closing on
// true would wipe the target before the dialog shows.
function onResetShow(v: boolean) {
	if (!v) {
		resetTarget.value = null
	}
}

async function changeOwn() {
	if (!selfNew.value) return
	try {
		await userApi.changeOwnPassword(selfOld.value, selfNew.value)
		selfPwOpen.value = false
		selfOld.value = ''
		selfNew.value = ''
		message.success(t('common.confirm'))
	} catch (e) {
		message.error(String(e))
	}
}

const checksumOptions = [
	{ label: 'sha256', value: 'sha256' },
	{ label: 'md5', value: 'md5' },
]

const userColumns = computed<DataTableColumns<UserRecord>>(() => [
	{ title: 'ID', key: 'id', width: 60 },
	{ title: t('login.username'), key: 'username' },
	{
		title: t('build.role'),
		key: 'role',
		width: 100,
		render: (u) => h(NTag, { size: 'small', type: u.role === 'admin' ? 'warning' : 'default' }, { default: () => u.role }),
	},
	{
		title: '',
		key: 'actions',
		width: 140,
		render: (u) =>
			h(NSpace, { size: 4 }, {
				default: () => [
					h(
						NButton,
						{ size: 'tiny', quaternary: true, onClick: () => ((resetTarget.value = u), (resetPassword.value = '')) },
						{ default: () => t('build.resetPassword') },
					),
				],
			}),
	},
])
</script>

<template>
	<div>
		<div class="header-row">
			<h2>{{ t('settings.title') }}</h2>
		</div>

		<NCard size="small" :title="t('settings.system')" class="settings-card">
			<div class="form-grid">
				<label>{{ t('settings.maxParallel') }}</label>
				<NInput v-model:value="settings.max_parallel_batches" style="width: 120px" />
				<div></div>
				<div class="hint">{{ t('settings.maxParallelHint') }}</div>

				<label>{{ t('settings.checksum') }}</label>
				<NSelect v-model:value="settings.checksum" :options="checksumOptions" style="width: 160px" />
				<div></div>
				<div class="hint">{{ t('settings.checksumHint') }}</div>

				<label>{{ t('settings.whitelist') }}</label>
				<textarea v-model="settings.publish_whitelist" class="n-input" rows="3" spellcheck="false" />
				<div></div>
				<div class="hint">{{ t('settings.whitelistHint') }}</div>

				<label>{{ t('settings.tmpDir') }}</label>
				<input v-model="settings.build_tmp_dir" class="n-input mono" spellcheck="false" />
				<div></div>
				<div class="hint">{{ t('settings.tmpDirHint') }}</div>
			</div>
			<div style="margin-top: 12px">
				<NButton type="primary" :loading="saving" @click="save">{{ t('expect.save') }}</NButton>
			</div>
		</NCard>

		<NCard size="small" :title="t('settings.users')" class="settings-card">
			<template #header-extra>
				<NSpace>
					<NButton size="small" quaternary @click="selfPwOpen = true">{{ t('build.changeOwnPassword') }}</NButton>
					<NButton size="small" type="primary" @click="newUserOpen = true">{{ t('build.newUser') }}</NButton>
				</NSpace>
			</template>
			<NDataTable :columns="userColumns" :data="users" size="small" />
		</NCard>

		<NModal :show="newUserOpen" preset="dialog" :title="t('build.newUser')" :show-icon="false" @update:show="newUserOpen = $event">
			<NInput v-model:value="newUsername" :placeholder="t('login.username')" style="margin-bottom: 8px" />
			<NInput v-model:value="newPassword" type="password" :placeholder="t('login.password')" style="margin-bottom: 8px" />
			<NSelect
				v-model:value="newRole"
				:options="[
					{ label: 'member', value: 'member' },
					{ label: 'admin', value: 'admin' },
				]"
			/>
			<template #action>
				<NButton @click="newUserOpen = false">{{ t('common.cancel') }}</NButton>
				<NButton type="primary" @click="createUser">{{ t('common.confirm') }}</NButton>
			</template>
		</NModal>

		<NModal :show="resetTarget !== null" preset="dialog" :title="t('build.resetPassword')" :show-icon="false" @update:show="onResetShow">
			<NInput v-model:value="resetPassword" type="password" :placeholder="t('login.password')" />
			<template #action>
				<NButton @click="resetTarget = null">{{ t('common.cancel') }}</NButton>
				<NButton type="primary" @click="doReset">{{ t('common.confirm') }}</NButton>
			</template>
		</NModal>

		<NModal :show="selfPwOpen" preset="dialog" :title="t('build.changeOwnPassword')" :show-icon="false" @update:show="selfPwOpen = $event">
			<NInput v-model:value="selfOld" type="password" :placeholder="t('settings.oldPassword')" style="margin-bottom: 8px" />
			<NInput v-model:value="selfNew" type="password" :placeholder="t('settings.newPassword')" />
			<template #action>
				<NButton @click="selfPwOpen = false">{{ t('common.cancel') }}</NButton>
				<NButton type="primary" @click="changeOwn">{{ t('common.confirm') }}</NButton>
			</template>
		</NModal>
	</div>
</template>

<style scoped>
.header-row h2 {
	font-size: 20px;
	font-weight: 650;
	margin-bottom: 12px;
}
.settings-card {
	max-width: 860px;
	margin-bottom: 16px;
}
.form-grid {
	display: grid;
	grid-template-columns: 170px 1fr;
	gap: 10px 14px;
	align-items: center;
}
.form-grid label {
	color: rgba(0, 0, 0, 0.55);
	font-size: 13px;
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
}
.mono {
	font-family: 'JetBrains Mono', Consolas, monospace;
}
.hint {
	color: rgba(0, 0, 0, 0.45);
	font-size: 12.5px;
}
</style>
