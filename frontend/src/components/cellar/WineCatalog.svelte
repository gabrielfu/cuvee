<script lang="ts">
  import type { Wine } from "$lib/api/wines";
	import WineCard from "./WineCard.svelte";
	import WineCatalogNavbar from "./WineCatalogNavbar.svelte";

  export let wines: Wine[];
  export let images: string[];
</script>

<main>
  <h1 class="fira-sans font-bold text-3xl px-2">MY CELLAR </h1>

  <WineCatalogNavbar />

  {#each wines as wine, i}
  <WineCard wine={wine} image={images[i]} />
  {/each}
</main>