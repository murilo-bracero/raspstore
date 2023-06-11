<script lang="ts">
  import type { ActionResult } from '@sveltejs/kit';
  import { NotificationType, toast } from '../../../stores/toast';
  import { goto } from '$app/navigation';

  async function handleSubmit(this: any) {
    const data = new FormData(this);

    const response = await fetch(this.action, {
      method: 'POST',
      body: data
    });

    const body = (await response.json()) as ActionResult;

    if (body.type === 'success') {
      goto('/');
      return;
    }

    toast({
      message: 'username or password must be provided',
      type: NotificationType.ERROR
    });
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
      class="rounded-lg border-2 border-black p-1"
      type="password"
      id="password"
      name="password"
    />
    <button
      class="my-3 self-end rounded-full border-2 border-black bg-sky-400 px-3 py-1 font-bold"
      type="submit">Log in</button
    >
  </form>
</main>
