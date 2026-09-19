import { ref, computed } from 'vue'
import { projectApi, type ProjectRecord } from '../api/http'

// Current-project state (requirement 1.6: the project is the top-level
// container; every section works inside the selected project). Kept in
// localStorage so a refresh stays in the same project.
const PROJECT_KEY = 'embedflow.projectId'

const projects = ref<ProjectRecord[]>([])
const currentId = ref<number>(Number(localStorage.getItem(PROJECT_KEY)) || 0)

export function useProject() {
	const current = computed(() => projects.value.find((p) => p.id === currentId.value) ?? null)

	async function load() {
		try {
			projects.value = await projectApi.list()
		} catch {
			projects.value = []
		}
		// Drop a stale selection; default to the first project.
		if (!projects.value.find((p) => p.id === currentId.value)) {
			currentId.value = projects.value[0]?.id ?? 0
			persist()
		}
	}

	function select(id: number) {
		currentId.value = id
		persist()
	}

	async function create(name: string, note: string) {
		const res = await projectApi.create(name, note)
		await load()
		select(res.id)
		return res.id
	}

	async function remove(id: number) {
		await projectApi.remove(id)
		if (currentId.value === id) {
			currentId.value = 0
		}
		await load()
	}

	function persist() {
		if (currentId.value) {
			localStorage.setItem(PROJECT_KEY, String(currentId.value))
		} else {
			localStorage.removeItem(PROJECT_KEY)
		}
	}

	return { projects, current, currentId, load, select, create, remove }
}
