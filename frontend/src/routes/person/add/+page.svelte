<script lang="ts">
	import PersonForm from '$lib/components/PersonForm.svelte';
</script>

<div class="page">
	<h1>Добавить деятеля</h1>
	<PersonForm />
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
