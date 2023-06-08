<script>
  import '../../app.css';
  import Header from '../../components/Header.svelte';
  import Footer from '../../components/Footer.svelte';
  import Toast from '../../components/Toast.svelte';
</script>

<Toast />
<Header />
<slot />
<Footer />

<style lang="postcss">
  :global(html) {
    background-color: theme(colors.white);
    font-family: 'consolas';
  }
</style>
