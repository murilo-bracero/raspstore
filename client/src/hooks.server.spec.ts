import { assert, describe, it } from 'vitest';
import { handle } from './hooks.server';
import { cookieKeys } from '$lib/config/cookies';
import { RequestEvent } from '@sveltejs/kit';

describe('test handle', () => {
  it('should handle public endpoints', async () => {
    const event = {
      url: {
        pathname: 'invalid'
      },
      request: {
        method: 'invalid'
      }
    };

    const testTable = [
      {
        method: 'GET',
        path: '/api/login',
        fail: false
      },
      {
        method: 'GET',
        path: '/api/login/callback',
        fail: false
      },
      {
        method: 'GET',
        path: '/login',
        fail: false
      },
      {
        method: 'POST',
        path: '/login',
        fail: true
      },
      {
        method: 'GET',
        path: '/',
        fail: true
      }
    ];

    for (const tt of testTable) {
      event.url.pathname = tt.path;
      event.request.method = tt.method;

      const resolve = () => {
        if (tt.fail) {
          assert.fail('should not resolve');
        }
      };

      try {
        await handle({ event, resolve });
      } catch (e) {
        if (!tt.fail) {
          assert.fail('should resolve');
        }
      }
    }
  });

  it('should success if route is not public and tokens were provided', async () => {
    const event = {
      url: {
        pathname: '/'
      },
      request: {
        method: 'GET'
      },
      cookies: {
        get: (key: string) => {
          return {
            [cookieKeys.accessToken]: 'accessToken',
            [cookieKeys.idToken]: 'idToken',
            [cookieKeys.refreshToken]: 'refreshToken'
          }[key];
        }
      }
    };

    const resolve = () => {
      return { status: 200 };
    };

    await handle({ event, resolve });
  });

  it('should fail if route is not public and no tokens', async () => {
    const event = {
      url: {
        pathname: '/'
      },
      request: {
        method: 'GET'
      },
      cookies: {
        get: (key: string) => {
          return {}[key];
        }
      }
    };

    const resolve = () => {
      assert.fail('should not resolve');
    };

    try {
      await handle({ event, resolve });
    } catch (e: any) {
      assert.equal(e.status, 307);
      assert.equal(e.location, '/login');
    }
  });

  it('should fail if resolve returns 401 and tokens were provided', async () => {
    const event = {
      url: {
        pathname: '/'
      },
      request: {
        method: 'GET'
      },
      cookies: {
        _cookies: {
          [cookieKeys.accessToken]: 'accessToken',
          [cookieKeys.idToken]: 'idToken',
          [cookieKeys.refreshToken]: 'refreshToken'
        },
        get: (key: string) => {
          return event.cookies._cookies[key];
        },
        delete: (key: string) => {
          delete event.cookies._cookies[key];
        }
      }
    };

    const resolve = () => {
      return { status: 401 };
    };

    try {
      await handle({ event, resolve });

      assert.fail('should not resolve');
    } catch (e: any) {
      assert.equal(e.status, 307);
      assert.equal(e.location, '/login');
    }
  });
});
