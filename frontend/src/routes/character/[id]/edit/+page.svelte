<script lang="ts">
	import { page } from '$app/stores';
	import CharacterForm from '$lib/components/CharacterForm.svelte';
	import { charactersApi } from '$lib/api';
	import type { Character } from '../../../../app.d';

	let character = $state<Character | null>(null);
	let loading = $state(true);

	const id = $derived(Number($page.params.id));

	$effect(() => {
		if (!id) return;
		load(id);
	});

	async function load(charId: number) {
		character = await charactersApi.get(charId);
		loading = false;
	}
</script>

<div class="page">
	<h1>Редактирование персонажа</h1>
	{#if loading}
		<div class="loading">Загрузка…</div>
	{:else if character}
		<CharacterForm {character} />
	{:else}
		<div class="loading">Персонаж не найден</div>
	{/if}
</div>

<style>
	.page {
		max-width: 720px;
		margin: 0 auto;
	}
	h1 {
		font-size: 24px;
		margin-bottom: 20px;
	}
	.loading {
		text-align: center;
		padding: 60px;
		color: var(--text-dim);
	}
</style>
