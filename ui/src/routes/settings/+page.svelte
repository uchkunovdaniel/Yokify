<script lang="ts">
	import { toast } from 'svelte-sonner';
	import { Button } from '#lib/components/ui/button/index.js';
	import { Input } from '#lib/components/ui/input/index.js';
	import { Label } from '#lib/components/ui/label/index.js';
	import { Switch } from '#lib/components/ui/switch/index.js';
	import { Badge } from '#lib/components/ui/badge/index.js';
	import { Separator } from '#lib/components/ui/separator/index.js';
	import * as RadioGroup from '#lib/components/ui/radio-group/index.js';
	import PageHeader from '#lib/components/page-header.svelte';
	import ThemeSwitcher from '#lib/components/theme-switcher.svelte';
	import { page } from '$app/state';

	let engineChoice = $state('podman');

	let domain = $state('localhost');
	let autoStopDefault = $state(false);
	let git = $state([
		{ name: 'GitHub', account: '', repos: 0, connected: false },
		{ name: 'GitLab', account: '', repos: 0, connected: false },
		{ name: 'Bitbucket', account: '', repos: 0, connected: false }
	]);

	const sections = ['Engine', 'Git', 'Defaults', 'Users', 'Updates'];

	const active = $derived(page.url.hash || `#${sections[0].toLowerCase()}`);
</script>

<PageHeader>
	{#snippet title()}
		<h1 class="type-heading">Settings</h1>
	{/snippet}
	{#snippet actions()}
		<ThemeSwitcher />
	{/snippet}
</PageHeader>

<div class="flex flex-wrap gap-8 p-6">
	<nav aria-label="Settings sections" class="flex w-40 flex-col gap-1">
		{#each sections as s (s)}
			{@const href = `#${s.toLowerCase()}`}
			<a
				{href}
				aria-current={href === active ? 'true' : undefined}
				class="rounded-md px-2 py-1.5 {href === active
					? 'bg-accent text-foreground'
					: 'text-muted-foreground hover:bg-accent hover:text-foreground'}"
			>
				{s}
			</a>
		{/each}
	</nav>

	<div class="flex max-w-2xl min-w-0 flex-1 flex-col gap-8">
		<section id="engine" class="flex flex-col gap-4">
			<div>
				<h2 class="type-heading">Container engine</h2>
				<p class="mt-1 text-muted-foreground">
					Builds and runs your functions. Changing it applies to the next deploy.
				</p>
			</div>
			<RadioGroup.Root bind:value={engineChoice} class="grid gap-3 sm:grid-cols-2">
				{#each [{ v: 'podman', t: 'Podman', d: 'Rootless, no background daemon.' }, { v: 'docker', t: 'Docker', d: 'Uses Docker Engine if it is installed.' }] as o (o.v)}
					<Label
						for="engine-{o.v}"
						class="flex cursor-pointer items-start gap-3 rounded-xl border p-4 has-[[data-state=checked]]:border-teal has-[[data-state=checked]]:bg-teal-tint"
					>
						<RadioGroup.Item value={o.v} id="engine-{o.v}" class="mt-0.5" />
						<span class="flex flex-col gap-1">
							<span class="flex items-center gap-2 font-medium">
								{o.t}
								{#if o.v === 'podman'}<Badge variant="secondary">Default</Badge>{/if}
							</span>
							<span class="font-normal text-muted-foreground">{o.d}</span>
						</span>
					</Label>
				{/each}
			</RadioGroup.Root>
			<div class="flex items-center justify-between gap-4">
				<span class="flex items-center gap-2 text-muted-foreground">
					<span class="size-1.5 rounded-full bg-ink-3"></span>
					Not connected
				</span>
				<Button
					variant="outline"
					onclick={() =>
						toast('Not connected yet', { description: 'The Yōkify server is not running.' })}
				>
					Test connection
				</Button>
			</div>
		</section>

		<Separator />

		<section id="git" class="flex flex-col gap-2">
			<h2 class="type-heading">Git accounts</h2>
			<p class="mb-2 text-muted-foreground">
				Connect an account to deploy its repositories and redeploy on every push.
			</p>
			{#each git as g (g.name)}
				<div class="flex flex-wrap items-center gap-3 border-t py-3">
					<div class="min-w-40 flex-1">
						<div class="font-medium">{g.name}</div>
						<div class="type-label">
							{g.connected
								? `Connected as ${g.account} · ${g.repos} repositories`
								: 'Not connected'}
						</div>
					</div>
					{#if g.connected}
						<Button variant="outline" size="sm" onclick={() => (g.connected = false)}>
							Disconnect
						</Button>
					{:else}
						<Button
							variant="outline"
							size="sm"
							onclick={() => toast(`Opening ${g.name} to authorise Yōkify`)}
						>
							Connect {g.name}
						</Button>
					{/if}
				</div>
			{/each}
		</section>

		<Separator />

		<form
			id="defaults"
			class="flex flex-col gap-5"
			onsubmit={(e) => {
				e.preventDefault();
				toast('Saved defaults', { description: 'New functions use these settings.' });
			}}
		>
			<div>
				<h2 class="type-heading">Defaults for new functions</h2>
				<p class="mt-1 text-muted-foreground">
					Each function can change these in its own settings.
				</p>
			</div>
			<div class="flex max-w-sm flex-col gap-2">
				<Label for="domain">Base address</Label>
				<Input id="domain" class="type-mono" bind:value={domain} />
				<span class="type-label">
					Functions are reachable at <span class="type-mono">{domain}:PORT</span>.
				</span>
			</div>
			<div class="flex items-start justify-between gap-6">
				<div>
					<Label for="autostop-default">Stop idle functions</Label>
					<p class="mt-1 text-muted-foreground">
						Off by default. Saves memory on small machines; the first request after a stop waits for
						the start.
					</p>
				</div>
				<Switch id="autostop-default" bind:checked={autoStopDefault} />
			</div>
			<div><Button type="submit">Save</Button></div>
		</form>
	</div>
</div>
