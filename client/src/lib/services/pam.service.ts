import { error } from '@sveltejs/kit';
import { LoginForm, LoginResponse } from '../stores/login';
import { coreURLs } from '../config/urls';

export async function pamSignIn(form: LoginForm): Promise<LoginResponse> {
  const response = await fetch(coreURLs.loginPAM, {
    method: 'POST',
    headers: { Authorization: 'Basic ' + btoa(`${form.username}:${form.password}`) }
  });

  if (response.status !== 200) {
    throw error(response.status, response.statusText);
  }

  return response.json();
}
