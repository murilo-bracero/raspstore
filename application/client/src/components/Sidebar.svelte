<script>
  import { clickOutside } from '../directives/clickOutsideDirective';

  export let open = false;

  function handleCloseSidebar() {
    open = false;
  }
</script>

<aside
  class="absolute left-[-100%] z-50 h-full w-[80%] border-r-2 border-black bg-white transition-[left] duration-300 ease-in-out"
  class:open
  use:clickOutside
  on:click_outside={handleCloseSidebar}
>
  <button class="mx-7 my-5 flex items-center justify-center" on:click={handleCloseSidebar}>
    <svg xmlns="http://www.w3.org/2000/svg" height="20" viewBox="0 -960 960 960" width="20">
      <path
        d="m249-183-66-66 231-231-231-231 66-66 231 231 231-231 66 66-231 231 231 231-66 66-231-231-231 231Z"
      />
    </svg>
  </button>
  <nav class="container flex flex-col">
    <ul class="relative m-0 list-none px-1">
      <li class="relative">
        <a
          class="flex cursor-pointer items-center justify-between border-b-2 border-t-2 border-black px-6 py-3"
          href="/link"
        >
          <p class="font-bold">Settings</p>
          <span>
            <svg class="w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 -960 960 960">
              <path
                d="m366-55-21-134q-13-4-29.5-13.5T288-221l-125 58L47-368l114-83q-1-6-1.5-14.5T159-480q0-6 .5-14.5T161-509L47-593l116-203 127 57q10-8 26-17t29-13l21-137h228l21 136q13 5 29.5 13.5T672-739l126-57 115 203-115 82q1 7 2 15.5t1 15.5q0 7-1 15t-2 15l115 82-116 205-126-58q-11 9-26.5 18.5T615-189L594-55H366Zm112-295q54 0 92-38t38-92q0-54-38-92t-92-38q-54 0-92 38t-38 92q0 54 38 92t92 38Z"
              />
            </svg>
          </span>
        </a>
      </li>
    </ul>
  </nav>
</aside>

<style>
  .open {
    left: 0;
  }
</style>
