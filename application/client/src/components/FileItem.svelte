<script lang="ts">
  import type { FileData } from '../stores/file';

  export let fileData: FileData | undefined;
</script>

<section class="my-1 flex w-[85%] flex-row rounded-xl border-2 border-black p-4">
  <figure>
    <svg
      class="mt-2 fill-red-500"
      xmlns="http://www.w3.org/2000/svg"
      height="20"
      viewBox="0 -960 960 960"
      width="20"
    >
      <path
        d="M353-448h37v-83h48q15.725 0 26.362-10.638Q475-552.275 475-568v-48q0-15.725-10.638-26.362Q453.725-653 438-653h-85v205Zm37-120v-48h48v48h-48Zm122 120h84q15 0 26-10.638 11-10.637 11-26.362v-131q0-15.725-11-26.362Q611-653 596-653h-84v205Zm37-37v-131h47v131h-47Zm129 37h37v-83h50v-37h-50v-48h50v-37h-87v205ZM286-192q-39.05 0-66.525-27.475Q192-246.95 192-286v-542q0-39.463 27.475-67.231Q246.95-923 286-923h542q39.463 0 67.231 27.769Q923-867.463 923-828v542q0 39.05-27.769 66.525Q867.463-192 828-192H286ZM132-37q-39.463 0-67.231-27.769Q37-92.537 37-132v-636h95v636h636v95H132Z"
      />
    </svg>
  </figure>
  <div class="flex flex-col items-start justify-center pl-3">
    <p class=" w-60 truncate text-lg font-bold">{fileData?.filename}</p>
    <div class="flex w-60 flex-col items-start justify-center">
      <p class="w-60 truncate"><span class="font-bold">Owner: </span>{fileData?.owner.username}</p>
      <p class="truncate"><span class="font-bold">Size: </span>{fileData?.size}</p>
    </div>
  </div>
</section>
