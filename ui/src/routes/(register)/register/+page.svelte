<script lang="ts">
  import { enhance } from '$app/forms';
  import { NotificationType, toast } from '$lib/stores/toast.js';

  export let form;

  $: {
    if (form && !form.success) {
      toast({
        message: form.message || 'Could not register user',
        type: NotificationType.ERROR
      });
    }
  }
</script>

<main class="flex h-screen items-center justify-center">
  <form
    class="flex flex-col items-start justify-center rounded-lg border-2 border-black p-6"
    method="POST"
    action="?/register"
    use:enhance
  >
    <label for="username">Username:</label>
    <input
      class="mb-3 rounded-lg border-2 border-black p-1"
      type="text"
      name="username"
      id="username"
    />
    <label for="name">Name (Optional):</label>
    <input class="mb-3 rounded-lg border-2 border-black p-1" type="text" name="name" id="name" />
    <label for="password">Password:</label>
    <input
      class="mb-3 rounded-lg border-2 border-black p-1"
      type="password"
      name="password"
      id="password"
    />
    <label for="confirm-password">Confirm Password:</label>
    <input
      class="mb-3 rounded-lg border-2 border-black p-1"
      type="password"
      name="confirm-password"
      id="confirm-password"
    />
    <button
      class="my-3 w-full self-center rounded-full border-2 border-black bg-sky-300 px-3 py-1 font-bold"
      type="submit"
    >
      Sign up
    </button>

    <a href="/login" class="self-center text-sky-600 underline">Already have an account?</a>
  </form>
</main>
