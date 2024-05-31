<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import * as Dialog from "$lib/components/ui/dialog";
  import * as Select from "$lib/components/ui/select";
  import Label from "$lib/components/ui/label/label.svelte";
  import { WandSparkles } from "lucide-svelte";
  
  import { suggestRegion, listVintageCharts, listVintageChartRegions } from "$lib/api/vintageCharts";
  import type { Wine } from "$lib/api/wines";
  import type { Selected } from "bits-ui";
    import { getRegion } from "$lib/api/regions";

  export let wine: Wine;
  let vcs: string[] = [];
  let allRegions: Record<string, string[]> = {};
  let selectedRegions: Record<string, Selected<string>> = {};
  listVintageCharts().then((_vcs) => {
    vcs = _vcs.map((vc) => vc.symbol);
    _vcs.forEach((vc) => {
      listVintageChartRegions(vc.symbol).then((regions) => {
        allRegions[vc.symbol] = regions;
      });

      getRegion(wine.id, vc.symbol).then((region) => {
        selectedRegions[vc.symbol] = {value: region.region, label: region.region};
      }).catch((_) => {
        selectedRegions[vc.symbol] = {value: "", label: ""};
      });
    });
  });

  function suggest(vc: string, wine: Wine) {
    suggestRegion(vc, wine).then((resp) => {
      selectedRegions[vc] = {value: resp.region, label: resp.region};
    }).catch((error) => {
      console.error(error);
    })
  }
</script>

<Dialog.Root>
  <Dialog.Trigger class="text-sm text-muted-foreground font-light underline mb-4">
    Add Rating
  </Dialog.Trigger>
  <Dialog.Content class="sm:max-w-[600px] p-8">

    <!-- <form method="POST" id="rating-form" action="?/updateRatings"> -->

    {#each vcs as vc}
    <div class="grid grid-flow-col items-center">
      <Label>{vc}</Label>
      <Select.Root bind:selected={selectedRegions[vc]}>
        <Select.Trigger class="w-[425px] h-8">
          <Select.Value placeholder="Pick a region" />
        </Select.Trigger>
        <Select.Content class="overflow-y-scroll max-h-[10rem]">
          {#each allRegions[vc] as region}
            <Select.Item value={region}>{region}</Select.Item>
          {/each}
        </Select.Content>
      </Select.Root>

      <div 
        class="hover:cursor-pointer" 
        data-tooltip="Ask AI to suggest a suitable region" 
        on:click={() => suggest(vc, wine)}
      >
        <WandSparkles class="inline ml-2" size="20" />
      </div>
    </div>
    {/each}

    <Button type="submit" class="mt-4 h-8 max-w-[96px]" on:click={() => {
      console.log(selectedRegions);
    }}>Save</Button>
    <!-- </form> -->

  </Dialog.Content>
</Dialog.Root>