<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { tick } from 'svelte';
	import { peopleApi, coverApi } from '$lib/api';
	import { displayName, genderLabel, roleLabel } from '$lib/constants';
	import type { Person } from '../../../app.d';
	import ArrowLeft from 'lucide-svelte/icons/arrow-left';
	import Trash2 from 'lucide-svelte/icons/trash-2';
	import Edit from 'lucide-svelte/icons/pencil';
	import X from 'lucide-svelte/icons/x';
	import ChevronLeft from 'lucide-svelte/icons/chevron-left';
	import ChevronRight from 'lucide-svelte/icons/chevron-right';

	let person = $state<Person | null>(null);
	let loading = $state(true);
	let mainImgUrl = $state('');
	let thumbUrls = $state<Record<string, string>>({});
	let titleCovers = $state<Record<number, string>>({});
	let lightboxIdx = $state(-1);
	let stripEl = $state<HTMLElement | null>(null);
	let stripCanLeft = $state(false);
	let stripCanRight = $state(false);
	let loadedId = 0;

	const id = $derived(Number($page.params.id));

	$effect(() => {
		if (id === loadedId) return;
		loadedId = id;
		load();
	});

	async function load() {
		loading = true;
		person = await peopleApi.get(id);
		thumbUrls = {};
		titleCovers = {};
		if (person) {
			if (person.mainImage) {
				mainImgUrl = await coverApi.dataURL(person.mainImage);
			}
			for (const f of person.images || []) {
				await ensureThumb(f);
			}
			for (const t of person.titles || []) {
				if (t.cover && titleCovers[t.id] === undefined) {
					titleCovers[t.id] = await coverApi.dataURL(t.cover);
				}
			}
		}
		loading = false;
		await tick();
		updateStripArrows();
	}

	async function ensureThumb(file: string) {
		if (thumbUrls[file] !== undefined) return;
		thumbUrls[file] = '';
		thumbUrls[file] = await coverApi.dataURL(file);
	}

	function lbStep(dir: number) {
		if (!person || person.images.length === 0) return;
		lightboxIdx = (lightboxIdx + dir + person.images.length) % person.images.length;
	}

	function updateStripArrows() {
		const el = stripEl;
		if (!el) return;
		stripCanLeft = el.scrollLeft > 4;
		stripCanRight = el.scrollLeft + el.clientWidth < el.scrollWidth - 4;
	}

	function scrollStrip(dir: number) {
		const el = stripEl;
		if (!el) return;
		el.scrollBy({ left: dir * el.clientWidth * 0.8, behavior: 'smooth' });
	}

	async function deletePerson() {
		if (!person) return;
		if (!confirm('Удалить этого деятеля?')) return;
		await peopleApi.delete(id);
		goto('/people');
	}

	function titlesLabel(n: number): string {
		const m10 = n % 10;
		const m100 = n % 100;
		const word = m10 === 1 && m100 !== 11 ? 'тайтле' : 'тайтлах';
		return `Встречается в ${n} ${word}`;
	}
</script>

<svelte:window
	onkeydown={(e) => {
		if (e.key === 'Escape') {
			lightboxIdx = -1;
		} else if (e.key === 'ArrowLeft' && lightboxIdx >= 0) {
			lbStep(-1);
		} else if (e.key === 'ArrowRight' && lightboxIdx >= 0) {
			lbStep(1);
		}
	}}
/>

{#if loading}
	<div class="loading">Загрузка…</div>
{:else if !person}
	<div class="loading">Деятель не найден</div>
{:else}
	<div class="page">
		<div class="back">
			<button class="btn" onclick={() => history.back()}>
				<ArrowLeft size={16} />
				Назад
			</button>
			<button class="btn" onclick={() => goto(`/person/${id}/edit`)}>
				<Edit size={14} />
				Редактировать
			</button>
			<button class="btn btn-danger" onclick={deletePerson}><Trash2 size={14} /> Удалить</button>
		</div>

		<div class="detail-header">
			<div class="poster">
				{#if mainImgUrl}
					<img src={mainImgUrl} alt="" />
				{:else}
					<div class="poster-placeholder">{displayName(person.names)[0]}</div>
				{/if}
			</div>
			<div class="header-info">
				<div class="type-badge">Деятель · {roleLabel(person.role)}</div>
				<h1 class="detail-title">{person.names.map((n) => n.value).join(' / ')}</h1>
				{#if person.age}
					<div class="info-line">Возраст: {person.age}</div>
				{/if}
				{#if person.birthDate}
					<div class="info-line">Дата рождения: {person.birthDate}</div>
				{/if}
				{#if person.deathDate}
					<div class="info-line">Дата смерти: {person.deathDate}</div>
				{/if}
				{#if person.gender}
					<div class="info-line">Пол: {genderLabel(person.gender)}</div>
				{/if}
				{#if person.titles.length > 0}
					<div class="in-titles">{titlesLabel(person.titles.length)}</div>
				{/if}
			</div>
		</div>

		{#if person.description}
			<section class="section">
				<h3>Биография</h3>
				<p class="synopsis">{person.description}</p>
			</section>
		{/if}

		{#if person.titles.length > 0}
			<section class="section">
				<h3>Тайтлы</h3>
				<div class="cards">
					{#each person.titles as t (t.id)}
						<button class="rel-card" onclick={() => goto(`/title/${t.id}`)}>
							<span class="rel-poster">
								{#if titleCovers[t.id]}
									<img src={titleCovers[t.id]} alt="" />
								{:else}
									<span class="rel-ph">{(t.name || '?').slice(0, 1).toUpperCase()}</span>
								{/if}
								<span class="rel-name">{t.name || 'Без названия'}</span>
							</span>
						</button>
					{/each}
				</div>
			</section>
		{/if}

		{#if person.images.length > 0}
			<section class="section">
				<h3>Дополнительные изображения</h3>
				<div class="strip">
					<button class="strip-arrow" class:off={!stripCanLeft} aria-label="Прокрутить влево" onclick={() => scrollStrip(-1)}>
						<ChevronLeft size={18} />
					</button>
					<div class="strip-row" bind:this={stripEl} onscroll={updateStripArrows}>
						{#each person.images as f, i}
							<button class="strip-thumb" onclick={() => (lightboxIdx = i)}>
								<img src={thumbUrls[f]} alt="" />
							</button>
						{/each}
					</div>
					<button class="strip-arrow" class:off={!stripCanRight} aria-label="Прокрутить вправо" onclick={() => scrollStrip(1)}>
						<ChevronRight size={18} />
					</button>
				</div>
			</section>
		{/if}

		{#if lightboxIdx >= 0 && person}
			<div class="lightbox">
				<button class="lb-backdrop" aria-label="Закрыть" onclick={() => (lightboxIdx = -1)}></button>
				<img class="lb-img" src={thumbUrls[person.images[lightboxIdx]]} alt="" />
				{#if person.images.length > 1}
					<button class="lb-arrow lb-prev" aria-label="Предыдущее" onclick={() => lbStep(-1)}><ChevronLeft size={26} /></button>
					<button class="lb-arrow lb-next" aria-label="Следующее" onclick={() => lbStep(1)}><ChevronRight size={26} /></button>
					<span class="lb-counter">{lightboxIdx + 1} / {person.images.length}</span>
				{/if}
				<button class="lb-close" aria-label="Закрыть" onclick={() => (lightboxIdx = -1)}><X size={18} /></button>
			</div>
		{/if}
	</div>
{/if}

<style>
	.page {
		max-width: 960px;
		margin: 0 auto;
	}
	.loading {
		text-align: center;
		padding: 60px;
		color: var(--text-dim);
	}
	.back {
		display: flex;
		gap: 8px;
		margin-bottom: 20px;
		flex-wrap: wrap;
	}

	.detail-header {
		display: flex;
		gap: 24px;
		margin-bottom: 24px;
	}
	.poster {
		width: 200px;
		flex-shrink: 0;
	}
	.poster img {
		width: 100%;
		border-radius: var(--radius);
		box-shadow: var(--shadow);
	}
	.poster-placeholder {
		width: 100%;
		aspect-ratio: 2 / 3;
		border-radius: var(--radius);
		background: var(--bg-elev-2);
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 64px;
		font-weight: 800;
		color: var(--accent);
	}
	.header-info {
		flex: 1;
		min-width: 0;
	}
	.type-badge {
		font-size: 12px;
		color: var(--text-dim);
		margin-bottom: 4px;
		text-transform: uppercase;
		letter-spacing: 0.04em;
	}
	.detail-title {
		font-size: 24px;
		font-weight: 700;
		line-height: 1.3;
		margin-bottom: 6px;
	}
	.info-line {
		font-size: 14px;
		color: var(--text-dim);
		margin-bottom: 6px;
	}
	.in-titles {
		font-size: 13px;
		color: var(--text-faint);
	}

	.section {
		margin-bottom: 22px;
	}
	.section h3 {
		font-size: 16px;
		font-weight: 600;
		margin-bottom: 10px;
	}
	.synopsis {
		color: var(--text-dim);
		line-height: 1.6;
		font-size: 14px;
	}

	.cards {
		display: flex;
		flex-wrap: wrap;
		gap: 14px;
	}
	.rel-card {
		padding: 0;
		background: none;
		border: none;
		cursor: pointer;
	}
	.rel-poster {
		position: relative;
		display: flex;
		align-items: center;
		justify-content: center;
		width: 118px;
		aspect-ratio: 2 / 3;
		border-radius: var(--radius-sm);
		overflow: hidden;
		border: 1px solid var(--border);
		background: var(--bg-elev-2);
	}
	.rel-poster img {
		width: 100%;
		height: 100%;
		object-fit: cover;
		display: block;
		transition: transform 0.15s;
	}
	.rel-card:hover .rel-poster img {
		transform: scale(1.05);
	}
	.rel-ph {
		font-size: 40px;
		font-weight: 800;
		color: var(--accent);
	}
	.rel-name {
		position: absolute;
		left: 0;
		right: 0;
		bottom: 0;
		padding: 22px 8px 8px;
		font-size: 12px;
		font-weight: 600;
		line-height: 1.25;
		color: #fff;
		text-align: center;
		background: linear-gradient(to top, rgba(0, 0, 0, 0.8), rgba(0, 0, 0, 0));
		pointer-events: none;
	}

	.strip {
		display: flex;
		align-items: center;
		gap: 8px;
	}
	.strip-arrow {
		flex: none;
		width: 34px;
		height: 34px;
		border-radius: 50%;
		border: 1px solid var(--border);
		background: var(--bg-elev);
		color: var(--text-dim);
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
	}
	.strip-arrow:hover {
		color: var(--text);
		background: var(--bg-elev-2);
	}
	.strip-arrow.off {
		visibility: hidden;
	}
	.strip-row {
		display: flex;
		gap: 8px;
		overflow-x: auto;
		scrollbar-width: none;
	}
	.strip-row::-webkit-scrollbar {
		display: none;
	}
	.strip-thumb {
		flex: none;
		padding: 0;
		border: 1px solid var(--border);
		border-radius: var(--radius-sm);
		overflow: hidden;
		cursor: zoom-in;
		background: var(--bg-elev);
	}
	.strip-thumb img {
		display: block;
		width: 90px;
		height: 120px;
		object-fit: cover;
	}
	.lightbox {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.88);
		z-index: 100;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	.lb-backdrop {
		position: absolute;
		inset: 0;
		background: transparent;
		border: none;
		padding: 0;
		cursor: zoom-out;
	}
	.lb-img {
		position: relative;
		max-width: 88vw;
		max-height: 86vh;
		border-radius: var(--radius);
		box-shadow: var(--shadow);
		pointer-events: none;
	}
	.lb-arrow {
		position: absolute;
		top: 50%;
		transform: translateY(-50%);
		width: 46px;
		height: 46px;
		border-radius: 50%;
		border: none;
		background: rgba(0, 0, 0, 0.55);
		color: #fff;
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
	}
	.lb-arrow:hover {
		background: rgba(0, 0, 0, 0.8);
	}
	.lb-prev {
		left: 18px;
	}
	.lb-next {
		right: 18px;
	}
	.lb-close {
		position: absolute;
		top: 16px;
		right: 16px;
		width: 38px;
		height: 38px;
		border-radius: 50%;
		border: none;
		background: rgba(0, 0, 0, 0.55);
		color: #fff;
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
	}
	.lb-counter {
		position: absolute;
		bottom: 18px;
		left: 50%;
		transform: translateX(-50%);
		background: rgba(0, 0, 0, 0.55);
		color: #fff;
		font-size: 13px;
		padding: 4px 12px;
		border-radius: 99px;
	}

	@media (max-width: 640px) {
		.detail-header {
			flex-direction: column;
			align-items: center;
			text-align: center;
		}
		.poster {
			width: 160px;
		}
	}
</style>
