<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  let input: HTMLInputElement;

  const dispatcher = createEventDispatcher();

  function handleUploadClick() {
    input.click();
  }

  function handleUploadEvent() {
    if (!input.files) {
      return;
    }

    dispatcher('upload', {
      file: input.files[0]
    });
  }
</script>

<footer class="container fixed bottom-0 flex min-w-full flex-col items-center">
  <input type="file" class="hidden" bind:this={input} on:change={handleUploadEvent} />
  <button
    on:click={handleUploadClick}
    class="z-10 -mb-5 mr-2 inline-flex h-12 w-12 items-center justify-center rounded-full border-2 border-black bg-sky-500"
  >
    <svg xmlns="http://www.w3.org/2000/svg" height="30" viewBox="0 -960 960 960" width="30">
      <path
        d="M220-160q-24 0-42-18t-18-42v-143h60v143h520v-143h60v143q0 24-18 42t-42 18H220Zm230-153v-371L330-564l-43-43 193-193 193 193-43 43-120-120v371h-60Z"
      />
    </svg>
  </button>
  <div class="z-0 mx-3 h-12 w-96 rounded-t-lg border-2 border-b-0 border-black bg-white" />
</footer>
