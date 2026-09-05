<script lang="ts">
	import CharacterForm from '$lib/components/CharacterForm.svelte';
</script>

<div class="page">
	<h1>Добавить персонажа</h1>
	<CharacterForm />
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
</style>
