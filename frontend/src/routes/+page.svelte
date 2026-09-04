<script lang="ts">
	import TitleCard from '$lib/components/TitleCard.svelte';
	import { titlesApi, metaApi } from '$lib/api';
	import { TYPES, ALL_CATEGORIES, STATUSES, typeLabel, categoryLabel } from '$lib/constants';
	import Search from 'lucide-svelte/icons/search';
	import Grid from 'lucide-svelte/icons/layout-grid';
	import List from 'lucide-svelte/icons/list';
	import Sliders from 'lucide-svelte/icons/sliders-horizontal';
	import { onMount } from 'svelte';

	let titles = $state<any[]>([]);
	let loading = $state(true);
	let view = $state<'grid' | 'list'>('grid');
	let showFilters = $state(false);

	let sort = $state('created');
	let typeFilter = $state('');
	let catFilter = $state('');
	let statusFilter = $state('');
	let tagFilter = $state('');
	let search = $state('');
	let allTags = $state<string[]>([]);

	let debounce: any;

	onMount(async () => {
		allTags = (await metaApi.allTags()) || [];
		await load();
	});

	async function load() {
		loading = true;
		try {
			titles = (await titlesApi.list({
				sort,
				type: typeFilter,
				category: catFilter,
				status: statusFilter,
				search,
				tags: tagFilter ? [tagFilter] : []
			})) || [];
		} catch (e) {
			console.error(e);
			titles = [];
		}
		loading = false;
	}

	function onSearch(v: string) {
		search = v;
		clearTimeout(debounce);
		debounce = setTimeout(load, 250);
	}

	function setSort(v: string) {
		sort = v;
		load();
	}
	function setFilter(key: 'type' | 'cat' | 'status' | 'tag', v: string) {
		if (key === 'type') typeFilter = v;
		if (key === 'cat') catFilter = v;
		if (key === 'status') statusFilter = v;
		if (key === 'tag') tagFilter = v;
		load();
	}

	const sortOptions = [
		{ id: 'created', label: 'По добавлению' },
		{ id: 'updated', label: 'По обновлению' },
		{ id: 'title', label: 'По алфавиту' },
		{ id: 'score', label: 'По оценке' }
	];
</script>

<div class="page">
	<div class="topbar">
		<h1>Каталог</h1>
		<div class="count">{titles.length} тайтлов</div>
	</div>

	<div class="toolbar">
		<div class="search-box">
			<Search size={16} />
			<input
				class="search-input"
				placeholder="Поиск по названию или автору…"
				value={search}
				oninput={(e) => onSearch(e.currentTarget.value)}
			/>
		</div>
		<select class="select sort" value={sort} onchange={(e) => setSort(e.currentTarget.value)}>
			{#each sortOptions as o}
				<option value={o.id}>{o.label}</option>
			{/each}
		</select>
		<button class="btn btn-icon" class:active={showFilters} onclick={() => (showFilters = !showFilters)} title="Фильтры">
			<Sliders size={16} />
		</button>
		<div class="view-toggle">
			<button class="btn btn-icon" class:active={view === 'grid'} onclick={() => (view = 'grid')} title="Сетка">
				<Grid size={16} />
			</button>
			<button class="btn btn-icon" class:active={view === 'list'} onclick={() => (view = 'list')} title="Список">
				<List size={16} />
			</button>
		</div>
	</div>

	{#if showFilters}
		<div class="filters">
			<div class="filter-group">
				<span class="label">Тип</span>
				<div class="chips">
					<button class="chip" class:active={typeFilter === ''} onclick={() => setFilter('type', '')}>Все</button>
					{#each TYPES as t}
						<button class="chip" class:active={typeFilter === t.id} onclick={() => setFilter('type', t.id)}>
							{t.label}
						</button>
					{/each}
				</div>
			</div>
			<div class="filter-group">
				<span class="label">Категория</span>
				<div class="chips">
					<button class="chip" class:active={catFilter === ''} onclick={() => setFilter('cat', '')}>Все</button>
					{#each ALL_CATEGORIES as c}
						<button class="chip" class:active={catFilter === c.id} onclick={() => setFilter('cat', c.id)}>
							{c.label}
						</button>
					{/each}
				</div>
			</div>
			<div class="filter-group">
				<span class="label">Мой статус</span>
				<div class="chips">
					<button class="chip" class:active={statusFilter === ''} onclick={() => setFilter('status', '')}>Все</button>
					{#each STATUSES as s}
						<button class="chip" class:active={statusFilter === s.id} onclick={() => setFilter('status', s.id)}>
							{s.label}
						</button>
					{/each}
				</div>
			</div>
			{#if allTags.length > 0}
				<div class="filter-group">
					<span class="label">Тег</span>
					<div class="chips">
						<button class="chip" class:active={tagFilter === ''} onclick={() => setFilter('tag', '')}>Все</button>
						{#each allTags as t}
							<button class="chip" class:active={tagFilter === t} onclick={() => setFilter('tag', t)}>{t}</button>
						{/each}
					</div>
				</div>
			{/if}
		</div>
	{/if}

	{#if loading}
		<div class="empty">Загрузка…</div>
	{:else if titles.length === 0}
		<div class="empty">
			<div class="empty-icon">📚</div>
			<h2>Пока пусто</h2>
			<p>Добавьте свой первый тайтл — книгу, мангу, фильм или сериал.</p>
			<a href="/add" class="btn btn-primary">Добавить тайтл</a>
		</div>
	{:else if view === 'grid'}
		<div class="grid">
			{#each titles as t (t.id)}
				<TitleCard title={t} view="grid" />
			{/each}
		</div>
	{:else}
		<div class="list-view">
			{#each titles as t (t.id)}
				<TitleCard title={t} view="list" />
			{/each}
		</div>
	{/if}
</div>

<style>
	.page {
		max-width: 1400px;
		margin: 0 auto;
	}
	.topbar {
		display: flex;
		align-items: baseline;
		gap: 14px;
		margin-bottom: 18px;
	}
	h1 {
		font-size: 26px;
		font-weight: 700;
	}
	.count {
		color: var(--text-faint);
		font-size: 13px;
	}

	.toolbar {
		display: flex;
		gap: 8px;
		margin-bottom: 14px;
		flex-wrap: wrap;
	}
	.search-box {
		flex: 1;
		min-width: 200px;
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 0 12px;
		background: var(--bg-elev);
		border-radius: var(--radius-sm);
		color: var(--text-dim);
	}
	.search-input {
		flex: 1;
		background: none;
		border: none;
		padding: 9px 0;
		color: var(--text);
	}
	.search-input:focus {
		outline: none;
	}
	.sort {
		width: auto;
		min-width: 150px;
	}
	.view-toggle {
		display: flex;
		gap: 2px;
		background: var(--bg-elev);
		border-radius: var(--radius-sm);
		padding: 2px;
	}
	.view-toggle .btn-icon {
		height: 32px;
		width: 32px;
		background: none;
	}
	.btn-icon.active {
		background: var(--bg-elev-2);
		color: var(--accent);
	}

	.filters {
		background: var(--bg-elev);
		border-radius: var(--radius);
		padding: 14px 16px;
		margin-bottom: 18px;
		display: flex;
		flex-direction: column;
		gap: 12px;
	}
	.filter-group {
		display: flex;
		flex-direction: column;
		gap: 6px;
	}
	.chips {
		display: flex;
		gap: 6px;
		flex-wrap: wrap;
	}
	.chip {
		padding: 4px 11px;
		border-radius: 99px;
		font-size: 12px;
		background: var(--bg-elev-2);
		color: var(--text-dim);
		transition: background 0.12s, color 0.12s;
	}
	.chip:hover {
		color: var(--text);
	}
	.chip.active {
		background: var(--accent);
		color: #fff;
	}

	.grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
		gap: 18px;
	}
	.list-view {
		display: flex;
		flex-direction: column;
		gap: 6px;
	}

	.empty {
		text-align: center;
		padding: 60px 20px;
		color: var(--text-dim);
	}
	.empty-icon {
		font-size: 56px;
		margin-bottom: 12px;
	}
	.empty h2 {
		color: var(--text);
		margin-bottom: 6px;
	}
	.empty p {
		margin-bottom: 18px;
	}

	@media (max-width: 640px) {
		.grid {
			grid-template-columns: repeat(auto-fill, minmax(110px, 1fr));
			gap: 12px;
		}
	}
</style>
