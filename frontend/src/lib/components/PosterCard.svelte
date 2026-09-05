<script lang="ts">
	import { coverApi } from '../api';
	import { onMount } from 'svelte';

	let {
		href,
		image = '',
		title,
		subtitle = ''
	}: { href: string; image?: string; title: string; subtitle?: string } = $props();

	let imgUrl = $state('');
	let loaded = $state(false);

	onMount(async () => {
		if (image) {
			imgUrl = await coverApi.dataURL(image);
		}
		loaded = true;
	});
</script>

<a {href} class="card">
	<div class="poster">
		{#if imgUrl}
			<img src={imgUrl} alt="" />
		{:else if loaded}
			<div class="poster-placeholder">{title[0]}</div>
		{/if}
	</div>
	<div class="info">
		<div class="name">{title}</div>
		{#if subtitle}
			<div class="subtitle">{subtitle}</div>
		{/if}
	</div>
</a>

<style>
	.card {
		display: flex;
		flex-direction: column;
		text-decoration: none;
		color: inherit;
		transition: transform 0.12s;
	}
	.card:hover {
		transform: translateY(-3px);
	}
	.poster {
		width: 100%;
		aspect-ratio: 2 / 3;
		border-radius: var(--radius);
		overflow: hidden;
		background: var(--bg-elev);
		box-shadow: var(--shadow);
	}
	.poster img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}
	.poster-placeholder {
		width: 100%;
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 42px;
		font-weight: 800;
		background: linear-gradient(135deg, var(--accent-soft), var(--bg-elev-2));
		color: var(--accent);
	}
	.info {
		padding: 8px 2px 0;
	}
	.name {
		font-size: 13px;
		font-weight: 600;
		overflow: hidden;
		text-overflow: ellipsis;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		line-height: 1.3;
	}
	.subtitle {
		font-size: 11px;
		color: var(--text-dim);
		margin-top: 3px;
	}
</style>
