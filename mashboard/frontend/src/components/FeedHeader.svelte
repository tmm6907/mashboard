<script>
    let props = $props();
</script>

<div class="flex items-center justify-between">
    <h2 class="text-2xl font-black">{props.header || "Feed"}</h2>
    <div class="flex items-center space-x-1">
        <button class="btn btn-sm btn-ghost sort" aria-label="Sort"
            ><i class="fa fa-sort"></i></button
        >
        <button class="btn btn-sm btn-ghost filter" aria-label="Filter"
            ><i class="fa fa-filter"></i></button
        >
    </div>
</div>
