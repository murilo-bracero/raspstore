import { writable } from 'svelte/store';

export interface PageData {
  size: number;
  totalElements: number;
  page: number;
  next: string;
  content: FileData[];
}

export interface FileData {
  fileId: string;
  path: string;
  filename: string;
  size: number;
  owner: FileUserData;
  editors: FileUserData[];
  viewers: FileUserData[];
  createdAt: string;
  updatedAt: string;
  createdBy: FileUserData;
  updatedBy: FileUserData;
}

export interface FileUserData {
  userId: string;
  username: string;
}
