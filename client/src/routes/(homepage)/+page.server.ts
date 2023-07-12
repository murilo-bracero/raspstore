import { fail, redirect } from '@sveltejs/kit';
import { getPageFiles } from '../../service/file-info-service';
import type { PageData } from '../../stores/file';
import { uploadFile } from '../../service/fs-service';

export async function load({ cookies }): Promise<PageData> {
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
