import React, { useState, useEffect } from 'react';
import type { User, PaginatedResponse } from '../../services/api';
import { rolesApi, usersApi } from '../../services/api';

export default function RolesPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [newRole, setNewRole] = useState('');
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(10);
  const [total, setTotal] = useState(0);

  const token = localStorage.getItem('token');

  const fetchUsers = async (pageNum: number) => {
    setLoading(true);
    setError('');
    try {
      const response: PaginatedResponse<User> = await usersApi.getAll(pageNum, limit);
      setUsers(response.data);
      setTotal(response.total);
      setPage(pageNum);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch users');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers(1);
  }, []);

  const handleAssignRole = async (userId: string) => {
    if (!token || !newRole) {
      setError('Please select a role and ensure you are authenticated');
      return;
    }

    try {
      setLoading(true);
      await rolesApi.assignRole(userId, newRole, token);
      setNewRole('');
      setSelectedUser(null);
      await fetchUsers(page);
      alert('Role assigned successfully');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to assign role');
    } finally {
      setLoading(false);
    }
  };

  const totalPages = Math.ceil(total / limit);

  return (
    <div className="p-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Role Management</h1>
        <p className="text-gray-600">Manage user roles and permissions</p>
      </div>

      {error && (
        <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">
          {error}
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Users List */}
        <div className="lg:col-span-2">
          <div className="bg-white rounded-lg shadow">
            <div className="px-6 py-4 border-b border-gray-200">
              <h2 className="text-lg font-bold text-gray-900">Users</h2>
            </div>

            {loading && !selectedUser ? (
              <div className="px-6 py-8 text-center">
                <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
                <p className="mt-2 text-gray-600">Loading users...</p>
              </div>
            ) : (
              <>
                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead className="bg-gray-50 border-b border-gray-200">
                      <tr>
                        <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase">
                          Username
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase">
                          Email
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase">
                          Current Role
                        </th>
                        <th className="px-6 py-3 text-right text-xs font-semibold text-gray-700 uppercase">
                          Action
                        </th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-200">
                      {users.length === 0 ? (
                        <tr>
                          <td colSpan={4} className="px-6 py-8 text-center text-gray-500">
                            No users found
                          </td>
                        </tr>
                      ) : (
                        users.map((user) => (
                          <tr
                            key={user.id}
                            className={`hover:bg-gray-50 transition ${
                              selectedUser?.id === user.id ? 'bg-blue-50' : ''
                            }`}
                          >
                            <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                              {user.username}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                              {user.email}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap">
                              <span className="px-3 py-1 text-xs font-semibold rounded-full bg-blue-100 text-blue-800">
                                {user.role}
                              </span>
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-right">
                              <button
                                onClick={() =>
                                  setSelectedUser(selectedUser?.id === user.id ? null : user)
                                }
                                className={`text-sm font-medium py-2 px-4 rounded transition ${
                                  selectedUser?.id === user.id
                                    ? 'bg-blue-600 text-white'
                                    : 'text-blue-600 hover:text-blue-900'
                                }`}
                              >
                                {selectedUser?.id === user.id ? 'Selected' : 'Select'}
                              </button>
                            </td>
                          </tr>
                        ))
                      )}
                    </tbody>
                  </table>
                </div>

                {/* Pagination */}
                {totalPages > 1 && (
                  <div className="px-6 py-4 border-t border-gray-200 flex justify-between items-center">
                    <div className="text-sm text-gray-600">
                      {total > 0
                        ? `Showing ${(page - 1) * limit + 1} to ${Math.min(
                            page * limit,
                            total
                          )} of ${total}`
                        : 'No users'}
                    </div>
                    <div className="flex gap-2">
                      <button
                        onClick={() => fetchUsers(page - 1)}
                        disabled={page === 1}
                        className="px-4 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        Previous
                      </button>
                      <button
                        onClick={() => fetchUsers(page + 1)}
                        disabled={page === totalPages}
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
        </div>

        {/* Role Assignment Panel */}
        <div>
          <div className="bg-white rounded-lg shadow sticky top-8">
            <div className="px-6 py-4 border-b border-gray-200">
              <h2 className="text-lg font-bold text-gray-900">Assign Role</h2>
            </div>

            <div className="p-6">
              {selectedUser ? (
                <div className="space-y-4">
                  <div className="p-4 bg-blue-50 rounded-lg">
                    <p className="text-sm text-gray-600 mb-1">Selected User:</p>
                    <p className="text-lg font-bold text-gray-900">{selectedUser.username}</p>
                    <p className="text-xs text-gray-500">{selectedUser.email}</p>
                  </div>

                  <div className="p-4 bg-gray-50 rounded-lg">
                    <p className="text-sm text-gray-600 mb-2">Current Role:</p>
                    <p className="text-base font-semibold text-blue-600">{selectedUser.role}</p>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Select New Role
                    </label>
                    <select
                      value={newRole}
                      onChange={(e) => setNewRole(e.target.value)}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none"
                    >
                      <option value="">Choose a role...</option>
                      <option value="user">User</option>
                      <option value="admin">Admin</option>
                      <option value="manager">Manager</option>
                      <option value="viewer">Viewer</option>
                    </select>
                  </div>

                  <button
                    onClick={() => handleAssignRole(selectedUser.id)}
                    disabled={!newRole || loading}
                    className="w-full bg-green-600 hover:bg-green-700 disabled:bg-gray-400 text-white font-semibold py-2 px-4 rounded-lg transition"
                  >
                    {loading ? 'Assigning...' : 'Assign Role'}
                  </button>

                  <button
                    onClick={() => {
                      setSelectedUser(null);
                      setNewRole('');
                    }}
                    className="w-full border border-gray-300 text-gray-700 hover:bg-gray-50 font-semibold py-2 px-4 rounded-lg transition"
                  >
                    Cancel
                  </button>
                </div>
              ) : (
                <div className="text-center py-8">
                  <p className="text-gray-500">Select a user from the list to assign a role</p>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Role Information */}
      <div className="mt-8 bg-white rounded-lg shadow p-6">
        <h2 className="text-lg font-bold text-gray-900 mb-4">Available Roles</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <div className="p-4 border border-gray-200 rounded-lg">
            <h3 className="font-bold text-gray-900 mb-1">Admin</h3>
            <p className="text-sm text-gray-600">Full access to all features</p>
          </div>
          <div className="p-4 border border-gray-200 rounded-lg">
            <h3 className="font-bold text-gray-900 mb-1">Manager</h3>
            <p className="text-sm text-gray-600">Manage organizations and documents</p>
          </div>
          <div className="p-4 border border-gray-200 rounded-lg">
            <h3 className="font-bold text-gray-900 mb-1">User</h3>
            <p className="text-sm text-gray-600">Create and view documents</p>
          </div>
          <div className="p-4 border border-gray-200 rounded-lg">
            <h3 className="font-bold text-gray-900 mb-1">Viewer</h3>
            <p className="text-sm text-gray-600">Read-only access</p>
          </div>
        </div>
      </div>
    </div>
  );
}
