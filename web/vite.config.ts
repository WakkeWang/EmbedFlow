import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'

// https://vite.dev/config/
export default defineConfig({
	plugins: [vue()],
	server: {
		proxy: {
			// Dev proxy: API and WS ride the same origin as the dev server
			// (production serves the built assets from the Go binary instead).
			'/api': 'http://127.0.0.1:18601',
			'/ws': { target: 'ws://127.0.0.1:18601', ws: true },
		},
	},
})
