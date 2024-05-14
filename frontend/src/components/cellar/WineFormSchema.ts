import { z } from "zod";

function validateVintage(vintage: string): boolean {
  return vintage == "NV" || /^\d{4}$/.test(vintage);
}

function numberErrorMessage(name: string, min: number, max: number): string {
  return `${name} must be a number between ${min} and ${max}.`;
}

function stringLengthErrorMessage(name: string, min: number, max: number): string {
  return `${name} must be between ${min} and ${max} characters long.`;
}

function stringField(name: string, min: number, max: number, withMessage?: boolean): z.ZodString {
  const options = withMessage ? { message: stringLengthErrorMessage(name, min, max) } : {};
  return z.string().min(min, options).max(max, options);
}

export const wineFormSchema = z
  .object({
    name: stringField("Name", 1, 255, true),
    vintage: stringField("Vintage", 1, 4, false),
    format: stringField("Format", 1, 255, true),
    country: stringField("Country", 1, 255, true),
    region: stringField("Region", 1, 255, true),
    purchases: z.array(
      z.object({
        quantity: z.coerce.number().int().min(1, { message: numberErrorMessage("Quantity", 1, 9999) }).max(9999, { message: numberErrorMessage("Quantity", 1, 9999) }),
        price: z.coerce.number().min(0, { message: "Price must be 0 or above." } ),
        date: z.string()
      })
    ).min(1, { message: "At least one purchase must be provided." }),
    imageUrl: z.string().nullable()
  })
  .refine((data) => validateVintage(data.vintage), {
    message: "Vintage must be either a year or 'NV'.",
    path: ["vintage"]
  });

export type WineFormSchema = typeof wineFormSchema;
