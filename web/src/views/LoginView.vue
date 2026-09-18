<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { NCard, NForm, NFormItem, NButton, useMessage } from 'naive-ui'
import { useI18n } from 'vue-i18n'
import { ofetch } from 'ofetch'
import { setToken } from '../api/http'

const { t } = useI18n()
const router = useRouter()
const message = useMessage()

const username = ref('')
const password = ref('')
const loading = ref(false)

async function submit() {
	loading.value = true
	try {
		// M1 interim: the server issues a session token on login. The
		// endpoint lands with the auth slice; bootstrap flow uses the same
		// shape (requirement 1.5: server-side session, cookie or token).
		const res = await ofetch<{ token: string }>('/api/login', {
			method: 'POST',
			body: { username: username.value, password: password.value },
		})
		setToken(res.token)
		router.push('/devices')
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
	background: #101014;
}
.login-card {
	width: 360px;
}
input.n-input {
	width: 100%;
	padding: 6px 10px;
	border-radius: 3px;
	border: 1px solid rgba(255, 255, 255, 0.24);
	background: transparent;
	color: inherit;
}
</style>
