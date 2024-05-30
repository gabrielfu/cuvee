import { baseUrl } from "./utils";
import type { Wine } from "./wines";

export type VintageChart = {
  name: string;
  symbol: string;
}

export type Rating = {
  region: string;
  vintage: string;
  score: string;
  maturity: string;
  notes: string;
};

export type RatingWithSymbol = Rating & {
  symbol: string;
}

export type SuggestRegionResponse = {
  region: string;
}

export const listVintageCharts = async (): Promise<VintageChart[]> => {
  const response = await fetch(`${baseUrl}/api/v1/vintage_charts`);
  if (response.ok) {
    const data = await response.json();
    return data as VintageChart[];
  }
  throw new Error("Failed to fetch vintage charts");
}

export const listVintageChartRegions = async (vc: string): Promise<string[]> => {
  const response = await fetch(`${baseUrl}/api/v1/vintage_charts/${vc}/regions`);
  if (response.ok) {
    const data = await response.json();
    return data as string[];
  }
  throw new Error("Failed to fetch regions");
}

export const getRating = async (vc: string, region: string, vintage: string): Promise<Rating> => {
  const response = await fetch(`${baseUrl}/api/v1/vintage_charts/${vc}/ratings?region=${region}&vintage=${vintage}`);
  if (response.ok) {
    const data = await response.json();
    return data as Rating;
  }
  throw new Error("Failed to fetch wines");
}

export const suggestRegion = async (vc: string, wine: Wine): Promise<SuggestRegionResponse> => {
  const response = await fetch(`${baseUrl}/api/v1/vintage_charts/${vc}/suggest`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      name: wine.name,
      vintage: wine.vintage,
      country: wine.country,
      region: wine.region,
    }),
  });
  if (response.ok) {
    const data = await response.json();
    return data as SuggestRegionResponse;
  }
  throw new Error("Failed to fetch region suggestion");
}
