import { z } from "zod";

function validateVintage(vintage: string): boolean {
  return vintage == "NV" || /^\d{4}$/.test(vintage);
}

export const wineFormSchema = z
  .object({
    name: z.string().min(1).max(255),
    vintage: z.string().min(1).max(4),
    format: z.string().min(1).max(255),
    country: z.string().min(1).max(255),
    region: z.string().min(1).max(255),
    purchases: z.array(
      z.object({
        quantity: z.coerce.number().int().min(1).max(1000),
        price: z.coerce.number().min(0),
        date: z.string()
      })
    ),
    imageUrl: z.string().nullable()
  })
  .refine((data) => validateVintage(data.vintage), {
    message: "Vintage must be either a year or 'NV'.",
    path: ["vintage"]
  });

export type WineFormSchema = typeof wineFormSchema;
