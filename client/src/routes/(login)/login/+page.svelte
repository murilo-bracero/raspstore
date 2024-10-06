<script lang="ts">
  import { goto } from '$app/navigation';
  import { NotificationType, toast } from '$lib/stores/toast';

  async function handleSubmit(this: any) {
    toast({
      message: 'Local login is not implemented yet, but new maintainers are welcome.',
      type: NotificationType.ERROR
    });
  }

  async function handleOidcClick() {
    const res = await fetch('/api/login');

    const body = await res.json();

    goto(body.authUrl);
  }
</script>

<main class="flex h-screen items-center justify-center">
  <form
    class="flex flex-col items-start justify-center rounded-lg border-2 border-black p-6"
    on:submit|preventDefault={handleSubmit}
  >
    <label class="mb-1 p-1 font-bold" for="username">Username:</label>
    <input class="rounded-lg border-2 border-black p-1" type="text" id="username" name="username" />
    <label class="mb-1 mt-3 p-1 font-bold" for="password">Password:</label>
    <input
      class="mb-3 rounded-lg border-2 border-black p-1"
      type="password"
      id="password"
      name="password"
    />
    <button
      class="my-3 w-full cursor-not-allowed self-center rounded-full border-2 border-black bg-sky-200 px-3 py-1 font-bold opacity-50"
      type="submit">PAM Sign In</button
    >

    <button
      class="my-3 w-full self-center rounded-full border-2 border-black bg-sky-600 px-3 py-1 font-bold"
      type="button"
      on:click={handleOidcClick}>OIDC Sign In</button
    >
  </form>
</main>
