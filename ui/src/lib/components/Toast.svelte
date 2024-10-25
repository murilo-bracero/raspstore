<script>
  import { NotificationType, notifications } from '../stores/toast';
</script>

<section class="fixed left-0 top-0 z-[100] flex flex-col items-center justify-start">
  {#each $notifications as notification (notification.message)}
    <div
      class="w-screen border-b-2 border-black p-2"
      class:error={notification.type === NotificationType.ERROR}
      class:success={notification.type === NotificationType.SUCCESS}
    >
      <p class="font-bold">{notification.message}</p>
    </div>
  {/each}
</section>

<style>
  .error {
    @apply bg-yellow-500;
  }

  .success {
    @apply bg-green-500;
  }
</style>
