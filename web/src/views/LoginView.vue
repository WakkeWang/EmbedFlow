<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { NCard, NForm, NFormItem, NButton, useMessage } from 'naive-ui'
import { useI18n } from 'vue-i18n'
import { ofetch } from 'ofetch'
import { setToken } from '../api/http'
import { useProject } from '../store/project'

const { t } = useI18n()
const router = useRouter()
const message = useMessage()
const { projects, load, create } = useProject()

const username = ref('')
const password = ref('')
const loading = ref(false)
// First-run form: when no project exists yet, offer to create one on login.
const newProjectName = ref('')
const hasProjects = ref(true)

async function submit() {
	loading.value = true
	try {
		const res = await ofetch<{ token: string }>('/api/login', {
			method: 'POST',
			body: { username: username.value, password: password.value },
		})
		setToken(res.token)
		await load()
		if (projects.value.length === 0 && newProjectName.value.trim()) {
			await create(newProjectName.value.trim(), '')
		}
		hasProjects.value = projects.value.length > 0
		router.push('/')
	} catch {
		message.error(t('login.error'))
	} finally {
		loading.value = false
	}
}
</script>

<template>
	<div class="login-wrap">
		<NCard :title="t('login.title')" class="login-card">
			<NForm @submit.prevent="submit">
				<NFormItem :label="t('login.username')">
					<input v-model="username" class="n-input" data-testid="username" />
				</NFormItem>
				<NFormItem :label="t('login.password')">
					<input
						v-model="password"
						type="password"
						class="n-input"
						data-testid="password"
					/>
				</NFormItem>
				<NFormItem v-if="!hasProjects" :label="t('project.name')">
					<input v-model="newProjectName" class="n-input" :placeholder="t('project.emptyHint')" />
				</NFormItem>
				<NButton type="primary" attr-type="submit" :loading="loading" block>
					{{ t('login.submit') }}
				</NButton>
			</NForm>
		</NCard>
	</div>
</template>

<style scoped>
.login-wrap {
	display: flex;
	align-items: center;
	justify-content: center;
	height: 100vh;
	background: linear-gradient(160deg, #eef2fb 0%, #f6f7f9 55%, #e9eef8 100%);
}
.login-card {
	width: 380px;
	border-radius: 16px;
}
input.n-input {
	width: 100%;
	padding: 8px 12px;
	border-radius: 6px;
	border: 1px solid rgba(0, 0, 0, 0.16);
	background: transparent;
	color: inherit;
}
</style>
