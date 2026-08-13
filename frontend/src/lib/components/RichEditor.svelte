<script lang="ts">
	import { onMount } from 'svelte';
	import { Editor } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import { SvelteRenderer } from 'svelte-tiptap';

	let { content = '', onUpdate }: { content?: string; onUpdate?: (html: string) => void } = $props();

	let editor: Editor | null = $state(null);
	let el: HTMLElement | null = $state(null);

	onMount(() => {
		editor = new Editor({
			element: el!,
			extensions: [StarterKit],
			content: content || '<p></p>',
			onUpdate: ({ editor }) => {
				onUpdate?.(editor.getHTML());
			}
		});
		return () => editor?.destroy();
	});
</script>

<div class="editor" bind:this={el}></div>

<style>
	.editor {
		min-height: 120px;
		padding: 10px 14px;
		border-radius: var(--radius-sm);
		background: var(--bg-elev-2);
		border: 1px solid var(--border);
		outline: none;
		color: var(--text);
		font-size: 14px;
		line-height: 1.6;
	}
	.editor:focus-within {
		border-color: var(--accent);
	}
	.editor :global(p) {
		margin: 0 0 8px;
	}
	.editor :global(p:last-child) {
		margin-bottom: 0;
	}
	.editor :global(strong) {
		font-weight: 700;
	}
	.editor :global(em) {
		font-style: italic;
	}
	.editor :global(ul),
	.editor :global(ol) {
		padding-left: 20px;
		margin: 4px 0;
	}
	.editor :global(blockquote) {
		border-left: 3px solid var(--accent);
		padding-left: 12px;
		margin: 8px 0;
		color: var(--text-dim);
	}
</style>
