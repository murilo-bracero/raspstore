import { error } from '@sveltejs/kit';
import type { PageData } from '../stores/file';
import { coreURLs } from '../config/urls';

export async function getPageFiles(token: string): Promise<PageData> {
  const response = await fetch(coreURLs.files, {
    headers: {
      Authorization: `Bearer ${token}`
    }
  });

  if (response.status !== 200) {
    throw error(response.status, response.statusText);
  }

  return response.json();
}
