<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
	NGrid,
	NGridItem,
	NCard,
	NButton,
	NTag,
	NModal,
	NEmpty,
} from 'naive-ui'
import { useI18n } from 'vue-i18n'
import { deviceApi, type Device } from '../api/http'

const { t } = useI18n()
const router = useRouter()

const devices = ref<Device[]>([])
const loading = ref(true)
// DS-2A: busy popup shows occupier user + session kind + elapsed duration.
const busyInfo = ref<{
	deviceId: number
	user: string
	kind: string
	since: string
} | null>(null)

onMounted(load)

async function load() {
	loading.value = true
	try {
		devices.value = await deviceApi.list()
	} catch {
		// Server may not expose /api/devices yet (T5 server part); show empty.
		devices.value = []
	} finally {
		loading.value = false
	}
}

function openSession(d: Device) {
	router.push(`/terminal/device-${d.id}`)
}
</script>

<template>
	<div>
		<h2>{{ t('device.list') }}</h2>
		<NEmpty
			v-if="!loading && devices.length === 0"
			:description="t('device.emptyHint')"
		>
			<template #icon></template>
		</NEmpty>
		<NGrid v-else :cols="4" :x-gap="12" :y-gap="12" responsive="screen" item-responsive>
			<NGridItem v-for="d in devices" :key="d.id" span="1 s:1 m:2 l:1">
				<!-- DS-1A: card per device; name + status light + busy badge first,
				     actions second, details folded. -->
				<NCard :title="d.name" size="small">
					<template #header-extra>
						<NTag size="small" type="success">{{ t('device.idle') }}</NTag>
					</template>
					<div class="device-project">{{ d.project }}</div>
					<template #footer>
						<NButton type="primary" size="small" @click="openSession(d)">
							{{ t('device.open') }}
						</NButton>
					</template>
				</NCard>
			</NGridItem>
		</NGrid>

		<!-- DS-2A: busy modal with full occupier info. -->
		<NModal
			:show="busyInfo !== null"
			preset="dialog"
			:title="t('device.busy')"
			:positive-text="t('common.confirm')"
			:negative-text="t('common.cancel')"
			@positive-click="busyInfo = null"
			@negative-click="busyInfo = null"
		>
			<span v-if="busyInfo">{{ t('device.sessionOf', { user: busyInfo.user, kind: busyInfo.kind }) }}</span>
		</NModal>
	</div>
</template>

<style scoped>
.device-project {
	color: rgba(255, 255, 255, 0.52);
	font-size: 12px;
}
</style>
