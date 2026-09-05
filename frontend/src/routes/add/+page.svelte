<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { titlesApi, coverApi, shelvesApi, charactersApi } from '$lib/api';
	import {
		TYPES,
		ALL_CATEGORIES,
		READ_CATEGORIES,
		WATCH_CATEGORIES,
		STATUSES,
		RELEASE_STATUSES,
		NAME_KINDS,
		CREATOR_ROLES,
		READ_PROGRESS,
		WATCH_PROGRESS,
		progressForType,
		displayName
	} from '$lib/constants';
	import type { Title, Shelf, Character } from '../../app.d';
	import Plus from 'lucide-svelte/icons/plus';
	import Minus from 'lucide-svelte/icons/minus';
	import ImagePlus from 'lucide-svelte/icons/image-plus';
	import Link from 'lucide-svelte/icons/link';
	import X from 'lucide-svelte/icons/x';
	import ChevronLeft from 'lucide-svelte/icons/chevron-left';
	import ChevronRight from 'lucide-svelte/icons/chevron-right';

	let type = $state('read');
	let category = $state('book');
	let names = $state([{ kind: 'original', value: '' }]);
	let cover = $state('');
	let coverUrl = $state('');
	let images = $state<string[]>([]);
	let imageUrls = $state<Record<string, string>>({});
	let imagesUploading = $state(false);
	let synopsis = $state('');
	let creators = $state<{ role: string; name: string }[]>([]);
	let genres = $state<string[]>([]);
	let tags = $state<string[]>([]);
	let score = $state(0);
	let status = $state('planned');
	let releaseStatus = $state('ongoing');
	let customList = $state('');
	let progress: Record<string, number> = $state({ volumes: 0, chapters: 0, pages: 0, seasons: 0, episodes: 0, minutes: 0 });
	let coverUrlInput = $state('');
	let coverUploading = $state(false);
	let saving = $state(false);
	let error = $state('');
	let shelves = $state<Shelf[]>([]);
	let shelfId = $state(0);
	let relations = $state<{ relatedId: number; label: string; reverseLabel: string }[]>([]);
	let allTitles = $state<Title[]>([]);
	let charRefs = $state<{ id: number }[]>([]);
	let allCharacters = $state<Character[]>([]);

	const typeShelves = $derived(shelves.filter((s) => s.kind === type));

	onMount(async () => {
		shelves = (await shelvesApi.list()) || [];
		allTitles = (await titlesApi.list()) || [];
		allCharacters = (await charactersApi.list('name')) || [];
	});

	$effect(() => {
		if (type === 'read' && !READ_CATEGORIES.find((c) => c.id === category)) {
			category = 'book';
		}
		if (type === 'watch' && !WATCH_CATEGORIES.find((c) => c.id === category)) {
			category = 'movie';
		}
		if (shelfId !== 0 && !typeShelves.find((s) => s.id === shelfId)) {
			shelfId = 0;
		}
	});

	async function onCoverFile(e: Event) {
		const input = e.target as HTMLInputElement;
		const file = input.files?.[0];
		if (!file) return;
		coverUploading = true;
		const reader = new FileReader();
		reader.onload = async () => {
			cover = await coverApi.uploadDataURL(reader.result as string);
			coverUrl = await coverApi.dataURL(cover);
			coverUploading = false;
		};
		reader.readAsDataURL(file);
	}

	async function onCoverURL() {
		if (!coverUrlInput.trim()) return;
		coverUploading = true;
		try {
			cover = await coverApi.uploadFromURL(coverUrlInput.trim());
			coverUrl = await coverApi.dataURL(cover);
			coverUrlInput = '';
		} catch (e) {
			error = 'Не удалось загрузить обложку по ссылке';
		}
		coverUploading = false;
	}

	async function onImagesInput(e: Event) {
		const input = e.target as HTMLInputElement;
		const files = Array.from(input.files ?? []);
		input.value = '';
		if (files.length === 0) return;
		imagesUploading = true;
		for (const file of files) {
			const dataUrl = await new Promise<string>((resolve) => {
				const reader = new FileReader();
				reader.onload = () => resolve(reader.result as string);
				reader.readAsDataURL(file);
			});
			const name = await coverApi.uploadDataURL(dataUrl);
			if (name && !images.includes(name)) {
				images.push(name);
				imageUrls[name] = dataUrl;
			}
		}
		imagesUploading = false;
	}

	function removeImage(i: number) {
		images.splice(i, 1);
	}

	function moveImage(i: number, dir: number) {
		const j = i + dir;
		if (j < 0 || j >= images.length) return;
		const f = images[i];
		images[i] = images[j];
		images[j] = f;
	}

	async function save() {
		const hasName = names.some((n) => n.value.trim() !== '');
		if (!hasName) {
			error = 'Укажите хотя бы одно название';
			return;
		}
		saving = true;
		error = '';
		try {
			const t: Title = {
				id: 0,
				type,
				category,
				names: names.filter((n) => n.value.trim()),
				cover,
				images: [...images],
				synopsis,
				creators: creators.filter((c) => c.name.trim()),
				genres: genres.filter((g) => g.trim()),
				tags: tags.filter((t) => t.trim()),
				relations: relations.filter((r) => r.relatedId > 0),
				characters: charRefs.filter((c) => c.id > 0),
				score,
				status,
				releaseStatus,
				customList,
				progress: {
					volumes: progress.volumes,
					chapters: progress.chapters,
					pages: progress.pages,
					seasons: progress.seasons,
					episodes: progress.episodes,
					minutes: progress.minutes,
					totalChapters: progress.totalChapters,
					totalEpisodes: progress.totalEpisodes
				},
				notes: '',
				spineColor: '',
				createdAt: '',
				updatedAt: ''
			};
			const id = await titlesApi.save(t);
			if (shelfId !== 0) {
				const sh = shelves.find((s) => s.id === shelfId);
				if (sh) await shelvesApi.setItems(shelfId, [...(sh.titleIds || []), id]);
			}
			goto('/title/' + id);
		} catch (e: any) {
			error = e.message || 'Ошибка сохранения';
		}
		saving = false;
	}

	function addName() { names.push({ kind: 'original', value: '' }); }
	function addCreator() { creators.push({ role: 'author', name: '' }); }
	function addGenre() { genres.push(''); }
	function addTag() { tags.push(''); }
	function addRelation() { relations.push({ relatedId: 0, label: '', reverseLabel: '' }); }

	function addCharRef() { charRefs.push({ id: 0 }); }
	function charCandidates(i: number) {
		const picked = charRefs.map((c, j) => (j === i ? 0 : c.id));
		return allCharacters.filter((c) => !picked.includes(c.id));
	}

	function relationCandidates(i: number) {
		const picked = relations.map((r, j) => (j === i ? 0 : r.relatedId));
		return allTitles.filter((t) => !picked.includes(t.id));
	}
</script>

<div class="page">
	<h1>Добавить тайтл</h1>

	{#if error}
		<div class="error">{error}</div>
	{/if}

	<div class="form">
		<div class="form-row">
			<div class="field">
				<span class="label">Тип</span>
				<select class="select" bind:value={type}>
					{#each TYPES as t}<option value={t.id}>{t.label}</option>{/each}
				</select>
			</div>
			<div class="field">
				<span class="label">Категория</span>
				<select class="select" bind:value={category}>
					{#if type === 'read'}
						{#each READ_CATEGORIES as c}<option value={c.id}>{c.label}</option>{/each}
					{:else}
						{#each WATCH_CATEGORIES as c}<option value={c.id}>{c.label}</option>{/each}
					{/if}
				</select>
			</div>
		</div>

		<div class="field">
			<span class="label">Названия</span>
			{#each names as n, i}
				<div class="row-2">
					<select class="select sm" bind:value={n.kind}>
						{#each NAME_KINDS as k}<option value={k.id}>{k.label}</option>{/each}
					</select>
					<input class="input" placeholder="Название…" bind:value={n.value} />
					{#if names.length > 1}
						<button class="btn btn-icon sm" onclick={() => names.splice(i, 1)}><Minus size={14} /></button>
					{/if}
				</div>
			{/each}
			<button class="btn sm" onclick={addName}><Plus size={14} /> Добавить название</button>
		</div>

		<div class="field">
			<span class="label">Обложка</span>
			<div class="cover-area">
				<div class="cover-preview">
					{#if coverUrl}
						<img src={coverUrl} alt="Обложка" />
					{:else}
						<div class="cover-placeholder"><ImagePlus size={32} /></div>
					{/if}
				</div>
				<div class="cover-actions">
					<label class="btn">
						<ImagePlus size={14} />
						{coverUploading ? 'Загрузка…' : 'С устройства'}
						<input type="file" accept="image/*" hidden onchange={onCoverFile} />
					</label>
					<div class="row-2">
						<input class="input" placeholder="Ссылка на изображение…" bind:value={coverUrlInput} />
						<button class="btn btn-primary sm" onclick={onCoverURL} disabled={coverUploading}>
							<Link size={14} />
						</button>
					</div>
				</div>
			</div>
		</div>

		<div class="field">
			<span class="label">Изображения</span>
			{#if images.length > 0}
				<div class="gallery-edit">
					{#each images as f, i}
						<div class="gallery-thumb">
							<img src={imageUrls[f]} alt="" />
							<div class="thumb-actions">
								{#if i > 0}
									<button class="btn btn-icon sm" onclick={() => moveImage(i, -1)}><ChevronLeft size={13} /></button>
								{/if}
								{#if i < images.length - 1}
									<button class="btn btn-icon sm" onclick={() => moveImage(i, 1)}><ChevronRight size={13} /></button>
								{/if}
								<button class="btn btn-icon sm" onclick={() => removeImage(i)}><X size={13} /></button>
							</div>
						</div>
					{/each}
				</div>
			{/if}
			<label class="btn sm">
				<ImagePlus size={14} />
				{imagesUploading ? 'Загрузка…' : 'Добавить изображения'}
				<input type="file" accept="image/*" multiple hidden onchange={onImagesInput} />
			</label>
		</div>

		<div class="field">
			<span class="label">Описание</span>
			<textarea class="textarea" rows={4} bind:value={synopsis} placeholder="Сюжет, синопсис…"></textarea>
		</div>

		<div class="field">
			<span class="label">Создатели</span>
			{#each creators as c, i}
				<div class="row-3">
					<select class="select sm" bind:value={c.role}>
						{#each CREATOR_ROLES as r}<option value={r.id}>{r.label}</option>{/each}
					</select>
					<input class="input" placeholder="Имя…" bind:value={c.name} />
					<button class="btn btn-icon sm" onclick={() => creators.splice(i, 1)}><Minus size={14} /></button>
				</div>
			{/each}
			<button class="btn sm" onclick={addCreator}><Plus size={14} /> Добавить создателя</button>
		</div>

		<div class="form-row">
			<div class="field">
				<span class="label">Жанры</span>
				{#each genres as g, i}
					<div class="row-2">
						<input class="input" placeholder="Жанр…" bind:value={genres[i]} />
						<button class="btn btn-icon sm" onclick={() => genres.splice(i, 1)}><Minus size={14} /></button>
					</div>
				{/each}
				<button class="btn sm" onclick={addGenre}><Plus size={14} /> Жанр</button>
			</div>
			<div class="field">
				<span class="label">Теги</span>
				{#each tags as t, i}
					<div class="row-2">
						<input class="input" placeholder="Тег…" bind:value={tags[i]} />
						<button class="btn btn-icon sm" onclick={() => tags.splice(i, 1)}><Minus size={14} /></button>
					</div>
				{/each}
				<button class="btn sm" onclick={addTag}><Plus size={14} /> Тег</button>
			</div>
		</div>

		<div class="form-row">
			<div class="field">
				<span class="label">Оценка (0–5)</span>
				<input type="number" class="input" min="0" max="5" step="0.5" bind:value={score} />
			</div>
			<div class="field">
				<span class="label">Мой статус</span>
				<select class="select" bind:value={status}>
					{#each STATUSES as s}<option value={s.id}>{s.label}</option>{/each}
				</select>
			</div>
		</div>

		<div class="field">
			<span class="label">Статус тайтла</span>
			<select class="select" bind:value={releaseStatus}>
				{#each RELEASE_STATUSES as s}<option value={s.id}>{s.label}</option>{/each}
			</select>
		</div>

		{#if typeShelves.length > 0}
			<div class="field">
				<span class="label">Полка</span>
				<select class="select" bind:value={shelfId}>
					<option value={0}>Без полки</option>
					{#each typeShelves as s}<option value={s.id}>{s.name}</option>{/each}
				</select>
			</div>
		{/if}

		<div class="field">
			<span class="label">Свой список</span>
			<input class="input" placeholder="Название пользовательского списка…" bind:value={customList} />
		</div>

		<div class="field">
			<span class="label">Персонажи</span>
			{#each charRefs as char, i}
				<div class="row-2">
					<select class="select sm" bind:value={char.id}>
						<option value={0} disabled hidden>Персонаж…</option>
						{#each charCandidates(i) as c (c.id)}
							<option value={c.id}>{displayName(c.names)}</option>
						{/each}
					</select>
					<button class="btn btn-icon sm" onclick={() => charRefs.splice(i, 1)}><Minus size={14} /></button>
				</div>
			{/each}
			<button class="btn sm" onclick={addCharRef}><Plus size={14} /> Добавить персонажа</button>
			<span class="rel-hint">Список персонажей создаётся в разделе «Персонажи»; здесь выбираются уже добавленные.</span>
		</div>

		<div class="field">
			<span class="label">Связи</span>
			{#each relations as rel, i}
				<div class="rel-edit">
					<div class="row-3">
						<select class="select sm" bind:value={rel.relatedId}>
							<option value={0} disabled hidden>Тайтл…</option>
							{#each relationCandidates(i) as t (t.id)}
								<option value={t.id}>{displayName(t.names)}</option>
							{/each}
						</select>
						<input class="input" placeholder="Метка здесь (например, «Источник»)…" bind:value={rel.label} />
						<button class="btn btn-icon sm" onclick={() => relations.splice(i, 1)}><Minus size={14} /></button>
					</div>
					<input class="input" placeholder="Метка на другом тайтле (например, «Экранизация»)…" bind:value={rel.reverseLabel} />
				</div>
			{/each}
			<button class="btn sm" onclick={addRelation}><Plus size={14} /> Добавить связь</button>
		</div>

		<div class="field">
			<span class="label">Прогресс</span>
			<div class="progress-fields">
				{#each progressForType(type) as f}
					<div class="prog-item">
						<span class="prog-label">{f.label}</span>
						<div class="prog-row">
							<button class="btn btn-icon sm" onclick={() => (progress[f.id] = Math.max(0, progress[f.id] - 1))}><Minus size={14} /></button>
							<input type="number" class="input sm" min="0" bind:value={progress[f.id]} />
							<button class="btn btn-icon sm" onclick={() => (progress[f.id] = progress[f.id] + 1)}><Plus size={14} /></button>
						</div>
					</div>
				{/each}
				{#if type === 'read'}
					<div class="prog-item">
						<span class="prog-label">Вышло глав</span>
						<input type="number" class="input sm" min="0" bind:value={progress.totalChapters} />
					</div>
				{:else}
					<div class="prog-item">
						<span class="prog-label">Вышло серий</span>
						<input type="number" class="input sm" min="0" bind:value={progress.totalEpisodes} />
					</div>
				{/if}
			</div>
		</div>

		<div class="actions">
			<button class="btn" onclick={() => history.back()}>Отмена</button>
			<button class="btn btn-primary" onclick={save} disabled={saving}>
				{saving ? 'Сохранение…' : 'Сохранить'}
			</button>
		</div>
	</div>
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
	.error {
		background: rgba(224, 72, 77, 0.12);
		color: #e0484d;
		padding: 10px 14px;
		border-radius: var(--radius-sm);
		margin-bottom: 14px;
		font-size: 13px;
	}
	.form {
		background: var(--bg-elev);
		padding: 20px 24px;
		border-radius: var(--radius);
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
	.rel-edit {
		display: flex;
		flex-direction: column;
		gap: 6px;
		margin-bottom: 8px;
	}
	.rel-hint {
		display: block;
		font-size: 11px;
		color: var(--text-dim);
		margin-top: 6px;
	}
	.sm {
		width: auto;
		min-width: 110px;
		padding: 7px 10px;
		font-size: 12px;
	}
	.btn.sm {
		padding: 6px 10px;
		font-size: 12px;
	}
	.btn-icon.sm {
		width: 32px;
		height: 32px;
		min-width: 32px;
	}
	.input.sm {
		width: 70px;
		text-align: center;
	}

	.cover-area {
		display: flex;
		gap: 16px;
		margin-bottom: 8px;
	}
	.cover-preview {
		width: 120px;
		height: 170px;
		border-radius: var(--radius);
		overflow: hidden;
		flex-shrink: 0;
		background: var(--bg-elev-2);
		display: flex;
		align-items: center;
		justify-content: center;
	}
	.cover-preview img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}
	.cover-placeholder {
		color: var(--text-faint);
	}
	.cover-actions {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.gallery-edit {
		display: flex;
		gap: 10px;
		flex-wrap: wrap;
		margin-bottom: 8px;
	}
	.gallery-thumb {
		width: 84px;
	}
	.gallery-thumb img {
		width: 84px;
		height: 112px;
		object-fit: cover;
		display: block;
		border-radius: 6px;
		border: 1px solid var(--border);
	}
	.thumb-actions {
		display: flex;
		gap: 2px;
		margin-top: 4px;
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

	.actions {
		display: flex;
		justify-content: flex-end;
		gap: 8px;
		margin-top: 20px;
		padding-top: 16px;
		border-top: 1px solid var(--border);
	}

	@media (max-width: 640px) {
		.form-row {
			grid-template-columns: 1fr;
		}
		.cover-area {
			flex-direction: column;
			align-items: center;
		}
	}
</style>
