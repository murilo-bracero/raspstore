import { User, UserManager } from 'oidc-client-ts';
import { writable } from 'svelte/store';

export interface LoginForm {
  username: string;
  password: string;
}

export interface LoginResponse {
  accessToken: string;
  refreshToken: string;
}

export const isAuthenticated = writable<Boolean>(false);
export const userStore = writable<User | undefined>();
