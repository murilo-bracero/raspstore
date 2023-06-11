<script lang="ts">
  import FileItem from '../../components/FileItem.svelte';
  import SearchBar from '../../components/SearchBar.svelte';
  import type { PageData } from '../../stores/file';

  export let data: PageData;
</script>

<main class="mb-24 flex flex-col items-center justify-center overflow-x-hidden">
  <article class="my-2">
    <SearchBar />
  </article>
  {#each data.content as file}
    <FileItem bind:fileData={file} />
  {/each}
</main>

<style>
  main {
    overflow: auto;
  }
</style>
