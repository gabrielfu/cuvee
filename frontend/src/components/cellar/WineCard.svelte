<script lang="ts">
  import * as Card from "$lib/components/ui/card";
  import * as Collapsible from "$lib/components/ui/collapsible";
  import * as Table from "$lib/components/ui/table";
  import * as Accordion from "$lib/components/ui/accordion/index.js";
  import { Label } from "$lib/components/ui/label";

  import type { Wine } from "$lib/api/wines";
  import WineCardEditDialog from "./WineCardEditDialog.svelte";
  import {
    CircleDollarSign,
    Wine as WineSign,
    Calendar,
    ChevronDown,
    ChevronUp
  } from "lucide-svelte";
  import { getRating, listVintageCharts, type RatingWithSymbol } from "$lib/api/vintageCharts";
  import { getRegion, type Region } from "$lib/api/regions";
  import RegionAddDialog from "./RegionAddDialog.svelte";

  export let wine: Wine;
  let collapsibleOpen = false;

  let ratings: RatingWithSymbol[] = [];
  function getRatings(wine: Wine) {
    ratings = [];
    listVintageCharts().then((vcs) => {
      vcs.forEach((vc) => {
        getRegion(wine.id, vc.symbol)
        .then((region: Region) => {
          if (region) {
            getRating(vc.symbol, region.region, wine.vintage).then((rating) => {
              const ratingWithSymbol = { ...rating, symbol: vc.symbol, region: region.region };
              ratings = [...ratings, ratingWithSymbol];
            })
          }
        })
        .catch((_) => {});
      });
    });
  }
  getRatings(wine);

</script>

<Card.Card class="bg-transparent border-x-transparent border-t-transparent text-card-foreground">
  <Card.Content>
    <div class="grid grid-cols-[minmax(224px,1fr),3fr,2fr] pt-4">
      <img src={wine.imageUrl} alt={wine.name} class="h-[224px] max-h-[224px] max-w-[224px]" />
      <div>
        <div class="grid grid-cols-2 items-center justify-between">
          <div>
            <Card.Title class="font-medium mt-4 mb-2">
              {wine.name}
            </Card.Title>

            <p class="text-sm text-muted-foreground mb-2">
              {wine.vintage}&#xff5c;{wine.format}&#xff5c;{wine.region}&#xff5c;{wine.country}
            </p>

            <WineCardEditDialog {wine} />
          </div>
        </div>

        <Card.Card class="bg-card w-[160px] mx-2 px-4 py-2 mb-2">
          <p class="text-sm">
            <WineSign class="inline mr-2" size="20" />{wine.summary.quantity} bottle{wine.summary
              .quantity > 1
              ? "s"
              : ""}
          </p>
          <p class="text-sm">
            <CircleDollarSign class="inline mr-2" size="20" />{wine.summary.price.toFixed(2)}
          </p>
        </Card.Card>

        <Collapsible.Root class="pl-4 pt-2" bind:open={collapsibleOpen}>
          <Collapsible.Trigger class="text-sm text-muted-foreground underline">
            {#if collapsibleOpen}
              Collapse purchases<ChevronUp class="inline" size="20" />
            {:else}
              Expand purchases<ChevronDown class="inline" size="20" />
            {/if}
          </Collapsible.Trigger>
          <Collapsible.Content>
            <Table.Root class="w-auto pt-2">
              <Table.Header>
                <Table.Row class="leading-none">
                  <Table.Head class="h-auto py-1 text-sm text-muted-foreground">
                    <WineSign class="inline mr-2" size="20" />
                  </Table.Head>
                  <Table.Head class="h-auto py-1 text-sm text-muted-foreground">
                    <CircleDollarSign class="inline mr-2" size="20" />
                  </Table.Head>
                  <Table.Head class="h-auto py-1 text-sm text-muted-foreground">
                    <Calendar class="inline mr-2" size="20" />
                  </Table.Head>
                </Table.Row>
              </Table.Header>
              <Table.Body>
                {#each [...wine.purchases].sort().reverse() as entry}
                  <Table.Row class="leading-none">
                    <Table.Cell class="py-1 text-sm">{entry.quantity}</Table.Cell>
                    <Table.Cell class="py-1 text-sm">{entry.price}</Table.Cell>
                    <Table.Cell class="py-1 text-sm">
                      {entry.date.toISOString().slice(0, 10)}
                    </Table.Cell>
                  </Table.Row>
                {/each}
              </Table.Body>
            </Table.Root>
          </Collapsible.Content>
        </Collapsible.Root>
      </div>

      <div>
        <Accordion.Root class="mb-2">
          {#each ratings as rating}
          <Accordion.Item value={rating.symbol}>
            <Accordion.Trigger class="hover:no-underline py-2">
              <Label class="rounded-sm w-[40px] text-center bg-accent shadow-sm py-1 text-xs text-accent-foreground">
                {rating.symbol}
              </Label>
              <span class={rating.score ? "" : "font-light text-sm text-muted-foreground"}>{rating.score || "No rating"}</span>
            </Accordion.Trigger>
            <Accordion.Content>
              <div class="mb-2">
                <Label>Region</Label>
                <p class="text-sm text-muted-foreground">{rating.region}</p>
                {#if rating.maturity}
                  <Label>Maturity</Label>
                  <p class="text-sm text-muted-foreground">{rating.maturity}</p>
                {/if}
                {#if rating.notes}
                  <Label>Notes</Label>
                  <p class="text-sm text-muted-foreground">{rating.notes}</p>
                {/if}
              </div>
              <Label class="text-sm text-muted-foreground font-light underline">Delete</Label>
            </Accordion.Content>
          </Accordion.Item>
          {/each}
        </Accordion.Root> 
        <RegionAddDialog {wine} />
      </div>
    </div>
  </Card.Content>
</Card.Card>
