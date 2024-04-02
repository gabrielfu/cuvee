<script lang="ts">
  import { Button, buttonVariants } from "$lib/components/ui/button";
  import * as Dialog from "$lib/components/ui/dialog";
  import * as Form from "$lib/components/ui/form";
	import { Input } from "$lib/components/ui/input";
	import * as Table from "$lib/components/ui/table";
	import { CirclePlus } from "lucide-svelte";
  import { tick } from "svelte";

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
  $formData.purchases = [...$formData.purchases, {quantity: 1, price: 0, date: new Date()}];

	function addPurchase() {
		$formData.purchases = [...$formData.purchases, {quantity: 1, price: 0, date: new Date()}];

		tick().then(() => {
			const urlInputs = Array.from(
				document.querySelectorAll<HTMLElement>("#wine-form input[name='purchases']")
			);
			const lastInput = urlInputs[urlInputs.length - 1];
			lastInput && lastInput.focus();
		});
	}

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

  <Dialog.Content class="min-w-[800px] sm:max-w-[425px] overflow-y-scroll max-h-screen">
    <Dialog.Header>
      <Dialog.Title>Add Wine</Dialog.Title>
      <Dialog.Description>
        Add a wine to your cellar.
      </Dialog.Description>
    </Dialog.Header>

    <form method="POST" use:enhance id="wine-form">
      <Form.Field {form} name="name">
        <Form.Control let:attrs>
          <Form.Label>Name</Form.Label>
          <Input {...attrs} bind:value={$formData.name} />
        </Form.Control>
        <Form.FieldErrors />
      </Form.Field>

      <Form.Field {form} name="vintage">
        <Form.Control let:attrs>
          <Form.Label>Vintage</Form.Label>
          <Input {...attrs} bind:value={$formData.vintage} />
        </Form.Control>
        <Form.FieldErrors />
      </Form.Field>

      <Form.Field {form} name="format">
        <Form.Control let:attrs>
          <Form.Label>Format</Form.Label>
          <Input {...attrs} bind:value={$formData.format} />
        </Form.Control>
        <Form.FieldErrors />
      </Form.Field>

      <Form.Field {form} name="country">
        <Form.Control let:attrs>
          <Form.Label>Coutry</Form.Label>
          <Input {...attrs} bind:value={$formData.country} />
        </Form.Control>
        <Form.FieldErrors />
      </Form.Field>

      <Form.Field {form} name="region">
        <Form.Control let:attrs>
          <Form.Label>Region</Form.Label>
          <Input {...attrs} bind:value={$formData.region} />
        </Form.Control>
        <Form.FieldErrors />
      </Form.Field>

      <div class="my-4">
        <Form.Fieldset {form} name="purchases">
          <Form.Legend>Purchases</Form.Legend>
          <Table.Table>
            <Table.Header>
              <Table.Row>
                <Table.Head class="h-auto text-sm">Quantity</Table.Head>
                <Table.Head class="h-auto text-sm">Price</Table.Head>
                <Table.Head class="h-auto text-sm">Date</Table.Head>
              </Table.Row>
            </Table.Header>
            <Table.Body>
              {#each $formData.purchases as _, i}
                <Table.Row>
                  <Table.Cell>
                    <Form.Control let:attrs>
                      <Input type="number" min=1 {...attrs} bind:value={$formData.purchases[i].quantity} />
                    </Form.Control>
                  </Table.Cell>
                  <Table.Cell>
                    <Form.Control let:attrs>
                      <Input type="number" step="any" class="[&::-webkit-inner-spin-button]:appearance-none" {...attrs} bind:value={$formData.purchases[i].price} />
                    </Form.Control>
                  </Table.Cell>
                  <Table.Cell>
                    <Form.Control let:attrs>
                      <Input type="date" {...attrs} bind:value={$formData.purchases[i].date} />
                    </Form.Control>
                  </Table.Cell>
                </Table.Row>

              {/each}
            </Table.Body>
          </Table.Table>
        </Form.Fieldset>
        <Button type="button" variant="outline" size="sm" class="mt-2" on:click={addPurchase}>
          Add Purchase
        </Button>
      </div>

      <Form.Button>Submit</Form.Button>
    </form>

  </Dialog.Content>
</Dialog.Root>
