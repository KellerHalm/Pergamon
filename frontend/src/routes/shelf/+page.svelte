<script lang="ts">
	import { onMount } from 'svelte';
	import { shelvesApi, titlesApi, coverApi } from '$lib/api';
	import { TYPES, displayName, categoryLabel } from '$lib/constants';
	import type { Title, Shelf } from '../../app.d';
	import Plus from 'lucide-svelte/icons/plus';
	import Trash2 from 'lucide-svelte/icons/trash-2';
	import X from 'lucide-svelte/icons/x';

	let shelves = $state<Shelf[]>([]);
	let titlesMap = $state<Record<number, Title>>({});
	let loading = $state(true);
	let activeShelf: number | null = $state(null);
	let showAddShelf = $state(false);
	let newShelfName = $state('');
	let newShelfKind = $state('read');
	let coverCache: Record<string, string> = $state({});

	onMount(async () => {
		await load();
	});

	async function load() {
		loading = true;
		shelves = (await shelvesApi.list()) || [];
		if (shelves.length > 0 && activeShelf === null) {
			activeShelf = shelves[0].id;
		}
		const allIds = shelves.flatMap((s) => s.titleIds || []);
		const unique = [...new Set(allIds)];
		const map: Record<number, Title> = {};
		for (const id of unique) {
			map[id] = await titlesApi.get(id);
		}
		titlesMap = map;
		for (const id of unique) {
			const t = map[id];
			if (t?.cover) {
				const url = await coverApi.dataURL(t.cover);
				if (url) coverCache[t.cover] = url;
			}
		}
		loading = false;
	}

	async function getCover(filename: string): Promise<string> {
		if (coverCache[filename]) return coverCache[filename];
		const url = await coverApi.dataURL(filename);
		coverCache[filename] = url || '';
		return coverCache[filename];
	}

	async function addShelf() {
		if (!newShelfName.trim()) return;
		try {
			const id = await shelvesApi.save({
				id: 0,
				name: newShelfName.trim(),
				kind: newShelfKind,
				position: shelves.length,
				titleIds: [],
				createdAt: ''
			});
			newShelfName = '';
			showAddShelf = false;
			await load();
			activeShelf = id;
		} catch (e: any) {
			alert('Ошибка добавления полки: ' + (e?.message || JSON.stringify(e)));
		}
	}

	async function deleteShelf(id: number) {
		if (!confirm('Удалить полку?')) return;
		await shelvesApi.delete(id);
		if (activeShelf === id) {
			activeShelf = shelves.find((s) => s.id !== id)?.id || null;
		}
		await load();
	}

	function shelfTitles(shelf: Shelf): Title[] {
		return shelf.titleIds.map((id) => titlesMap[id]).filter(Boolean);
	}

	function spinalColor(title: Title): string {
		return title.spineColor || '#5b6470';
	}
</script>

<div class="page">
	<div class="topbar">
		<h1>Виртуальная полка</h1>
		<button class="btn btn-primary" onclick={() => (showAddShelf = !showAddShelf)}>
			<Plus size={14} /> Новая полка
		</button>
	</div>

	{#if showAddShelf}
		<div class="add-shelf">
			<input class="input" placeholder="Название полки" bind:value={newShelfName} onkeydown={(e) => e.key === 'Enter' && addShelf()} />
			<select class="select sm" bind:value={newShelfKind}>
				{#each TYPES as t}<option value={t.id}>{t.label}</option>{/each}
			</select>
			<button class="btn btn-primary" onclick={addShelf}>Создать</button>
			<button class="btn" onclick={() => (showAddShelf = false)}><X size={14} /></button>
		</div>
	{/if}

	{#if shelves.length === 0}
		<div class="empty">
			<div class="empty-icon">📚</div>
			<h2>Нет полок</h2>
			<p>Создайте полку, чтобы расставить тайтлы</p>
		</div>
	{:else}
		<div class="shelf-tabs">
			{#each shelves as sh}
				<button class="tab" class:active={activeShelf === sh.id} onclick={() => (activeShelf = sh.id)}>
					{sh.name}
					<span class="tab-count">{sh.titleIds?.length ?? 0}</span>
				</button>
				<button class="tab-del" onclick={() => deleteShelf(sh.id)} title="Удалить">
					<Trash2 size={12} />
				</button>
			{/each}
		</div>

		{#if activeShelf}
			{@const shelf = shelves.find((s) => s.id === activeShelf)}
			{#if shelf}
				{#if shelfTitles(shelf).length === 0}
					<div class="empty-sm">Пустая полка — перетащите сюда тайтлы</div>
				{:else if shelf.kind === 'read'}
					<div class="bookshelf">
						<div class="wood-plank top"></div>
						<div class="books">
							{#each shelfTitles(shelf) as t, i (t.id)}
								<a href="/title/{t.id}" class="book-spine" style="background:{spinalColor(t)}; height:{60 + Math.random() * 40}px">
									<span class="spine-text">{displayName(t.names)}</span>
								</a>
							{/each}
						</div>
						<div class="wood-plank bottom"></div>
					</div>
				{:else}
					<div class="mediashelf">
						{#each shelfTitles(shelf) as t (t.id)}
							<a href="/title/{t.id}" class="dvd-case">
								<div class="dvd-spine" style="background:{spinalColor(t)}">
									<span class="dvd-text">{displayName(t.names)}</span>
								</div>
								<div class="dvd-cover">
									{#if coverCache[t.cover]}
										<img src={coverCache[t.cover]} alt="" />
									{:else}
										<span class="dvd-placeholder">D</span>
									{/if}
								</div>
							</a>
						{/each}
					</div>
				{/if}
			{/if}
		{/if}
	{/if}
	</div>

<style>
	.page {
		max-width: 1200px;
		margin: 0 auto;
	}
	.topbar {
		display: flex;
		align-items: center;
		gap: 12px;
		margin-bottom: 16px;
	}
	h1 {
		font-size: 24px;
		font-weight: 700;
	}
	.empty,
	.empty-sm {
		text-align: center;
		padding: 40px;
		color: var(--text-dim);
	}
	.empty-icon {
		font-size: 48px;
		margin-bottom: 8px;
	}

	.add-shelf {
		display: flex;
		gap: 8px;
		margin-bottom: 16px;
		background: var(--bg-elev);
		padding: 14px;
		border-radius: var(--radius);
	}
	.sm {
		width: auto;
		min-width: 100px;
	}

	.shelf-tabs {
		display: flex;
		gap: 4px;
		margin-bottom: 20px;
		flex-wrap: wrap;
	}
	.tab {
		display: inline-flex;
		align-items: center;
		gap: 6px;
		padding: 8px 14px;
		border-radius: var(--radius-sm);
		background: var(--bg-elev);
		font-size: 13px;
		font-weight: 500;
		color: var(--text-dim);
		transition: all 0.12s;
	}
	.tab:hover {
		color: var(--text);
	}
	.tab.active {
		background: var(--accent);
		color: #fff;
	}
	.tab-count {
		font-size: 11px;
		opacity: 0.7;
	}
	.tab-del {
		padding: 4px;
		color: var(--text-faint);
		border-radius: 4px;
	}
	.tab-del:hover {
		color: #e0484d;
	}

	.bookshelf {
		margin-bottom: 10px;
	}
	.books {
		display: flex;
		align-items: flex-end;
		gap: 3px;
		padding: 0 24px 8px;
		min-height: 120px;
	}
	.book-spine {
		width: 22px;
		min-height: 60px;
		border-radius: 2px;
		cursor: pointer;
		text-decoration: none;
		display: flex;
		align-items: flex-end;
		justify-content: center;
		transition: transform 0.15s, filter 0.15s;
		position: relative;
		box-shadow: 1px 0 3px rgba(0, 0, 0, 0.3);
	}
	.book-spine::before {
		content: '';
		position: absolute;
		inset: 0;
		background: linear-gradient(90deg, rgba(255, 255, 255, 0.15) 0%, transparent 30%, rgba(0, 0, 0, 0.15) 100%);
		border-radius: 2px;
	}
	.book-spine:hover {
		transform: translateY(-6px);
		filter: brightness(1.15);
	}
	.spine-text {
		writing-mode: vertical-rl;
		text-orientation: mixed;
		font-size: 8px;
		font-weight: 700;
		color: rgba(255, 255, 255, 0.9);
		text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
		padding: 4px 1px;
		overflow: hidden;
		max-height: 90%;
	}

	.wood-plank {
		height: 18px;
		border-radius: 3px;
		background: var(--shelf-wood);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3), inset 0 2px 4px rgba(255, 255, 255, 0.1);
	}
	.wood-plank.bottom {
		margin: 0 16px;
	}
	.wood-plank.top {
		margin: 0 16px 2px;
		height: 8px;
		box-shadow: inset 0 -2px 4px rgba(0, 0, 0, 0.2);
	}

	.mediashelf {
		display: flex;
		flex-wrap: wrap;
		gap: 16px;
		padding: 0 20px 20px;
	}
	.dvd-case {
		display: flex;
		flex-direction: column;
		width: 100px;
		cursor: pointer;
		text-decoration: none;
		transition: transform 0.15s;
	}
	.dvd-case:hover {
		transform: translateY(-4px);
	}
	.dvd-spine {
		height: 8px;
		border-radius: 2px 2px 0 0;
		display: flex;
		align-items: center;
		justify-content: center;
		box-shadow: 0 1px 4px rgba(0, 0, 0, 0.3);
	}
	.dvd-text {
		display: none;
	}
	.dvd-cover {
		flex: 1;
		min-height: 140px;
		background: var(--bg-elev-2);
		border-radius: 0 0 4px 4px;
		display: flex;
		align-items: center;
		justify-content: center;
		overflow: hidden;
		box-shadow: var(--shadow);
	}
	.dvd-cover img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}
	.dvd-placeholder {
		font-size: 36px;
		font-weight: 800;
		color: var(--text-faint);
	}

	@media (max-width: 640px) {
		.books {
			padding: 0 12px 6px;
			overflow-x: auto;
		}
		.book-spine {
			width: 18px;
		}
	}
</style>
