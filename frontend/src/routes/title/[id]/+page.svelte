<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import RichEditor from '$lib/components/RichEditor.svelte';
	import { titlesApi, coverApi, shelvesApi } from '$lib/api';
	import {
		STATUSES,
		RELEASE_STATUSES,
		READ_CATEGORIES,
		WATCH_CATEGORIES,
		TYPES,
		NAME_KINDS,
		CREATOR_ROLES,
		READ_PROGRESS,
		WATCH_PROGRESS,
		statusLabel,
		statusColor,
		releaseStatusLabel,
		releaseStatusColor,
		categoryLabel,
		typeLabel,
		progressForType
	} from '$lib/constants';
	import type { Title, Shelf } from '../../../app.d';
	import ArrowLeft from 'lucide-svelte/icons/arrow-left';
	import Trash2 from 'lucide-svelte/icons/trash-2';
	import Plus from 'lucide-svelte/icons/plus';
	import Minus from 'lucide-svelte/icons/minus';
	import Edit from 'lucide-svelte/icons/pencil';

	let title = $state<Title | null>(null);
	let loading = $state(true);
	let editing = $state(false);
	let coverUrl = $state('');
	let notesHtml = $state('');

	let editTitle = $state<any>(null);
	let shelves = $state<Shelf[]>([]);
	let editShelfId = $state(0);
	let editShelfInitial = $state(0);
	let statusMenu = $state('');
	let notesTimer: any;

	const id = $derived(Number($page.params.id));
	const editShelfOptions = $derived(shelves.filter((s) => s.kind === (title?.type ?? '')));

	$effect(() => {
		if (!editTitle) return;
		if (editTitle.type === 'read' && !READ_CATEGORIES.find((c) => c.id === editTitle.category)) {
			editTitle.category = 'book';
		}
		if (editTitle.type === 'watch' && !WATCH_CATEGORIES.find((c) => c.id === editTitle.category)) {
			editTitle.category = 'movie';
		}
	});

	onMount(async () => {
		await load();
	});

	async function load() {
		loading = true;
		title = await titlesApi.get(id);
		if (title) {
			notesHtml = title.notes || '';
			if (title.cover) {
				coverUrl = await coverApi.dataURL(title.cover);
			}
		}
		loading = false;
	}

	async function adjustProgress(field: string, delta: number) {
		title = await titlesApi.adjustProgress(id, field, delta);
	}

	async function saveStatuses() {
		if (!title) return;
		await titlesApi.save({ ...title });
	}

	function pickStatus(field: 'status' | 'releaseStatus', value: string) {
		if (!title) return;
		title[field] = value;
		statusMenu = '';
		saveStatuses();
	}

	async function saveNotes(html: string) {
		notesHtml = html;
		clearTimeout(notesTimer);
		notesTimer = setTimeout(async () => {
			if (!title) return;
			title.notes = html;
			await titlesApi.save({ ...title });
		}, 500);
	}

	async function startEdit() {
		editTitle = JSON.parse(JSON.stringify(title));
		editTitle.names = [...editTitle.names];
		editTitle.creators = [...editTitle.creators];
		editTitle.genres = [...editTitle.genres];
		editTitle.tags = [...editTitle.tags];
		shelves = (await shelvesApi.list()) || [];
		editShelfId = shelves.find((s) => s.titleIds.includes(id))?.id ?? 0;
		editShelfInitial = editShelfId;
		editing = true;
	}

	async function saveEdit() {
		if (!editTitle) return;
		await titlesApi.save(editTitle);
		if (editShelfId !== editShelfInitial) {
			const current = shelves.filter((s) => s.titleIds.includes(id));
			for (const s of current) {
				if (s.id !== editShelfId) {
					await shelvesApi.setItems(s.id, s.titleIds.filter((x) => x !== id));
				}
			}
			if (editShelfId !== 0) {
				const target = shelves.find((s) => s.id === editShelfId);
				if (target && !target.titleIds.includes(id)) {
					await shelvesApi.setItems(target.id, [...target.titleIds, id]);
				}
			}
		}
		editing = false;
		await load();
	}

	async function deleteTitle() {
		if (!title) return;
		if (!confirm('Удалить этот тайтл?')) return;
		await titlesApi.delete(id);
		goto('/');
	}

	function addName() {
		editTitle.names.push({ kind: 'original', value: '' });
	}
	function removeName(i: number) {
		editTitle.names.splice(i, 1);
	}
	function addCreator() {
		editTitle.creators.push({ role: 'author', name: '' });
	}
	function removeCreator(i: number) {
		editTitle.creators.splice(i, 1);
	}
	function addGenre() {
		editTitle.genres.push('');
	}
	function removeGenre(i: number) {
		editTitle.genres.splice(i, 1);
	}
	function addTag() {
		editTitle.tags.push('');
	}
	function removeTag(i: number) {
		editTitle.tags.splice(i, 1);
	}

	function onCoverInput(e: Event) {
		const input = e.target as HTMLInputElement;
		const file = input.files?.[0];
		if (!file) return;
		const reader = new FileReader();
		reader.onload = async () => {
			const dataUrl = reader.result as string;
			editTitle.cover = await coverApi.uploadDataURL(dataUrl);
			coverUrl = await coverApi.dataURL(editTitle.cover);
		};
		reader.readAsDataURL(file);
	}
</script>

<svelte:window
	onkeydown={(e) => e.key === 'Escape' && (statusMenu = '')}
	onclick={() => (statusMenu = '')}
/>

{#if loading}
	<div class="loading">Загрузка…</div>
{:else if !title}
	<div class="loading">Тайтл не найден</div>
{:else}
	<div class="page">
		<div class="back">
			<button class="btn" onclick={() => history.back()}>
				<ArrowLeft size={16} />
				Назад
			</button>
			<button class="btn" onclick={startEdit}><Edit size={14} /> Редактировать</button>
			<button class="btn btn-danger" onclick={deleteTitle}><Trash2 size={14} /> Удалить</button>
		</div>

		{#if editing}
			{@const e = editTitle}
			<div class="edit-panel">
				<h2>Редактирование</h2>

				<div class="form-row">
					<div class="field">
						<span class="label">Тип</span>
						<select class="select" bind:value={e.type}>
							{#each TYPES as t}<option value={t.id}>{t.label}</option>{/each}
						</select>
					</div>
					<div class="field">
						<span class="label">Категория</span>
						<select class="select" bind:value={e.category}>
							{#if e.type === 'read'}
								{#each READ_CATEGORIES as c}<option value={c.id}>{c.label}</option>{/each}
							{:else}
								{#each WATCH_CATEGORIES as c}<option value={c.id}>{c.label}</option>{/each}
							{/if}
						</select>
					</div>
				</div>

				<div class="field">
					<span class="label">Названия</span>
					{#each e.names as n, i}
						<div class="row-2">
							<select class="select sm" bind:value={n.kind}>
								{#each NAME_KINDS as k}<option value={k.id}>{k.label}</option>{/each}
							</select>
							<input class="input" placeholder="Название…" bind:value={n.value} />
							<button class="btn btn-icon sm" onclick={() => removeName(i)}><Minus size={14} /></button>
						</div>
					{/each}
					<button class="btn sm" onclick={addName}><Plus size={14} /> Название</button>
				</div>

				<div class="field">
					<span class="label">Обложка</span>
					<div class="cover-preview">
						{#if coverUrl}<img src={coverUrl} alt="" />{/if}
						<input type="file" accept="image/*" onchange={onCoverInput} />
					</div>
				</div>

				<div class="field">
					<span class="label">Описание</span>
					<textarea class="textarea" rows={4} bind:value={e.synopsis}></textarea>
				</div>

				<div class="field">
					<span class="label">Создатели</span>
					{#each e.creators as c, i}
						<div class="row-3">
							<select class="select sm" bind:value={c.role}>
								{#each CREATOR_ROLES as r}<option value={r.id}>{r.label}</option>{/each}
							</select>
							<input class="input" placeholder="Имя…" bind:value={c.name} />
							<button class="btn btn-icon sm" onclick={() => removeCreator(i)}><Minus size={14} /></button>
						</div>
					{/each}
					<button class="btn sm" onclick={addCreator}><Plus size={14} /> Создатель</button>
				</div>

				<div class="field">
					<span class="label">Жанры</span>
					{#each e.genres as g, i}
						<div class="row-2">
							<input class="input" placeholder="Жанр…" bind:value={e.genres[i]} />
							<button class="btn btn-icon sm" onclick={() => removeGenre(i)}><Minus size={14} /></button>
						</div>
					{/each}
					<button class="btn sm" onclick={addGenre}><Plus size={14} /> Жанр</button>
				</div>

				<div class="field">
					<span class="label">Теги</span>
					{#each e.tags as t, i}
						<div class="row-2">
							<input class="input" placeholder="Тег…" bind:value={e.tags[i]} />
							<button class="btn btn-icon sm" onclick={() => removeTag(i)}><Minus size={14} /></button>
						</div>
					{/each}
					<button class="btn sm" onclick={addTag}><Plus size={14} /> Тег</button>
				</div>

				<div class="form-row">
					<div class="field">
						<span class="label">Оценка (0–5)</span>
						<input type="number" class="input" min="0" max="5" step="0.5" bind:value={e.score} />
					</div>
					<div class="field">
						<span class="label">Мой статус</span>
						<select class="select" bind:value={e.status}>
							{#each STATUSES as s}<option value={s.id}>{s.label}</option>{/each}
						</select>
					</div>
				</div>

				<div class="field">
					<span class="label">Статус тайтла</span>
					<select class="select" bind:value={e.releaseStatus}>
						<option value="">Не указан</option>
						{#each RELEASE_STATUSES as s}<option value={s.id}>{s.label}</option>{/each}
					</select>
				</div>

				{#if editShelfOptions.length > 0}
					<div class="field">
						<span class="label">Полка</span>
						<select class="select" bind:value={editShelfId}>
							<option value={0}>Без полки</option>
							{#each editShelfOptions as s}<option value={s.id}>{s.name}</option>{/each}
						</select>
					</div>
				{/if}

				{#if e.customList}
					<div class="field">
						<span class="label">Свой список</span>
						<input class="input" bind:value={e.customList} />
					</div>
				{/if}

				<div class="field">
					<span class="label">Прогресс</span>
					<div class="progress-fields">
						{#each progressForType(e.type) as f}
							<div class="prog-item">
								<span class="prog-label">{f.label}</span>
								<div class="prog-row">
									<button class="btn btn-icon sm" onclick={() => (e.progress[f.id] = Math.max(0, e.progress[f.id] - 1))}><Minus size={14} /></button>
									<input type="number" class="input sm" min="0" bind:value={e.progress[f.id]} />
									<button class="btn btn-icon sm" onclick={() => (e.progress[f.id] = e.progress[f.id] + 1)}><Plus size={14} /></button>
								</div>
							</div>
						{/each}
						{#if e.type === 'read'}
							<div class="prog-item">
								<span class="prog-label">Вышло глав</span>
								<input type="number" class="input sm" min="0" bind:value={e.progress.totalChapters} />
							</div>
						{:else}
							<div class="prog-item">
								<span class="prog-label">Вышло серий</span>
								<input type="number" class="input sm" min="0" bind:value={e.progress.totalEpisodes} />
							</div>
						{/if}
					</div>
				</div>

				<div class="edit-actions">
					<button class="btn" onclick={() => (editing = false)}>Отмена</button>
					<button class="btn btn-primary" onclick={saveEdit}>Сохранить</button>
				</div>
			</div>
		{:else}
			<div class="detail-header">
				<div class="poster">
					{#if coverUrl}
						<img src={coverUrl} alt="" />
					{:else}
						<div class="poster-placeholder">{title.names.length > 0 ? title.names[0].value[0] : '?'}</div>
					{/if}
				</div>
				<div class="header-info">
					<div class="type-badge">{typeLabel(title.type)} · {categoryLabel(title.category)}</div>
					<h1 class="detail-title">{title.names.map((n: any) => n.value).join(' / ')}</h1>
					{#if title.creators.length > 0}
						<div class="creators">
							{title.creators.map((c: any) => `${c.name} (${c.role})`).join(', ')}
						</div>
					{/if}
					<div class="badges">
						<div class="badge-menu">
							<button
								class="badge"
								style="background:{statusColor(title.status)}22;color:{statusColor(title.status)}"
								onclick={(e) => {
									e.stopPropagation();
									statusMenu = statusMenu === 'status' ? '' : 'status';
								}}
							>
								{statusLabel(title.status)}
							</button>
							{#if statusMenu === 'status'}
								<div class="menu">
									{#each STATUSES as s}
										<button
											class="menu-item"
											class:active={title.status === s.id}
											onclick={(e) => {
												e.stopPropagation();
												pickStatus('status', s.id);
											}}
										>
											<span class="dot" style="background:{s.color}"></span>
											{s.label}
										</button>
									{/each}
								</div>
							{/if}
						</div>
						<div class="badge-menu">
								<button
									class="badge"
									style="background:{releaseStatusColor(title.releaseStatus)}22;color:{releaseStatusColor(title.releaseStatus)}"
									onclick={(e) => {
										e.stopPropagation();
										statusMenu = statusMenu === 'release' ? '' : 'release';
									}}
								>
									{releaseStatusLabel(title.releaseStatus) || 'Не указан'}
								</button>
								{#if statusMenu === 'release'}
									<div class="menu">
										<button
											class="menu-item"
											class:active={title.releaseStatus === ''}
											onclick={(e) => {
												e.stopPropagation();
												pickStatus('releaseStatus', '');
											}}
										>
											<span class="dot" style="background:{releaseStatusColor('')}"></span>
											Не указан
										</button>
										{#each RELEASE_STATUSES as s}
											<button
												class="menu-item"
												class:active={title.releaseStatus === s.id}
												onclick={(e) => {
													e.stopPropagation();
													pickStatus('releaseStatus', s.id);
												}}
											>
												<span class="dot" style="background:{s.color}"></span>
												{s.label}
											</button>
										{/each}
									</div>
								{/if}
							</div>
						{#if title.score > 0}
							<span class="badge score">★ {title.score.toFixed(1)}</span>
						{/if}
					</div>
					{#if title.genres.length > 0}
						<div class="genre-list">
							{#each title.genres as g}<span class="badge">{g}</span>{/each}
						</div>
					{/if}
					{#if title.tags.length > 0}
						<div class="tag-list">
							{#each title.tags as t}<span class="badge dim">#{t}</span>{/each}
						</div>
					{/if}
				</div>
			</div>

			{#if title.synopsis}
				<section class="section">
					<h3>Описание</h3>
					<p class="synopsis">{title.synopsis}</p>
				</section>
			{/if}

				<section class="section">
					<h3>Прогресс</h3>
					<div class="progress-grid">
						{#each progressForType(title.type) as f}
							<div class="prog-block">
								<span class="prog-name">{f.label}</span>
								<div class="prog-controls">
									<button class="btn btn-icon" onclick={() => adjustProgress(f.id, -1)}><Minus size={18} /></button>
									<span class="prog-val">{(title.progress as any)[f.id]}</span>
									<button class="btn btn-icon" onclick={() => adjustProgress(f.id, 1)}><Plus size={18} /></button>
								</div>
							</div>
						{/each}
						{#if title.type === 'read'}
							<div class="prog-block">
								<span class="prog-name">Вышло глав</span>
								<div class="prog-controls">
									<button class="btn btn-icon" onclick={() => adjustProgress('totalChapters', -1)}><Minus size={18} /></button>
									<span class="prog-val">{title.progress.totalChapters}</span>
									<button class="btn btn-icon" onclick={() => adjustProgress('totalChapters', 1)}><Plus size={18} /></button>
								</div>
							</div>
						{:else}
							<div class="prog-block">
								<span class="prog-name">Вышло серий</span>
								<div class="prog-controls">
									<button class="btn btn-icon" onclick={() => adjustProgress('totalEpisodes', -1)}><Minus size={18} /></button>
									<span class="prog-val">{title.progress.totalEpisodes}</span>
									<button class="btn btn-icon" onclick={() => adjustProgress('totalEpisodes', 1)}><Plus size={18} /></button>
								</div>
							</div>
						{/if}
					</div>
				</section>

			<section class="section">
				<h3>Заметки</h3>
				<RichEditor content={notesHtml} onUpdate={saveNotes} />
			</section>
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
	.creators {
		font-size: 13px;
		color: var(--text-dim);
		margin-bottom: 10px;
	}
	.badges {
		display: flex;
		gap: 6px;
		margin-bottom: 10px;
		flex-wrap: wrap;
	}
	.badge.score {
		background: rgba(245, 166, 35, 0.15);
		color: #f5a623;
	}
	.badge.dim {
		background: var(--bg-elev-2);
		color: var(--text-dim);
	}
	.badge-menu {
		position: relative;
	}
	.menu {
		position: absolute;
		top: calc(100% + 6px);
		left: 0;
		z-index: 20;
		min-width: 170px;
		background: var(--bg-elev);
		border: 1px solid var(--border);
		border-radius: var(--radius-sm);
		box-shadow: var(--shadow);
		padding: 4px;
		display: flex;
		flex-direction: column;
		gap: 2px;
	}
	.menu-item {
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 6px 10px;
		border-radius: var(--radius-sm);
		font-size: 12px;
		text-align: left;
		color: var(--text);
	}
	.menu-item:hover {
		background: var(--bg-elev-2);
	}
	.menu-item.active {
		color: var(--accent);
		font-weight: 600;
	}
	.dot {
		width: 8px;
		height: 8px;
		border-radius: 50%;
		flex: none;
	}
	.genre-list,
	.tag-list {
		display: flex;
		gap: 4px;
		flex-wrap: wrap;
		margin-bottom: 4px;
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

	.progress-grid {
		display: flex;
		gap: 16px;
		flex-wrap: wrap;
	}
	.prog-block {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 6px;
	}
	.prog-name {
		font-size: 12px;
		color: var(--text-dim);
		text-transform: uppercase;
	}
	.prog-controls {
		display: flex;
		align-items: center;
		gap: 10px;
	}
	.prog-val {
		font-size: 22px;
		font-weight: 700;
		min-width: 40px;
		text-align: center;
	}

	.edit-panel {
		background: var(--bg-elev);
		border-radius: var(--radius);
		padding: 20px 24px;
	}
	.edit-panel h2 {
		margin-bottom: 16px;
	}
	.form-row {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 14px;
	}
	.row-2 {
		display: flex;
		gap: 6px;
		margin-bottom: 6px;
	}
	.row-3 {
		display: flex;
		gap: 6px;
		margin-bottom: 6px;
	}
	.row-2 .input,
	.row-3 .input {
		flex: 1;
	}
	.sm {
		width: auto;
		min-width: 100px;
		padding: 6px 10px;
		font-size: 12px;
	}
	.btn.sm {
		padding: 5px 10px;
		font-size: 12px;
	}
	.btn-icon.sm {
		width: 30px;
		height: 30px;
		min-width: 30px;
	}
	.cover-preview {
		display: flex;
		gap: 12px;
		align-items: center;
	}
	.cover-preview img {
		height: 80px;
		border-radius: 8px;
	}
	.progress-fields {
		display: flex;
		gap: 20px;
		flex-wrap: wrap;
	}
	.prog-item {
		display: flex;
		flex-direction: column;
		gap: 4px;
	}
	.prog-label {
		font-size: 11px;
		color: var(--text-dim);
		text-transform: uppercase;
	}
	.prog-row {
		display: flex;
		align-items: center;
		gap: 4px;
	}
	.input.sm {
		width: 70px;
		text-align: center;
	}
	.edit-actions {
		display: flex;
		justify-content: flex-end;
		gap: 8px;
		margin-top: 18px;
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
		.form-row {
			grid-template-columns: 1fr;
		}
	}
</style>
