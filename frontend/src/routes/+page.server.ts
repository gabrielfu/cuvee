import { superValidate, message } from "sveltekit-superforms";
import { zod } from "sveltekit-superforms/adapters";
import type { PageServerLoad } from './$types';
import { fail, type Actions } from "@sveltejs/kit";
import { listWines } from '$lib/api/wines';
import { wineFormSchema } from '../components/cellar/WineFormSchema';

export const load: PageServerLoad = async () => {
	const wines = await listWines();
	const images: string[] = [
		"https://images.vivino.com/thumbs/uDMfPG10R1efm5H1do--Ow_pb_600x600.png",
		"https://images.vivino.com/thumbs/Eh506eKdSXGvndyHHfqPug_pb_600x600.png",
		"https://images.vivino.com/thumbs/MCs7Ix2zS56U3-vs6GOR4A_pb_600x600.png",
	];

	return { 
		wines, 
		images,
    form: await superValidate(zod(wineFormSchema)),
	};
};

export const actions: Actions = {
  default: async (event) => {
    const form = await superValidate(event, zod(wineFormSchema));
    if (!form.valid) {
      return fail(400, {
        form,
      });
    }

    const response = await fetch(`http://localhost:8080/wines`, { 
      method: "POST", 
      body: JSON.stringify(form.data), 
      headers: { "Content-Type": "application/json" },
    });
    if (!response.ok) {
      return fail(500, {
        form,
        error: await response.text(),
      });
    }

    return message(form, { text: "Success!"});
  },
};

export const ssr = true;
