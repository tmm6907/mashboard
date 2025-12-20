<script lang="ts">
    import Feed from "src/components/Feed.svelte";
    import SidePanel from "src/components/SidePanel.svelte";
    import { GetFeedItems } from "../../wailsjs/go/main/App.js";
</script>

<Feed
    header={"Saved"}
    saved={true}
    request={(offset: number, category: string = "", saved: boolean) => {
        return GetFeedItems({ offset, category, saved });
    }}
/>
<SidePanel />
