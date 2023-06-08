<script lang="ts">
  import { goto } from '$app/navigation';
  import { clickOutside } from '../directives/clickOutsideDirective';
  import { uploadFile } from '../service/fs-service';
  import { NotificationType, toast } from '../stores/toast';

  export let open = true;

  let file: any;
  let path: string;
  let choosenFile = '';

  function handleCloseClick() {
    open = false;
  }

  function handleFileChange(event: any) {
    file = event.target.files[0];
    choosenFile = file.name;
  }

  function handleSubmitForm(event: any) {
    event.preventDefault();

    const formData = createFormData();

    uploadFile(formData)
      .then(() => {
        handleCloseClick();
        toast({
          message: 'Uploaded successfully',
          type: NotificationType.SUCCESS
        });
      })
      .catch((err) => {
        if (err.status === 401) {
          goto('/login');
          return;
        }

        toast({
          message: 'Could not upload file',
          type: NotificationType.ERROR
        });
      });
  }

  function createFormData(): FormData {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('path', path);
    return formData;
  }
</script>

<section
  class="bg-red absolute bottom-[-13rem] z-50 h-52 w-screen border-t-2 border-black bg-white transition-[bottom] duration-300 ease-in-out"
  class:open
  use:clickOutside
  on:click_outside={handleCloseClick}
>
  <form class="flex flex-col p-6" on:submit={handleSubmitForm}>
    <div class="mb-3 flex flex-row items-center justify-start">
      <label
        class="flex-3 mr-2 rounded-lg border-2 border-black bg-sky-400 px-3 py-1 font-bold"
        for="file_input"
      >
        Choose file
      </label>
      <p class="flex-1 truncate">{choosenFile ? choosenFile : 'No files'}</p>
      <input type="file" id="file_input" on:change={handleFileChange} hidden />
    </div>
    <input
      type="text"
      bind:value={path}
      placeholder="Folder path i.e. /folder1"
      class="rounded-lg border-2 border-black p-1"
    />
    <div class="flex items-center justify-center">
      <button
        type="submit"
        class="mt-3 rounded-lg border-2 border-black bg-sky-400 px-3 py-1 font-bold">Upload</button
      >
    </div>
  </form>
</section>

<style>
  .open {
    bottom: 0;
  }
</style>
