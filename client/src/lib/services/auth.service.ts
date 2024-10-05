import { BaseClient, Issuer, TokenSet } from 'openid-client';
import { settings } from '../config/oidc';

export class AuthService {
  private static _instance: AuthService;

  static get instance(): AuthService {
    if (!AuthService._instance) {
      AuthService._instance = new AuthService();
    }

    return AuthService._instance;
  }

  private client: BaseClient | undefined;

  private constructor() {}

  async getClient() {
    if (!this.client) {
      const iss = await Issuer.discover(settings.authority);

      this.client = new iss.Client({
        client_id: settings.client_id,
        client_secret: settings.client_secret,
        redirect_uris: [settings.redirect_uri],
        response_types: [settings.response_type],
        post_logout_redirect_uris: [settings.post_logout_redirect_uri]
      });
    }

    return this.client;
  }

  async refresh(token: string): Promise<TokenSet> {
    const client = await this.getClient();
    return client.refresh(token);
  }
}
