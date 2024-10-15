import { cookieKeys } from '$lib/config/cookies';
import { AuthService } from '$lib/services/auth.service';
import { clearAuthCookies, parseFromCookies, populateAuthCookies } from '$lib/utils/cookies.util';
import { redirect, type RequestEvent } from '@sveltejs/kit';

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
      (path) => event.url.pathname === path.path && path.method === event.request.method
    )
  ) {
    return resolve(event);
  }

  const tokenSet = parseFromCookies(event.cookies);

  if (!tokenSet.accessToken) {
    throw redirect(307, '/login');
  }

  const response = await resolve(event);

  if (response.status === 401 && tokenSet.refreshToken) {
    await refreshUserToken(event);

    return handle({ event, resolve });
  }

  return response;
};

async function refreshUserToken(event: RequestEvent) {
  const newTokenSet = await AuthService.instance
    .refresh(event.cookies.get(cookieKeys.refreshToken)!)
    .catch(() => {
      clearAuthCookies(event.cookies);

      throw redirect(307, '/login');
    });

  populateAuthCookies(event.cookies, {
    accessToken: newTokenSet.access_token!,
    idToken: newTokenSet.id_token!,
    refreshToken: newTokenSet.refresh_token!
  });
}
