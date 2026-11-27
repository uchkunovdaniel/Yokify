<script lang="ts">
	import { page } from '$app/state';
	import * as Sidebar from '#lib/components/ui/sidebar/index.js';
	import logo from '#lib/assets/yokify.svg';

	const nav = [
		{ title: 'Functions', href: '/' },
		{ title: 'Deployments', href: '/deployments' },
		{ title: 'Logs', href: '/logs' },
		{ title: 'Settings', href: '/settings' }
	];

	function isActive(href: string) {
		const path = page.url.pathname;
		return href === '/' ? path === '/' || path.startsWith('/functions') : path.startsWith(href);
	}
</script>

<Sidebar.Root collapsible="offcanvas">
	<Sidebar.Header class="px-4 pt-4 pb-2">
		<a href="/" class="flex items-center gap-2" aria-label="Yōkify home">
			<img src={logo} alt="logo" width="24" height="24" class="size-6" />
			<span class="wordmark text-[17px] leading-none text-foreground">Yōkify</span>
		</a>
	</Sidebar.Header>
	<Sidebar.Content>
		<Sidebar.Group>
			<Sidebar.GroupContent>
				<Sidebar.Menu class="gap-2">
					{#each nav as item (item.href)}
						<Sidebar.MenuItem>
							<Sidebar.MenuButton class="h-12" isActive={isActive(item.href)}>
								{#snippet child({ props })}
									<a href={item.href} {...props}>
										<span>{item.title}</span>
									</a>
								{/snippet}
							</Sidebar.MenuButton>
							<!--{#if item.count}-->
							<!--	<Sidebar.MenuBadge class="tabular text-ink-3">{item.count}</Sidebar.MenuBadge>-->
							<!--{/if}-->
						</Sidebar.MenuItem>
					{/each}
				</Sidebar.Menu>
			</Sidebar.GroupContent>
		</Sidebar.Group>
		<Sidebar.Separator class="mx-4" />
		<div class="type-label px-4 py-2 leading-[18px]">
			<!--			<div>{engine.name} {engine.version} · {engine.host}</div>-->
			<!--			<div class="tabular">{engine.containers} containers · {engine.memory}</div>-->
		</div>
	</Sidebar.Content>
</Sidebar.Root>
