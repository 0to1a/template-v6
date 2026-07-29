import path from 'node:path';

import tailwindcss from '@tailwindcss/vite';
import react from '@vitejs/plugin-react';
import { tanstackRouter } from '@tanstack/router-plugin/vite';
import { defineConfig } from 'vite';

// Generic regex proxy; new services need no extra entry
const backendTarget = 'http://localhost:8080';
const connectPathPattern = String.raw`^/(?:[a-z][a-z0-9_]*\.)+[A-Za-z0-9]+Service/`;

export default defineConfig({
	plugins: [
		tanstackRouter({
			target: 'react',
			autoCodeSplitting: true,
			routesDirectory: './src/routes',
			generatedRouteTree: './src/routeTree.gen.ts'
		}),
		react(),
		tailwindcss()
	],
	resolve: {
		alias: {
			'@': path.resolve(import.meta.dirname, './src')
		}
	},
	build: {
		outDir: 'dist'
	},
	server: {
		proxy: {
			[connectPathPattern]: backendTarget,
			'/health': backendTarget
		}
	},
	test: {
		environment: 'jsdom',
		globals: false,
		include: ['src/**/*.{test,spec}.{js,ts,jsx,tsx}']
	}
});
