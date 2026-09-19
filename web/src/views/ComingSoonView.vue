<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

// Structural placeholder for sections whose entities land with M2/M4/M5.
const props = defineProps<{ milestone: string }>()
const { t } = useI18n()
const n = computed(() => props.milestone.replace(/^M/i, ''))
</script>

<template>
	<div class="placeholder">
		<div class="placeholder-icon">
			<span>{{ props.milestone }}</span>
		</div>
		<div class="placeholder-title">{{ t('comingsoon.title', { n }) }}</div>
		<div class="placeholder-body">{{ t('comingsoon.body', { n }) }}</div>
	</div>
</template>

<style scoped>
.placeholder {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	min-height: 50vh;
	text-align: center;
}
.placeholder-icon {
	width: 64px;
	height: 64px;
	border-radius: 18px;
	background: linear-gradient(135deg, rgba(47, 107, 255, 0.14), rgba(47, 107, 255, 0.04));
	border: 1px solid rgba(47, 107, 255, 0.28);
	color: #2f6bff;
	font-size: 18px;
	font-weight: 700;
	display: flex;
	align-items: center;
	justify-content: center;
	margin-bottom: 18px;
}
.placeholder-title {
	font-size: 16px;
	font-weight: 600;
	margin-bottom: 8px;
}
.placeholder-body {
	color: rgba(0, 0, 0, 0.45);
	font-size: 13.5px;
	max-width: 460px;
	line-height: 1.7;
	padding: 0 16px;
}
</style>
