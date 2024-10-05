import { cookieKeys } from '$lib/config/cookies';
import { AuthService } from '$lib/services/auth.service';
import { type Cookies, redirect, type RequestEvent } from '@sveltejs/kit';

const PUBLIC_PATHS = [
  {
    method: 'GET',
    path: '/api/login'
  },
  {
    method: 'GET',
    path: '/api/login/callback'
  },
  {
    method: 'GET',
    path: '/login'
  }
];

export const handle = async ({ event, resolve }) => {
  if (
    PUBLIC_PATHS.some(
      (path) => event.url.pathname.startsWith(path.path) && path.method === event.request.method
    )
  ) {
    return resolve(event);
  }

  if (!hasAccessTokenCookie(event)) {
    throw redirect(307, '/login');
  }

  const client = await AuthService.instance.getClient();

  await client
    .userinfo(event.cookies.get(cookieKeys.accessToken)!)
    .then((user) => {
      event.locals.user = user;
    })
    .catch(async (err) => {
      if (err.status === 401 && hasRefreshTokenCookie(event)) {
        await refreshUserToken(event);
        return handle({ event, resolve });
      }

      throw redirect(307, '/login');
    });

  const response = await resolve(event);

  if (response.status === 401 && hasRefreshTokenCookie(event)) {
    await refreshUserToken(event);

    return handle({ event, resolve });
  }

  return response;
};

function hasAccessTokenCookie(event: RequestEvent) {
  return event.cookies.get(cookieKeys.accessToken);
}

async function refreshUserToken(event: RequestEvent) {
  const newTokenSet = await AuthService.instance
    .refresh(event.cookies.get(cookieKeys.refreshToken)!)
    .catch((err) => {
      console.error('error refreshing token', err);

      clearAuthCookies(event.cookies);

      throw redirect(307, '/login');
    });

  setCookie(event.cookies, cookieKeys.accessToken, newTokenSet.access_token!);

  setCookie(event.cookies, cookieKeys.refreshToken, newTokenSet.refresh_token!);

  setCookie(event.cookies, cookieKeys.idToken, newTokenSet.id_token!);
}

function hasRefreshTokenCookie(event: RequestEvent) {
  return event.cookies.get(cookieKeys.refreshToken);
}

function setCookie(cookies: Cookies, key: string, value: string): void {
  cookies.set(key, value, {
    httpOnly: true,
    secure: true,
    sameSite: true,
    path: '/'
  });
}

function clearAuthCookies(cookies: Cookies) {
  Object.keys(cookieKeys).forEach((key) => {
    cookies.delete(key);
  });
}
