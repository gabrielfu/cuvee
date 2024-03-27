import { z } from "zod";

export const wineFormSchema = z.object({
    name: z.string().min(1).max(255),
    // vintage: z.string().min(1).max(10),
    // format: z.string().min(1).max(255),
    // country: z.string().min(1).max(255),
    // region: z.string().min(1).max(255),
    // purchases: z.array(z.object({
    //     quantity: z.number().int().min(1).max(1000),
    //     price: z.number().int().min(0),
    //     date: z.date(),
    // })),
});
 
export type WineFormSchema = typeof wineFormSchema;
