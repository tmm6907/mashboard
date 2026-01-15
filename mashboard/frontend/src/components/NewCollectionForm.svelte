<script lang="ts">
    import { appState } from "src/stores/appState.svelte.ts";
    import { CreateNewCollection } from "../../wailsjs/go/main/App.js";

    function resetForm() {
        const nameField = document.getElementById("new-collection-name");
        nameField.textContent = "";
        appState.showNewCollectionForm = false;
    }

    async function createNewCollection(e: SubmitEvent) {
        e.preventDefault();
        let form = e.target as HTMLFormElement;
        if (!form) console.error("couldn't find form");
        let data = new FormData(form);
        let res = await CreateNewCollection(
            `${data.get("new-collection-name")}`,
        );
        if (res.error) {
            console.error(res.error);
            return;
        }
        resetForm();
    }
</script>

{#if appState.showNewCollectionForm}
    <dialog
        class="bg-base-300 min-w-[32ch] rounded border border-primary p-4"
        open={appState.showNewCollectionForm}
    >
        <form onsubmit={async (e) => await createNewCollection(e)}>
            <h4 class="mb-4">Create new collection</h4>
            <label class="label mb-4" for="new-collection-name"
                >Name <input
                    id="new-collection-name"
                    class="input"
                    name="new-collection-name"
                    type="text"
                /></label
            >
            <div class="flex space-x-4 justify-center">
                <button
                    class="btn btn-sm bg-neutral-200 text-black"
                    type="button"
                    onclick={() => resetForm()}
                >
                    Cancel
                </button>
                <button type="submit" class="btn btn-sm btn-primary"
                    >Create</button
                >
            </div>
        </form>
    </dialog>
{/if}
