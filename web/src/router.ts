import { createRouter, createWebHistory } from 'vue-router'

export const router = createRouter({
	history: createWebHistory(),
	routes: [
		{ path: '/login', component: () => import('./views/LoginView.vue') },
		{
			path: '/',
			component: () => import('./views/LayoutView.vue'),
			children: [
				{ path: '', redirect: '/devices' },
				{ path: 'devices', component: () => import('./views/DevicesView.vue') },
				{ path: 'expect', component: () => import('./views/ExpectView.vue') },
				{ path: 'history', component: () => import('./views/HistoryView.vue') },
				{ path: 'terminal/:sessionId', component: () => import('./views/TerminalView.vue') },
			],
		},
	],
})
