import { error } from '@sveltejs/kit';
import type { FileData, PageData } from '../stores/file';
import { getToken } from './token-service';

export async function getFiles(): Promise<FileData[]> {
  const token = getToken();

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

  return response.json().then((body: PageData) => body.content);
}
