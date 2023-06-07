<script lang="ts">
  import FileItem from '../../components/FileItem.svelte';
  import SearchBar from '../../components/SearchBar.svelte';
  import type { PageData } from '../../stores/file';
  import { files } from '../../stores/file';
  import { onMount } from 'svelte';

  onMount(async () => {
    fetch(import.meta.env.VITE_FILES_SERVICE_URL, {
      headers: {
        Authorization:
          'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjI2ODYwMDc2OTcsInVpZCI6ImM3NGQ3NzIwLTAwMjYtNDQ2Ni1iNTlmLWQxYjRhN2Y2ODg2ZiJ9.XPqmqcaJCH5dPyzL-BJmRKrpkKBqdLaWkA5P6Rg-cyw'
      }
    })
      .then((res) => res.json())
      .then((data: PageData) => {
        files.set(data.content);
      })
      .catch((err) => console.log(err));
  });
</script>

<main class="flex flex-col items-center justify-center">
  <article class="my-2">
    <SearchBar />
  </article>
  {#each $files as file}
    <FileItem bind:fileData={file} />
  {/each}
</main>
