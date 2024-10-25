import { coreURLs } from '$lib/config/urls';
import { http, HttpResponse } from 'msw';
import { setupServer } from 'msw/node';
import { assert, beforeAll, describe, expect, it } from 'vitest';
import { uploadFile } from './fs.service';

coreURLs.upload = 'http://files/v1/uploads';

const handler = http.post(coreURLs.upload, ({ request }) => {
  if (request.headers.get('Authorization') !== 'Bearer access_token') {
    return new HttpResponse(null, { status: 401 });
  }

  return HttpResponse.json({ content: [] }, { status: 201 });
});

const server = setupServer(handler);

beforeAll(() => {
  server.listen({
    onUnhandledRequest: 'error'
  });
});

describe('file.service.ts testing', async () => {
  it('should call upload API endpoint with 201 response', async () => {
    const token = 'access_token';
    const formData = new FormData();
    formData.append('file', new Blob());

    await uploadFile(formData, token).catch(assert.fail);
  });

  it('should call API and throw error if status is not 201', async () => {
    const token = 'access_token_1';
    const formData = new FormData();
    formData.append('file', new Blob());

    await uploadFile(formData, token)
      .then(() => assert.fail())
      .catch((err) => assert.equal(err.status, 401));
  });
});
