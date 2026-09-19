import { createRouter, createWebHistory } from 'vue-router'

// Route tree mirrors the requirement-1.6 organization: dashboard + four
// sections. Deploy holds the M1 entities (devices / flash rules / session
// history); Build holds the M2 entities (items / batches / records); Test
// and Release are structural placeholders (M4/M5).
export const router = createRouter({
	history: createWebHistory(),
	routes: [
		{ path: '/login', component: () => import('./views/LoginView.vue') },
		{
			path: '/',
			component: () => import('./views/LayoutView.vue'),
			children: [
				{ path: '', component: () => import('./views/DashboardView.vue') },
				{
					path: 'build',
					component: () => import('./views/BuildLayout.vue'),
					children: [
						{ path: '', redirect: '/build/items' },
						{ path: 'items', component: () => import('./views/BuildItemsView.vue') },
						{ path: 'batches', component: () => import('./views/BatchesView.vue') },
						{ path: 'records', component: () => import('./views/BuildRecordsView.vue') },
					],
				},
				{
					path: 'deploy',
					component: () => import('./views/DeployLayout.vue'),
					children: [
						{ path: '', redirect: '/deploy/devices' },
						{ path: 'devices', component: () => import('./views/DevicesView.vue') },
						{ path: 'rules', component: () => import('./views/DeployRulesView.vue') },
						{ path: 'deployments', component: () => import('./views/DeploymentsView.vue') },
						{ path: 'flash', component: () => import('./views/ExpectView.vue') },
						{ path: 'history', component: () => import('./views/HistoryView.vue') },
					],
				},
				{ path: 'test', component: () => import('./views/ComingSoonView.vue'), props: { milestone: 'M4' } },
				{ path: 'release', component: () => import('./views/ComingSoonView.vue'), props: { milestone: 'M5' } },
				{ path: 'settings', component: () => import('./views/SettingsView.vue') },
				{ path: 'devices', redirect: '/deploy/devices' },
				{ path: 'expect', redirect: '/deploy/flash' },
				{ path: 'history', redirect: '/deploy/history' },
				{ path: 'terminal/:sessionId', component: () => import('./views/TerminalView.vue') },
			],
		},
	],
})
