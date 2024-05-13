import { baseUrl } from "./utils";

export type ImageResult = {
  link: string;
  height: number;
  width: number;
};

export const searchImages = async (
  name: string,
  vintage: string,
  country: string,
  region: string
): Promise<ImageResult[]> => {
  return [
    {
      link: "https://images.vivino.com/thumbs/DSTc4YH6SnanmH_f3jXjog_pb_x600.png",
      height: 600,
      width: 175
    },
    {
      link: "https://images.vivino.com/thumbs/5NkgHoNcTDueYur5LMxhZA_pb_x600.png",
      height: 600,
      width: 180
    },
    {
      link: "https://preview.redd.it/attis-mar-albari%C3%B1o-good-or-gimmick-v0-k1zbf8bk0cza1.jpg?width=1080\u0026crop=smart\u0026auto=webp\u0026s=a936097fe1d5540c0aa3a978db28a5a3ba640ce1",
      height: 2053,
      width: 1080
    },
    {
      link: "https://images.vivino.com/thumbs/TA2gJ9Y4Sbm6aOD5eRGGpQ_pb_x600.png",
      height: 600,
      width: 179
    },
    {
      link: "https://thumbs.vivino.com/wineries/16843/logos/1Byj999XSNeQ0gqic_kIDA_500x500.jpeg",
      height: 500,
      width: 500
    },
    {
      link: "https://images.vivino.com/thumbs/wmKYr_dURYSUbZYQgMHqCA_pb_x600.png",
      height: 600,
      width: 175
    },
    {
      link: "https://preview.redd.it/attis-mar-albari%C3%B1o-good-or-gimmick-v0-k1zbf8bk0cza1.jpg?width=640\u0026crop=smart\u0026auto=webp\u0026s=7cdc0f19aeecc3757252ba867ecd04c27685b2a0",
      height: 1216,
      width: 640
    },
    {
      link: "https://images.vivino.com/thumbs/x8wYCC5STbKx5O_ujOOiDg_pb_600x600.png",
      height: 600,
      width: 600
    }
  ];

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
};
