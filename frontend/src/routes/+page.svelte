<script lang="ts">
  import Main from "../components/Main.svelte";

  import type { PageData } from './$types';
  export let data: PageData;
</script>

<svelte:head>
  <title>cuvée</title>
</svelte:head>

<Main wines={data.wines} images={data.images}/>
