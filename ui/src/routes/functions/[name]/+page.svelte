<script lang="ts">
	import { untrack } from 'svelte';
	import { page } from '$app/state';
	import { toast } from 'svelte-sonner';
	import { Button } from '#lib/components/ui/button/index.js';
	import { Badge } from '#lib/components/ui/badge/index.js';
	import { Input } from '#lib/components/ui/input/index.js';
	import { Switch } from '#lib/components/ui/switch/index.js';
	import { Label } from '#lib/components/ui/label/index.js';
	import { Separator } from '#lib/components/ui/separator/index.js';
	import * as Tabs from '#lib/components/ui/tabs/index.js';
	import * as Table from '#lib/components/ui/table/index.js';
	import * as Breadcrumb from '#lib/components/ui/breadcrumb/index.js';
	import * as Select from '#lib/components/ui/select/index.js';
	import PageHeader from '#lib/components/page-header.svelte';
	import StatusDot from '#lib/components/status-dot.svelte';
	import { sourceLabel, type Fn, type Deployment, type LogLine } from '#lib/types.js';
	import * as Empty from '#lib/components/ui/empty/index.js';
	import LogView from '#lib/components/log-view.svelte';

	let fn = $state<Fn | undefined>(undefined);
	let history = $state<Deployment[]>([]);
	let log = $state({ lines: [] as LogLine[], paused: true });

	// Build settings suggested by Yōkipack. "Detected" until someone edits the value.
	const detected = { image: '', install: '', build: '', port: '', run: '' };

	let build = $state(untrack(() => ({ ...detected })));
	const fields = [
		{ key: 'image', label: 'Image' },
		{ key: 'install', label: 'Install' },
		{ key: 'build', label: 'Build' },
		{ key: 'port', label: 'Port' },
		{ key: 'run', label: 'Run' }
	] as const;

	const dockerfile = $derived(
		[
			`FROM ${build.image}`,
			'WORKDIR /app',
			'COPY package.json package-lock.json ./',
			`RUN ${build.install}`,
			'COPY . .',
			`RUN ${build.build}`,
			`ENV PORT=${build.port}`,
			`EXPOSE ${build.port}`,
			`CMD ${build.run}`
		].join('\n')
	);

	let autoStop = $state(false);
	let idleAfter = $state('15');

	let env = $state<{ key: string; value: string; secret: boolean; shown: boolean }[]>([]);
</script>

{#if fn}
	<PageHeader>
		{#snippet title()}
			<Breadcrumb.Root>
				<Breadcrumb.List class="type-heading">
					<Breadcrumb.Item>
						<Breadcrumb.Link href="/" class="text-muted-foreground">Functions</Breadcrumb.Link>
					</Breadcrumb.Item>
					<Breadcrumb.Separator>/</Breadcrumb.Separator>
					<Breadcrumb.Item>
						<Breadcrumb.Page class="font-semibold">{fn.name}</Breadcrumb.Page>
					</Breadcrumb.Item>
				</Breadcrumb.List>
			</Breadcrumb.Root>
		{/snippet}
		{#snippet actions()}
			{#if fn.address}
				<Button href="http://{fn.address}" class="bg-teal text-primary-foreground hover:bg-teal/85">
					Open app
				</Button>
			{/if}
			<Button
				onclick={() =>
					toast('Deploy started', {
						description: `${fn.name} · building from ${sourceLabel(fn.source)}`
					})}
			>
				Deploy
			</Button>
		{/snippet}
	</PageHeader>

	<Tabs.Root value="overview" class="gap-0">
		<div class="border-b px-6">
			<Tabs.List variant="line" class="h-11">
				<Tabs.Trigger value="overview">Overview</Tabs.Trigger>
				<Tabs.Trigger value="build">Build</Tabs.Trigger>
				<Tabs.Trigger value="logs">Logs</Tabs.Trigger>
				<Tabs.Trigger value="environment">Environment</Tabs.Trigger>
			</Tabs.List>
		</div>

		<!-- Overview -->
		<Tabs.Content
			value="overview"
			class="grid gap-8 p-6 lg:grid-cols-[minmax(0,1fr)_minmax(0,1fr)]"
		>
			<section class="flex flex-col gap-3">
				<h2 class="type-heading">Details</h2>
				<dl class="grid grid-cols-[120px_minmax(0,1fr)] gap-x-4 gap-y-3">
					<dt class="text-muted-foreground">Status</dt>
					<dd><StatusDot status={fn.status} /></dd>
					<dt class="text-muted-foreground">Source</dt>
					<dd>{sourceLabel(fn.source)}</dd>
					<dt class="text-muted-foreground">Image</dt>
					<dd class="type-mono">{fn.image}</dd>
					<dt class="text-muted-foreground">Address</dt>
					<dd class="type-mono">
						{#if fn.address}<a
								href="http://{fn.address}"
								class="text-teal underline-offset-4 hover:underline"
							>
								{fn.address}
							</a>{:else}—{/if}
					</dd>
					<dt class="text-muted-foreground">Requests · 24 h</dt>
					<dd class="tabular">{fn.requests24h?.toLocaleString('en-US') ?? '—'}</dd>
					<dt class="text-muted-foreground">Engine</dt>
					<dd>Podman</dd>
				</dl>

				<Separator class="my-3" />

				<div class="flex items-start justify-between gap-6">
					<div>
						<Label for="autostop" class="type-heading">Stop when idle</Label>
						<p class="mt-1 max-w-[52ch] text-muted-foreground">
							Stops the container after a period with no requests and starts it on the next one. The
							first request then waits for the start.
						</p>
					</div>
					<Switch id="autostop" bind:checked={autoStop} />
				</div>
				{#if autoStop}
					<div class="flex items-center gap-3">
						<Label for="idle">Stop after</Label>
						<Select.Root type="single" bind:value={idleAfter}>
							<Select.Trigger id="idle" class="w-32">{idleAfter} minutes</Select.Trigger>
							<Select.Content>
								{#each ['5', '15', '30', '60'] as m (m)}
									<Select.Item value={m}>{m} minutes</Select.Item>
								{/each}
							</Select.Content>
						</Select.Root>
					</div>
				{/if}
			</section>

			<section class="flex flex-col gap-3">
				<h2 class="type-heading">Deployments</h2>
				<Table.Root>
					<Table.Body>
						{#each history as d (d.commit)}
							<Table.Row class="h-12">
								<Table.Cell class="w-28"><StatusDot status={d.result} /></Table.Cell>
								<Table.Cell class="w-full max-w-0">
									<div class="truncate">
										<span class="type-mono text-muted-foreground">{d.commit}</span>
										{d.message}
									</div>
									<div class="type-label">{d.author} · {d.when} · {d.duration}</div>
								</Table.Cell>
								<Table.Cell class="text-right">
									{#if d.live}
										<Badge variant="secondary">Live</Badge>
									{:else if d.result === 'Success'}
										<Button
											variant="outline"
											size="sm"
											onclick={() =>
												toast(`Rolled back to ${d.commit}`, {
													description: 'No rebuild needed. The old image is running again.'
												})}
										>
											Roll back
										</Button>
									{/if}
								</Table.Cell>
							</Table.Row>
						{:else}
							<Table.Row>
								<Table.Cell class="text-muted-foreground">No deployments yet.</Table.Cell>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			</section>
		</Tabs.Content>

		<!-- Build: settings suggested by Yōkipack -->
		<Tabs.Content value="build" class="grid gap-8 p-6 lg:grid-cols-[minmax(0,1.1fr)_minmax(0,1fr)]">
			<form
				class="flex flex-col"
				onsubmit={(e) => {
					e.preventDefault();
					toast('Saved build settings', { description: 'Deploy to apply' });
				}}
			>
				<p class="mb-4 max-w-[60ch] text-muted-foreground">
					Yōkify found a {fn.runtime} project. These are suggestions; change anything before you deploy.
				</p>
				{#each fields as f (f.key)}
					<div class="grid grid-cols-[96px_minmax(0,1fr)_72px] items-center gap-3 border-b py-2.5">
						<Label for={f.key}>{f.label}</Label>
						<Input id={f.key} class="type-mono h-7" bind:value={build[f.key]} />
						{#if build[f.key] === detected[f.key]}
							<Badge class="bg-teal-tint text-teal">Detected</Badge>
						{:else}
							<span class="type-label">Edited</span>
						{/if}
					</div>
				{/each}
				<div class="mt-5 flex gap-2">
					<Button type="submit">Save</Button>
					<Button variant="outline" type="button" onclick={() => (build = { ...detected })}>
						Reset to detected
					</Button>
				</div>
			</form>
			<section class="flex flex-col gap-2">
				<h2 class="type-label tracking-wider uppercase">Dockerfile preview</h2>
				<pre class="type-mono overflow-x-auto rounded-lg border bg-surface p-4">{dockerfile}</pre>
			</section>
		</Tabs.Content>

		<!-- Logs -->
		<Tabs.Content value="logs" class="flex flex-col gap-3 p-6">
			<div class="flex items-center justify-between">
				<span class="flex items-center gap-2 text-muted-foreground">
					<span class="size-1.5 rounded-full {log.paused ? 'bg-ink-3' : 'bg-running'}"></span>
					{log.paused ? 'Paused' : 'Live'}
				</span>
				<Button variant="outline" size="sm" onclick={() => (log.paused = !log.paused)}>
					{log.paused ? 'Resume' : 'Pause'}
				</Button>
			</div>
			<LogView lines={log.lines} class="h-[420px]" />
		</Tabs.Content>

		<!-- Environment -->
		<Tabs.Content value="environment" class="flex max-w-3xl flex-col gap-4 p-6">
			<p class="text-muted-foreground">
				Variables are passed to the container on the next deploy. Secrets stay hidden until you show
				them.
			</p>
			<Table.Root>
				<Table.Header>
					<Table.Row class="hover:bg-transparent">
						<Table.Head class="type-label">Name</Table.Head>
						<Table.Head class="type-label">Value</Table.Head>
						<Table.Head><span class="sr-only">Actions</span></Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each env as v (v.key)}
						<Table.Row class="h-11">
							<Table.Cell class="type-mono">{v.key}</Table.Cell>
							<Table.Cell class="type-mono text-muted-foreground">
								{v.shown ? v.value : '••••••••••••'}
							</Table.Cell>
							<Table.Cell class="text-right">
								{#if v.secret}
									<Button variant="ghost" size="sm" onclick={() => (v.shown = !v.shown)}>
										{v.shown ? 'Hide' : 'Show'}
									</Button>
								{/if}
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
			<div><Button variant="outline">Add variable</Button></div>
		</Tabs.Content>
	</Tabs.Root>
{:else}
	<PageHeader>
		{#snippet title()}<h1 class="type-heading">Function not found</h1>{/snippet}
	</PageHeader>
	<Empty.Root class="m-6 border">
		<Empty.Header>
			<Empty.Title class="type-heading">No function called “{page.params.name}”</Empty.Title>
			<Empty.Description>
				It may have been removed, or the server is not connected yet.
			</Empty.Description>
		</Empty.Header>
		<Empty.Content><Button variant="outline" href="/">Back to functions</Button></Empty.Content>
	</Empty.Root>
{/if}
