<script lang="ts">
	import type { Snippet } from 'svelte';
	import * as Sidebar from '#lib/components/ui/sidebar/index.js';

	// The bar at the top of every page: title (or breadcrumb) on the left, actions on the right.
	let { title, actions }: { title: Snippet; actions?: Snippet } = $props();
</script>

<header class="flex min-h-15 flex-wrap items-center gap-3 border-b px-6 py-3">
	<Sidebar.Trigger class="-ml-2 md:hidden" />
	<div class="min-w-0 flex-1">{@render title()}</div>
	{#if actions}
		<div class="flex flex-wrap items-center gap-2">{@render actions()}</div>
	{/if}
</header>
