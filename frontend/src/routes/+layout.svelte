<script lang="ts">
  import "../app.css";
  import Header from "../components/Header.svelte";
  import Footer from "../components/Footer.svelte";
  import { ArrowUp } from "lucide-svelte";

  let y: number;
  let innerHeight = 0;
  let innerWidth = 0;

  function goTop() {
    document.body.scrollIntoView();
  }
</script>

<div
  class="container relative flex flex-col max-w-[1400px] mx-auto w-full text-sm sm:text-base min-h-screen"
>
  <div
    class={"fixed bottom-10 right-10 duration-200 flex z-10 " +
      (y > 0 ? " opacity-full pointer-events-auto" : " pointer-events-none opacity-0")}
  >
    <button
      on:click={goTop}
      class="ml-auto rounded-full bg-secondary text-secondary-foreground px-3 sm:px-4 hover:bg-primary hover:text-primary-foreground cursor-pointer aspect-square grid place-items-center"
    >
      <ArrowUp />
    </button>
  </div>
  <Header />
  <slot />
  <Footer />
</div>

<svelte:window bind:scrollY={y} bind:innerHeight bind:innerWidth />
