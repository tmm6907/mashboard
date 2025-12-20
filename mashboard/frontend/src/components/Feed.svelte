<script>
    import FeedHeader from "./FeedHeader.svelte";
    import { feed } from "src/stores/feeds.svelte.ts";
    import FeedItem from "./FeedItem.svelte";
    import { onDestroy, onMount } from "svelte";
    let props = $props();
    let scrollY = $state(0);
    let isNearBottom = $state(false);
    let hasMoreData = $state(false);

    let checkInterval = $state(null);
    let isChecking = $state(false);

    $effect(() => {
        if (scrollY > 0) {
            isNearBottom =
                scrollY >
                document.documentElement.scrollHeight -
                    window.innerHeight -
                    200;
        }
    });

    const checkForNewFeeds = async () => {
        if (isChecking) return;
        isChecking = true;
        try {
            let res = await props.request(0, feed.category, feed.saved);
            if (res.error) {
                console.error(res.error);
                return;
            }
            if (
                JSON.stringify(res.data) !==
                JSON.stringify(feed.items.slice(0, 25))
            ) {
                feed.setNewFeeds(res.data);
            }
        } catch (error) {
            console.error("Error checking feeds:", error);
        } finally {
            isChecking = false;
        }
    };

    const checkMoreFeeds = async () => {
        feed.offset += 25;
        feed.isLoading = true;
        let res = await props.request(feed.offset, feed.category, feed.saved);
        feed.isLoading = false;
        if (res.error) {
            console.error(res.error);
            return;
        }
        hasMoreData = res.data != undefined;
        if (hasMoreData) feed.addMoreFeeds(res.data);
    };

    onMount(async () => {
        feed.items = [];
        feed.newItems = [];
        feed.offset = 0;
        feed.saved = props.saved;
        isNearBottom = props.saved;
        feed.isLoading = true;
        let res = await props.request(feed.offset, feed.category, feed.saved);
        feed.isLoading = false;
        if (res.error) {
            console.error(res.error);
            return;
        }
        feed.setFeeds(res.data);
        await checkMoreFeeds();
        checkInterval = setInterval(checkForNewFeeds, 300000);
    });
    onDestroy(() => {
        clearInterval(checkInterval);
    });
</script>

<svelte:window bind:scrollY />

<section class="flex justify-center p-4">
    <div class="flex flex-col space-y-4 w-[75ch] mb-24">
        <FeedHeader header={props.header} />
        {#each feed.items as feedItem (feedItem.id)}
            <FeedItem item={feedItem} />
        {/each}
    </div>

    {#if feed.newItems.length > 0}
        <div class="fixed top-8 left-1/2 transform -translate-x-1/2 z-10">
            <button
                onclick={() => {
                    feed.refresh();
                    window.scrollTo({ top: 0, behavior: "smooth" });
                }}
                class="btn btn-primary"
            >
                New Feeds Available
            </button>
        </div>
    {/if}
    {#if isNearBottom && hasMoreData}
        <div class="fixed bottom-8 left-1/2 transform -translate-x-1/2 z-10">
            <button
                onclick={async () => await checkMoreFeeds()}
                disabled={feed.isLoading}
                class="btn btn-secondary"
            >
                {feed.isLoading ? "Loading..." : "Show More"}
            </button>
        </div>
    {/if}
</section>

<style></style>
