import { superValidate, setError } from "sveltekit-superforms";
import { zod } from "sveltekit-superforms/adapters";
import type { PageServerLoad, Actions } from "./$types";
import { fail } from "@sveltejs/kit";
import { listWines } from "$lib/api/wines";
import { wineFormSchema } from "../components/cellar/WineFormSchema";
import { baseUrl } from "$lib/api/utils";

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

    const response = await fetch(`${baseUrl}/api/v1/wines`, {
      method: "POST",
      body: JSON.stringify(form.data),
      headers: { "Content-Type": "application/json" },
    });
    if (!response.ok) {
      try {
        const error = await response.json();
        if (error.type == "validation") {
          let message = "Invalid fields:\n";
          for (const innerError of error.error) {
            message += `${innerError.field}: ${innerError.message}\n`;
          }
          return setError(form, "", message);
        }
        return setError(form, "", error.error);
      } catch (e) {
        return setError(
          form,
          "",
          "Encountered error when trying to create a wine. Please try again later."
        );
      }
    }

    return { form };
  }
};

export const ssr = true;
