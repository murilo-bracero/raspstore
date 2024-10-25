import { json } from '@sveltejs/kit';
import { settings } from '../../../lib/config/oidc';
import { AuthService } from '$lib/services/auth.service';

export async function GET() {
  const client = await AuthService.instance.getClient();

  const authUrl = client.authorizationUrl({
    scope: settings.scope
  });

  return json({
    authUrl: authUrl
  });
}
