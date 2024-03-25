const baseUrl = 'http://localhost:8080';

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
};

export const listWines = async (): Promise<Wine[]> => {
    //
    return fakeWines
    //
    // const response = await fetch(`${baseUrl}/wines`);
    // if (response.ok) {
    //     return response.json();
    // }
    // throw new Error('Failed to fetch wines');
}

export const getWine = async (id: string): Promise<Wine> => {
    const response = await fetch(`${baseUrl}/wines/${id}`);
    if (response.ok) {
        return response.json();
    }
    throw new Error('Failed to fetch wine');
}

export const fakeWines: Wine[] = [
    {
      id: "1",
      name: "Château Lafite Rothschild",
      vintage: "2018",
      format: "750ml",
      country: "France",
      region: "Bordeaux",
      purchases: [
        { quantity: 1, price: 5000, date: new Date("2024-01-01") },
      ],
      summary: { quantity: 1, price: 5000 },
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
        { quantity: 1, price: 820, date: new Date("2024-02-06") },
      ],
      summary: { quantity: 8, price: 780 },
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
        { quantity: 2, price: 560, date: new Date("2024-03-04") },
      ],
      summary: { quantity: 3, price: 573.33 },
    },
  ];