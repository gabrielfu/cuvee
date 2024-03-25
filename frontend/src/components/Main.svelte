<script lang="ts">
	import type { Wine } from "$lib/api/wines";
  import * as Card from "$lib/components/ui/card";
  import * as Collapsible from "$lib/components/ui/collapsible";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import * as Select from "$lib/components/ui/select";
  import * as Table from "$lib/components/ui/table";

  export let wines: Wine[];
  export let images: string[];

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
  let sortBy = { value: "default", label: sortBys["default"] };
</script>

<main>
  <p class="fira-sans font-bold text-3xl px-2">MY CELLAR </p>

  <div class="flex items-center justify-between">

    <Input class="max-w-sm h-[36px]" type="text" id="search" placeholder="Search" />

    <div class="flex items-center space-x-2 my-4 px-2">
      <Label class="font-normal pr-1">Sort By</Label>
      <Select.Root selected={sortBy}>
        <Select.Trigger class="w-[180px] h-[36px]">
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

  {#each wines as wine, i}
  <Card.Card class="bg-transparent border-x-transparent border-t-transparent text-card-foreground">
    <Card.Content>
      <div class="grid grid-cols-[minmax(224px,1fr),3fr,2fr] pt-4">
        <img src={images[i]} alt={wine.name} width="224" />
        <div>
          <Card.Title class="font-medium mt-4 mb-2">
            {wine.name}
          </Card.Title>

          <p class="text-sm text-muted-foreground mb-4">{wine.vintage}&#xff5c;{wine.format}&#xff5c;{wine.region}&#xff5c;{wine.country}</p>

          <Card.Card class="bg-card w-[160px] mx-2 px-4 py-2 mb-2">
            <p class="text-sm">
              <i class="fa-solid fa-bottle-droplet text-sm" /> &nbsp;{wine.summary.quantity} bottle{wine.summary.quantity > 1 ? 's' : ''}
            </p>
            <p class="text-sm">
              <i class="fa-solid fa-dollar-sign text-sm" /> &nbsp;{wine.summary.price.toFixed(2)}
            </p>
          </Card.Card>

          <Collapsible.Root class="pl-4 pt-2">
            <Collapsible.Trigger class="text-sm text-muted-foreground">Expand purchases <i class="fa-solid fa-chevron-down text-sm" /></Collapsible.Trigger>
            <Collapsible.Content>
              <Table.Root class="w-auto pt-2">
                <Table.Header>
                  <Table.Row class="leading-none">
                    <Table.Head class="h-auto py-1 text-sm text-muted-foreground"><i class="fa-solid fa-bottle-droplet text-sm" /></Table.Head>
                    <Table.Head class="h-auto py-1 text-sm text-muted-foreground"><i class="fa-solid fa-dollar-sign text-sm" /></Table.Head>
                    <Table.Head class="h-auto py-1 text-sm text-muted-foreground"><i class="fa-regular fa-calendar text-sm" /></Table.Head>
                  </Table.Row>
                </Table.Header>
                <Table.Body>
                  {#each [...wine.purchases].sort().reverse() as entry}
                  <Table.Row class="leading-none">
                    <Table.Cell class="py-1 text-sm">{entry.quantity}</Table.Cell>
                    <Table.Cell class="py-1 text-sm">{entry.price}</Table.Cell>
                    <Table.Cell class="py-1 text-sm">{entry.date.toISOString().slice(0, 10)}</Table.Cell>
                  </Table.Row>
                  {/each}
                </Table.Body>
              </Table.Root>
            </Collapsible.Content>
          </Collapsible.Root>
        </div>
        <div>

        </div>
      </div>
    </Card.Content>
  </Card.Card>
  {/each}
</main>