import type { ListQuery, Shelf, SyncResult, Title, WebDAVConfig } from '../app.d';

const api = () => window.go.main.App;

export const titlesApi = {
	list: (q: ListQuery = {}) => api().ListTitles(q),
	get: (id: number) => api().GetTitle(id),
	save: (t: Title) => api().SaveTitle(t),
	delete: (id: number) => api().DeleteTitle(id),
	adjustProgress: (id: number, field: string, delta: number) => api().AdjustProgress(id, field, delta)
};

export const metaApi = {
	allTags: () => api().AllTags(),
	allGenres: () => api().AllGenres()
};

export const shelvesApi = {
	list: () => api().ListShelves() as Promise<Shelf[]>,
	save: (s: Shelf) => api().SaveShelf(s),
	delete: (id: number) => api().DeleteShelf(id),
	setItems: (shelfId: number, ids: number[]) => api().SetShelfItems(shelfId, ids)
};

export const settingsApi = {
	get: (k: string) => api().GetSetting(k),
	set: (k: string, v: string) => api().SetSetting(k, v)
};

export const coverApi = {
	uploadDataURL: (d: string) => api().UploadCoverDataURL(d),
	uploadFromURL: (u: string) => api().UploadCoverFromURL(u),
	dataURL: (f: string) => api().CoverDataURL(f),
	color: (f: string) => api().CoverColor(f)
};

export const syncApi = {
	now: () => api().SyncNow() as Promise<SyncResult>,
	test: (c: WebDAVConfig) => api().TestWebDAV(c),
	saveConfig: (c: WebDAVConfig) => api().SaveWebDAVConfig(c),
	getConfig: () => api().GetWebDAVConfig() as Promise<WebDAVConfig>,
	exportJSON: () => api().ExportJSON(),
	importJSON: (p: string) => api().ImportJSONPath(p),
	exportSQLite: (p: string) => api().ExportSQLitePath(p),
	importSQLite: (p: string) => api().ImportSQLitePath(p)
};
