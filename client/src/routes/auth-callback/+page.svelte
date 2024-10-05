<script lang="ts">
  import { User, UserManager } from 'oidc-client-ts';
  import { settings } from '../../config/oidc';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { NotificationType, toast } from '$lib/stores/toast';
  import { isAuthenticated, userStore } from '$lib/stores/login';

  onMount(async () => {
    const um = new UserManager(settings);

    console.log(um.metadataService.getTokenEndpoint());

    um.signinRedirectCallback().then(successfullSignIn).catch(failedSignIn);
  });

  function successfullSignIn(user: User) {
    userStore.set(user);
    isAuthenticated.set(true);

    goto('/');
  }

  function failedSignIn() {
    toast({
      message: 'Could not sign in',
      type: NotificationType.ERROR
    });
    goto('/login');
  }
</script>

<main class="w-100 flex items-center justify-center">
  <div class="my-10 border-2 border-black bg-sky-400 p-6 opacity-100 shadow-[8px_6px_0_2px_#000]">
    <p class="font-bold tracking-wide">Authenticating...</p>
  </div>
</main>
