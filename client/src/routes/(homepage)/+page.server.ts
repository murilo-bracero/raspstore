import { cookieKeys } from '$lib/config/cookies.js';
import { getPageFiles } from '$lib/services/file.service';
import { uploadFile } from '$lib/services/fs.service';
import { type PageData } from '$lib/stores/file';
import { ActionFailure, fail } from '@sveltejs/kit';

export async function load({ cookies }): Promise<PageData | ActionFailure> {
  const token = cookies.get(cookieKeys.accessToken);

  if (!token) {
    return fail(401);
  }

  return getPageFiles(token).catch((err) => fail(err.status, err));
}

export const actions = {
  upload: async ({ cookies, request }) => {
    const token = cookies.get(cookieKeys.accessToken);

    if (!token) {
      return fail(401);
    }

    const data = await request.formData();

    uploadFile(data, token);
  }
};
