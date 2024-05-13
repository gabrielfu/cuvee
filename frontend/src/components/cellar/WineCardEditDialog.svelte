<script lang="ts">
  import type { Wine } from "$lib/api/wines";
  import { Button, buttonVariants } from "$lib/components/ui/button";
  import * as Dialog from "$lib/components/ui/dialog";
  import Input from "$lib/components/ui/input/input.svelte";
  import Label from "$lib/components/ui/label/label.svelte";
  import { Pencil } from "lucide-svelte";

  export let wine: Wine;
</script>

<Dialog.Root>
  <Dialog.Trigger
    class={buttonVariants({
      variant: "outline",
      class:
        "px-2 ml-2 mb-2 w-20 h-8 border-0 text-sm bg-secondary text-secondary-foreground hover:bg-primary hover:text-primary-foreground"
    })}
  >
    Edit<Pencil class="inline ml-2" size="20" />
  </Dialog.Trigger>
  <Dialog.Content class="sm:max-w-[425px]">
    <Dialog.Header>
      <Dialog.Title>Edit profile</Dialog.Title>
      <Dialog.Description>
        Make changes to your profile here. Click save when you're done.
      </Dialog.Description>
    </Dialog.Header>
    <div class="grid gap-4 py-4 justify-start items-start">
      <div class="grid grid-cols-4 items-center gap-4">
        <Label for="name" class="text-sm text-muted-foreground text-right">Name</Label>
        <Input type="text" id="name" name="name" value={wine.name} class="text-sm col-span-3" />
      </div>
      <div class="grid grid-cols-4 items-center gap-4">
        <Label for="vintage" class="text-sm text-muted-foreground text-right">Vintage</Label>
        <Input
          type="text"
          id="vintage"
          name="vintage"
          value={wine.vintage}
          class="text-sm col-span-3"
        />
      </div>
      <div class="grid grid-cols-4 items-center gap-4">
        <Label for="format" class="text-sm text-muted-foreground text-right">Format</Label>
        <Input
          type="text"
          id="format"
          name="format"
          value={wine.format}
          class="text-sm col-span-3"
        />
      </div>
      <div class="grid grid-cols-4 items-center gap-4">
        <Label for="price" class="text-sm text-muted-foreground text-right">Price</Label>
        <Input
          type="number"
          id="price"
          name="price"
          value={wine.summary.price}
          class="text-sm col-span-3"
        />
      </div>
      <div class="grid grid-cols-4 items-center gap-4">
        <Label for="date" class="text-sm text-muted-foreground text-right">Date</Label>
        <Input
          type="date"
          id="date"
          name="date"
          value={new Date().toISOString().slice(0, 10)}
          class="text-sm col-span-3"
        />
      </div>
    </div>
    <Dialog.Footer>
      <Button type="submit" class="text-sm h-8 bg-primary text-primary-foreground rounded-md p-4"
        >Save</Button
      >
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
