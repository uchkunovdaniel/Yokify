<script lang="ts">
	import './layout.css';
	import * as Sidebar from '#lib/components/ui/sidebar/index.js';
	import { Toaster } from '#lib/components/ui/sonner/index.js';
	import * as Tooltip from '#lib/components/ui/tooltip/index.js';
	import AppSidebar from '#lib/components/app-sidebar.svelte';
	import { ModeWatcher } from 'mode-watcher';
	import logo from '#lib/assets/yokify.svg';

	let { children } = $props();
</script>

<ModeWatcher defaultMode="dark" />

<svelte:head>
	<link rel="icon" href={logo} />
	<title>Yōkify</title>
</svelte:head>

<Tooltip.Provider delayDuration={150}>
	<Sidebar.Provider>
		<AppSidebar />
		<Sidebar.Inset>
			{@render children()}
		</Sidebar.Inset>
	</Sidebar.Provider>
</Tooltip.Provider>
<Toaster position="bottom-right" />
