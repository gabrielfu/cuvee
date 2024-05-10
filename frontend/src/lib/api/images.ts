import { baseUrl } from "./utils";

export type ImageResult = {
    link: string;
    height: number;
    width: number;
};

export const searchImages = async (name: string, vintage: string, country: string, region: string): Promise<ImageResult[]> => {
    let url = `${baseUrl}/images/search`;
    url += `?name=${encodeURIComponent(name)}`;
    url += `&vintage=${encodeURIComponent(vintage)}`;
    url += `&country=${encodeURIComponent(country)}`;
    url += `&region=${encodeURIComponent(region)}`;
    const response = await fetch(url);
    if (response.ok) {
        const result = await response.json();
        return result.items;
    }
    throw new Error('Failed to fetch images');
}
