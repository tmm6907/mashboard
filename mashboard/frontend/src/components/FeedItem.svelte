<script lang="ts">
    import type { FeedItemData } from "src/stores/feeds.svelte.ts";
    import { getTimeAgo } from "src/lib/utils.ts";
    import { appState } from "src/stores/appState.svelte.ts";
    import { HandleSaveFeedItem } from "../../wailsjs/go/main/App.js";

    let props = $props();
    let item: FeedItemData = props.item;
    let clipboardTimeout = $state(null);

    function copyLink(e: MouseEvent, link: string) {
        navigator.clipboard
            .writeText(link)
            .then(() => {
                const btn = e.target as HTMLButtonElement;
                const clipboard = btn.getElementsByClassName("clipboard")[0];
                clipboard.classList.toggle("hidden");
                console.log("copied ", link);
                console.log(clipboard);
                clipboardTimeout = setTimeout(
                    () => clipboard.classList.toggle("hidden"),
                    3000,
                );
            })
            .catch((err) => console.error(err));
        clearInterval(clipboardTimeout);
    }
</script>

<div class=" *:not-last:mb-4 bg-base-300 border-0 rounded shadow-sm py-2">
    <div>
        <div class="flex justify-between px-4 py-2 text-sm">
            <div>
                <span class="text-secondary"
                    ><a href={item.feedName}>{item.feedName}</a></span
                >
                |
                <span class="text-base-content/70"
                    >{getTimeAgo(item.pubDate)}</span
                >
            </div>
            <div class="flex space-x-1">
                {#if true || false}
                    <button class="btn btn-sm btn-ghost" aria-label="read feed"
                        ><i class="fa fa-check"></i></button
                    >
                {/if}
                <button
                    class="btn btn-sm btn-ghost"
                    aria-label="bookmark"
                    onclick={async () => {
                        item.saved = !item.saved;
                        let res = await HandleSaveFeedItem({
                            id: item.id,
                            value: item.saved,
                        });
                        if (res.error) {
                            item.saved = !item.saved;
                            console.error(res.error);
                            return;
                        }
                    }}
                    ><i
                        class={item.saved
                            ? "fa fa-bookmark text-primary"
                            : "fa-regular fa-bookmark"}
                    ></i></button
                >
            </div>
        </div>
        <h3 class="text-lg px-4">{item.title}</h3>
    </div>
    {#if item.image}
        <figure class="bg-black">
            <img
                loading="lazy"
                class="mx-auto h-96"
                src={item.image}
                alt={item.altText || item.title}
            />
        </figure>
    {:else}
        <div class="overflow-hidden">
            <p class="text-base-content/70 px-4 line-clamp-3 max-w-[75cih]">
                {item.description}
            </p>
        </div>
    {/if}
    <div class="flex justify-end space-x-2 px-4">
        <button class="feed-btn btn btn-ghost"
            ><i class="fa fa-plus"></i>Add</button
        >
        <button
            class="relative feed-btn btn btn-ghost"
            onclick={(e) => copyLink(e, item.link)}
        >
            <i class="fa fa-link"></i>Share
            <span
                class="clipboard hidden absolute -top-20 -left-1/2 p-4 bg-primary z-10 whitespace-nowrap"
                >Copied to clipboard!</span
            >
        </button>
        <a
            href={`#feeds/items/${item.id}`}
            class="feed-btn btn btn-ghost"
            onclick={() => appState.currentPage.push(`#feeds/items/${item.id}`)}
            >View</a
        >
    </div>
</div>

<style>
</style>
