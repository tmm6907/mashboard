<script>
    import CompactFeedItem from "./CompactFeedItem.svelte";

    let props = $props();
    let suggestedFeeds = props.suggestedFeeds;
</script>

<div class="bg-base-300 rounded h-72 p-4">
    <h2 class="text-lg">Suggested Feeds</h2>
    <div class="flex">
        {#if suggestedFeeds}
            {#each suggestedFeeds as feed}
                <CompactFeedItem
                    item={{
                        title: "item with longer name on Thursday, 6th",
                        feedName: "feed1",
                        createdAt: "2025-12-11 10:30:00",
                    }}
                />
            {/each}
        {:else}
            <div class="flex w-full py-4 justify-center">
                <span class="text-base-content/70">None found</span>
            </div>
        {/if}
    </div>
</div>
