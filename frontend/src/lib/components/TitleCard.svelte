<script lang="ts">
	import type { Title } from '../../app.d';
	import { displayName, statusLabel, statusColor, progressForType, releaseStatusLabel, releaseStatusColor } from '../constants';
	import { coverApi } from '../api';
	import { onMount } from 'svelte';

	let { title, view = 'grid' }: { title: Title; view?: 'grid' | 'list' } = $props();

	let coverUrl = $state('');
	let loaded = $state(false);

	onMount(async () => {
		if (title.cover) {
			coverUrl = await coverApi.dataURL(title.cover);
		}
		loaded = true;
	});

	const progressText = (t: Title): string => {
		const fields = progressForType(t.type);
		const parts: string[] = [];
		for (const f of fields) {
			const v = (t.progress as any)[f.id] as number;
			if (v > 0) parts.push(`${v} ${f.label.toLowerCase()}`);
		}
		return parts.slice(0, 2).join(' · ');
	};
</script>

<a href="/title/{title.id}" class="card {view}">
	{#if view === 'grid'}
		<div class="poster">
			{#if coverUrl}
				<img src={coverUrl} alt="" />
			{:else if loaded}
				<div class="poster-placeholder">{displayName(title.names)[0]}</div>
			{/if}
			{#if title.releaseStatus}
				<span class="release-badge" style="background:{releaseStatusColor(title.releaseStatus)}e6">
					{releaseStatusLabel(title.releaseStatus)}
				</span>
			{/if}
			<span class="status-dot" style="background:{statusColor(title.status)}"></span>
		</div>
		<div class="info">
			<div class="title-text">{displayName(title.names)}</div>
			<div class="meta">
				{#if title.score > 0}
					<span class="score">★ {title.score.toFixed(1)}</span>
				{/if}
				{#if progressText(title)}
					<span class="prog">{progressText(title)}</span>
				{/if}
			</div>
		</div>
	{:else}
		<div class="poster-mini">
			{#if coverUrl}
				<img src={coverUrl} alt="" />
			{:else if loaded}
				<div class="poster-placeholder sm">{displayName(title.names)[0]}</div>
			{/if}
		</div>
		<div class="list-info">
			<div class="title-text">{displayName(title.names)}</div>
			<div class="meta">
				<span class="badge-sm" style="background:{statusColor(title.status)}22;color:{statusColor(title.status)}">
					{statusLabel(title.status)}
				</span>
				{#if title.releaseStatus}
					<span class="badge-sm" style="background:{releaseStatusColor(title.releaseStatus)}22;color:{releaseStatusColor(title.releaseStatus)}">
						{releaseStatusLabel(title.releaseStatus)}
					</span>
				{/if}
				{#if title.score > 0}<span class="score">★ {title.score.toFixed(1)}</span>{/if}
				{#if progressText(title)}<span class="prog">{progressText(title)}</span>{/if}
			</div>
		</div>
	{/if}
</a>

<style>
	.card {
		display: flex;
		flex-direction: column;
		text-decoration: none;
		color: inherit;
		transition: transform 0.12s;
	}
	.card.grid:hover {
		transform: translateY(-3px);
	}

	.poster {
		width: 100%;
		aspect-ratio: 2 / 3;
		border-radius: var(--radius);
		overflow: hidden;
		background: var(--bg-elev);
		box-shadow: var(--shadow);
		position: relative;
	}
	.poster img,
	.poster-mini img {
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
	.poster-placeholder.sm {
		font-size: 20px;
	}
	.status-dot {
		position: absolute;
		top: 8px;
		right: 8px;
		width: 12px;
		height: 12px;
		border-radius: 50%;
		border: 2px solid var(--bg);
	}
	.release-badge {
		position: absolute;
		top: 8px;
		left: 8px;
		padding: 2px 8px;
		border-radius: 99px;
		font-size: 10px;
		font-weight: 700;
		color: #fff;
		text-shadow: 0 1px 2px rgba(0, 0, 0, 0.35);
		box-shadow: 0 1px 4px rgba(0, 0, 0, 0.35);
	}

	.info {
		padding: 8px 2px 0;
	}
	.title-text {
		font-size: 13px;
		font-weight: 600;
		overflow: hidden;
		text-overflow: ellipsis;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		line-height: 1.3;
	}
	.meta {
		display: flex;
		gap: 8px;
		margin-top: 4px;
		font-size: 11px;
		color: var(--text-dim);
		flex-wrap: wrap;
	}
	.score {
		color: #f5a623;
		font-weight: 600;
	}

	.card.list {
		flex-direction: row;
		align-items: center;
		gap: 12px;
		padding: 10px;
		border-radius: var(--radius);
		background: var(--bg-elev);
		transition: background 0.12s;
	}
	.card.list:hover {
		background: var(--bg-elev-2);
	}
	.poster-mini {
		width: 40px;
		height: 56px;
		border-radius: 6px;
		overflow: hidden;
		flex-shrink: 0;
		background: var(--bg-elev-2);
	}
	.list-info {
		flex: 1;
		min-width: 0;
	}
	.badge-sm {
		padding: 1px 7px;
		border-radius: 99px;
		font-size: 10px;
		font-weight: 600;
	}
</style>
