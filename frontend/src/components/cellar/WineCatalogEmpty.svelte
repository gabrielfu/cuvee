<script lang="ts">
  import * as Card from "$lib/components/ui/card";
  import WineAddDialog from "./WineAddDialog.svelte";
  import type { WineFormSchema } from "./WineFormSchema";
  import type { SuperValidated, Infer } from "sveltekit-superforms";

  export let form: SuperValidated<Infer<WineFormSchema>>;
</script>

<Card.Card class="bg-transparent border-transparent text-card-foreground text-center pt-4 mb-4">
  <Card.Content>
    <p class="my-4 text-lg">Your cellar is empty...</p>
    <WineAddDialog data={form} displayText="Click here to add a wine" />
  </Card.Content>
</Card.Card>
