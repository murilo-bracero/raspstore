export function getToken(): string | null {
  return localStorage.getItem(import.meta.env.VITE_TOKEN_KEY);
}

export function setToken(token: string) {
  localStorage.setItem(import.meta.env.VITE_TOKEN_KEY, token);
}
