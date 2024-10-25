import { error } from '@sveltejs/kit';
import { coreURLs } from '../config/urls';

export async function uploadFile(body: FormData, token: string) {
  const response = await fetch(coreURLs.upload, {
    method: 'POST',
    body: body,
    headers: {
      Authorization: `Bearer ${token}`
    }
  });

  if (response.status === 201) {
    return;
  }

  throw error(response.status, response.statusText);
}
