import { baseUrl } from "./utils";

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

export const listVintageCharts = async (): Promise<VintageChart[]> => {
  const response = await fetch(`${baseUrl}/vintage_charts`);
  if (response.ok) {
    const data = await response.json();
    return data as VintageChart[];
  }
  throw new Error("Failed to fetch vintage charts");
}

export const getRating = async (vc: string, region: string, vintage: string): Promise<Rating> => {
  const response = await fetch(`${baseUrl}/vintage_charts/${vc}/ratings?region=${region}&vintage=${vintage}`);
  if (response.ok) {
    const data = await response.json();
    return data as Rating;
  }
  throw new Error("Failed to fetch wines");
}
