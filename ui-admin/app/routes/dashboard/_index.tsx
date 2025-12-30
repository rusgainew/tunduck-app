import React, { useState, useEffect } from 'react';
import { usersApi, organizationsApi, documentsApi } from '../../services/api';

export default function DashboardIndexPage() {
  const [stats, setStats] = useState({
    totalUsers: 0,
    totalOrganizations: 0,
    totalDocuments: 0,
    loading: true,
    error: '',
  });

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const [usersRes, orgsRes, docsRes] = await Promise.all([
          usersApi.getAll(1, 1),
          organizationsApi.getAll(),
          documentsApi.getAll(),
        ]);

        setStats({
          totalUsers: usersRes.total,
          totalOrganizations: orgsRes.length,
          totalDocuments: docsRes.length,
          loading: false,
          error: '',
        });
      } catch (err) {
        setStats((prev) => ({
          ...prev,
          error: err instanceof Error ? err.message : 'Failed to load stats',
          loading: false,
        }));
      }
    };

    fetchStats();
  }, []);

  const statCards = [
    {
      title: 'Total Users',
      value: stats.totalUsers,
      icon: '👥',
      color: 'blue',
      href: '/dashboard/users',
    },
    {
      title: 'Organizations',
      value: stats.totalOrganizations,
      icon: '🏢',
      color: 'green',
      href: '/dashboard/organizations',
    },
    {
      title: 'Documents',
      value: stats.totalDocuments,
      icon: '📄',
      color: 'purple',
      href: '/dashboard/documents',
    },
  ];

  const getColorClasses = (color: string) => {
    const colors: Record<string, string> = {
      blue: 'from-blue-500 to-blue-600',
      green: 'from-green-500 to-green-600',
      purple: 'from-purple-500 to-purple-600',
      red: 'from-red-500 to-red-600',
    };
    return colors[color] || colors.blue;
  };

  return (
    <div className="p-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Dashboard Overview</h1>
        <p className="text-gray-600">Welcome to the Tunduct ESF Management System</p>
      </div>

      {stats.error && (
        <div className="mb-4 p-4 bg-yellow-50 border border-yellow-200 rounded-lg text-yellow-700">
          {stats.error}
        </div>
      )}

      {/* Statistics Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        {statCards.map((card) => (
          <a
            key={card.title}
            href={card.href}
            className="group cursor-pointer"
          >
            <div
              className={`bg-gradient-to-br ${getColorClasses(
                card.color
              )} rounded-lg shadow-lg p-6 text-white group-hover:shadow-xl transition transform group-hover:-translate-y-1`}
            >
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-semibold">{card.title}</h3>
                <span className="text-4xl">{card.icon}</span>
              </div>
              <p className="text-4xl font-bold">
                {stats.loading ? '...' : card.value}
              </p>
              <p className="text-sm text-white/80 mt-2">Click to manage</p>
            </div>
          </a>
        ))}
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* System Features */}
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-bold text-gray-900 mb-4">System Features</h2>
          <ul className="space-y-3">
            <li className="flex items-start gap-3">
              <span className="text-green-600 font-bold text-lg mt-0.5">✓</span>
              <div>
                <p className="font-semibold text-gray-900">User Management</p>
                <p className="text-sm text-gray-600">Create and manage system users with different roles</p>
              </div>
            </li>
            <li className="flex items-start gap-3">
              <span className="text-green-600 font-bold text-lg mt-0.5">✓</span>
              <div>
                <p className="font-semibold text-gray-900">Organization Management</p>
                <p className="text-sm text-gray-600">Manage multiple organizations with separate databases</p>
              </div>
            </li>
            <li className="flex items-start gap-3">
              <span className="text-green-600 font-bold text-lg mt-0.5">✓</span>
              <div>
                <p className="font-semibold text-gray-900">Document Management</p>
                <p className="text-sm text-gray-600">Create, edit, and publish ESF documents</p>
              </div>
            </li>
            <li className="flex items-start gap-3">
              <span className="text-green-600 font-bold text-lg mt-0.5">✓</span>
              <div>
                <p className="font-semibold text-gray-900">Role-Based Access Control</p>
                <p className="text-sm text-gray-600">Fine-grained permissions with role management</p>
              </div>
            </li>
          </ul>
        </div>

        {/* Recent Activity / Help */}
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-bold text-gray-900 mb-4">Quick Start Guide</h2>
          <div className="space-y-4">
            <div className="border-l-4 border-blue-500 pl-4">
              <h3 className="font-semibold text-gray-900">1. Add Users</h3>
              <p className="text-sm text-gray-600">Create user accounts and assign roles from the Users page</p>
            </div>
            <div className="border-l-4 border-green-500 pl-4">
              <h3 className="font-semibold text-gray-900">2. Create Organizations</h3>
              <p className="text-sm text-gray-600">Set up organizations with separate databases</p>
            </div>
            <div className="border-l-4 border-purple-500 pl-4">
              <h3 className="font-semibold text-gray-900">3. Manage Documents</h3>
              <p className="text-sm text-gray-600">Create and publish ESF documents within organizations</p>
            </div>
            <div className="border-l-4 border-red-500 pl-4">
              <h3 className="font-semibold text-gray-900">4. Assign Roles</h3>
              <p className="text-sm text-gray-600">Control access with granular role-based permissions</p>
            </div>
          </div>
        </div>
      </div>

      {/* System Status */}
      <div className="mt-8 bg-gradient-to-r from-green-50 to-blue-50 border border-green-200 rounded-lg p-6">
        <div className="flex items-center gap-4">
          <div className="flex-1">
            <h3 className="text-lg font-bold text-gray-900 mb-1">System Status</h3>
            <p className="text-gray-600">All systems operational</p>
          </div>
          <div className="flex items-center gap-2">
            <span className="flex h-3 w-3 bg-green-600 rounded-full animate-pulse"></span>
            <span className="text-green-700 font-semibold">Healthy</span>
          </div>
        </div>
      </div>
    </div>
  );
}
