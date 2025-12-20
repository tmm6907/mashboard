<script lang="ts">
    import { appState } from "src/stores/appState.svelte.ts";
    import Search from "./Search.svelte";
    import { collections } from "src/stores/collections.svelte.ts";
    import { searchField } from "src/stores/search.svelte.ts";
    import {
        SearchForFeed,
        GetCollections,
    } from "../../wailsjs/go/main/App.js";
    import type { SearchResult } from "src/lib/types.ts";
    import NewCollectionForm from "./NewCollectionForm.svelte";
    import { onMount } from "svelte";

    $effect(() => {
        if (searchField.text.length > 0) {
            SearchForFeed(searchField.text).then((res) => {
                if (res.error) {
                    console.error(res.error);
                    return;
                }
                let data: SearchResult[] = res.data;
                console.log("search results", data);
                searchField.results = data;
            });
        }
    });

    onMount(async () => {
        let res = await GetCollections();
        if (res.error) {
            console.error(res.error);
            return;
        }
        collections.items = res.data;
        console.log(collections.items);
    });
</script>

<div class="bg-base-300 box w-full">
    <div class="mashboard-logo"></div>
    <ul class="my-4 px-4">
        <li><Search id="feed-search" placeholder="Search Feeds" /></li>
    </ul>
    <ul class="menu mb-2 w-full">
        <li>
            <a href="#/" onclick={() => appState.currentPage.push("#/")}>
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-5 w-5"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
                    />
                </svg>
                Home
            </a>
        </li>
        <li>
            <a
                href="#saved/"
                onclick={() => appState.currentPage.push("#saved/")}
            >
                <i class="fa-solid fa-bookmark"></i>
                Saved
            </a>
        </li>
        <li>
            <a
                href="#explore/"
                onclick={() => appState.currentPage.push("#explore/")}
            >
                <i class="fa-solid fa-compass"></i>
                Explore
            </a>
        </li>
    </ul>
    <div class="relative flex items-center">
        <h3 class="px-4">Collections</h3>
        <button
            class="btn btn-sm btn-ghost"
            aria-label="New collection"
            onclick={() => (appState.showNewCollectionForm = true)}
            title="New collection"><i class="fa fa-add"></i></button
        >
        <div class="fixed top-1/3 left-1/2 transform -translate-x-1/2 z-10">
            <NewCollectionForm />
        </div>
    </div>
    <ul class="menu w-full">
        {#each collections.items as item (item.id)}
            <li>
                <a
                    href={`#collections/${item.name.toLowerCase()}`}
                    onclick={() =>
                        appState.currentPage.push(
                            `#collections/${item.name.toLowerCase()}`,
                        )}
                >
                    <i class="fa fa-file"></i>
                    {item.name}
                    <!-- {#if item.count > 0}
                        {item.count}
                    {/if} -->
                </a>
            </li>
        {/each}
    </ul>
</div>

<style>
    .mashboard-logo {
        width: auto;
        height: 48px;
        margin: 1rem;
        background-image: url("../assets/images/mashboard.svg");
        background-repeat: no-repeat;
        background-size: contain;
        background-position: start;
    }
</style>
