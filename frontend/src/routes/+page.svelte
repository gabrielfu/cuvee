<script lang="ts">
  import WineCatalog from "../components/cellar/WineCatalog.svelte";

  import type { PageData } from './$types';
  export let data: PageData;
</script>

<svelte:head>
  <title>cuvée</title>
</svelte:head>

<WineCatalog wines={data.wines} images={data.images} form={data.form}/>
