import { cookieKeys } from '$lib/config/cookies.js';
import { getPageFiles } from '$lib/services/file.service';
import { uploadFile } from '$lib/services/fs.service';
import { type PageData } from '$lib/stores/file';
import { fail, redirect } from '@sveltejs/kit';

export async function load({ cookies, request, locals }): Promise<PageData> {
  const token = cookies.get(cookieKeys.accessToken);

  if (!token) {
    throw redirect(307, '/login');
  }

  return getPageFiles(token).catch((err) => {
    if (err.status === 401) {
      //TODO:try refresh
      throw redirect(307, '/login');
    }

    throw err;
  });
}

export const actions = {
  upload: async ({ cookies, request }) => {
    const token = cookies.get(cookieKeys.accessToken);

    if (!token) {
      //throw redirect(307, '/login');
      return;
    }

    const data = await request.formData();

    uploadFile(data, token).catch((err) => {
      if (err.status === 401) {
        //TODO: try refresh
        throw redirect(307, '/login');
      }

      return fail(err.status, { error: true });
    });
  }
};
