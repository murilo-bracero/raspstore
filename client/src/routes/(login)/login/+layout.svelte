<script>
  import '../../../app.css';
  import Toast from '../../../components/Toast.svelte';
</script>

<Toast />
<slot />

<style lang="postcss">
  :global(html) {
    background-color: theme(colors.white);
    font-family: 'consolas';
  }
</style>
