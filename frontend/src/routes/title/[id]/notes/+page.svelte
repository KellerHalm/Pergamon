<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import RichEditor from '$lib/components/RichEditor.svelte';
	import { titlesApi, notesApi } from '$lib/api';
	import type { Note, Title } from '../../../../app.d';
	import ArrowLeft from 'lucide-svelte/icons/arrow-left';
	import Plus from 'lucide-svelte/icons/plus';
	import Pencil from 'lucide-svelte/icons/pencil';
	import Trash2 from 'lucide-svelte/icons/trash-2';
	import NotebookPen from 'lucide-svelte/icons/notebook-pen';

	let title = $state<Title | null>(null);
	let notes = $state<Note[]>([]);
	let loading = $state(true);
	let formOpen = $state(false);
	let formKey = $state(0);
	let editId = $state(0);
	let heading = $state('');
	let content = $state('');
	let saving = $state(false);

	const id = $derived(Number($page.params.id));

	onMount(async () => {
		title = await titlesApi.get(id);
		await loadNotes();
		loading = false;
	});

	async function loadNotes() {
		notes = (await notesApi.list(id)) || [];
	}

	function openAdd() {
		editId = 0;
		heading = '';
		content = '';
		formKey++;
		formOpen = true;
	}

	function openEdit(n: Note) {
		editId = n.id;
		heading = n.heading;
		content = n.content;
		formKey++;
		formOpen = true;
	}

	async function saveNote() {
		saving = true;
		await notesApi.save({ id: editId, titleId: id, heading: heading.trim(), content });
		formOpen = false;
		saving = false;
		await loadNotes();
	}

	async function removeNote(n: Note) {
		if (!confirm('Удалить эту заметку?')) return;
		await notesApi.delete(n.id);
		await loadNotes();
	}

	function noteDate(n: Note): string {
		return n.updatedAt && n.updatedAt !== n.createdAt ? n.updatedAt : n.createdAt || '';
	}

	function fmtDate(s: string): string {
		if (!s) return '';
		const d = new Date(s.includes('T') ? s : s.replace(' ', 'T') + 'Z');
		if (isNaN(d.getTime())) return s;
		return d.toLocaleString('ru-RU', {
			day: 'numeric',
			month: 'long',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}
</script>

<div class="page">
	<div class="top">
		<button class="btn" onclick={() => goto(`/title/${id}`)}>
			<ArrowLeft size={16} />
			Назад
		</button>
		<button class="btn btn-primary" onclick={openAdd}><Plus size={14} /> Добавить заметку</button>
	</div>

	<div class="head">
		<NotebookPen size={22} />
		<h1>Заметки</h1>
		{#if notes.length > 0}<span class="count">{notes.length}</span>{/if}
	</div>
	{#if title}
		<div class="subtitle">{title.names.map((n) => n.value).join(' / ')}</div>
	{/if}

	{#if formOpen}
		<div class="edit-panel">
			<h2>{editId === 0 ? 'Новая заметка' : 'Редактирование заметки'}</h2>
			<div class="field">
				<span class="label">Заголовок</span>
				<input class="input" placeholder="Заголовок заметки…" bind:value={heading} />
			</div>
			<div class="field">
				<span class="label">Текст</span>
				{#key formKey}
					<RichEditor content={content} onUpdate={(html) => (content = html)} />
				{/key}
			</div>
			<div class="edit-actions">
				<button class="btn" onclick={() => (formOpen = false)}>Отмена</button>
				<button class="btn btn-primary" disabled={saving} onclick={saveNote}>Сохранить</button>
			</div>
		</div>
	{/if}

	{#if loading}
		<div class="loading">Загрузка…</div>
	{:else if notes.length === 0 && !formOpen}
		<div class="empty">
			<NotebookPen size={40} />
			<p>Заметок пока нет</p>
			<button class="btn btn-primary" onclick={openAdd}><Plus size={14} /> Добавить заметку</button>
		</div>
	{:else}
		<div class="notes">
			{#each notes as n (n.id)}
				<article class="note">
					<div class="note-head">
						<div class="note-meta">
							<h3>{n.heading || 'Без названия'}</h3>
							<span class="date">
								{fmtDate(noteDate(n))}{n.updatedAt !== n.createdAt ? ' · изменено' : ''}
							</span>
						</div>
						<div class="note-actions">
							<button class="btn btn-icon sm" title="Редактировать" onclick={() => openEdit(n)}>
								<Pencil size={14} />
							</button>
							<button class="btn btn-icon sm btn-danger" title="Удалить" onclick={() => removeNote(n)}>
								<Trash2 size={14} />
							</button>
						</div>
					</div>
					{#if n.content}
						<div class="note-content">
							{@html n.content}
						</div>
					{/if}
				</article>
			{/each}
		</div>
	{/if}
</div>

<style>
	.page {
		max-width: 960px;
		margin: 0 auto;
	}
	.top {
		display: flex;
		justify-content: space-between;
		gap: 8px;
		margin-bottom: 20px;
		flex-wrap: wrap;
	}
	.head {
		display: flex;
		align-items: center;
		gap: 10px;
		color: var(--accent);
		margin-bottom: 4px;
	}
	.head h1 {
		font-size: 22px;
		font-weight: 700;
		color: var(--text);
	}
	.count {
		background: var(--bg-elev-2);
		color: var(--text-dim);
		border-radius: 99px;
		font-size: 12px;
		font-weight: 600;
		padding: 2px 10px;
	}
	.subtitle {
		color: var(--text-dim);
		font-size: 13px;
		margin-bottom: 20px;
	}
	.loading {
		text-align: center;
		padding: 60px;
		color: var(--text-dim);
	}

	.edit-panel {
		background: var(--bg-elev);
		border-radius: var(--radius);
		padding: 20px 24px;
		margin-bottom: 22px;
	}
	.edit-panel h2 {
		margin-bottom: 16px;
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
	.edit-actions {
		display: flex;
		justify-content: flex-end;
		gap: 8px;
		margin-top: 18px;
	}

	.empty {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 12px;
		padding: 60px 20px;
		color: var(--text-dim);
		background: var(--bg-elev);
		border-radius: var(--radius);
	}
	.empty p {
		font-size: 15px;
	}

	.notes {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}
	.note {
		background: var(--bg-elev);
		border-radius: var(--radius);
		padding: 16px 20px;
	}
	.note-head {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		gap: 10px;
		margin-bottom: 8px;
	}
	.note-meta h3 {
		font-size: 16px;
		font-weight: 600;
	}
	.date {
		font-size: 12px;
		color: var(--text-faint, var(--text-dim));
	}
	.note-actions {
		display: flex;
		gap: 4px;
		flex: none;
	}
	.btn-icon.sm {
		width: 30px;
		height: 30px;
		min-width: 30px;
	}
	.note-content {
		color: var(--text);
		font-size: 14px;
		line-height: 1.6;
		word-break: break-word;
	}
	.note-content :global(p) {
		margin: 0 0 8px;
	}
	.note-content :global(p:last-child) {
		margin-bottom: 0;
	}
	.note-content :global(ul),
	.note-content :global(ol) {
		padding-left: 20px;
		margin: 4px 0;
	}
	.note-content :global(blockquote) {
		border-left: 3px solid var(--accent);
		padding-left: 12px;
		margin: 8px 0;
		color: var(--text-dim);
	}
</style>
