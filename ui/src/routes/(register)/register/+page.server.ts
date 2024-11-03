import { coreURLs } from '$lib/config/urls';
import { fail } from '@sveltejs/kit';

export const actions = {
  register: async ({ request }) => {
    const data = await request.formData();
    const username = data.get('username') as string;
    const name = data.get('name') as string;
    const password = data.get('password') as string;
    const confirmPassword = data.get('confirm-password') as string;

    if (password !== confirmPassword) {
      return fail(400, { success: false, message: 'Passwords do not match' });
    }

    if (!username || !password) {
      return fail(400, { success: false, message: 'Username and password fields are required' });
    }

    return fetch(coreURLs.register, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        username,
        name,
        password
      })
    })
      .then(async (res) => {
        if (!res.ok) {
          const { status } = res;

          if (status === 422) {
            return fail(400, {
              success: false,
              message: 'User registration is disabled'
            });
          }
        }
        return { success: true };
      })
      .catch(() => {
        return fail(500, { success: false, message: 'Something went wrong' });
      });
  }
};
