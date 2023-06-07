<script>
  import '../../app.css';
  import Header from '../../components/Header.svelte';
  import Footer from '../../components/Footer.svelte';
</script>

<Header />
<slot />
<Footer />

<style lang="postcss">
  :global(html) {
    background-color: theme(colors.white);
    font-family: 'consolas';
  }
</style>
