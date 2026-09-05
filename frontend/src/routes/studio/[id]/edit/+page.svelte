<script lang="ts">
	import { page } from '$app/stores';
	import StudioForm from '$lib/components/StudioForm.svelte';
	import { studiosApi } from '$lib/api';
	import type { Studio } from '../../../../app.d';

	let studio = $state<Studio | null>(null);
	let loading = $state(true);

	const id = $derived(Number($page.params.id));

	$effect(() => {
		if (!id) return;
		load(id);
	});

	async function load(studioId: number) {
		studio = await studiosApi.get(studioId);
		loading = false;
	}
</script>

<div class="page">
	<h1>Редактирование студии</h1>
	{#if loading}
		<div class="loading">Загрузка…</div>
	{:else if studio}
		<StudioForm {studio} />
	{:else}
		<div class="loading">Студия не найдена</div>
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
