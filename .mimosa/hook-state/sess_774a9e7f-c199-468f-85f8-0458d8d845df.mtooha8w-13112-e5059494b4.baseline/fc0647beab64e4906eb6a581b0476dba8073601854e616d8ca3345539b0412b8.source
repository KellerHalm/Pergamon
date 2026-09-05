import type { Character, ListQuery, Note, Person, Relation, Shelf, Studio, SyncResult, Title, WebDAVConfig } from '../app.d';

const api = () => window.go.main.App;

export const titlesApi = {
	list: (q: ListQuery = {}) => api().ListTitles(q),
	get: (id: number) => api().GetTitle(id),
	save: (t: Title) => api().SaveTitle(t),
	updateIncomingRelations: (id: number, rels: Relation[]) => api().UpdateIncomingRelations(id, rels),
	delete: (id: number) => api().DeleteTitle(id),
	adjustProgress: (id: number, field: string, delta: number) => api().AdjustProgress(id, field, delta)
};

export const notesApi = {
	list: (titleId: number) => api().ListNotes(titleId) as Promise<Note[]>,
	save: (n: Note) => api().SaveNote(n),
	delete: (id: number) => api().DeleteNote(id)
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

export const charactersApi = {
	list: (sort: string = '') => api().ListCharacters(sort) as Promise<Character[]>,
	get: (id: number) => api().GetCharacter(id) as Promise<Character>,
	save: (c: Character) => api().SaveCharacter(c),
	delete: (id: number) => api().DeleteCharacter(id)
};

export const studiosApi = {
	list: (sort: string = '') => api().ListStudios(sort) as Promise<Studio[]>,
	get: (id: number) => api().GetStudio(id) as Promise<Studio>,
	save: (s: Studio) => api().SaveStudio(s),
	delete: (id: number) => api().DeleteStudio(id)
};

export const peopleApi = {
	list: (sort: string = '') => api().ListPeople(sort) as Promise<Person[]>,
	get: (id: number) => api().GetPerson(id) as Promise<Person>,
	save: (p: Person) => api().SavePerson(p),
	delete: (id: number) => api().DeletePerson(id)
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
