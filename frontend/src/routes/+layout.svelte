<script lang="ts">
  import "../app.css";
</script>

<slot />

<style global>
  @font-face {
      font-family: 'Schuss';
      font-style: normal;
      font-weight: 400;
      src: url('../static/fonts/SchussSerifPro-Regular.woff2') format("woff2");
  }

  :global(html) {
    font-family: 'Schuss';
    font-size: 20px;
    font-style: italic;
  }
</style>