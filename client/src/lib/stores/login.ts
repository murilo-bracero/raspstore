import { writable } from 'svelte/store';

export interface LoginForm {
  username: string;
  password: string;
}

export interface LoginResponse {
  accessToken: string;
  refreshToken: string;
}
