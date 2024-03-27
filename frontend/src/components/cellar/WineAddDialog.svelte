<script lang="ts">
	import type { Wine } from "$lib/api/wines";
  import { Button, buttonVariants } from "$lib/components/ui/button";
  import * as Dialog from "$lib/components/ui/dialog";
  import * as Form from "$lib/components/ui/form";
	import Input from "$lib/components/ui/input/input.svelte";
	import Label from "$lib/components/ui/label/label.svelte";
	import { CirclePlus } from "lucide-svelte";

  import { wineFormSchema, type WineFormSchema } from "./WineFormSchema";
  import {
    type SuperValidated,
    type Infer,
    superForm,
  } from "sveltekit-superforms";
  import { zodClient } from "sveltekit-superforms/adapters";

  export let data: SuperValidated<Infer<WineFormSchema>>;

  const form = superForm(data, {
    validators: zodClient(wineFormSchema),
  })
  const { form: formData, enhance } = form;
</script>

<Dialog.Root>
  <Dialog.Trigger class={
    buttonVariants({ 
      variant: "outline", 
      class: "ml-2 mb-2 w-20 h-10 border-0 text-sm bg-secondary text-secondary-foreground hover:bg-primary hover:text-primary-foreground",
    })
  }>
      Add<CirclePlus class="inline ml-2" size=20 />
  </Dialog.Trigger>

  <Dialog.Content class="sm:max-w-[425px]">
    <Dialog.Header>
      <Dialog.Title>Add Wine</Dialog.Title>
      <Dialog.Description>
        Add a wine to your cellar.
      </Dialog.Description>
    </Dialog.Header>

    <form method="POST" use:enhance>
      <Form.Field {form} name="name">
        <Form.Control let:attrs>
          <Form.Label>Name</Form.Label>
          <Input {...attrs} bind:value={$formData.name} />
        </Form.Control>
        <Form.FieldErrors />
      </Form.Field>
      <Form.Button>Submit</Form.Button>
    </form>

  </Dialog.Content>
</Dialog.Root>
