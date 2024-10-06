<script lang="ts">
  import FileItem from '$lib/components/FileItem.svelte';
  import Footer from '$lib/components/Footer.svelte';
  import Header from '$lib/components/Header.svelte';
  import SearchBar from '$lib/components/SearchBar.svelte';
  import { type PageData } from '$lib/stores/file';
  import { NotificationType, toast } from '$lib/stores/toast';

  export let data: PageData;

  async function onDocumentDrop(e: DragEvent) {
    if (!e.dataTransfer) {
      return;
    }

    if (!e.dataTransfer.files || e.dataTransfer.files.length < 1) {
      return;
    }

    const files = [...e.dataTransfer.items].filter((item) => item.kind === 'file');

    if (files.length < 1) {
      return;
    }

    upload(files[0].getAsFile() as File);
  }

  async function handleUpload(e: any) {
    const { file } = e.detail;

    upload(file);
  }

  function upload(file: File) {
    const formData = createFormData(file);

    const response = fetch('?/upload', {
      method: 'POST',
      body: formData
    })
      .then((response) => {
        if (!response.ok) {
          toast({
            message: 'Could not upload file',
            type: NotificationType.ERROR
          });
        }

        return response.json();
      })
      .catch((error) => {
        toast({
          message: 'Could not upload file',
          type: NotificationType.ERROR
        });
      });

    if (!response) {
      return;
    }

    toast({
      message: 'Document uploaded successfully',
      type: NotificationType.SUCCESS
    });
  }

  function createFormData(file: File): FormData {
    const formData = new FormData();
    formData.append('file', file);
    return formData;
  }
</script>

<svelte:window on:drop|preventDefault={onDocumentDrop} on:dragover|preventDefault />

<Header />
<main class="mb-24 flex flex-col items-center justify-center overflow-x-hidden">
  <article class="my-2">
    <SearchBar />
  </article>
  {#each data.content as file}
    <FileItem bind:fileData={file} />
  {/each}
</main>
<Footer on:upload={handleUpload} />

<style lang="postcss">
  :global(html) {
    background-color: theme(colors.white);
    font-family: 'consolas';
  }

  main {
    overflow: auto;
  }
</style>
