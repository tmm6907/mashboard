<script>
    import RecentlyViewed from "./RecentlyViewed.svelte";
    import SuggestedFeeds from "./SuggestedFeeds.svelte";
</script>

<section
    class="fixed top-4 right-4 flex-col space-y-8 hidden w-[calc(25vw-1rem)] 2xl:block overflow-hidden"
    aria-label="side panel"
>
    <RecentlyViewed />
    <SuggestedFeeds />
</section>
