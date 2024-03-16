<script lang="ts">
  import * as Card from "$lib/components/ui/card";

  type Vintage = number | "NV";
  type Wine = {
    name: string;
    vintage: Vintage;
    country: string;
    region: string;
    quantity: number;
    price: number;
    purchase_date: Date;
  };

  let wines: Wine[] = [
    {
      name: "Château Lafite Rothschild",
      vintage: 2018,
      country: "France",
      region: "Bordeaux",
      quantity: 1,
      price: 1000,
      purchase_date: new Date("2024-01-01"),
    },
    {
      name: "Albert Bichot Chassagne-Montrachet Blanc",
      vintage: 2019,
      country: "France",
      region: "Burgundy",
      quantity: 2,
      price: 200,
      purchase_date: new Date("2024-02-02"),
    },
    {
      name: "Louis Roederer",
      vintage: "NV",
      country: "France",
      region: "Champagne",
      quantity: 3,
      price: 150,
      purchase_date: new Date("2024-03-03"),
    },
  ];

  let images: string[] = [
    "https://images.vivino.com/thumbs/uDMfPG10R1efm5H1do--Ow_pb_600x600.png",
    "https://images.vivino.com/thumbs/Eh506eKdSXGvndyHHfqPug_pb_600x600.png",
    "https://images.vivino.com/thumbs/MCs7Ix2zS56U3-vs6GOR4A_pb_600x600.png",
  ];
</script>

<main>
  <h1>Wine Cellar</h1>
  {#each wines as wine, i}
  <Card.Card class="bg-card text-card-foreground">
    <Card.Content>
      <div class="grid grid-cols-[1fr,3fr] pt-4 pb-2">
        <img src={images[i]} alt={wine.name} width="256" />
        <div>
          <Card.Title class="font-normal fira-sans my-4">
            {wine.name.toUpperCase()}
          </Card.Title>
          <p class="mb-4">{wine.vintage}&#xff5c;{wine.region}&#xff5c;{wine.country}</p>
          <p class="fira-sans text-sm">PURCHASED</p>
          <p>
            <i class="fa-solid fa-bottle-droplet text-sm" /> &nbsp;{wine.quantity} bottle{wine.quantity > 1 ? 's' : ''}
          </p>
          <p>
            <i class="fa-solid fa-dollar-sign text-sm" /> &nbsp;HKD {wine.price}
          </p>
          <p>
            <i class="fa-regular fa-calendar text-sm" /> &nbsp;{wine.purchase_date.toISOString().slice(0, 10)}
          </p>
        </div>
      </div>
    </Card.Content>
  </Card.Card>
  {/each}
</main>