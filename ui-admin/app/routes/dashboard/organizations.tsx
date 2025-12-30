
import React, { useState, useEffect } from 'react';
import type { Organization, OrganizationPaginatedResponse } from '../../services/api';
import { organizationsApi } from '../../services/api';

export default function OrganizationsPage() {
  const [organizations, setOrganizations] = useState<Organization[]>([]);
  const [pagination, setPagination] = useState({
    page: 1,
    page_size: 10,
    total_items: 0,
    total_pages: 1,
    has_prev: false,
    has_next: false,
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [showForm, setShowForm] = useState(false);
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    token: '',
    dbName: '',
  });

  const fetchOrganizations = async (page: number = 1) => {
    setLoading(true);
    setError('');
    try {
      const response: OrganizationPaginatedResponse = await organizationsApi.getPaginated(
        page,
        pagination.page_size
      );
      setOrganizations(response.data);
      setPagination(response.pagination);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch organizations');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchOrganizations(1);
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const token = localStorage.getItem('token');
    if (!token) {
      setError('Not authenticated');
      return;
    }

    try {
      setLoading(true);
      await organizationsApi.create(formData, token);
      setFormData({ name: '', description: '', token: '', dbName: '' });
      setShowForm(false);
      await fetchOrganizations(1);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create organization');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!window.confirm('Are you sure you want to delete this organization?')) return;

    const token = localStorage.getItem('token');
    if (!token) {
      setError('Not authenticated');
      return;
    }

    try {
      setLoading(true);
      await organizationsApi.delete(id, token);
      await fetchOrganizations(pagination.page);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete organization');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="p-8">
      <div className="mb-8 flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 mb-2">Organizations</h1>
          <p className="text-gray-600">Manage ESF organizations</p>
        </div>
        <button
          onClick={() => setShowForm(!showForm)}
          className="bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-4 rounded-lg transition"
        >
          {showForm ? 'Cancel' : 'Create Organization'}
        </button>
      </div>

      {error && (
        <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">
          {error}
        </div>
      )}

      {showForm && (
        <div className="mb-8 bg-white p-6 rounded-lg shadow">
          <h2 className="text-xl font-bold mb-4">Create New Organization</h2>
          <form onSubmit={handleSubmit} className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Organization Name
              </label>
              <input
                type="text"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Database Name
              </label>
              <input
                type="text"
                value={formData.dbName}
                onChange={(e) => setFormData({ ...formData, dbName: e.target.value })}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none"
                required
              />
            </div>

            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Description
              </label>
              <textarea
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none"
                rows={3}
              />
            </div>

            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Token
              </label>
              <input
                type="text"
                value={formData.token}
                onChange={(e) => setFormData({ ...formData, token: e.target.value })}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none"
                required
              />
            </div>

            <button
              type="submit"
              disabled={loading}
              className="md:col-span-2 bg-green-600 hover:bg-green-700 disabled:bg-gray-400 text-white font-semibold py-2 px-4 rounded-lg transition"
            >
              {loading ? 'Creating...' : 'Create Organization'}
            </button>
          </form>
        </div>
      )}

      {loading && !showForm ? (
        <div className="text-center py-8">
          <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          <p className="mt-2 text-gray-600">Loading organizations...</p>
        </div>
      ) : (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {organizations.length === 0 ? (
              <div className="col-span-full text-center py-8 text-gray-500">
                No organizations found
              </div>
            ) : (
              organizations.map((org) => (
                <div key={org.id} className="bg-white p-6 rounded-lg shadow hover:shadow-lg transition">
                  <h3 className="text-lg font-bold text-gray-900 mb-2">{org.name}</h3>
                  <p className="text-gray-600 text-sm mb-3">{org.description}</p>
                  <div className="space-y-2 mb-4 text-xs text-gray-500">
                    <p>
                      <strong>Database:</strong> {org.dbName}
                    </p>
                    <p>
                      <strong>Token:</strong> {org.token.substring(0, 20)}...
                    </p>
                    <p>
                      <strong>Created:</strong> {new Date(org.createdAt).toLocaleDateString()}
                    </p>
                  </div>
                  <div className="flex gap-2">
                    <button className="flex-1 text-blue-600 hover:text-blue-900 font-medium text-sm py-2 rounded bg-blue-50 hover:bg-blue-100 transition">
                      Edit
                    </button>
                    <button
                      onClick={() => handleDelete(org.id)}
                      className="flex-1 text-red-600 hover:text-red-900 font-medium text-sm py-2 rounded bg-red-50 hover:bg-red-100 transition"
                    >
                      Delete
                    </button>
                  </div>
                </div>
              ))
            )}
          </div>

          {/* Pagination */}
          {pagination.total_pages > 1 && (
            <div className="mt-6 flex items-center justify-between">
              <div className="text-sm text-gray-600">
                Showing {(pagination.page - 1) * pagination.page_size + 1} to{' '}
                {Math.min(pagination.page * pagination.page_size, pagination.total_items)} of{' '}
                {pagination.total_items}
              </div>
              <div className="flex gap-2">
                <button
                  onClick={() => fetchOrganizations(pagination.page - 1)}
                  disabled={!pagination.has_prev}
                  className="px-4 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Previous
                </button>
                <div className="flex items-center gap-1">
                  {Array.from({ length: pagination.total_pages }, (_, i) => i + 1).map((p) => (
                    <button
                      key={p}
                      onClick={() => fetchOrganizations(p)}
                      className={`px-3 py-2 rounded-lg text-sm font-medium transition ${
                        p === pagination.page
                          ? 'bg-blue-600 text-white'
                          : 'border border-gray-300 text-gray-700 hover:bg-gray-50'
                      }`}
                    >
                      {p}
                    </button>
                  ))}
                </div>
                <button
                  onClick={() => fetchOrganizations(pagination.page + 1)}
                  disabled={!pagination.has_next}
                  className="px-4 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Next
                </button>
              </div>
            </div>
          )}
        </>
      )}
    </div>
  );
}
