<script lang="ts">
	import { onMount } from 'svelte';
	import { syncApi, settingsApi } from '$lib/api';
	import { ACCENTS } from '$lib/constants';
	import type { WebDAVConfig } from '../../app.d';
	import RefreshCw from 'lucide-svelte/icons/refresh-cw';
	import Download from 'lucide-svelte/icons/download';
	import Upload from 'lucide-svelte/icons/upload';
	import Cloud from 'lucide-svelte/icons/cloud';

	let webdav = $state<WebDAVConfig>({ url: '', username: '', password: '', remoteDir: '/Pergamon' });
	let syncResult = $state('');
	let syncing = $state(false);
	let syncOnStart = $state(false);
	let syncOnShutdown = $state(false);
	let theme = $state('dark');
	let accent = $state('#7c5cff');
	let loaded = $state(false);

	onMount(async () => {
		try {
			const cfg = await syncApi.getConfig();
			webdav = cfg;
		} catch (e) {}
		syncOnStart = (await settingsApi.get('sync_on_startup')) === '1';
		syncOnShutdown = (await settingsApi.get('sync_on_shutdown')) === '1';
		theme = (await settingsApi.get('theme')) || 'dark';
		accent = (await settingsApi.get('accent')) || '#7c5cff';
		loaded = true;
	});

	function errText(e: unknown): string {
		if (e instanceof Error) return e.message;
		return String(e);
	}

	async function saveWebDAV() {
		await syncApi.saveConfig(webdav);
		syncResult = 'Настройки сохранены';
		setTimeout(() => (syncResult = ''), 3000);
	}

	async function testWebDAV() {
		syncResult = 'Проверка…';
		try {
			await syncApi.test(webdav);
			syncResult = '✓ Подключение успешно';
		} catch (e: any) {
			syncResult = '✗ Ошибка: ' + errText(e);
		}
	}

	async function syncNow() {
		syncing = true;
		syncResult = '';
		try {
			const res = await syncApi.now();
			syncResult = res.message;
		} catch (e: any) {
			syncResult = 'Ошибка: ' + errText(e);
		}
		syncing = false;
	}

	async function toggleAutoStart() {
		syncOnStart = !syncOnStart;
		await settingsApi.set('sync_on_startup', syncOnStart ? '1' : '0');
	}

	async function toggleAutoShutdown() {
		syncOnShutdown = !syncOnShutdown;
		await settingsApi.set('sync_on_shutdown', syncOnShutdown ? '1' : '0');
	}

	async function setTheme(t: string) {
		theme = t;
		localStorage.setItem('theme', t);
		document.documentElement.setAttribute('data-theme', t);
		await settingsApi.set('theme', t);
	}

	async function setAccent(c: string) {
		accent = c;
		localStorage.setItem('accent', c);
		document.documentElement.style.setProperty('--accent', c);
		await settingsApi.set('accent', c);
	}

	async function exportJSON() {
		try {
			const path = await syncApi.exportJSON();
			syncResult = 'Экспорт: ' + path;
		} catch (e: any) {
			syncResult = 'Ошибка: ' + errText(e);
		}
	}

	async function importJSON() {
		const input = document.createElement('input');
		input.type = 'file';
		input.accept = '.json';
		input.onchange = async () => {
			const file = input.files?.[0];
			if (!file) return;
			try {
				const text = await file.text();
				const json = JSON.parse(text);
				await syncApi.importJSON(JSON.stringify(json, null, 2));
				syncResult = 'Импорт завершён';
			} catch (e: any) {
				syncResult = 'Ошибка: ' + errText(e);
			}
		};
		input.click();
	}
</script>

{#if loaded}
	<div class="page">
		<h1>Настройки</h1>

		<section class="section">
			<h2><Cloud size={18} /> Синхронизация WebDAV</h2>
			<p class="hint">Поддержка Яндекс.Диска, Nextcloud и других WebDAV-хранилищ</p>

			<div class="field">
				<span class="label">URL WebDAV</span>
				<input class="input" placeholder="https://webdav.yandex.ru" bind:value={webdav.url} />
			</div>
			<div class="form-row">
				<div class="field">
					<span class="label">Логин</span>
					<input class="input" placeholder="user@yandex.ru" bind:value={webdav.username} />
				</div>
				<div class="field">
					<span class="label">Пароль приложения</span>
					<input class="input" type="password" placeholder="OAuth-токен или пароль" bind:value={webdav.password} />
				</div>
			</div>
			<div class="field">
				<span class="label">Папка на диске</span>
				<input class="input" placeholder="/Pergamon" bind:value={webdav.remoteDir} />
			</div>

			<div class="actions">
				<button class="btn" onclick={testWebDAV}>Проверить</button>
				<button class="btn btn-primary" onclick={saveWebDAV}>Сохранить</button>
				<button class="btn" onclick={syncNow} disabled={syncing}>
					<RefreshCw size={14} class={syncing ? 'spin' : ''} />
					{syncing ? 'Синхронизация…' : 'Синхронизировать'}
				</button>
			</div>

			<div class="auto-sync">
				<label class="checkbox-row">
					<input type="checkbox" bind:checked={syncOnStart} onchange={toggleAutoStart} />
					<span>Авто-синхронизация при запуске</span>
				</label>
				<label class="checkbox-row">
					<input type="checkbox" bind:checked={syncOnShutdown} onchange={toggleAutoShutdown} />
					<span>Авто-синхронизация при закрытии</span>
				</label>
			</div>
		</section>

		{#if syncResult}
			<div class="status">{syncResult}</div>
		{/if}

		<section class="section">
			<h2>Бэкап</h2>
			<div class="backup-actions">
				<button class="btn" onclick={exportJSON}>
					<Download size={14} /> Экспорт JSON
				</button>
				<button class="btn" onclick={importJSON}>
					<Upload size={14} /> Импорт JSON
				</button>
			</div>
		</section>

		<section class="section">
			<h2>Внешний вид</h2>
			<div class="field">
				<span class="label">Тема</span>
				<div class="theme-row">
					<button class="btn theme-btn" class:active={theme === 'dark'} onclick={() => setTheme('dark')}>🌙 Тёмная</button>
					<button class="btn theme-btn" class:active={theme === 'light'} onclick={() => setTheme('light')}>☀️ Светлая</button>
				</div>
			</div>
			<div class="field">
				<span class="label">Акцентный цвет</span>
				<div class="accent-row">
					{#each ACCENTS as c}
						<button
							class="accent-dot"
							class:active={accent === c}
							style="background:{c}"
							onclick={() => setAccent(c)}
						></button>
					{/each}
				</div>
			</div>
		</section>
	</div>
{/if}

<style>
	.page {
		max-width: 640px;
		margin: 0 auto;
	}
	h1 {
		font-size: 24px;
		margin-bottom: 24px;
	}
	.section {
		background: var(--bg-elev);
		border-radius: var(--radius);
		padding: 18px 20px;
		margin-bottom: 18px;
	}
	.section h2 {
		display: flex;
		align-items: center;
		gap: 8px;
		font-size: 16px;
		font-weight: 600;
		margin-bottom: 12px;
	}
	.hint {
		font-size: 12px;
		color: var(--text-faint);
		margin-bottom: 14px;
	}
	.form-row {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 12px;
	}
	.actions {
		display: flex;
		gap: 8px;
		margin-top: 14px;
		flex-wrap: wrap;
	}
	.auto-sync {
		margin-top: 14px;
		display: flex;
		flex-direction: column;
		gap: 6px;
	}
	.checkbox-row {
		display: flex;
		align-items: center;
		gap: 8px;
		font-size: 13px;
		color: var(--text-dim);
		cursor: pointer;
	}
	.status {
		padding: 10px 14px;
		border-radius: var(--radius-sm);
		background: var(--accent-soft);
		color: var(--accent);
		font-size: 13px;
		margin-bottom: 16px;
	}
	.backup-actions {
		display: flex;
		gap: 8px;
		flex-wrap: wrap;
	}
	.theme-row {
		display: flex;
		gap: 8px;
	}
	.theme-btn.active {
		background: var(--accent);
		color: #fff;
	}
	.accent-row {
		display: flex;
		gap: 8px;
		flex-wrap: wrap;
	}
	.accent-dot {
		width: 28px;
		height: 28px;
		border-radius: 50%;
		border: 2px solid transparent;
		transition: border-color 0.12s, transform 0.12s;
	}
	.accent-dot:hover {
		transform: scale(1.15);
	}
	.accent-dot.active {
		border-color: var(--text);
		box-shadow: 0 0 0 2px var(--bg-elev);
	}
	.spin {
		animation: spin 1s linear infinite;
	}
	@keyframes spin {
		from { transform: rotate(0deg); }
		to { transform: rotate(360deg); }
	}

	@media (max-width: 640px) {
		.form-row {
			grid-template-columns: 1fr;
		}
	}
</style>
