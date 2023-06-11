<script>
  import { search } from '../stores/search';

  function handleSearchClick() {
    search.subscribe(({ predicate }) => console.log(predicate));
  }
</script>

<fieldset class="flex w-80 flex-row rounded-xl border-2 border-black bg-white py-1 pr-5">
  <button class="flex-4 mx-3" on:click={handleSearchClick}>
    <svg xmlns="http://www.w3.org/2000/svg" height="23" viewBox="0 -960 960 960" width="23">
      <path
        d="M795-93 526-361q-29 22.923-68.459 35.962Q418.082-312 372-312q-115.162 0-195.081-80Q97-472 97-585t80-193q80-80 193.5-80t193 80Q643-698 643-584.85q0 44.85-12.5 83.35Q618-463 593-429l270 268-68 68ZM371.353-406q74.897 0 126.272-52.25T549-585q0-74.5-51.522-126.75T371.353-764q-75.436 0-127.895 52.25Q191-659.5 191-585t52.311 126.75Q295.623-406 371.353-406Z"
      />
    </svg>
  </button>
  <input class="flex-1 text-xl" type="text" bind:value={$search.predicate} />
</fieldset>
