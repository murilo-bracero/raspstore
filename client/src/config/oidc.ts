import * as oidc from 'oidc-client-ts';

const UIHost = import.meta.env.UI_HOST || 'http://localhost:5173';
const clientId = process.env.RS_CLIENT_ID || import.meta.env.VITE_CLIENT_ID;
const clientSecret = process.env.RS_CLIENT_SECRET || import.meta.env.VITE_CLIENT_SECRET;
const additionalScopes =
  process.env.RS_ADDITIONAL_SCOPES || import.meta.env.VITE_ADDITIONAL_SCOPES || '';
const authorityURL = process.env.RS_AUTHORITY_URL || import.meta.env.VITE_AUTHORITY_URL;

export const settings = {
  authority: authorityURL,
  client_id: clientId,
  client_secret: clientSecret,
  redirect_uri: UIHost + '/auth-callback',
  post_logout_redirect_uri: UIHost + '/login',
  response_type: 'code',
  scope: 'openid email roles ' + additionalScopes,
  automaticSilentRenew: true,
  filterProtocolClaims: true,
  loadUserInfo: true,
  response_mode: 'fragment'
} as oidc.UserManagerSettings;
