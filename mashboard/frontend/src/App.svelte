<script lang="ts">
    import Nav from "./components/Nav.svelte";
    import Collections from "./pages/Collections.svelte";
    import Explore from "./pages/Explore.svelte";
    import FeedItem from "./pages/FeedItem.svelte";
    import Home from "./pages/Home.svelte";
    import Saved from "./pages/Saved.svelte";
    import { appState } from "./stores/appState.svelte.ts";
    import { GetFeedItem } from "../wailsjs/go/main/App.js";
    import { feed } from "./stores/feeds.svelte.ts";

    $effect(() => {
        console.log(
            $state.snapshot(
                appState.currentPage[appState.currentPage.length - 1],
            ),
        );
        if (appState.currentPage.length == 0) return;
        if (
            appState.currentPage[appState.currentPage.length - 1].startsWith(
                "#feeds/items/",
            )
        ) {
            let dynID = appState.currentPage[
                appState.currentPage.length - 1
            ].replace("#feeds/items/", "");
            if (dynID !== "") {
                console.log("Dynamic ID init", dynID);
                appState.dynamicID = Number(dynID);
                let req = { id: appState.dynamicID, saved: false };
                console.log("request ", req);
                GetFeedItem(req).then((res) => {
                    if (res.error) {
                        console.error(res.error);
                        return;
                    }
                    let data = res.data;
                    feed.currentItem = data;
                    if (
                        !appState.recentPosts.some(
                            (p) => p.feedId === feed.currentItem.feedId,
                        )
                    ) {
                        if (appState.recentPosts.length >= 5) {
                            appState.recentPosts.shift();
                        }
                        appState.recentPosts.push(feed.currentItem);
                    }
                });
            }
        } else if (
            appState.currentPage[appState.currentPage.length - 1].startsWith(
                "#feeds/",
            )
        ) {
            let dynID = appState.currentPage[
                appState.currentPage.length - 1
            ].replace("#feeds/", "");
            if (dynID !== "") {
                appState.dynamicID = Number(dynID);
            }
        } else if (
            appState.currentPage[appState.currentPage.length - 1].startsWith(
                "#collections/",
            )
        ) {
            let collectionType = appState.currentPage[
                appState.currentPage.length - 1
            ].replace("#collections/", "");
            if (collectionType !== "") {
                appState.collectionType = collectionType;
            }
        } else {
            appState.dynamicID = null;
            appState.collectionType = "";
        }
    });
</script>

<main class="main-content relative md h-dvh">
    <Nav />
    {#if appState.currentPage[appState.currentPage.length - 1] == "#/"}
        <Home />
    {:else if appState.currentPage[appState.currentPage.length - 1] == "#saved/"}
        <Saved />
    {:else if appState.currentPage[appState.currentPage.length - 1] == "#explore/"}
        <Explore />
    {:else if appState.currentPage[appState.currentPage.length - 1].startsWith("#feeds/items/")}
        <FeedItem />
    {:else if appState.currentPage[appState.currentPage.length - 1].startsWith("#collections/items/")}
        <Collections />
    {:else}
        <div>
            PAGE NOT FOUND {appState.currentPage[
                appState.currentPage.length - 1
            ]}
        </div>
    {/if}
</main>

<style>
    .main-content {
        display: grid;
    }
    @media (max-width: 64rem) {
        .main-content {
            grid-template-columns: 1fr;
        }
    }

    @media (min-width: 64rem) {
        .main-content {
            grid-template-columns: 1fr 2fr 1fr;
        }
    }

    /*@media (min-width: 64rem) {
        .main-content {
            grid-template-columns: 1fr 2fr 1fr;
        }
    }*/
</style>
