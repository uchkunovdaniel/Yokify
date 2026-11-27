<script lang="ts">
	import type { LogLine } from '#lib/types.js';

	// A log well: surface background, mono 12/18, time in ink-3. Status colour only on the level.
	let { lines, class: className = '' }: { lines: LogLine[]; class?: string } = $props();

	const levelClass = { info: 'text-ink-3', warn: 'text-building', error: 'text-failed' };
</script>

<div
	class="type-mono overflow-auto rounded-lg border bg-surface px-3.5 py-3 {className}"
	role="log"
	aria-live="polite"
>
	{#each lines as l, i (i)}
		<div class="flex gap-4 whitespace-pre">
			<time class="text-ink-3">{l.t}</time>
			<span class="w-10 shrink-0 uppercase {levelClass[l.level]}">{l.level}</span>
			<span class={l.level === 'error' ? 'text-failed' : 'text-foreground'}>{l.text}</span>
		</div>
	{/each}
</div>
