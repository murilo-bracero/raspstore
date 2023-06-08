import { error } from '@sveltejs/kit';
import { getToken } from './token-service';

export async function uploadFile(body: FormData) {
  const token = getToken();

  if (token == null) {
    throw new Error('Unauthorized');
  }

  const response = await fetch(import.meta.env.VITE_FS_SERVICE_URL, {
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
