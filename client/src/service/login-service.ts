import { error } from '@sveltejs/kit';
import type { LoginForm, LoginResponse } from '../stores/login';

export async function login(form: LoginForm): Promise<LoginResponse> {
  const response = await fetch(import.meta.env.VITE_LOGIN_URL, {
    method: 'POST',
    headers: { Authorization: 'Basic ' + btoa(`${form.username}:${form.password}`) }
  });

  if (response.status !== 200) {
    throw error(response.status, response.statusText);
  }

  return response.json();
}
