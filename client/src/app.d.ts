// See https://kit.svelte.dev/docs/types#app

import { UserinfoResponse } from 'openid-client';

// for information about these interfaces
declare global {
  namespace App {
    // interface Error {}
    interface Locals {
      user: UserinfoResponse;
    }
    // interface PageData {}
    // interface Platform {}
  }
}

export {};
