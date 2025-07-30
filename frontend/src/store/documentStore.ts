import { create } from 'zustand'
import { httpClient } from '../lib/httpClient';
import type { Document } from '../types';

interface DocumentStore {
  documents: Document[];
  fetchDocuments: (projectId: string) => Promise<void>;
  uploadDocument: (projectId: string, file: File) => Promise<void>;
}

export const useDocumentStore = create<DocumentStore>((set, get) => ({
  documents: [],

  fetchDocuments: async (projectId: string) => {
    const response = await httpClient.get<Document[]>(`${import.meta.env.VITE_API_URL}/document/list-documents?projectId=${projectId}`);
    if (response.status === 200) {
      set({ documents: response.data })
    }
    else {
      throw new Error("Failed to fetch documents")
    }
  },
  uploadDocument: async (projectId: string, file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('projectId', projectId);

    const response = await httpClient.post(`${import.meta.env.VITE_API_URL}/document/add-document`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });

    if (response.status !== 200) {
      throw new Error("Failed to upload document");
    }
  },
}))






