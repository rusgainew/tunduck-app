import React, { useState } from 'react';
import { Link, Outlet, useNavigate } from 'react-router';
import { authApi } from '../services/api';

export default function DashboardLayout() {
  const navigate = useNavigate();
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const user = localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user') || '{}') : null;

  const handleLogout = async () => {
    const token = localStorage.getItem('token');
    if (token) {
      try {
        await authApi.logout(token);
      } catch (err) {
        console.error('Logout error:', err);
      }
    }
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    navigate('/login');
  };

  const menuItems = [
    {
      name: 'Users',
      icon: '👥',
      href: '/dashboard/users',
      description: 'Manage system users',
    },
    {
      name: 'Organizations',
      icon: '🏢',
      href: '/dashboard/organizations',
      description: 'Manage ESF organizations',
    },
    {
      name: 'Documents',
      icon: '📄',
      href: '/dashboard/documents',
      description: 'Manage ESF documents',
    },
    {
      name: 'Roles',
      icon: '🔐',
      href: '/dashboard/roles',
      description: 'Manage user roles',
    },
  ];

  return (
    <div className="flex h-screen bg-gray-100">
      {/* Sidebar */}
      <div
        className={`${
          sidebarOpen ? 'w-64' : 'w-20'
        } bg-gradient-to-b from-blue-700 to-blue-900 text-white transition-all duration-300 flex flex-col`}
      >
        {/* Logo */}
        <div className="h-16 flex items-center justify-between px-4 border-b border-blue-600">
          {sidebarOpen && <h1 className="text-xl font-bold">Tunduct</h1>}
          <button
            onClick={() => setSidebarOpen(!sidebarOpen)}
            className="p-2 hover:bg-blue-600 rounded-lg transition"
          >
            {sidebarOpen ? '◀' : '▶'}
          </button>
        </div>

        {/* Menu Items */}
        <nav className="flex-1 px-2 py-4 space-y-2">
          {menuItems.map((item) => (
            <Link
              key={item.href}
              to={item.href}
              className="flex items-center gap-3 px-4 py-3 rounded-lg hover:bg-blue-600 transition group"
              title={!sidebarOpen ? item.name : ''}
            >
              <span className="text-2xl">{item.icon}</span>
              {sidebarOpen && (
                <div>
                  <p className="font-semibold text-sm">{item.name}</p>
                  <p className="text-xs text-blue-200">{item.description}</p>
                </div>
              )}
            </Link>
          ))}
        </nav>

        {/* User Info */}
        <div className="px-2 py-4 border-t border-blue-600 space-y-2">
          {sidebarOpen && (
            <div className="px-4 py-3 bg-blue-600 rounded-lg">
              <p className="text-xs text-blue-200">Logged in as</p>
              <p className="font-semibold text-sm truncate">{user?.username || 'User'}</p>
              <p className="text-xs text-blue-300 capitalize">{user?.role || 'user'}</p>
            </div>
          )}
          <button
            onClick={handleLogout}
            className="w-full flex items-center justify-center gap-2 px-4 py-2 bg-red-600 hover:bg-red-700 rounded-lg transition font-semibold text-sm"
          >
            <span>🚪</span>
            {sidebarOpen && 'Logout'}
          </button>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Top Bar */}
        <div className="bg-white border-b border-gray-200 px-8 py-4 flex items-center justify-between">
          <div>
            <h2 className="text-2xl font-bold text-gray-900">Dashboard</h2>
            <p className="text-sm text-gray-600">Welcome back, {user?.username}!</p>
          </div>
          <div className="flex items-center gap-4">
            <div className="px-4 py-2 bg-gray-100 rounded-lg">
              <p className="text-xs text-gray-600">Role</p>
              <p className="font-semibold text-gray-900 capitalize">{user?.role}</p>
            </div>
            <button
              onClick={handleLogout}
              className="px-4 py-2 bg-red-100 hover:bg-red-200 text-red-700 font-semibold rounded-lg transition"
            >
              Logout
            </button>
          </div>
        </div>

        {/* Page Content */}
        <div className="flex-1 overflow-auto">
          <Outlet />
        </div>
      </div>
    </div>
  );
}
