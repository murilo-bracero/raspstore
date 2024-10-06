import { pamSignIn } from '$lib/services/pam.service';
import { type LoginForm } from '$lib/stores/login';
import { fail } from '@sveltejs/kit';

export const actions = {
  default: async ({ cookies, request }) => {
    const data = await request.formData();
    const loginForm = parseRequest(data);

    if (!isFormValid(loginForm)) {
      return fail(400, { invalid: true });
    }

    const response = await pamSignIn(loginForm);

    cookies.set('jwt-token', response.accessToken, {
      httpOnly: true,
      sameSite: 'strict',
      secure: true,
      path: '/'
    });
  }
};

function parseRequest(data: FormData): LoginForm {
  const { username, password } = Object.fromEntries(data);

  return { username: username as string, password: password as string };
}

function isFormValid({ username, password }: LoginForm): boolean {
  return [username, password].filter((field) => !field || field.trim() === '').length === 0;
}
