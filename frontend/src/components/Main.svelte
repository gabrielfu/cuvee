<script lang="ts">
  import * as Card from "$lib/components/ui/card";
  import * as Carousel from "$lib/components/ui/carousel";
  import { Label } from "$lib/components/ui/label";
  import { Switch } from "$lib/components/ui/switch";

  type Vintage = number | "NV";
  type Purchase = {
    quantity: number;
    price: number;
    date: Date;
  };
  type Wine = {
    name: string;
    vintage: Vintage;
    country: string;
    region: string;
    purchases: Purchase[];
  };

  let wines: Wine[] = [
    {
      name: "Château Lafite Rothschild",
      vintage: 2018,
      country: "France",
      region: "Bordeaux",
      purchases: [
        { quantity: 1, price: 5000, date: new Date("2024-01-01") },
      ],
    },
    {
      name: "Albert Bichot Chassagne-Montrachet Blanc",
      vintage: 2019,
      country: "France",
      region: "Burgundy",
      purchases: [
        { quantity: 2, price: 800, date: new Date("2024-02-02") },
        { quantity: 1, price: 750, date: new Date("2024-02-03") },
        { quantity: 1, price: 780, date: new Date("2024-02-04") },
        { quantity: 3, price: 720, date: new Date("2024-02-05") },
        { quantity: 1, price: 820, date: new Date("2024-02-06") },
      ],
    },
    {
      name: "Louis Roederer",
      vintage: "NV",
      country: "France",
      region: "Champagne",
      purchases: [
        { quantity: 1, price: 600, date: new Date("2024-03-03") },
        { quantity: 2, price: 560, date: new Date("2024-03-04") },
      ],
    },
  ];

  let images: string[] = [
    "https://images.vivino.com/thumbs/uDMfPG10R1efm5H1do--Ow_pb_600x600.png",
    "https://images.vivino.com/thumbs/Eh506eKdSXGvndyHHfqPug_pb_600x600.png",
    "https://images.vivino.com/thumbs/MCs7Ix2zS56U3-vs6GOR4A_pb_600x600.png",
  ];

  const aggregatePurchase = (purchases: Purchase[]) => {
    const totalQuantity = purchases.reduce((acc, p) => acc + p.quantity, 0);
    const totalPrice = purchases.reduce((acc, p) => acc + p.price * p.quantity, 0);
    return { quantity: totalQuantity, price: totalPrice / totalQuantity };
  };

  let detailedView = false;
</script>

<main>
  <p class="fira-sans font-bold text-2xl px-2">MY CELLAR </p>

  <div class="flex items-center space-x-2 my-4 px-2">
    <Switch id="detailed-view" bind:checked={detailedView} />
    <Label for="detailed-view">Purchase Details</Label>
  </div>

  {#each wines as wine, i}
  <Card.Card class="bg-card text-card-foreground">
    <Card.Content>
      <div class="grid grid-cols-[1fr,3fr] pt-4">
        <img src={images[i]} alt={wine.name} width="224" />
        <div>
          <Card.Title class="font-medium my-4">
            {wine.name.toUpperCase()}
          </Card.Title>
          <p class="mb-4">{wine.vintage}&#xff5c;{wine.region}&#xff5c;{wine.country}</p>
          <p class="text-sm mb-2">PURCHASES</p>
          {#if !detailedView}
          {@const p = aggregatePurchase(wine.purchases)}
          <Card.Card class="w-1/4 mx-2 px-4 py-2">
            <p class="text-sm">
              <i class="fa-solid fa-bottle-droplet text-sm" /> &nbsp;{p.quantity} bottle{p.quantity > 1 ? 's' : ''}
            </p>
            <p class="text-sm">
              <i class="fa-solid fa-dollar-sign text-sm" /> &nbsp;HKD {p.price.toFixed(2)}
            </p>
          </Card.Card>
          {:else}
          <Carousel.Root class="mx-12">
            <Carousel.Content>
              {#each wine.purchases as p}
              <Carousel.Item class="basis-1/4 pl-4">
                <Card.Card class="px-4 py-2">
                  <p class="text-sm">
                    <i class="fa-solid fa-bottle-droplet text-sm" /> &nbsp;{p.quantity} bottle{p.quantity > 1 ? 's' : ''}
                  </p>
                  <p class="text-sm">
                    <i class="fa-solid fa-dollar-sign text-sm" /> &nbsp;HKD {p.price}
                  </p>
                  <p class="text-sm">
                    <i class="fa-regular fa-calendar text-sm" /> &nbsp;{p.date.toISOString().slice(0, 10)}
                  </p>
                </Card.Card>
              </Carousel.Item>
              {/each}
            </Carousel.Content>
            <Carousel.Previous />
            <Carousel.Next />
          </Carousel.Root>
          {/if}
        </div>
      </div>
    </Card.Content>
  </Card.Card>
  {/each}
</main>