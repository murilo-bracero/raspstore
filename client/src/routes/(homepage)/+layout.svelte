<script lang="ts">
  import '../../app.css';
  import Header from '../../components/Header.svelte';
  import Footer from '../../components/Footer.svelte';
  import Toast from '../../components/Toast.svelte';
  import type { ActionResult } from '@sveltejs/kit';
  import { NotificationType, toast } from '../../stores/toast';

  async function handleUpload(e: any) {
    const { actionUrl, formData } = e.detail;

    const response = await fetch(actionUrl, {
      method: 'POST',
      body: formData
    });

    const body = (await response.json()) as ActionResult;

    if (body.type === 'success') {
      toast({
        message: 'File uploaded successfully',
        type: NotificationType.SUCCESS
      });
      return;
    }

    toast({
      message: 'Could not upload file',
      type: NotificationType.ERROR
    });
  }
</script>

<Toast />
<Header />
<slot />
<Footer on:upload={handleUpload} />

<style lang="postcss">
  :global(html) {
    background-color: theme(colors.white);
    font-family: 'consolas';
  }
</style>
