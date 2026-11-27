<script lang="ts">
	import { Button } from '#lib/components/ui/button/index.js';
	import { Input } from '#lib/components/ui/input/index.js';
	import * as Select from '#lib/components/ui/select/index.js';
	import * as ToggleGroup from '#lib/components/ui/toggle-group/index.js';
	import PageHeader from '#lib/components/page-header.svelte';
	import LogView from '#lib/components/log-view.svelte';
	import type { Fn, LogLine } from '#lib/types.js';

	let functions = $state<Fn[]>([]);
	let buildLog = $state<LogLine[]>([]);
	let fnName = $state('');
	let live = $state({ lines: [] as LogLine[], paused: true });
	let source = $state('runtime');
	let level = $state('all');
	let query = $state('');

	const lines = $derived(
		(source === 'build' ? buildLog : live.lines).filter(
			(l) =>
				(level === 'all' || l.level === level) &&
				l.text.toLowerCase().includes(query.trim().toLowerCase())
		)
	);

	function download() {
		const text = lines.map((l) => `${l.t} ${l.level.toUpperCase()} ${l.text}`).join('\n');
		const a = document.createElement('a');
		a.href = URL.createObjectURL(new Blob([text], { type: 'text/plain' }));
		a.download = `${fnName}-${source}.log`;
		a.click();
		URL.revokeObjectURL(a.href);
	}
</script>

<PageHeader>
	{#snippet title()}
		<h1 class="type-heading">
			Logs{#if fnName}
				<span class="font-normal text-muted-foreground">· {fnName}</span>{/if}
		</h1>
	{/snippet}
	{#snippet actions()}
		<Button
			variant="outline"
			onclick={() => (live.paused = !live.paused)}
			disabled={source === 'build'}
		>
			{live.paused ? 'Resume' : 'Pause'}
		</Button>
		<Button variant="outline" onclick={download}>Download</Button>
	{/snippet}
</PageHeader>

<section class="flex flex-1 flex-col gap-3 p-6">
	<div class="flex flex-wrap items-center gap-2">
		<Select.Root type="single" bind:value={fnName}>
			<Select.Trigger class="w-44" aria-label="Function">{fnName || 'No functions'}</Select.Trigger>
			<Select.Content>
				{#each functions as f (f.name)}
					<Select.Item value={f.name}>{f.name}</Select.Item>
				{/each}
			</Select.Content>
		</Select.Root>
		<ToggleGroup.Root type="single" variant="outline" bind:value={source} aria-label="Log source">
			<ToggleGroup.Item value="runtime">Runtime</ToggleGroup.Item>
			<ToggleGroup.Item value="build">Build</ToggleGroup.Item>
		</ToggleGroup.Root>
		<ToggleGroup.Root type="single" variant="outline" bind:value={level} aria-label="Level">
			<ToggleGroup.Item value="all">All</ToggleGroup.Item>
			<ToggleGroup.Item value="warn">Warnings</ToggleGroup.Item>
			<ToggleGroup.Item value="error">Errors</ToggleGroup.Item>
		</ToggleGroup.Root>
		<label class="ml-auto w-full max-w-72">
			<span class="sr-only">Search logs</span>
			<Input type="search" placeholder="Search logs" bind:value={query} />
		</label>
	</div>
	<div class="flex items-center gap-2 text-muted-foreground">
		<span
			class="size-1.5 rounded-full {source === 'runtime' && !live.paused
				? 'bg-running'
				: 'bg-ink-3'}"
		></span>
		{source === 'build' ? 'Last build' : live.paused ? 'Paused' : 'Live'} ·
		<span class="tabular">{lines.length}</span> lines
	</div>
	<LogView {lines} class="min-h-[480px] flex-1" />
</section>
