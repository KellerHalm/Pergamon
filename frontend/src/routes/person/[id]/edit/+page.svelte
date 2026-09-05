<script lang="ts">
	import { page } from '$app/stores';
	import PersonForm from '$lib/components/PersonForm.svelte';
	import { peopleApi } from '$lib/api';
	import type { Person } from '../../../../app.d';

	let person = $state<Person | null>(null);
	let loading = $state(true);

	const id = $derived(Number($page.params.id));

	$effect(() => {
		if (!id) return;
		load(id);
	});

	async function load(personId: number) {
		person = await peopleApi.get(personId);
		loading = false;
	}
</script>

<div class="page">
	<h1>Редактирование деятеля</h1>
	{#if loading}
		<div class="loading">Загрузка…</div>
	{:else if person}
		<PersonForm {person} />
	{:else}
		<div class="loading">Деятель не найден</div>
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
