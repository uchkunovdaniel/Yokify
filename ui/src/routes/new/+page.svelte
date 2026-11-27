<script lang="ts">
	import { toast } from 'svelte-sonner';
	import { Button } from '#lib/components/ui/button/index.js';
	import { Input } from '#lib/components/ui/input/index.js';
	import { Label } from '#lib/components/ui/label/index.js';
	import { Badge } from '#lib/components/ui/badge/index.js';
	import * as RadioGroup from '#lib/components/ui/radio-group/index.js';
	import * as Select from '#lib/components/ui/select/index.js';
	import * as Breadcrumb from '#lib/components/ui/breadcrumb/index.js';
	import PageHeader from '#lib/components/page-header.svelte';

	let repo = $state('');
	let folder = $state('');
	let name = $state('');
	let port = $state('');
	let source = $state('git');
	let branch = $state('main');

	const repos: string[] = [];
	const portDetected = '';
</script>

<PageHeader>
	{#snippet title()}
		<Breadcrumb.Root>
			<Breadcrumb.List class="type-heading">
				<Breadcrumb.Item>
					<Breadcrumb.Link href="/" class="text-muted-foreground">Functions</Breadcrumb.Link>
				</Breadcrumb.Item>
				<Breadcrumb.Separator>/</Breadcrumb.Separator>
				<Breadcrumb.Item>
					<Breadcrumb.Page class="font-semibold">New function</Breadcrumb.Page>
				</Breadcrumb.Item>
			</Breadcrumb.List>
		</Breadcrumb.Root>
	{/snippet}
</PageHeader>

<form
	class="flex max-w-2xl flex-col gap-6 p-6"
	onsubmit={(e) => {
		e.preventDefault();
		toast('Not connected yet', { description: 'Deploying needs the Yōkify server.' });
	}}
>
	<section class="flex flex-col gap-3">
		<h2 class="type-heading">Where is the code?</h2>
		<RadioGroup.Root bind:value={source} class="grid gap-3 sm:grid-cols-2">
			{#each [{ v: 'git', t: 'Git repository', d: 'GitHub, GitLab or Bitbucket. Redeploys on push.' }, { v: 'folder', t: 'Folder on this machine', d: 'Same as yokify deploy ./path.' }] as o (o.v)}
				<Label
					for="src-{o.v}"
					class="flex cursor-pointer items-start gap-3 rounded-xl border p-4 has-[[data-state=checked]]:border-teal has-[[data-state=checked]]:bg-teal-tint"
				>
					<RadioGroup.Item value={o.v} id="src-{o.v}" class="mt-0.5" />
					<span class="flex flex-col gap-1">
						<span class="font-medium">{o.t}</span>
						<span class="font-normal text-muted-foreground">{o.d}</span>
					</span>
				</Label>
			{/each}
		</RadioGroup.Root>
	</section>

	{#if source === 'git'}
		<div class="grid gap-4 sm:grid-cols-[minmax(0,1fr)_160px]">
			<div class="flex flex-col gap-2">
				<Label for="repo">Repository</Label>
				<Select.Root type="single" bind:value={repo}>
					<Select.Trigger id="repo" class="w-full">
						{repo || 'Connect a Git account first'}
					</Select.Trigger>
					<Select.Content>
						{#each repos as r (r)}<Select.Item value={r}>{r}</Select.Item>{/each}
					</Select.Content>
				</Select.Root>
			</div>
			<div class="flex flex-col gap-2">
				<Label for="branch">Branch</Label>
				<Input id="branch" class="type-mono" bind:value={branch} />
			</div>
		</div>
	{:else}
		<div class="flex flex-col gap-2">
			<Label for="folder">Project folder</Label>
			<Input id="folder" class="type-mono" bind:value={folder} />
		</div>
	{/if}

	<div class="grid gap-4 sm:grid-cols-[minmax(0,1fr)_160px]">
		<div class="flex flex-col gap-2">
			<Label for="name">Name</Label>
			<Input id="name" class="type-mono" bind:value={name} />
			<span class="type-label">Lowercase letters, numbers and dashes.</span>
		</div>
		<div class="flex flex-col gap-2">
			<Label for="port" class="flex items-center gap-2">
				Port
				{#if port && port === portDetected}<Badge class="bg-teal-tint text-teal">
						Detected
					</Badge>{/if}
			</Label>
			<Input id="port" class="type-mono" inputmode="numeric" bind:value={port} />
		</div>
	</div>

	<p class="text-muted-foreground">
		Yōkify detects the rest (image, install, build and run) and shows it on the Build tab before the
		first deploy finishes.
	</p>
	<div class="flex gap-2">
		<Button type="submit">Deploy</Button>
		<Button variant="outline" href="/">Cancel</Button>
	</div>
</form>
