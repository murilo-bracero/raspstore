import { getPageFiles } from '$lib/services/file.service';
import { uploadFile } from '$lib/services/fs.service';
import { PageData } from '$lib/stores/file';
import { fail, redirect } from '@sveltejs/kit';

export async function load({ cookies, request, locals }): Promise<PageData> {
  console.log(locals);
  const token = cookies.get('jwt-token');

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
    const token = cookies.get('jwt-token');

    if (!token) {
      throw redirect(307, '/login');
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
