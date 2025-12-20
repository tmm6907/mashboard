<script lang="ts">
    import Fab from "src/components/Fab.svelte";
    import { getTimeAgo } from "src/lib/utils.ts";
    import { appState } from "src/stores/appState.svelte.ts";
    import { feed } from "src/stores/feeds.svelte.ts";
    import { BrowserOpenURL } from "../../wailsjs/runtime/runtime.js";
    import { onMount } from "svelte";
    import { HandleSaveFeedItem } from "../../wailsjs/go/main/App.js";

    onMount(() => {
        console.log("recent posts", appState.recentPosts);
        // if (appState.recentPosts.length > 5) {
        //     if (
        //         appState.recentPosts.filter((p) => p !== feeds.currentItem)
        //             .length == 0
        //     ) {
        //         appState.recentPosts.shift();
        //         appState.recentPosts.push(feeds.currentItem);
        //     }
        // } else {
        //     if (
        //         appState.recentPosts.filter((p) => p !== feeds.currentItem)
        //             .length == 0
        //     )
        //         appState.recentPosts.push(feeds.currentItem);
        // }
    });
</script>

<div class="col-span-2">
    <Fab
        actionBtns={[
            {
                text: "Add to Collection",
                action: () => {},
                icon: "fa fa-plus",
            },
            {
                text: "Share",
                action: () => {
                    window.location.href = feed.currentItem.link;
                },
                icon: "fa fa-link",
            },
            {
                text: "Back",
                action: () => {
                    if (appState.currentPage.length == 0) {
                        appState.currentPage.push("#/");
                    } else {
                        appState.currentPage.pop();
                    }
                },
                icon: "fa fa-angle-left",
            },
        ]}
    />
    {#if feed.currentItem}
        <div
            class="flex flex-col space-y-4 bg-base-300 border-0 rounded shadow-sm mx-auto my-8 max-w-[75ch]"
        >
            <div>
                <div class="flex justify-between px-4 py-2 text-sm">
                    <div>
                        <span class="text-secondary"
                            ><a href={feed.currentItem.feedName}
                                >{feed.currentItem.feedName}</a
                            ></span
                        >
                        |
                        <span class="text-base-content/70"
                            >{getTimeAgo(feed.currentItem.createdAt)}</span
                        >
                    </div>
                    <div class="flex space-x-1">
                        {#if true || false}
                            <button
                                class="btn btn-sm btn-ghost"
                                aria-label="read feed"
                                ><i class="fa fa-check"></i></button
                            >
                        {/if}
                        <button
                            class="btn btn-sm btn-ghost"
                            aria-label="bookmark"
                            onclick={async () => {
                                feed.currentItem.saved =
                                    !feed.currentItem.saved;
                                let res = await HandleSaveFeedItem({
                                    id: feed.currentItem.id,
                                    value: feed.currentItem.saved,
                                });
                                if (res.error) {
                                    feed.currentItem.saved =
                                        !feed.currentItem.saved;
                                    console.error(res.error);
                                    return;
                                }
                            }}
                            ><i
                                class={feed.currentItem.saved
                                    ? "fa fa-bookmark text-primary"
                                    : "fa-regular fa-bookmark"}
                            ></i></button
                        >
                    </div>
                </div>
                <h3 class="text-lg px-4">{feed.currentItem.title}</h3>
            </div>
            {#if feed.currentItem.image}
                <figure>
                    <img
                        src={feed.currentItem.image}
                        alt={feed.currentItem.altText || feed.currentItem.title}
                    />
                </figure>
            {/if}
            {#if feed.currentItem.link}
                <span class="px-4 text-lg"
                    >See more: <button
                        class="link max-w-[48ch] text-secondary break-all"
                        onclick={() => {
                            BrowserOpenURL(feed.currentItem.link);
                        }}
                        title={feed.currentItem.link}>Visit Site</button
                    ></span
                >
            {/if}
            <div class="max-h-48 mb-12 mt-10 overflow-auto">
                <p class="text-base-content/70 px-4">
                    {@html feed.currentItem.description}
                </p>
            </div>
        </div>
    {/if}
</div>

<style>
</style>
