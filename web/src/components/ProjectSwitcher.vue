<script setup lang="ts">
import { ref } from 'vue'
import { NSelect, NButton } from 'naive-ui'
import { useI18n } from 'vue-i18n'
import { useProject } from '../store/project'

// Project switcher: the top-level container selector (requirement 1.6).
const { t } = useI18n()
const { projects, currentId, select, create } = useProject()
const creating = ref(false)
const newName = ref('')

async function quickCreate() {
	if (!newName.value.trim()) return
	await create(newName.value.trim(), '')
	newName.value = ''
	creating.value = false
}
</script>

<template>
	<div class="project-switcher">
		<NSelect
			v-if="projects.length > 0"
			:value="currentId"
			:options="projects.map((p) => ({ label: p.name, value: p.id }))"
			size="small"
			:placeholder="t('nav.projectNone')"
			@update:value="(v: number) => select(v)"
		/>
		<div v-else class="none-hint">{{ t('nav.projectNone') }}</div>
		<NButton
			v-if="!creating"
			quaternary
			size="tiny"
			class="add-btn"
			:title="t('project.new')"
			@click="creating = true"
		>+</NButton>
		<div v-else class="quick-create">
			<input
				v-model="newName"
				class="n-input"
				:placeholder="t('project.name')"
				@keyup.enter="quickCreate"
				@keyup.esc="creating = false"
			/>
			<NButton size="tiny" type="primary" @click="quickCreate">OK</NButton>
		</div>
	</div>
</template>

<style scoped>
.project-switcher {
	display: flex;
	align-items: center;
	gap: 6px;
	padding: 2px 14px 14px;
}
.project-switcher :deep(.n-select) {
	flex: 1;
}
.none-hint {
	flex: 1;
	font-size: 13px;
	opacity: 0.55;
	padding: 4px 2px;
}
.add-btn {
	flex: 0 0 auto;
}
.quick-create {
	display: flex;
	gap: 4px;
	flex: 1;
}
input.n-input {
	width: 100%;
	padding: 4px 8px;
	border-radius: 6px;
	border: 1px solid rgba(255, 255, 255, 0.24);
	background: transparent;
	color: inherit;
	font-size: 13px;
}
</style>
