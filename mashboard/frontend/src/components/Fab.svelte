<script lang="ts">
    import type { ActionBtn } from "src/lib/types.ts";

    const colors = [
        "btn-primary",
        "btn-secondary",
        "btn-accent",
        "btn-warning",
        "btn-error",
        "btn-error",
        "btn-info",
        "btn-success-content",
        "btn-info-content",
        "btn-purple",
    ];
    let props = $props();
    let actionBtns: ActionBtn[] = $state();
    actionBtns = props.actionBtns;
    let isOpen = $state(false);

    function toggleFab() {
        isOpen = !isOpen;
    }

    function runAction(btn: ActionBtn) {
        isOpen = false;
        btn.action();
        console.log("acted");
    }
    function randColor() {
        return colors[Math.floor(Math.random() * colors.length)];
    }
</script>

<div class="fixed bottom-8 right-8">
    <!-- Main FAB Button -->
    <button
        onclick={toggleFab}
        title="Actions"
        class="w-16 h-16 btn btn-circle btn-secondary text-primary font-bold text-lg shadow-lg flex items-center justify-center transition-transform duration-300 z-50"
    >
        <i class="fa-solid fa-bolt"></i>
    </button>

    <!-- Speed Dial Buttons (open upwards) -->
    {#if isOpen}
        <div
            class="absolute bottom-20 right-0 flex flex-col gap-4 pointer-events-none"
        >
            {#each actionBtns as btn}
                <div
                    class="flex items-center gap-3 fab-speed-dial pointer-events-auto"
                >
                    <button
                        onclick={() => {
                            console.log("click");
                            runAction(btn);
                        }}
                        aria-label={btn.text}
                        class={`w-14 h-14 relative rounded-full btn btn-circle ${randColor()} font-bold shadow-lg flex items-center justify-center transition-all duration-200`}
                    >
                        <span
                            class="absolute top-auto -left-[calc(12ch+1rem)] w-[12ch] bg-base-300 text-base-content px-4 py-2 z-20 rounded text-sm"
                        >
                            {btn.text}
                        </span>
                        <i class={btn.icon ? btn.icon : ""}></i>
                    </button>
                </div>
            {/each}
        </div>
    {/if}
</div>

<style>
</style>
