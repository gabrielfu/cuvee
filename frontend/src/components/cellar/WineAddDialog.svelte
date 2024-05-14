<script lang="ts">
  import { Button, buttonVariants } from "$lib/components/ui/button";
  import * as Dialog from "$lib/components/ui/dialog";
  import * as Form from "$lib/components/ui/form";
  import { Input } from "$lib/components/ui/input";
  import * as Carousel from "$lib/components/ui/carousel";
  import * as Card from "$lib/components/ui/card";
  import { CirclePlus, X } from "lucide-svelte";
  import { tick } from "svelte";

  import { wineFormSchema, type WineFormSchema } from "./WineFormSchema";
  import { type SuperValidated, type Infer, superForm } from "sveltekit-superforms";
  import { zodClient } from "sveltekit-superforms/adapters";
  import { searchImages, type ImageResult } from "$lib/api/images";
  import DatePicker from "../utils/DatePicker.svelte";

  export let data: SuperValidated<Infer<WineFormSchema>>;
  export let displayText: string;

  const form = superForm(data, {
    dataType: "json",
    validators: zodClient(wineFormSchema),
    onUpdated: ({ form: f }) => {
      if (!f.valid) {
        console.error("Form is invalid:", f.errors);
      }
    }
  });
  const { form: formData, enhance, allErrors, errors } = form;
  $formData.purchases = [
    ...$formData.purchases,
    { quantity: 1, price: 0, date: new Date().toISOString().split("T")[0] }
  ];

  function addPurchase() {
    $formData.purchases = [
      ...$formData.purchases,
      { quantity: 1, price: 0, date: new Date().toISOString().split("T")[0] }
    ];

    tick().then(() => {
      const urlInputs = Array.from(
        document.querySelectorAll<HTMLElement>("#wine-form input[name='purchases']")
      );
      const lastInput = urlInputs[urlInputs.length - 1];
      lastInput && lastInput.focus();
    });
  }

  function removePurchase(i: number) {
    $formData.purchases = $formData.purchases.filter((_, index) => index !== i);
  }

  // function setPurchaseDate(event: Event, i: number) {
  //   const target = event.target as HTMLInputElement;
  //   const date = new Date(target.value);
  //   $formData.purchases[i].date = date.toISOString().split("T")[0];
  // }

  let imageCandidates: ImageResult[] = [];
  let imageError: string | null = null;
  let imageLoading = false;
  function searchWineImages() {
    imageLoading = true;
    searchImages($formData.name, $formData.vintage, $formData.country, $formData.region)
      .then((images) => {
        imageCandidates = images;
        imageError = null;
      })
      .catch((error) => {
        imageError = error.message;
      })
      .finally(() => {
        imageLoading = false;
        tick();
      });
  }

  function onImageClicked(image: ImageResult) {
    if ($formData.imageUrl == image.link) {
      $formData.imageUrl = null;
    } else {
      $formData.imageUrl = image.link;
    }
  }
</script>

<Dialog.Root>
  <Dialog.Trigger
    class={buttonVariants({
      variant: "outline",
      class:
        "ml-2 h-8 border-0 text-sm bg-secondary text-secondary-foreground hover:bg-primary hover:text-primary-foreground"
    })}
  >
    {displayText}<CirclePlus class="inline ml-2" size="20" />
  </Dialog.Trigger>

  <Dialog.Content class="min-w-[800px] sm:max-w-[425px] overflow-y-scroll max-h-screen p-8">
    <Dialog.Header>
      <Dialog.Title>Add Wine</Dialog.Title>
      <Dialog.Description>Add a wine to your cellar.</Dialog.Description>
    </Dialog.Header>

    <form method="POST" use:enhance id="wine-form" action="?/create">
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
        <Form.Description>Year number or "NV"</Form.Description>
        <Form.FieldErrors />
      </Form.Field>

      <Form.Field {form} name="format">
        <Form.Control let:attrs>
          <Form.Label>Format</Form.Label>
          <Input {...attrs} bind:value={$formData.format} />
        </Form.Control>
        <Form.Description>E.g., 750ml</Form.Description>
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
          {#each $formData.purchases as _, i}
            <div
              class="grid grid-cols-[1fr_48px] gap-4 rounded-md border-[1px] border-input px-4 py-4"
            >
              <Form.ElementField {form} name="purchases[{i}]">
                <Form.Control let:attrs>
                  <div class="flex gap-4 mb-2 items-start w-auto">
                    <Form.Label class="w-32 mt-2">Quantity</Form.Label>
                    <Input
                      type="number"
                      min="1"
                      {...attrs}
                      bind:value={$formData.purchases[i].quantity}
                    />
                  </div>
                </Form.Control>
                <Form.Control let:attrs>
                  <div class="flex gap-4 mb-2 items-start w-auto">
                    <Form.Label class="w-32 mt-2">Average Price</Form.Label>
                    <Input
                      type="number"
                      step="any"
                      class="[&::-webkit-inner-spin-button]:appearance-none"
                      {...attrs}
                      bind:value={$formData.purchases[i].price}
                    />
                  </div>
                </Form.Control>
                <Form.Control let:attrs>
                  <div class="flex gap-4 mb-2 items-start w-auto">
                    <Form.Label class="w-32 mt-2">Date</Form.Label>
                    <DatePicker 
                      class="flex h-8 w-full"
                      attrs={attrs} 
                      onValueChange={(v) => {
                        if (v) {
                          $formData.purchases[i].date = v.toString();
                        } else {
                          $formData.purchases[i].date = "";
                        }
                      }} 
                      bind:formValue={$formData.purchases[i].date}
                    />
                    <input hidden name={attrs.name} bind:value={$formData.purchases[i].date} />
                  </div>
                </Form.Control>
              </Form.ElementField>
              <Button
                class="w-8 h-8"
                on:click={() => removePurchase(i)}
                disabled={$formData.purchases.length == 1}
                ><X class="inline -m-2" size="20" /></Button
              >
            </div>
          {/each}
          <Form.FieldErrors />
        </Form.Fieldset>
        <Button type="button" variant="outline" size="sm" class="mt-2 h-8" on:click={addPurchase}>
          Add Purchase
        </Button>
      </div>

      <div class="my-4">
        <Form.Fieldset {form} name="imageUrl">
          <Form.Legend>Thumbnail</Form.Legend>
          <Button
            type="button"
            variant="outline"
            size="sm"
            class="mt-4 h-8"
            on:click={searchWineImages}
            disabled={imageLoading}
          >
            Search for an image online
          </Button>
          {#if imageError}
            <p class="text-red-600 text-center p-2">{imageError}</p>
          {:else}
            <Carousel.Root class="w-full max-w-lg">
              <Carousel.Content class="my-4 -ml-2">
                {#each imageCandidates as image (image.link)}
                  <Carousel.Item class=" basis-[200px]">
                    <div on:click={() => onImageClicked(image)}>
                      <Card.Root
                        class={"items-center text-center shadow-lg overflow-hidden border-2 " +
                          ($formData.imageUrl == image.link ? "border-red-600" : "")}
                      >
                        <Card.Content>
                          <img
                            src={image.link}
                            alt={image.link}
                            class="m-auto max-h-[200px] max-w-[200px]"
                          />
                        </Card.Content>
                      </Card.Root>
                    </div>
                  </Carousel.Item>
                {/each}
              </Carousel.Content>
            </Carousel.Root>
          {/if}
        </Form.Fieldset>
      </div>

      <Form.Button class="h-8">Submit</Form.Button>
      {#if $allErrors.length}
        <p class="invalid pt-4 text-red-600" style="white-space: pre-wrap">
          {$allErrors[0].messages[0]}
        </p>
      {/if}
    </form>
  </Dialog.Content>
</Dialog.Root>
