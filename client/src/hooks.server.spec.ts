import { assert, describe, it } from 'vitest';
import { handle } from './hooks.server';
import { cookieKeys } from '$lib/config/cookies';
import { AuthService } from '$lib/services/auth.service';
import { RequestEvent } from '@sveltejs/kit';
import { BaseClient, TokenSet } from 'openid-client';

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

  it('should success refresh tokens if resolve returns 401 and refresh/access tokens were provided', async () => {
    AuthService.instance.getClient = async () => {
      return {
        refresh: (rt: string) => {
          if (rt === 'refreshToken') {
            return {
              access_token: 'new_accessToken',
              id_token: 'new_idToken',
              refresh_token: 'new_refreshToken'
            } as TokenSet | undefined;
          }

          assert.fail('wrong refresh token');
        }
      } as unknown as BaseClient;
    };

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
        },
        set(key: string, value: string) {
          event.cookies._cookies[key] = value;
          assert.oneOf(key, Object.values(cookieKeys));
        }
      }
    };

    const resolve = (event: RequestEvent) => {
      if (event.cookies.get(cookieKeys.accessToken) === 'accessToken') {
        return { status: 401 };
      }

      return { status: 200 };
    };

    await handle({ event, resolve });
  });

  it('should fail and redirect to login if refresh tokens failed', async () => {
    AuthService.instance.getClient = async () => {
      return {
        refresh: (_: string) => {
          throw new Error('failed');
        }
      } as unknown as BaseClient;
    };

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
        },
        set(_: string, __: string) {
          assert.fail('should not set cookies');
        }
      }
    };

    const resolve = () => {
      return { status: 401 };
    };

    try {
      await handle({ event, resolve });
    } catch (e: any) {
      assert.equal(e.status, 307);
      assert.equal(e.location, '/login');
    }
  });
});
