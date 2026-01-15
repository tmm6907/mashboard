<script>
  import { searchField } from "src/stores/search.svelte.ts";

  let props = $props();
</script>

<div class="relative w-[32ch]">
  <label class="input rounded-none" for={props.id}>
    <svg
      class="h-[1em] opacity-50"
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
    >
      <g
        stroke-linejoin="round"
        stroke-linecap="round"
        stroke-width="2.5"
        fill="none"
        stroke="currentColor"
      >
        <circle cx="11" cy="11" r="8"></circle>
        <path d="m21 21-4.3-4.3"></path>
      </g>
    </svg>

    <div class="w-full">
      <input
        id={props.id}
        type="text"
        placeholder={props.placeholder || "Search here"}
        class=""
        bind:value={searchField.text}
      />
    </div>
  </label>
  <div
    class="fixed -bottom-12 left-0 flex flex-col w-full max-h-24 z-30 bg-base-300"
  >
    <ul>
      {#if searchField.results}
        {#each searchField.results as r}
          <li>
            <img src={r.image} alt="" />
            <span>{r.title}</span>
          </li>
        {/each}
      {:else}
        <li>
          Add {searchField.text}
        </li>
      {/if}
    </ul>
  </div>
</div>

<style></style>
