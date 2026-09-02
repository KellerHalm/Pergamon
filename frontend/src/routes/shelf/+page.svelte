<script lang="ts">
	import { onMount } from 'svelte';
	import { shelvesApi, titlesApi, coverApi } from '$lib/api';
	import { TYPES, displayName, categoryLabel } from '$lib/constants';
	import type { Title, Shelf } from '../../app.d';
	import Plus from 'lucide-svelte/icons/plus';
	import Trash2 from 'lucide-svelte/icons/trash-2';
	import X from 'lucide-svelte/icons/x';
	import House from 'lucide-svelte/icons/house';

	let shelves = $state<Shelf[]>([]);
	let titlesMap = $state<Record<number, Title>>({});
	let loading = $state(true);
	let activeShelf: number | null = $state(null);
	let showAddShelf = $state(false);
	let newShelfName = $state('');
	let newShelfKind = $state('read');
	let coverCache: Record<string, string> = $state({});
	let showPicker = $state(false);
	let pickerSearch = $state('');
	let pickerSelected = $state<number[]>([]);
	let allTitles = $state<Title[]>([]);
	let adding = $state(false);

	const mainShelf = $derived({
		id: 0,
		name: 'Все тайтлы',
		kind: 'main',
		position: -1,
		titleIds: allTitles.map((t) => t.id),
		createdAt: ''
	});

	const activeShelfObj = $derived(
		activeShelf ? shelves.find((s) => s.id === activeShelf) || mainShelf : mainShelf
	);
	const pickerTitles = $derived.by(() => {
		const kind = activeShelfObj && activeShelfObj.id !== 0 ? activeShelfObj.kind : '';
		const q = pickerSearch.trim().toLowerCase();
		return allTitles.filter(
			(t) => (!kind || t.type === kind) && (!q || displayName(t.names).toLowerCase().includes(q))
		);
	});

	onMount(async () => {
		await load();
	});

	async function load() {
		loading = true;
		shelves = (await shelvesApi.list()) || [];
		allTitles = (await titlesApi.list({})) || [];
		const map: Record<number, Title> = {};
		for (const t of allTitles) map[t.id] = t;
		titlesMap = map;
		for (const t of allTitles) {
			if (t.cover && !(t.cover in coverCache)) {
				const url = await coverApi.dataURL(t.cover);
				if (url) coverCache[t.cover] = url;
			}
		}
		if (activeShelf === null) activeShelf = 0;
		loading = false;
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
			activeShelf = shelves.find((s) => s.id !== id)?.id ?? 0;
		}
		await load();
	}

	function openPicker() {
		if (!activeShelfObj || activeShelfObj.id === 0) return;
		pickerSelected = [...(activeShelfObj.titleIds || [])];
		pickerSearch = '';
		showPicker = true;
	}

	function togglePick(id: number) {
		pickerSelected = pickerSelected.includes(id)
			? pickerSelected.filter((x) => x !== id)
			: [...pickerSelected, id];
	}

	async function applyPicker() {
		if (!activeShelfObj || activeShelfObj.id === 0) return;
		adding = true;
		try {
			await shelvesApi.setItems(activeShelfObj.id, pickerSelected);
			showPicker = false;
			await load();
		} catch (e: any) {
			alert('Ошибка сохранения полки: ' + (e?.message || JSON.stringify(e)));
		} finally {
			adding = false;
		}
	}

	function shelfTitles(shelf: Shelf): Title[] {
		return shelf.titleIds.map((id) => titlesMap[id]).filter(Boolean);
	}

	function spineHeight(t: Title): number {
		return 110 + ((t.id * 137) % 61);
	}

	function spinalColor(title: Title): string {
		return title.spineColor || '#5b6470';
	}
</script>

<div class="page">
	<div class="topbar">
		<h1>Виртуальная полка</h1>
		{#if activeShelfObj.id !== 0}
			<button class="btn" onclick={openPicker}>
				<Plus size={14} /> Добавить тайтл
			</button>
		{/if}
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

	<div class="shelf-tabs">
		<button class="tab main" class:active={activeShelf === 0 || activeShelf === null} onclick={() => (activeShelf = 0)}>
			<House size={13} />
			{mainShelf.name}
			<span class="tab-count">{allTitles.length}</span>
		</button>
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

	{#if shelfTitles(activeShelfObj).length === 0}
		<div class="empty-sm">
			{#if activeShelfObj.id === 0}
				<p>В каталоге пока нет тайтлов</p>
				<a href="/add" class="btn btn-primary">Добавить тайтл</a>
			{:else}
				Пустая полка — нажмите «Добавить тайтл», чтобы наполнить её
			{/if}
		</div>
	{:else}
		<div class="bookshelf">
			<div class="wood-plank top"></div>
			<div class="books">
				{#each shelfTitles(activeShelfObj) as t (t.id)}
					{#if t.type === 'read'}
						<a href="/title/{t.id}" class="book-spine" style="background:{spinalColor(t)}; height:{spineHeight(t)}px">
							<span class="spine-text">{displayName(t.names)}</span>
						</a>
					{:else}
						<a href="/title/{t.id}" class="cassette">
							<span class="cassette-strip" style="background:{spinalColor(t)}"></span>
							<span class="cassette-label">
								{#if coverCache[t.cover]}
									<img src={coverCache[t.cover]} alt="" />
								{:else}
									<span class="cassette-name">{displayName(t.names)}</span>
								{/if}
							</span>
							<span class="cassette-window">
								<span class="reel"></span>
								<span class="tape"></span>
								<span class="reel"></span>
							</span>
						</a>
					{/if}
				{/each}
			</div>
			<div class="wood-plank bottom"></div>
		</div>
	{/if}

	{#if showPicker}
		<div class="picker-overlay">
			<div class="picker">
				<div class="picker-head">
					<h3>Тайтлы{activeShelfObj.id !== 0 ? ` — «${activeShelfObj.name}»` : ''}</h3>
					<button class="btn" onclick={() => (showPicker = false)}><X size={14} /></button>
				</div>
				<input class="input" placeholder="Поиск…" bind:value={pickerSearch} />
				<div class="picker-list">
					{#each pickerTitles as t (t.id)}
						<button
							class="picker-row"
							class:sel={pickerSelected.includes(t.id)}
							onclick={() => togglePick(t.id)}
						>
							<span class="pick-check">{pickerSelected.includes(t.id) ? '✓' : ''}</span>
							<span class="pick-name">{displayName(t.names)}</span>
							<span class="pick-cat">{categoryLabel(t.category)}</span>
						</button>
					{:else}
						<div class="picker-empty">Ничего не найдено</div>
					{/each}
				</div>
				<div class="picker-foot">
					<button class="btn btn-primary" onclick={applyPicker} disabled={adding}>
						{adding ? 'Сохранение…' : 'Готово'}
					</button>
				</div>
			</div>
		</div>
	{/if}
</div>

<svelte:window onkeydown={(e) => e.key === 'Escape' && (showPicker = false)} />

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
	.empty-sm {
		text-align: center;
		padding: 40px;
		color: var(--text-dim);
	}
	.empty-sm p {
		margin-bottom: 12px;
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
	.tab.main {
		font-weight: 600;
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
		gap: 6px;
		padding: 12px 28px 14px;
		min-height: 210px;
		overflow-x: auto;
	}
	.book-spine {
		width: 36px;
		min-height: 110px;
		flex: none;
		border-radius: 3px;
		cursor: pointer;
		text-decoration: none;
		display: flex;
		align-items: flex-end;
		justify-content: center;
		transition: transform 0.15s, filter 0.15s;
		position: relative;
		box-shadow: 1px 0 4px rgba(0, 0, 0, 0.35);
	}
	.book-spine::before {
		content: '';
		position: absolute;
		inset: 0;
		background: linear-gradient(90deg, rgba(255, 255, 255, 0.15) 0%, transparent 30%, rgba(0, 0, 0, 0.15) 100%);
		border-radius: 3px;
	}
	.book-spine:hover {
		transform: translateY(-6px);
		filter: brightness(1.15);
	}
	.spine-text {
		writing-mode: vertical-rl;
		text-orientation: mixed;
		font-size: 11px;
		font-weight: 700;
		color: rgba(255, 255, 255, 0.9);
		text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
		padding: 6px 2px;
		overflow: hidden;
		max-height: 92%;
	}

	.cassette {
		flex: none;
		display: flex;
		flex-direction: column;
		gap: 6px;
		width: 178px;
		height: 106px;
		padding: 8px;
		border-radius: 6px;
		background: linear-gradient(180deg, #26282f 0%, #17181d 100%);
		box-shadow: 0 4px 10px rgba(0, 0, 0, 0.35), inset 0 1px 0 rgba(255, 255, 255, 0.07);
		cursor: pointer;
		text-decoration: none;
		transition: transform 0.15s;
	}
	.cassette:hover {
		transform: translateY(-5px);
	}
	.cassette-strip {
		display: block;
		flex: none;
		height: 4px;
		border-radius: 2px;
	}
	.cassette-label {
		flex: 1;
		min-height: 0;
		border-radius: 3px;
		background: #ece7d9;
		overflow: hidden;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	.cassette-label img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}
	.cassette-name {
		font-size: 10px;
		font-weight: 700;
		color: #2c2c31;
		padding: 4px 6px;
		text-align: center;
		overflow: hidden;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		line-clamp: 2;
		-webkit-box-orient: vertical;
	}
	.cassette-window {
		flex: none;
		display: flex;
		align-items: center;
		gap: 6px;
		height: 26px;
		padding: 0 8px;
		border-radius: 3px;
		background: #0c0d10;
		box-shadow: inset 0 1px 4px rgba(0, 0, 0, 0.8);
	}
	.reel {
		flex: none;
		width: 14px;
		height: 14px;
		border-radius: 50%;
		border: 3px solid #454a54;
		box-shadow: inset 0 0 0 2px #0c0d10;
	}
	.tape {
		flex: 1;
		height: 3px;
		border-radius: 2px;
		background: #2b2e35;
	}

	.wood-plank {
		height: 22px;
		border-radius: 3px;
		background: var(--shelf-wood);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3), inset 0 2px 4px rgba(255, 255, 255, 0.1);
	}
	.wood-plank.bottom {
		margin: 0 16px;
	}
	.wood-plank.top {
		margin: 0 16px 2px;
		height: 10px;
		box-shadow: inset 0 -2px 4px rgba(0, 0, 0, 0.2);
	}

	.picker-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.5);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 50;
	}
	.picker {
		width: min(480px, calc(100vw - 32px));
		max-height: 80vh;
		display: flex;
		flex-direction: column;
		gap: 10px;
		background: var(--bg-elev);
		border-radius: var(--radius);
		padding: 16px;
		box-shadow: var(--shadow);
	}
	.picker-head {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}
	.picker-head h3 {
		font-size: 15px;
		font-weight: 600;
	}
	.picker-list {
		flex: 1;
		overflow-y: auto;
		min-height: 120px;
		display: flex;
		flex-direction: column;
		gap: 4px;
	}
	.picker-row {
		display: flex;
		align-items: center;
		gap: 10px;
		padding: 8px 10px;
		border-radius: var(--radius-sm);
		background: var(--bg-elev-2);
		text-align: left;
		font-size: 13px;
		transition: filter 0.12s;
	}
	.picker-row:hover {
		filter: brightness(1.15);
	}
	.picker-row.sel {
		background: var(--accent);
		color: #fff;
	}
	.pick-check {
		width: 14px;
		font-weight: 700;
	}
	.pick-name {
		flex: 1;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
	.pick-cat {
		font-size: 11px;
		opacity: 0.7;
	}
	.picker-empty {
		text-align: center;
		color: var(--text-dim);
		padding: 24px;
	}
	.picker-foot {
		display: flex;
		justify-content: flex-end;
	}

	@media (max-width: 640px) {
		.books {
			padding: 10px 12px 8px;
		}
		.book-spine {
			width: 28px;
		}
		.cassette {
			width: 150px;
			height: 92px;
		}
	}
</style>
