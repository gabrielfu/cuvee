import { superValidate, setError } from "sveltekit-superforms";
import { zod } from "sveltekit-superforms/adapters";
import type { PageServerLoad, Actions } from "./$types";
import { fail } from "@sveltejs/kit";
import { listWines } from "$lib/api/wines";
import { wineFormSchema } from "../components/cellar/WineFormSchema";

export const load: PageServerLoad = async () => {
  const wines = await listWines();

  return {
    wines,
    form: await superValidate(zod(wineFormSchema))
  };
};

export const actions: Actions = {
  create: async (event) => {
    const form = await superValidate(event, zod(wineFormSchema));
    if (!form.valid) {
      return fail(400, {
        form
      });
    }

    const response = await fetch(`http://localhost:8080/wines`, {
      method: "POST",
      body: JSON.stringify(form.data),
      headers: { "Content-Type": "application/json" }
    });
    if (!response.ok) {
      const text = await response.text();
      return setError(form, "", text);
    }

    return { form };
  }
};

export const ssr = true;
