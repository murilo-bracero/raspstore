import { settings } from '../../../../lib/config/oidc';
import { Cookies, redirect } from '@sveltejs/kit';
import { AuthService } from '$lib/services/auth.service';
import { cookieKeys } from '$lib/config/cookies';

export async function GET({ url, cookies }) {
  const client = await AuthService.instance.getClient();

  const params = client.callbackParams(url.toString());

  const tokenSet = await client.callback(settings.redirect_uri, params, {
    response_type: settings.response_type
  });

  setCookie(cookies, cookieKeys.accessToken, tokenSet.access_token!);

  setCookie(cookies, cookieKeys.refreshToken, tokenSet.refresh_token!);

  setCookie(cookies, cookieKeys.idToken, tokenSet.id_token!);

  throw redirect(302, '/');
}

function setCookie(cookies: Cookies, key: string, value: string): void {
  cookies.set(key, value, {
    httpOnly: true,
    secure: true,
    sameSite: true,
    path: '/'
  });
}
