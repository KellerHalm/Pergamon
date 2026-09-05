<script lang="ts">
	import '../app.css';
	import { page } from '$app/stores';
	import Library from 'lucide-svelte/icons/library';
	import BookOpen from 'lucide-svelte/icons/book-open';
	import Users from 'lucide-svelte/icons/users';
	import Plus from 'lucide-svelte/icons/plus';
	import Settings from 'lucide-svelte/icons/settings';
	import { onMount } from 'svelte';

	let { children } = $props();

	let theme = $state<'dark' | 'light'>('dark');
	let mobileNav = false;

	onMount(() => {
		theme = (localStorage.getItem('theme') as 'dark' | 'light') || 'dark';
	});

	const navItems = [
		{ href: '/', label: 'Каталог', icon: BookOpen },
		{ href: '/shelf', label: 'Полка', icon: Library },
		{ href: '/characters', label: 'Персонажи', icon: Users },
		{ href: '/add', label: 'Добавить', icon: Plus },
		{ href: '/settings', label: 'Настройки', icon: Settings }
	];

	function isActive(href: string, path: string): boolean {
		if (href === '/') return path === '/' || path.startsWith('/title');
		if (href === '/characters') return path === '/characters' || path.startsWith('/character/');
		return path.startsWith(href);
	}

	function toggleTheme() {
		theme = theme === 'dark' ? 'light' : 'dark';
		localStorage.setItem('theme', theme);
		document.documentElement.setAttribute('data-theme', theme);
	}
</script>

<div class="app-shell">
	<aside class="sidebar">
		<div class="logo">
			<Library size={22} />
			<span>Медиатека</span>
		</div>
		<nav>
			{#each navItems as item}
				<a
					href={item.href}
					class:active={isActive(item.href, $page.url.pathname)}
				>
					<item.icon size={18} />
					<span>{item.label}</span>
				</a>
			{/each}
		</nav>
		<button class="theme-toggle" onclick={toggleTheme} title="Сменить тему">
			{theme === 'dark' ? '☀️ Светлая' : '🌙 Тёмная'}
		</button>
		<div class="version">v1.1</div>
	</aside>

	<nav class="mobile-nav">
		{#each navItems as item}
			<a
				href={item.href}
				class:active={isActive(item.href, $page.url.pathname)}
			>
				<item.icon size={20} />
				<span>{item.label}</span>
			</a>
		{/each}
	</nav>

	<main>
		{@render children()}
	</main>
</div>

<style>
	.app-shell {
		display: flex;
		min-height: 100vh;
	}

	.sidebar {
		width: 200px;
		background: var(--bg-elev);
		border-right: 1px solid var(--border);
		display: flex;
		flex-direction: column;
		padding: 16px 12px;
		position: sticky;
		top: 0;
		height: 100vh;
		flex-shrink: 0;
	}

	.logo {
		display: flex;
		align-items: center;
		gap: 10px;
		font-weight: 700;
		font-size: 16px;
		padding: 8px 12px 20px;
		color: var(--accent);
	}

	nav {
		display: flex;
		flex-direction: column;
		gap: 2px;
		flex: 1;
	}

	nav a {
		display: flex;
		align-items: center;
		gap: 10px;
		padding: 10px 12px;
		border-radius: var(--radius-sm);
		color: var(--text-dim);
		font-size: 14px;
		font-weight: 500;
		transition: background 0.15s, color 0.15s;
	}

	nav a:hover {
		background: var(--bg-elev-2);
		color: var(--text);
	}

	nav a.active {
		background: var(--accent-soft);
		color: var(--accent);
	}

	.theme-toggle {
		font-size: 12px;
		padding: 8px 12px;
		text-align: left;
		color: var(--text-dim);
		border-radius: var(--radius-sm);
	}
	.theme-toggle:hover {
		background: var(--bg-elev-2);
		color: var(--text);
	}

	.version {
		font-size: 10px;
		color: var(--text-faint);
		padding: 4px 12px;
	}

	.mobile-nav {
		display: none;
	}

	main {
		flex: 1;
		min-width: 0;
		padding: 24px 32px 100px;
		max-width: 100%;
	}

	@media (max-width: 768px) {
		.sidebar {
			display: none;
		}
		.mobile-nav {
			display: flex;
			position: fixed;
			bottom: 0;
			left: 0;
			right: 0;
			background: var(--bg-elev);
			border-top: 1px solid var(--border);
			z-index: 100;
			padding: 6px;
			justify-content: space-around;
		}
		.mobile-nav a {
			display: flex;
			flex-direction: column;
			align-items: center;
			gap: 2px;
			padding: 6px 12px;
			font-size: 10px;
			color: var(--text-dim);
			border-radius: var(--radius-sm);
		}
		.mobile-nav a.active {
			color: var(--accent);
		}
		main {
			padding: 16px;
			padding-bottom: 80px;
		}
	}
</style>
