import { cookieKeys } from '$lib/config/cookies';
import { Cookies } from '@sveltejs/kit';

interface CookiesTokenSet {
  accessToken: string;
  idToken: string;
  refreshToken: string;
}

export function parseFromCookies(cookies: Cookies): CookiesTokenSet {
  return {
    accessToken: cookies.get(cookieKeys.accessToken) || '',
    idToken: cookies.get(cookieKeys.idToken) || '',
    refreshToken: cookies.get(cookieKeys.refreshToken) || ''
  };
}

export function populateAuthCookies(cookies: Cookies, tokenSet: CookiesTokenSet) {
  setCookie(cookies, cookieKeys.accessToken, tokenSet.accessToken);
  setCookie(cookies, cookieKeys.idToken, tokenSet.idToken);
  setCookie(cookies, cookieKeys.refreshToken, tokenSet.refreshToken);
}

export function clearAuthCookies(cookies: Cookies) {
  Object.keys(cookieKeys).forEach((key) => {
    // @ts-ignore
    cookies.delete(key);
  });
}

function setCookie(cookies: Cookies, key: string, value: string): void {
  cookies.set(key, value, {
    httpOnly: true,
    secure: true,
    sameSite: true,
    path: '/'
  });
}
