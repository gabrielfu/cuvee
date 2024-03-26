<script lang="ts">
  import type { Wine } from "$lib/api/wines";
	import type { Selected } from "bits-ui";
	import WineCard from "./WineCard.svelte";
	import WineCatalogNavbar from "./WineCatalogNavbar.svelte";
	import { sortByOptions } from "$lib/sortBys";
	import { derived } from "svelte/store";

  export let wines: Wine[];
  export let images: string[];

  $: filteredWines = wines;

  let onSelectedChange = (e: Selected<keyof typeof sortByOptions>) => {
    switch (e.value) {
      case "name_asc":
        filteredWines = wines.sort((a, b) => a.name.localeCompare(b.name));
        break;
      case "name_desc":
        filteredWines = wines.sort((a, b) => b.name.localeCompare(a.name));
        break;
      case "vintage_asc":
        filteredWines = wines.sort((a, b) => a.vintage.localeCompare(b.vintage));
        break
      case "vintage_desc":
        filteredWines = wines.sort((a, b) => b.vintage.localeCompare(a.vintage));
        break
      case "price_asc":
        filteredWines = wines.sort((a, b) => a.summary.price - b.summary.price);
        break;
      case "price_desc":
        filteredWines = wines.sort((a, b) => b.summary.price - a.summary.price);
        break;
      case "newest":
        filteredWines = wines.sort((a, b) => {
          // last purchase
          let pa = [...a.purchases].sort((i, j) => j.date.getTime() - i.date.getTime())[0];
          let pb = [...b.purchases].sort((i, j) => j.date.getTime() - i.date.getTime())[0];
          // by descending date
          return pb.date.getTime() - pa.date.getTime();
        });
        break;
      case "oldest":
        filteredWines = wines.sort((a, b) => {
          // first purchase
          let pa = [...a.purchases].sort((i, j) => i.date.getTime() - j.date.getTime())[0];
          let pb = [...b.purchases].sort((i, j) => i.date.getTime() - j.date.getTime())[0];
          // by ascending date
          return pa.date.getTime() - pb.date.getTime();
        });
        break;
      default:
        filteredWines = wines.sort((a, b) => a.id.localeCompare(b.id));
        break;
    }
  };

  let search: string;
  $: if (search) {
    filteredWines = wines.filter(wine => wine.name.normalize("NFD").replace(/\p{Diacritic}/gu, "").toLowerCase().includes(search.normalize("NFD").replace(/\p{Diacritic}/gu, "").toLowerCase()));
  } else {
    filteredWines = wines;
  }
</script>

<main>
  <h1 class="fira-sans font-bold text-3xl px-2">MY CELLAR </h1>

  <WineCatalogNavbar onSelectedChange={onSelectedChange} bind:search={search} />

  {#each filteredWines as wine, i}
  <WineCard wine={wine} image={images[i]} />
  {/each}
</main>