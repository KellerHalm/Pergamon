<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { peopleApi, coverApi, titlesApi } from '$lib/api';
	import { NAME_KINDS, GENDERS, PEOPLE_ROLES, displayName } from '$lib/constants';
	import type { Person, Title } from '../../app.d';
	import Plus from 'lucide-svelte/icons/plus';
	import Minus from 'lucide-svelte/icons/minus';
	import ImagePlus from 'lucide-svelte/icons/image-plus';
	import Link from 'lucide-svelte/icons/link';
	import X from 'lucide-svelte/icons/x';
	import ChevronLeft from 'lucide-svelte/icons/chevron-left';
	import ChevronRight from 'lucide-svelte/icons/chevron-right';

	let { person = null }: { person?: Person | null } = $props();

	let names = $state<{ kind: string; value: string }[]>([{ kind: 'original', value: '' }]);
	let mainImage = $state('');
	let mainImageUrl = $state('');
	let mainUploading = $state(false);
	let mainUrlInput = $state('');
	let role = $state('author');
	let gender = $state('');
	let age = $state('');
	let birthDate = $state('');
	let deathDate = $state('');
	let description = $state('');
	let images = $state<string[]>([]);
	let imageUrls = $state<Record<string, string>>({});
	let imagesUploading = $state(false);
	let titleIds = $state<number[]>([]);
	let allTitles = $state<Title[]>([]);
	let saving = $state(false);
	let error = $state('');

	onMount(async () => {
		allTitles = (await titlesApi.list({ sort: 'title' })) || [];
		if (person) {
			names = person.names.length > 0 ? person.names.map((n) => ({ ...n })) : [{ kind: 'original', value: '' }];
			mainImage = person.mainImage || '';
			if (mainImage) mainImageUrl = await coverApi.dataURL(mainImage);
			role = person.role || 'author';
			gender = person.gender || '';
			age = person.age || '';
			birthDate = person.birthDate || '';
			deathDate = person.deathDate || '';
			description = person.description || '';
			images = [...(person.images || [])];
			titleIds = [...(person.titleIds || [])];
			for (const f of images) {
				imageUrls[f] = await coverApi.dataURL(f);
			}
		}
	});

	async function onMainFile(e: Event) {
		const input = e.target as HTMLInputElement;
		const file = input.files?.[0];
		if (!file) return;
		mainUploading = true;
		const reader = new FileReader();
		reader.onload = async () => {
			mainImage = await coverApi.uploadDataURL(reader.result as string);
			mainImageUrl = await coverApi.dataURL(mainImage);
			mainUploading = false;
		};
		reader.readAsDataURL(file);
	}

	async function onMainURL() {
		if (!mainUrlInput.trim()) return;
		mainUploading = true;
		try {
			mainImage = await coverApi.uploadFromURL(mainUrlInput.trim());
			mainImageUrl = await coverApi.dataURL(mainImage);
			mainUrlInput = '';
		} catch (e) {
			error = 'Не удалось загрузить изображение по ссылке';
		}
		mainUploading = false;
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

	function addName() {
		names.push({ kind: 'original', value: '' });
	}

	function titleCandidates(): Title[] {
		return allTitles.filter((t) => !titleIds.includes(t.id));
	}

	function pickTitle(e: Event) {
		const v = Number((e.currentTarget as HTMLSelectElement).value);
		if (v > 0 && !titleIds.includes(v)) {
			titleIds.push(v);
		}
	}

	async function save() {
		if (!names.some((n) => n.value.trim() !== '')) {
			error = 'Укажите хотя бы одно имя';
			return;
		}
		saving = true;
		error = '';
		try {
			const p: Person = {
				id: person?.id ?? 0,
				names: names.filter((n) => n.value.trim()),
				mainImage,
				age: age.trim(),
				birthDate: birthDate.trim(),
				deathDate: deathDate.trim(),
				gender,
				role,
				description,
				images: [...images],
				titles: [],
				titleIds: [...titleIds],
				createdAt: person?.createdAt ?? '',
				updatedAt: ''
			};
			const id = await peopleApi.save(p);
			goto('/person/' + id);
		} catch (e: any) {
			error = e.message || 'Ошибка сохранения';
		}
		saving = false;
	}
</script>

{#if error}
	<div class="error">{error}</div>
{/if}

<div class="form">
	<div class="field">
		<span class="label">Main изображение</span>
		<div class="cover-area">
			<div class="cover-preview">
				{#if mainImageUrl}
					<img src={mainImageUrl} alt="Main" />
				{:else}
					<div class="cover-placeholder"><ImagePlus size={32} /></div>
				{/if}
			</div>
			<div class="cover-actions">
				<label class="btn">
					<ImagePlus size={14} />
					{mainUploading ? 'Загрузка…' : 'С устройства'}
					<input type="file" accept="image/*" hidden onchange={onMainFile} />
				</label>
				<div class="row-2">
					<input class="input" placeholder="Ссылка на изображение…" bind:value={mainUrlInput} />
					<button class="btn btn-primary sm" onclick={onMainURL} disabled={mainUploading}>
						<Link size={14} />
					</button>
				</div>
			</div>
		</div>
	</div>

	<div class="field">
		<span class="label">Имена</span>
		{#each names as n, i}
			<div class="row-2">
				<select class="select sm" bind:value={n.kind}>
					{#each NAME_KINDS as k}<option value={k.id}>{k.label}</option>{/each}
				</select>
				<input class="input" placeholder="Имя…" bind:value={n.value} />
				{#if names.length > 1}
					<button class="btn btn-icon sm" onclick={() => names.splice(i, 1)}><Minus size={14} /></button>
				{/if}
			</div>
		{/each}
		<button class="btn sm" onclick={addName}><Plus size={14} /> Добавить имя</button>
	</div>

	<div class="form-row">
		<div class="field">
			<span class="label">Роль</span>
			<select class="select" bind:value={role}>
				{#each PEOPLE_ROLES as r}<option value={r.id}>{r.label}</option>{/each}
			</select>
		</div>
		<div class="field">
			<span class="label">Пол</span>
			<select class="select" bind:value={gender}>
				{#each GENDERS as g}<option value={g.id}>{g.label}</option>{/each}
			</select>
		</div>
	</div>

	<div class="field">
		<span class="label">Возраст</span>
		<input class="input" placeholder="Например: 54 года…" bind:value={age} />
	</div>

	<div class="form-row">
		<div class="field">
			<span class="label">Дата рождения</span>
			<input class="input" placeholder="Например: 1975 или 1975-03-01…" bind:value={birthDate} />
		</div>
		<div class="field">
			<span class="label">Дата смерти</span>
			<input class="input" placeholder="Пусто, если жив…" bind:value={deathDate} />
		</div>
	</div>

	<div class="field">
		<span class="label">Описание</span>
		<textarea class="textarea" rows={4} bind:value={description} placeholder="Биография, деятельность…"></textarea>
	</div>

	<div class="field">
		<span class="label">Тайтлы</span>
		{#if titleIds.length > 0}
			{#each titleIds as tid, i}
				<div class="row-2">
					<input class="input" disabled value={displayName(allTitles.find((t) => t.id === tid)?.names || [])} />
					<button class="btn btn-icon sm" onclick={() => titleIds.splice(i, 1)}><Minus size={14} /></button>
				</div>
			{/each}
		{:else}
			<span class="hint">Деятель пока не привязан ни к одному тайтлу.</span>
		{/if}
		{#if titleCandidates().length > 0}
			<div class="row-2">
				<select class="select sm" value={0} onchange={pickTitle}>
					<option value={0} disabled hidden>Добавить в тайтл…</option>
					{#each titleCandidates() as t (t.id)}
						<option value={t.id}>{displayName(t.names)}</option>
					{/each}
				</select>
			</div>
		{/if}
	</div>

	<div class="field">
		<span class="label">Дополнительные изображения</span>
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

	<div class="actions">
		<button class="btn" onclick={() => history.back()}>Отмена</button>
		<button class="btn btn-primary" onclick={save} disabled={saving}>
			{saving ? 'Сохранение…' : 'Сохранить'}
		</button>
	</div>
</div>

<style>
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
	.field {
		margin-bottom: 14px;
	}
	.label {
		display: block;
		font-size: 12px;
		color: var(--text-dim);
		text-transform: uppercase;
		letter-spacing: 0.04em;
		margin-bottom: 6px;
	}
	.row-2 {
		display: flex;
		gap: 6px;
		margin-bottom: 6px;
	}
	.row-2 .input {
		flex: 1;
	}
	.form-row {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 14px;
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
	.hint {
		font-size: 12px;
		color: var(--text-faint);
		display: block;
		margin-bottom: 6px;
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

	.actions {
		display: flex;
		justify-content: flex-end;
		gap: 8px;
		margin-top: 20px;
		padding-top: 16px;
		border-top: 1px solid var(--border);
	}

	@media (max-width: 640px) {
		.cover-area {
			flex-direction: column;
			align-items: center;
		}
		.form-row {
			grid-template-columns: 1fr;
		}
	}
</style>
