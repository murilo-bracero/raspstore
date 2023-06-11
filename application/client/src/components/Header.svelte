<script>
  import Sidebar from './Sidebar.svelte';

  let open = false;

  function toggleSideBar() {
    open = !open;
  }
</script>

<header class="container flex min-w-full flex-col items-center">
  <Sidebar bind:open />
  <div class="flex w-screen flex-row justify-between border-b-2 border-black px-4 py-2">
    <button on:click={toggleSideBar}>
      <figure class="w-12 self-end justify-self-end">
        <img
          class="rounded-full"
          src="https://newprofilepic2.photo-cdn.net//assets/images/article/profile.jpg"
          alt="profile"
        />
      </figure>
    </button>
  </div>
</header>
