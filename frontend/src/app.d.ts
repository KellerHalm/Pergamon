declare global {
	interface Window {
		go: {
			main: {
				App: {
					ListTitles(q: ListQuery): Promise<Title[]>;
					GetTitle(id: number): Promise<Title>;
					SaveTitle(t: Title): Promise<number>;
					UpdateIncomingRelations(titleId: number, rels: Relation[]): Promise<void>;
					DeleteTitle(id: number): Promise<void>;
					AdjustProgress(id: number, field: string, delta: number): Promise<Title>;
					ListNotes(titleId: number): Promise<Note[]>;
					SaveNote(n: Note): Promise<number>;
					DeleteNote(id: number): Promise<void>;
					AllTags(): Promise<string[]>;
					AllGenres(): Promise<string[]>;
					ListShelves(): Promise<Shelf[]>;
					SaveShelf(s: Shelf): Promise<number>;
					DeleteShelf(id: number): Promise<void>;
					SetShelfItems(shelfId: number, titleIds: number[]): Promise<void>;
					ListCharacters(sort: string): Promise<Character[]>;
					GetCharacter(id: number): Promise<Character>;
					SaveCharacter(c: Character): Promise<number>;
					DeleteCharacter(id: number): Promise<void>;
					GetSetting(key: string): Promise<string>;
					SetSetting(key: string, value: string): Promise<void>;
					UploadCoverDataURL(dataUrl: string): Promise<string>;
					UploadCoverFromURL(url: string): Promise<string>;
					CoverDataURL(filename: string): Promise<string>;
					CoverColor(filename: string): Promise<string>;
					SyncNow(): Promise<SyncResult>;
					TestWebDAV(cfg: WebDAVConfig): Promise<void>;
					SaveWebDAVConfig(cfg: WebDAVConfig): Promise<void>;
					GetWebDAVConfig(): Promise<WebDAVConfig>;
					ExportJSON(): Promise<string>;
					ImportJSONPath(path: string): Promise<void>;
					ExportSQLitePath(dest: string): Promise<void>;
					ImportSQLitePath(src: string): Promise<void>;
				};
			};
		};
	}
}

export interface Name {
	kind: string;
	value: string;
}

export interface Creator {
	role: string;
	name: string;
}

export interface Relation {
	relatedId: number;
	label: string;
	reverseLabel?: string;
	name?: string;
	cover?: string;
	status?: string;
}

export interface Progress {
	volumes: number;
	chapters: number;
	pages: number;
	seasons: number;
	episodes: number;
	minutes: number;
	totalChapters: number;
	totalEpisodes: number;
}

export interface Title {
	id: number;
	type: string;
	category: string;
	names: Name[];
	cover: string;
	images: string[];
	synopsis: string;
	creators: Creator[];
	genres: string[];
	tags: string[];
	relations: Relation[];
	reverseRelations?: Relation[];
	characters: CharacterRef[];
	score: number;
	status: string;
	releaseStatus: string;
	customList: string;
	progress: Progress;
	notes: string;
	spineColor: string;
	createdAt: string;
	updatedAt: string;
}

export interface CharacterRef {
	id: number;
	name?: string;
	mainImage?: string;
}

export interface TitleRef {
	id: number;
	name: string;
	cover: string;
	status: string;
}

export interface Character {
	id: number;
	names: Name[];
	mainImage: string;
	age: string;
	description: string;
	images: string[];
	titles: TitleRef[];
	titleIds: number[];
	createdAt: string;
	updatedAt: string;
}

export interface Note {
	id: number;
	titleId: number;
	heading: string;
	content: string;
	createdAt?: string;
	updatedAt?: string;
}

export interface Shelf {
	id: number;
	name: string;
	kind: string;
	position: number;
	titleIds: number[];
	createdAt: string;
}

export interface ListQuery {
	sort?: string;
	type?: string;
	category?: string;
	status?: string;
	tags?: string[];
	search?: string;
}

export interface SyncResult {
	direction: string;
	uploaded: number;
	message: string;
	time: string;
}

export interface WebDAVConfig {
	url: string;
	username: string;
	password: string;
	remoteDir: string;
}

export {};
