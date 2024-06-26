import { baseUrl } from "./utils";

export type PurchaseSummary = {
  quantity: number;
  price: number;
};

export type Purchase = {
  quantity: number;
  price: number;
  date: Date;
};

export type Wine = {
  id: string;
  name: string;
  vintage: string;
  format: string;
  country: string;
  region: string;
  purchases: Purchase[];
  summary: PurchaseSummary;
  imageUrl: string;
};

export const listWines = async (): Promise<Wine[]> => {
  // return fakeWines;
  // return [];
  const response = await fetch(`${baseUrl}/api/v1/wines`);
  if (response.ok) {
    const data = await response.json();
    // case dates
    data.forEach((wine: Wine) => {
      wine.purchases.forEach((purchase: Purchase) => {
        purchase.date = new Date(purchase.date);
      });
    });
    return data;
  }
  throw new Error("Failed to fetch wines");
};

export const getWine = async (id: string): Promise<Wine> => {
  const response = await fetch(`${baseUrl}/api/v1/wines/${id}`);
  if (response.ok) {
    return response.json();
  }
  throw new Error("Failed to fetch wine");
};

export const deleteWine = async (id: string): Promise<void> => {
  const response = await fetch(`${baseUrl}/api/v1/wines/${id}`, {
    method: "DELETE"
  });
  if (!response.ok) {
    throw new Error("Failed to delete wine");
  }
}

export const fakeWines: Wine[] = [
  {
    id: "1",
    name: "Château Lafite Rothschild",
    vintage: "2018",
    format: "750ml",
    country: "France",
    region: "Bordeaux",
    purchases: [{ quantity: 1, price: 5000, date: new Date("2024-01-01") }],
    summary: { quantity: 1, price: 5000 },
    imageUrl: "https://images.vivino.com/thumbs/uDMfPG10R1efm5H1do--Ow_pb_600x600.png"
  },
  {
    id: "2",
    name: "Albert Bichot Chassagne-Montrachet Blanc",
    vintage: "2019",
    format: "750ml",
    country: "France",
    region: "Burgundy",
    purchases: [
      { quantity: 2, price: 800, date: new Date("2024-02-02") },
      { quantity: 1, price: 750, date: new Date("2024-02-03") },
      { quantity: 1, price: 780, date: new Date("2024-02-04") },
      { quantity: 3, price: 720, date: new Date("2024-02-05") },
      { quantity: 1, price: 820, date: new Date("2024-02-06") }
    ],
    summary: { quantity: 8, price: 780 },
    imageUrl: "https://images.vivino.com/thumbs/Eh506eKdSXGvndyHHfqPug_pb_600x600.png"
  },
  {
    id: "3",
    name: "Louis Roederer",
    vintage: "NV",
    format: "750ml",
    country: "France",
    region: "Champagne",
    purchases: [
      { quantity: 1, price: 600, date: new Date("2024-03-03") },
      { quantity: 2, price: 560, date: new Date("2024-03-04") }
    ],
    summary: { quantity: 3, price: 573.33 },
    imageUrl: "https://images.vivino.com/thumbs/MCs7Ix2zS56U3-vs6GOR4A_pb_600x600.png"
  }
];
