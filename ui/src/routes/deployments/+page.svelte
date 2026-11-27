<script lang="ts">
	import { toast } from 'svelte-sonner';
	import { Button } from '#lib/components/ui/button/index.js';
	import { Badge } from '#lib/components/ui/badge/index.js';
	import * as Table from '#lib/components/ui/table/index.js';
	import PageHeader from '#lib/components/page-header.svelte';
	import StatusDot from '#lib/components/status-dot.svelte';
	import type { Deployment } from '#lib/types.js';

	let deployments = $state<Deployment[]>([]);
</script>

<PageHeader>
	{#snippet title()}
		<h1 class="type-heading">Deployments</h1>
	{/snippet}
</PageHeader>

<section class="overflow-x-auto p-6">
	<Table.Root class="min-w-[760px]">
		<Table.Header>
			<Table.Row class="hover:bg-transparent">
				<Table.Head class="type-label">Result</Table.Head>
				<Table.Head class="type-label">Function</Table.Head>
				<Table.Head class="type-label">Commit</Table.Head>
				<Table.Head class="type-label">Author</Table.Head>
				<Table.Head class="type-label">When</Table.Head>
				<Table.Head class="type-label text-right">Duration</Table.Head>
				<Table.Head><span class="sr-only">Actions</span></Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each deployments as d (d.fn + d.commit)}
				<Table.Row class="h-11">
					<Table.Cell><StatusDot status={d.result} /></Table.Cell>
					<Table.Cell>
						<a href="/functions/{d.fn}" class="font-medium underline-offset-4 hover:underline">
							{d.fn}
						</a>
					</Table.Cell>
					<Table.Cell class="max-w-80 truncate">
						<span class="type-mono text-muted-foreground">{d.commit}</span>
						{d.message}
					</Table.Cell>
					<Table.Cell class="text-muted-foreground">{d.author}</Table.Cell>
					<Table.Cell class="text-muted-foreground">{d.when}</Table.Cell>
					<Table.Cell class="text-right">{d.duration}</Table.Cell>
					<Table.Cell class="text-right">
						{#if d.live}
							<Badge variant="secondary">Live</Badge>
						{:else if d.result === 'Success'}
							<Button
								variant="outline"
								size="sm"
								onclick={() =>
									toast(`Rolled back ${d.fn} to ${d.commit}`, {
										description: 'No rebuild needed. The old image is running again.'
									})}
							>
								Roll back
							</Button>
						{:else if d.result === 'Failed'}
							<Button variant="ghost" size="sm" href="/logs">Build log</Button>
						{/if}
					</Table.Cell>
				</Table.Row>
			{:else}
				<Table.Row>
					<Table.Cell colspan={7} class="text-muted-foreground">No deployments yet.</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</section>
