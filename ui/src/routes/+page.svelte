<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '#lib/components/ui/button/index.js';
	import { Input } from '#lib/components/ui/input/index.js';
	import * as Select from '#lib/components/ui/select/index.js';
	import * as Table from '#lib/components/ui/table/index.js';
	import * as Empty from '#lib/components/ui/empty/index.js';
	import { Kbd } from '#lib/components/ui/kbd/index.js';
	import PageHeader from '#lib/components/page-header.svelte';
	import StatusDot from '#lib/components/status-dot.svelte';
	import logo from '#lib/assets/yokify.svg';

	// rows and error come from +page.ts (GET /api/containers).
	let { data } = $props();
	const all = $derived(data.rows);

	let query = $state('');
	let status = $state('all');
	const statuses = ['all', 'Running', 'Building', 'Failed', 'Stopped'];

	const rows = $derived(
		all.filter(
			(f) =>
				(status === 'all' || f.status === status) &&
				f.name.toLowerCase().includes(query.trim().toLowerCase())
		)
	);

	function onKeydown(e: KeyboardEvent) {
		const el = e.target as HTMLElement;
		if (el.closest('input, textarea, [contenteditable]')) return;
		if (e.key === 'n' && !e.metaKey && !e.ctrlKey && !e.altKey) goto('/new');
	}
</script>

<svelte:window onkeydown={onKeydown} />

<PageHeader>
	{#snippet title()}
		<h1 class="type-heading">Functions</h1>
	{/snippet}
	{#snippet actions()}
		<Button href="/new">New function <Kbd class="ml-1">N</Kbd></Button>
	{/snippet}
</PageHeader>

{#if all.length === 0}
	<Empty.Root class="m-6 w-[90%]">
		<Empty.Header>
			<Empty.Media>
				<img src={logo} alt="" width="96" height="96" class="size-24" />
			</Empty.Media>
			<Empty.Title class="type-heading">No functions yet</Empty.Title>
			<Empty.Description>
				{#if data.error}
					Can't reach the containers: {data.error}
				{:else}
					Deploy a Git repository or a project folder and Yōkify will build and run it on Podman.
				{/if}
			</Empty.Description>
		</Empty.Header>
		<Empty.Content class="flex flex-col items-center gap-3">
			<Button href="/new">New function</Button>
			<code class="type-mono rounded-md border bg-surface px-3 py-1.5">
				yokify deploy ./my-app <br /> (Coming soon!)
			</code>
		</Empty.Content>
	</Empty.Root>
{:else}
	<section class="flex flex-col gap-4 p-6">
		<div class="flex flex-wrap gap-2">
			<label class="min-w-56 flex-1">
				<span class="sr-only">Search functions</span>
				<Input type="search" placeholder="Search functions" bind:value={query} />
			</label>
			<Select.Root type="single" bind:value={status}>
				<Select.Trigger class="w-40" aria-label="Filter by status">
					{status === 'all' ? 'All statuses' : status}
				</Select.Trigger>
				<Select.Content>
					{#each statuses as s (s)}
						<Select.Item value={s}>{s === 'all' ? 'All statuses' : s}</Select.Item>
					{/each}
				</Select.Content>
			</Select.Root>
		</div>

		<div class="overflow-x-auto">
			<Table.Root class="min-w-[760px]">
				<Table.Header>
					<Table.Row class="hover:bg-transparent">
						<Table.Head class="type-label">Name</Table.Head>
						<Table.Head class="type-label">Status</Table.Head>
						<Table.Head class="type-label">Image</Table.Head>
						<Table.Head class="type-label">Address</Table.Head>
						<Table.Head class="type-label text-right">Requests · 24 h</Table.Head>
						<Table.Head class="type-label">Last deploy</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each rows as f (f.name)}
						<Table.Row class="h-11">
							<Table.Cell>
								<a
									href="/functions/{f.name}"
									class="font-medium underline-offset-4 hover:underline"
								>
									{f.name}
								</a>
							</Table.Cell>
							<Table.Cell><StatusDot status={f.status} /></Table.Cell>
							<Table.Cell class="type-mono text-muted-foreground">{f.image}</Table.Cell>
							<Table.Cell class="type-mono">
								{#if f.address}
									<a href="http://{f.address}" class="text-teal underline-offset-4 hover:underline">
										{f.address}
									</a>
								{:else}
									<span class="text-ink-3">—</span>
								{/if}
							</Table.Cell>
							<Table.Cell class="text-right">
								{f.requests24h === null ? '—' : f.requests24h.toLocaleString('en-US')}
							</Table.Cell>
							<Table.Cell class="text-muted-foreground">{f.lastDeploy}</Table.Cell>
						</Table.Row>
					{:else}
						<Table.Row class="hover:bg-transparent">
							<Table.Cell colspan={6} class="py-10 text-center text-muted-foreground">
								No functions match “{query}”{status !== 'all' ? ` with status ${status}` : ''}.
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</div>
	</section>
{/if}
