<script>
  import Toast from '$lib/components/Toast.svelte';
  import '../../../app.css';
</script>

<Toast />
<slot />

<style lang="postcss">
  :global(html) {
    background-color: theme(colors.white);
    font-family: 'consolas';
  }
</style>
