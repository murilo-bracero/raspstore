import { error } from '@sveltejs/kit';

export async function uploadFile(body: FormData, token: string) {
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
