<script lang="ts">
	import CharacterCard from '$lib/components/CharacterCard.svelte';
	import { charactersApi } from '$lib/api';
	import Plus from 'lucide-svelte/icons/plus';
	import { onMount } from 'svelte';

	let characters = $state<any[]>([]);
	let loading = $state(true);
	let sort = $state('created');

	let debounce: any;

	onMount(load);

	async function load() {
		loading = true;
		try {
			characters = (await charactersApi.list(sort)) || [];
		} catch (e) {
			console.error(e);
			characters = [];
		}
		loading = false;
	}

	function setSort(v: string) {
		sort = v;
		load();
	}

	const sortOptions = [
		{ id: 'created', label: 'По добавлению' },
		{ id: 'updated', label: 'По обновлению' },
		{ id: 'name', label: 'По алфавиту' }
	];

	function countLabel(n: number): string {
		const m10 = n % 10;
		const m100 = n % 100;
		if (m10 === 1 && m100 !== 11) return `${n} персонаж`;
		if (m10 >= 2 && m10 <= 4 && (m100 < 12 || m100 > 14)) return `${n} персонажа`;
		return `${n} персонажей`;
	}
</script>

<div class="page">
	<div class="topbar">
		<h1>Персонажи</h1>
		<div class="count">{countLabel(characters.length)}</div>
	</div>

	<div class="toolbar">
		<select class="select sort" value={sort} onchange={(e) => setSort(e.currentTarget.value)}>
			{#each sortOptions as o}
				<option value={o.id}>{o.label}</option>
			{/each}
		</select>
		<a href="/character/add" class="btn btn-primary">
			<Plus size={16} />
			Добавить персонажа
		</a>
	</div>

	{#if loading}
		<div class="empty">Загрузка…</div>
	{:else if characters.length === 0}
		<div class="empty">
			<div class="empty-icon">🧑‍🎤</div>
			<h2>Персонажей пока нет</h2>
			<p>Создайте первого персонажа и привяжите его к тайтлам.</p>
			<a href="/character/add" class="btn btn-primary">Добавить персонажа</a>
		</div>
	{:else}
		<div class="grid">
			{#each characters as c (c.id)}
				<CharacterCard character={c} />
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
		justify-content: space-between;
		align-items: center;
		flex-wrap: wrap;
	}
	.sort {
		width: auto;
		min-width: 150px;
	}
	.grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
		gap: 18px;
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
