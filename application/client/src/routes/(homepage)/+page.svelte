<script lang="ts">
  import FileItem from '../../components/FileItem.svelte';
  import SearchBar from '../../components/SearchBar.svelte';
  import { files } from '../../stores/file';
  import { getFiles } from '../../service/file-info-service';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';

  onMount(async () => {
    getFiles()
      .then((filesData) => files.set(filesData))
      .catch((error) => {
        if (error.status === 401) {
          goto('/login');
          return;
        }

        console.error(error);
      });
  });
</script>

<main class="mb-24 flex flex-col items-center justify-center overflow-x-hidden">
  <article class="my-2">
    <SearchBar />
  </article>
  {#each $files as file}
    <FileItem bind:fileData={file} />
  {/each}
</main>

<style>
  main {
    overflow: auto;
  }
</style>
