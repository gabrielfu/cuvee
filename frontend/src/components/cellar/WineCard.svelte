<script lang="ts">
  import * as Card from "$lib/components/ui/card";
  import * as Collapsible from "$lib/components/ui/collapsible";
  import * as Table from "$lib/components/ui/table";

	import type { Wine } from "$lib/api/wines";

  export let wine: Wine;
  export let image: string;
</script>

<Card.Card class="bg-transparent border-x-transparent border-t-transparent text-card-foreground">
  <Card.Content>
    <div class="grid grid-cols-[minmax(224px,1fr),3fr,2fr] pt-4">
      <img src={image} alt={wine.name} width="224" />
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