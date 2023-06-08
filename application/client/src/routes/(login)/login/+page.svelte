<script lang="ts">
  import { goto } from '$app/navigation';
  import { loginForm, type LoginForm, type LoginResponse } from '../../../stores/login';

  function handleLoginSubmit(event: any) {
    event.preventDefault();
    loginForm.subscribe(loginFormSubscriptionAction);
  }

  function loginFormSubscriptionAction(form: LoginForm) {
    if (!validateLoginForm(form)) {
      throw new Error('Login Form invalid');
    }

    fetch(import.meta.env.VITE_LOGIN_URL, {
      method: 'POST',
      headers: { Authorization: 'Basic ' + btoa(`${form.username}:${form.password}`) }
    })
      .then((res) => res.json())
      .then((res: LoginResponse) => {
        storeToken(res.accessToken);
        goto('/');
      })
      .catch((err) => console.log(err));
  }

  function validateLoginForm(form: LoginForm): boolean {
    return [form.password, form.username].filter((field) => field && field === '').length === 0;
  }

  function storeToken(token: string) {
    localStorage.setItem(import.meta.env.VITE_TOKEN_KEY, token);
  }
</script>

<main class="flex h-screen items-center justify-center">
  <form
    class="flex flex-col items-start justify-center rounded-lg border-2 border-black p-6"
    on:submit={handleLoginSubmit}
  >
    <label class="mb-1 p-1 font-bold" for="username">Username:</label>
    <input
      class="rounded-lg border-2 border-black p-1"
      type="text"
      id="username"
      bind:value={$loginForm.username}
    />
    <label class="mb-1 mt-3 p-1 font-bold" for="password">Password:</label>
    <input
      class="rounded-lg border-2 border-black p-1"
      type="password"
      id="password"
      bind:value={$loginForm.password}
    />
    <button
      class="my-3 self-end rounded-full border-2 border-black bg-sky-400 px-3 py-1 font-bold"
      type="submit">Log in</button
    >
  </form>
</main>
