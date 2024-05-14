<script lang="ts">
	import CalendarIcon from "lucide-svelte/icons/calendar";
	import {
		CalendarDate,
		type DateValue,
		getLocalTimeZone,
		parseDate,
		today,
	} from "@internationalized/date";
	import { cn } from "$lib/utils.js";
	import { buttonVariants } from "$lib/components/ui/button";
	import { Calendar } from "$lib/components/ui/calendar";
	import * as Popover from "$lib/components/ui/popover";

  function formatDate(date: Date): string {
    return date.toISOString().slice(0, 10);
  }

  export let attrs: any;
  export let formValue: string;
  let value: DateValue | undefined;
  $: value = formValue ? parseDate(formValue) : undefined;
  export let onValueChange: (v: DateValue | undefined) => void;
  let placeholder: DateValue = today(getLocalTimeZone());

  let className: string | null | undefined = undefined;
  export { className as class };
</script>

<Popover.Root>
  <Popover.Trigger
    {...attrs}
    class={cn(
      buttonVariants({ variant: "outline" }),
      "w-[280px] justify-start pl-4 text-left font-normal",
      !value && "text-muted-foreground",
      className
    )}
  >
    {value ? formatDate(value.toDate(getLocalTimeZone())) : "Pick a date"}
    <CalendarIcon class="ml-auto h-4 w-4 opacity-50" />
  </Popover.Trigger>
  <Popover.Content class="w-auto p-0" side="top">
    <Calendar
      {value}
      bind:placeholder
      minValue={new CalendarDate(1900, 1, 1)}
      maxValue={today(getLocalTimeZone())}
      calendarLabel="Date of birth"
      initialFocus
      onValueChange={onValueChange}
    />
  </Popover.Content>
</Popover.Root>