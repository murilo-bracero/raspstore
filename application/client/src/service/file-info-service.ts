import { error } from '@sveltejs/kit';
import type { PageData } from '../stores/file';

export async function getPageFiles(token: string): Promise<PageData> {
  if (token === null) {
    throw new Error('Unauthorized');
  }

  const response = await fetch(import.meta.env.VITE_FILES_SERVICE_URL, {
    headers: {
      Authorization: `Bearer ${token}`
    }
  });

  if (response.status !== 200) {
    throw error(response.status, response.statusText);
  }

  return response.json();
}
