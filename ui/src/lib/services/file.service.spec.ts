import { afterAll, afterEach, assert, beforeAll, describe, expect, it } from 'vitest';
import { getPageFiles } from './file.service';
import { http, HttpResponse } from 'msw';
import { setupServer } from 'msw/node';
import { coreURLs } from '$lib/config/urls';

const handler = http.get('http://files/v1/files', ({ request }) => {
  if (request.headers.get('Authorization') !== 'Bearer access_token') {
    return new HttpResponse(null, { status: 401 });
  }

  return HttpResponse.json({ content: [] });
});

coreURLs.files = 'http://files/v1/files';

const server = setupServer(handler);

beforeAll(() => {
  server.listen({
    onUnhandledRequest: 'error'
  });
});

afterAll(() => server.close());

afterEach(() => server.resetHandlers());

describe('file.service.ts testing', async () => {
  it('should call API and if success return a list of files', async () => {
    const token = 'access_token';

    const result = await getPageFiles(token);

    expect(result).toBeInstanceOf(Object);
    expect(result.content).toBeInstanceOf(Array);
  });

  it('should call API and throw error if status is not 200', async () => {
    const token = 'wrong_token';

    await getPageFiles(token)
      .then(() => assert.fail())
      .catch((err) => assert.equal(err.status, 401));
  });
});
