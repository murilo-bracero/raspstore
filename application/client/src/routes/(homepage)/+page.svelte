<script lang="ts">
  import { goto } from '$app/navigation';
  import FileItem from '../../components/FileItem.svelte';
  import SearchBar from '../../components/SearchBar.svelte';
  import type { PageData } from '../../stores/file';
  import { files } from '../../stores/file';
  import { onMount } from 'svelte';

  onMount(async () => {
    const token = getToken();

    if (token === null) {
      goto('/login');
      return;
    }

    fetch(import.meta.env.VITE_FILES_SERVICE_URL, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })
      .then((res) => res.json())
      .then((data: PageData) => {
        files.set(data.content);
      })
      .catch((err) => console.log(err));
  });

  function getToken(): string | null {
    return localStorage.getItem(import.meta.env.VITE_TOKEN_KEY);
  }
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
