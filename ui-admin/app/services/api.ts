/**
 * API Service Client
 * Клиент для работы с API tunduct-system
 */

// Safely get API base URL (works in both browser and server)
const API_BASE_URL = typeof window !== 'undefined' 
  ? (process.env.VITE_API_BASE_URL || 'http://localhost:8080/api')
  : 'http://localhost:8080/api';

// Types for API responses
export interface AuthResponse {
  token: string;
  user: UserInfo;
}

export interface UserInfo {
  id: string;
  username: string;
  email: string;
  fullName: string;
  phone: string;
  role: string;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface User extends UserInfo {}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
}

export interface Organization {
  id: string;
  name: string;
  description: string;
  token: string;
  dbName: string;
  createdAt: string;
  updatedAt: string;
}

export interface OrganizationPaginatedResponse {
  data: Organization[];
  pagination: {
    page: number;
    page_size: number;
    total_items: number;
    total_pages: number;
    has_next: boolean;
    has_prev: boolean;
  };
}

export interface EsfDocument {
  id: string;
  organizationId: string;
  title: string;
  content: string;
  status: string;
  createdAt: string;
  updatedAt: string;
}

export interface ErrorResponse {
  code: string;
  message: string;
  details?: string;
}

// Auth API
export const authApi = {
  register: async (data: {
    username: string;
    email: string;
    password: string;
    confirmPassword: string;
    fullName: string;
    phone: string;
  }): Promise<AuthResponse> => {
    const response = await fetch(`${API_BASE_URL}/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error(`Registration failed: ${response.statusText}`);
    return response.json();
  },

  login: async (username: string, password: string): Promise<AuthResponse> => {
    const response = await fetch(`${API_BASE_URL}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password }),
    });
    if (!response.ok) throw new Error(`Login failed: ${response.statusText}`);
    return response.json();
  },

  getCurrentUser: async (token: string): Promise<UserInfo> => {
    const response = await fetch(`${API_BASE_URL}/auth/me`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    if (!response.ok) throw new Error('Failed to fetch current user');
    return response.json();
  },

  logout: async (token: string): Promise<void> => {
    const response = await fetch(`${API_BASE_URL}/auth/logout`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
    });
    if (!response.ok) throw new Error('Logout failed');
  },
};

// Users API
export const usersApi = {
  getAll: async (page: number = 1, limit: number = 10): Promise<PaginatedResponse<User>> => {
    const response = await fetch(`${API_BASE_URL}/users?page=${page}&limit=${limit}`);
    if (!response.ok) throw new Error('Failed to fetch users');
    return response.json();
  },

  getById: async (id: string): Promise<User> => {
    const response = await fetch(`${API_BASE_URL}/users/${id}`);
    if (!response.ok) throw new Error('Failed to fetch user');
    return response.json();
  },

  update: async (id: string, data: Partial<User>, token: string): Promise<User> => {
    const response = await fetch(`${API_BASE_URL}/users/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error('Failed to update user');
    return response.json();
  },

  delete: async (id: string, token: string): Promise<void> => {
    const response = await fetch(`${API_BASE_URL}/users/${id}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token}` },
    });
    if (!response.ok) throw new Error('Failed to delete user');
  },
};

// Organizations API
export const organizationsApi = {
  getAll: async (): Promise<Organization[]> => {
    const response = await fetch(`${API_BASE_URL}/esf-organizations`);
    if (!response.ok) throw new Error('Failed to fetch organizations');
    return response.json();
  },

  getPaginated: async (
    page: number = 1,
    limit: number = 10,
    sort: string = 'created_at',
    order: string = 'desc'
  ): Promise<OrganizationPaginatedResponse> => {
    const params = new URLSearchParams({
      page: page.toString(),
      page_size: limit.toString(),
      sort,
      order,
    });
    const response = await fetch(`${API_BASE_URL}/esf-organizations/paginated?${params}`);
    if (!response.ok) throw new Error('Failed to fetch organizations');
    return response.json();
  },

  getById: async (id: string): Promise<Organization> => {
    const response = await fetch(`${API_BASE_URL}/esf-organizations/${id}`);
    if (!response.ok) throw new Error('Failed to fetch organization');
    return response.json();
  },

  create: async (data: Partial<Organization>, token: string): Promise<Organization> => {
    const response = await fetch(`${API_BASE_URL}/esf-organizations`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error('Failed to create organization');
    return response.json();
  },

  update: async (id: string, data: Partial<Organization>, token: string): Promise<Organization> => {
    const response = await fetch(`${API_BASE_URL}/esf-organizations/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error('Failed to update organization');
    return response.json();
  },

  delete: async (id: string, token: string): Promise<void> => {
    const response = await fetch(`${API_BASE_URL}/esf-organizations/${id}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token}` },
    });
    if (!response.ok) throw new Error('Failed to delete organization');
  },
};

// Documents API
export const documentsApi = {
  getAll: async (organizationId?: string): Promise<EsfDocument[]> => {
    let url = `${API_BASE_URL}/esf-documents`;
    if (organizationId) {
      url += `?organization_id=${organizationId}`;
    }
    const response = await fetch(url);
    if (!response.ok) throw new Error('Failed to fetch documents');
    return response.json();
  },

  getPaginated: async (
    page: number = 1,
    limit: number = 10,
    organizationId?: string
  ): Promise<any> => {
    const params = new URLSearchParams({
      page: page.toString(),
      page_size: limit.toString(),
    });
    if (organizationId) {
      params.append('organization_id', organizationId);
    }
    const response = await fetch(`${API_BASE_URL}/esf-documents/paginated?${params}`);
    if (!response.ok) throw new Error('Failed to fetch documents');
    return response.json();
  },

  getById: async (id: string): Promise<EsfDocument> => {
    const response = await fetch(`${API_BASE_URL}/esf-documents/${id}`);
    if (!response.ok) throw new Error('Failed to fetch document');
    return response.json();
  },

  create: async (data: Partial<EsfDocument>, token: string): Promise<EsfDocument> => {
    const response = await fetch(`${API_BASE_URL}/esf-documents`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error('Failed to create document');
    return response.json();
  },

  update: async (id: string, data: Partial<EsfDocument>, token: string): Promise<EsfDocument> => {
    const response = await fetch(`${API_BASE_URL}/esf-documents/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error('Failed to update document');
    return response.json();
  },

  delete: async (id: string, token: string): Promise<void> => {
    const response = await fetch(`${API_BASE_URL}/esf-documents/${id}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token}` },
    });
    if (!response.ok) throw new Error('Failed to delete document');
  },
};

// Roles API
export const rolesApi = {
  getUserRole: async (userId: string, token: string): Promise<{ role: string }> => {
    const response = await fetch(`${API_BASE_URL}/roles/${userId}`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    if (!response.ok) throw new Error('Failed to fetch user role');
    return response.json();
  },

  assignRole: async (userId: string, role: string, token: string): Promise<void> => {
    const response = await fetch(`${API_BASE_URL}/roles/assign`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ user_id: userId, role }),
    });
    if (!response.ok) throw new Error('Failed to assign role');
  },

  updateRole: async (userId: string, role: string, token: string): Promise<void> => {
    const response = await fetch(`${API_BASE_URL}/roles/${userId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ role }),
    });
    if (!response.ok) throw new Error('Failed to update role');
  },

  checkPermissions: async (userId: string, token: string): Promise<string[]> => {
    const response = await fetch(`${API_BASE_URL}/roles/permissions/${userId}`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    if (!response.ok) throw new Error('Failed to check permissions');
    return response.json();
  },

  listUsersByRole: async (token: string): Promise<any> => {
    const response = await fetch(`${API_BASE_URL}/roles/users`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    if (!response.ok) throw new Error('Failed to fetch users by role');
    return response.json();
  },
};
