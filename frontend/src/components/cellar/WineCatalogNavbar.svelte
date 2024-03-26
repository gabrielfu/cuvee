<script lang="ts">
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import * as Select from "$lib/components/ui/select";
	import type { Selected } from "bits-ui";

  const sortBys = {
    "default": "Default",
    "name_asc": "Name (asc)",
    "name_desc": "Name (desc)",
    "vintage_asc": "Vintage (asc)",
    "vintage_desc": "Vintage (desc)",
    "price_asc": "Price (asc)",
    "price_desc": "Price (desc)",
    "newest": "Newest",
    "oldest": "Oldest",
  };
  let selected = { value: "default", label: sortBys["default"] };

  export let onSelectedChange: (e: any) => void;
</script>

<div class="flex items-center justify-between">
  <Input class="max-w-sm h-[40px]" type="text" id="search" placeholder="Search" />

  <div class="flex items-center space-x-2 my-4 px-2">
    <Label class="font-normal pr-1">Sort By</Label>
    <Select.Root bind:selected={selected} onSelectedChange={onSelectedChange}>
      <Select.Trigger class="w-[180px] h-[40px]">
        <Select.Value class="text-sm" placeholder="" />
      </Select.Trigger>
      <Select.Content>
        {#each Object.entries(sortBys) as [value, label]}
        <Select.Item class="text-sm" value={value}>{label}</Select.Item>
        {/each}
      </Select.Content>
    </Select.Root>
  </div>
</div>