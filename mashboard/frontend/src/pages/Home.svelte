<script>
    import Feed from "src/components/Feed.svelte";
    import SidePanel from "src/components/SidePanel.svelte";
    import { GetFeedItems } from "../../wailsjs/go/main/App.js";
</script>

<Feed
    header={"Home"}
    request={(offset = 0, category = undefined, saved = undefined) => {
        return GetFeedItems({ offset, category, saved });
    }}
/>
<SidePanel />
