<script lang="ts">
    import type { FeedItemData } from "src/stores/feeds.svelte.ts";

    let props = $props();
    let item: FeedItemData = props.item;
</script>

{#if item}
    <li class="list-row w-full">
        <figure>
            <img
                src={item.image ||
                    "https://img.daisyui.com/images/stock/photo-1606107557195-0e29a4b5b4aa.webp   "}
                class="w-12 h-12"
                alt={item.altText || item.title}
            />
        </figure>
        <div>
            <span class="text-secondary">
                <a href={item.feedName} class="line-clamp-2">{item.feedName}</a>
            </span>
            <h4 class="line-clamp-2">{item.title}</h4>
        </div>
        <div>
            <button class="btn btn-sm btn-ghost" aria-label="bookmark">
                <i
                    class={item.saved
                        ? "fa fa-bookmark text-primary"
                        : "fa-regular fa-bookmark"}
                ></i>
            </button>
        </div>
    </li>
{/if}

<style>
</style>
