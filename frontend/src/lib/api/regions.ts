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
