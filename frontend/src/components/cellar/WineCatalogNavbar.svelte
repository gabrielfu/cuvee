<script lang="ts">
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import * as Select from "$lib/components/ui/select";
  import { sortByOptions } from "$lib/sortBys";
  import WineAddDialog from "./WineAddDialog.svelte";
  import type { WineFormSchema } from "./WineFormSchema";
  import type { SuperValidated, Infer } from "sveltekit-superforms";

  let selected = { value: "default", label: sortByOptions["default"] };
  export let onSelectedChange: (e: any) => void;

  export let search = "";

  export let form: SuperValidated<Infer<WineFormSchema>>;
</script>

<div class="flex items-center justify-between">
  <Input class="max-w-sm" type="text" id="search" placeholder="Search" bind:value={search} />

  <div class="flex items-center space-x-2 my-4 px-2">
    <Label class="font-normal pr-1">Sort By</Label>
    <Select.Root bind:selected {onSelectedChange}>
      <Select.Trigger class="w-[180px] h-[40px]">
        <Select.Value class="text-sm" placeholder="" />
      </Select.Trigger>
      <Select.Content>
        {#each Object.entries(sortByOptions) as [value, label]}
          <Select.Item class="text-sm" {value}>{label}</Select.Item>
        {/each}
      </Select.Content>
    </Select.Root>
  </div>

  <WineAddDialog data={form} displayText="Add" />
</div>
