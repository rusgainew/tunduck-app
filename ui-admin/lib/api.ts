import { api, apiClient } from './api-client';
import {
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  User,
  EsfOrganization,
  EsfDocument,
  PaginatedResponse,
} from './types';

// Auth API
export const authApi = {
  login: async (data: LoginRequest): Promise<AuthResponse> => {
    const response = await api.post<AuthResponse>('/api/auth/login', data);
    apiClient.setToken(response.data.token);
    return response.data;
  },

  register: async (data: RegisterRequest): Promise<AuthResponse> => {
    const response = await api.post<AuthResponse>('/api/auth/register', data);
    apiClient.setToken(response.data.token);
    return response.data;
  },

  logout: async (): Promise<void> => {
    await api.post('/api/auth/logout');
    if (typeof window !== 'undefined') {
      localStorage.removeItem('auth_token');
    }
  },

  getCurrentUser: async (): Promise<User> => {
    const response = await api.get<User>('/api/auth/me');
    return response.data;
  },
};

// Users API
export const usersApi = {
  getAll: async (params?: { page?: number; limit?: number }): Promise<PaginatedResponse<User>> => {
    const response = await api.get<PaginatedResponse<User>>('/api/users', { params });
    return response.data;
  },

  getById: async (id: string): Promise<User> => {
    const response = await api.get<User>(`/api/users/${id}`);
    return response.data;
  },
};

// Organizations API
export const organizationsApi = {
  getAll: async (): Promise<EsfOrganization[]> => {
    const response = await api.get<EsfOrganization[]>('/api/esf-organizations');
    return response.data;
  },

  getPaginated: async (params?: {
    page?: number;
    pageSize?: number;
    search?: string;
    sortBy?: string;
    sortOrder?: 'asc' | 'desc';
  }): Promise<PaginatedResponse<EsfOrganization>> => {
    const response = await api.get<PaginatedResponse<EsfOrganization>>(
      '/api/esf-organizations/paginated',
      { params }
    );
    return response.data;
  },

  getById: async (id: string): Promise<EsfOrganization> => {
    const response = await api.get<EsfOrganization>(`/api/esf-organizations/${id}`);
    return response.data;
  },

  create: async (data: Omit<EsfOrganization, 'id' | 'createdAt' | 'updatedAt'>): Promise<EsfOrganization> => {
    const response = await api.post<EsfOrganization>('/api/esf-organizations', data);
    return response.data;
  },

  update: async (id: string, data: Partial<EsfOrganization>): Promise<EsfOrganization> => {
    const response = await api.put<EsfOrganization>(`/api/esf-organizations/${id}`, data);
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/api/esf-organizations/${id}`);
  },
};

// Documents API
export const documentsApi = {
  getAll: async (orgId?: string): Promise<{ success: boolean; data: EsfDocument[]; count: number }> => {
    const params = orgId ? { org_id: orgId } : {};
    const response = await api.get<{ success: boolean; data: EsfDocument[]; count: number }>(
      '/api/esf-documents',
      { params }
    );
    return response.data;
  },

  getPaginated: async (params?: {
    page?: number;
    pageSize?: number;
    org_id?: string;
    search?: string;
    sortBy?: string;
    sortOrder?: 'asc' | 'desc';
  }): Promise<PaginatedResponse<EsfDocument>> => {
    const response = await api.get<PaginatedResponse<EsfDocument>>(
      '/api/esf-documents/paginated',
      { params }
    );
    return response.data;
  },

  getById: async (id: string, orgId?: string): Promise<EsfDocument> => {
    const params = orgId ? { org_id: orgId } : {};
    const response = await api.get<EsfDocument>(`/api/esf-documents/${id}`, { params });
    return response.data;
  },

  create: async (data: Partial<EsfDocument>, orgId?: string): Promise<EsfDocument> => {
    const params = orgId ? { org_id: orgId } : {};
    const response = await api.post<EsfDocument>('/api/esf-documents', data, { params });
    return response.data;
  },

  update: async (id: string, data: Partial<EsfDocument>, orgId?: string): Promise<EsfDocument> => {
    const params = orgId ? { org_id: orgId } : {};
    const response = await api.put<EsfDocument>(`/api/esf-documents/${id}`, data, { params });
    return response.data;
  },

  delete: async (id: string, orgId?: string): Promise<void> => {
    const params = orgId ? { org_id: orgId } : {};
    await api.delete(`/api/esf-documents/${id}`, { params });
  },
};
