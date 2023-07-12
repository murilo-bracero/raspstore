import { derived, writable } from 'svelte/store';

export interface Notification {
  message: string;
  type: NotificationType;
}

export enum NotificationType {
  ERROR = 'ERROR',
  SUCCESS = 'SUCCESS'
}

export const notifications = writable<Notification[]>([]);

export function toast(notification: Notification) {
  notifications.update((state) => {
    if (state.filter((n) => n.message === notification.message).length > 0) {
      return state;
    }
    setTimeout(removeToast, 3000);
    return [notification, ...state];
  });
}

function removeToast() {
  notifications.update((state) => {
    return [...state.slice(0, state.length - 1)];
  });
}
