<script>
    import { appState } from "src/stores/appState.svelte.ts";
    import CompactFeedItem from "./CompactFeedItem.svelte";
</script>

<div class="bg-base-300 rounded h-72 py-4">
    <h2 class="text-lg px-4">Recently Viewed</h2>
    <ul class="list overflow-auto max-h-[calc(100%-2rem)] overscroll-contain">
        {#each appState.recentPosts as post}
            <CompactFeedItem item={post} />
        {/each}
    </ul>
</div>
