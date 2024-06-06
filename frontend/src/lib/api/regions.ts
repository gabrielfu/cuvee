import { baseUrl } from "./utils";

export type Region = {
  id: string;
  wineId: string;
  symbol: string;
  region: string;
}

export const listRegions = async (wineId: string): Promise<Region[]> => {
  const response = await fetch(`${baseUrl}/api/v1/regions/wines/${wineId}`);
  if (response.ok) {
    const data = await response.json();
    return data as Region[];
  }
  throw new Error("Failed to fetch regions");
}

export const getRegion = async (wineId: string, symbol: string): Promise<Region> => {
  const response = await fetch(`${baseUrl}/api/v1/regions/wines/${wineId}/vintage_charts/${symbol}`);
  if (response.ok) {
    const data = await response.json();
    return data as Region;
  }
  throw new Error("Failed to fetch regions");
}

export const upsertRegion = async (wineId: string, symbol: string, region: string): Promise<void> => {
  const obj = { wineId, symbol, region };
  const response = await fetch(`${baseUrl}/api/v1/regions/wines/${wineId}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(obj)
  });
  if (!response.ok) {
    const message = await response.text();
    throw new Error(`Failed to upsert region: ${message}`);
  }
}

export const deleteRegion = async (wineId: string, symbol: string): Promise<void> => {
  const response = await fetch(`${baseUrl}/api/v1/regions/wines/${wineId}/vintage_charts/${symbol}`, {
    method: "DELETE"
  });
  if (!response.ok) {
    throw new Error("Failed to delete region");
  }
}
