import { create } from 'zustand'
import type { Project, ProjectCreateRequest } from '../types/project';
import { httpClient } from '../lib/httpClient';

interface ProjectStore {
  projects: Project[];
  fetchProjects: () => Promise<void>;
  createProject: (project: ProjectCreateRequest) => Promise<void>;
  uploadDocument: (projectId: string, file: File) => Promise<void>;
  askQuestion: (question: { question: string; projectId: string }) => Promise<void>;
}

export const useProjectStore = create<ProjectStore>((set, get) => ({
  projects: [],

  fetchProjects: async () => {
    const response = await httpClient.get<Project[]>(`${import.meta.env.VITE_API_URL}/nummer`)
    if (response.status === 200) {
      set({ projects: response.data })
    }
    else {
      throw new Error("Failed to fetch projects")
    }
  },

  createProject: async (newProject: ProjectCreateRequest) => {
    const response = await httpClient.post<Project>(`${import.meta.env.VITE_API_URL}/project/create-project`, newProject)
    if (response.status === 201) {
      set((state) => ({ projects: [...state.projects, response.data] }))
    } else {
      throw new Error("Failed to create project")
    }
  },

  uploadDocument: async (projectId: string, file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('projectId', projectId);

    const response = await httpClient.post(`${import.meta.env.VITE_API_URL}/project/add-document`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      params: { projectId },
    });

    if (response.status !== 200) {
      throw new Error("Failed to upload document");
    }
  },
  askQuestion: async (question) => {
    const response = await httpClient.post(`${import.meta.env.VITE_API_URL}/project/question`, question);
    if (response.status !== 200) {
      throw new Error("Failed to ask question");
    }
    return response.data;
  }
}))






