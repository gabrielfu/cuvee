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
    const result = await response.json();
    if (response.ok) {
        const items = result.items;
        if (items != null && items.length > 0) {
            return items;
        }
        throw new Error("No images found");
    }
    throw new Error(`Failed to fetch images: ${result.error}`);
}
