<script lang="ts">
	import type { Status } from '#lib/types.js';

	// Status colours show state only, at dot or text size. The word is always shown too.
	let { status }: { status: Status | 'Success' } = $props();

	const color: Record<string, string> = {
		Running: 'bg-running',
		Success: 'bg-running',
		Building: 'bg-building',
		Failed: 'bg-failed',
		Stopped: 'bg-ink-3'
	};
</script>

<span class="inline-flex items-center gap-2 font-medium text-foreground">
	<span
		class="size-1.5 shrink-0 rounded-full transition-colors duration-200 {color[status]}"
		aria-hidden="true"
	></span>
	{status}
</span>
