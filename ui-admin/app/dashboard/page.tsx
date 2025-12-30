"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { useAuthStore } from "@/lib/store";
import { usersApi, organizationsApi, documentsApi } from "@/lib/api";
import { Users, Building2, FileText, Activity } from "lucide-react";

export default function DashboardPage() {
  const router = useRouter();
  const { isAuthenticated, user } = useAuthStore();

  useEffect(() => {
    if (!isAuthenticated) {
      router.push("/login");
    }
  }, [isAuthenticated, router]);

  const { data: usersData } = useQuery({
    queryKey: ["users", { page: 1, limit: 10 }],
    queryFn: () => usersApi.getAll({ page: 1, limit: 10 }),
    enabled: isAuthenticated,
  });

  const { data: orgsData } = useQuery({
    queryKey: ["organizations"],
    queryFn: () => organizationsApi.getAll(),
    enabled: isAuthenticated,
  });

  const { data: docsData } = useQuery({
    queryKey: ["documents"],
    queryFn: () => documentsApi.getAll(),
    enabled: isAuthenticated,
  });

  if (!isAuthenticated) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  const stats = [
    {
      name: "Пользователи",
      value: usersData?.total || 0,
      icon: Users,
      color: "bg-blue-500",
      href: "/dashboard/users",
    },
    {
      name: "Организации",
      value: orgsData?.length || 0,
      icon: Building2,
      color: "bg-green-500",
      href: "/dashboard/organizations",
    },
    {
      name: "Документы",
      value: docsData?.count || 0,
      icon: FileText,
      color: "bg-purple-500",
      href: "/dashboard/documents",
    },
    {
      name: "Активность",
      value: "100%",
      icon: Activity,
      color: "bg-orange-500",
      href: "#",
    },
  ];

  return (
    <div className="space-y-6">
      {/* Welcome Section */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-2xl font-bold text-gray-900">
          Добро пожаловать, {user?.fullName}!
        </h2>
        <p className="mt-2 text-gray-600">
          Роль: <span className="font-semibold">{user?.role}</span>
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {stats.map((stat) => {
          const Icon = stat.icon;
          return (
            <a
              key={stat.name}
              href={stat.href}
              className="bg-white rounded-lg shadow hover:shadow-lg transition-shadow p-6"
            >
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-600">
                    {stat.name}
                  </p>
                  <p className="text-3xl font-bold text-gray-900 mt-2">
                    {stat.value}
                  </p>
                </div>
                <div className={`${stat.color} rounded-full p-3`}>
                  <Icon className="w-6 h-6 text-white" />
                </div>
              </div>
            </a>
          );
        })}
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Recent Users */}
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">
            Последние пользователи
          </h3>
          {usersData?.data && usersData.data.length > 0 ? (
            <div className="space-y-3">
              {usersData.data.slice(0, 5).map((user) => (
                <div
                  key={user.id}
                  className="flex items-center justify-between"
                >
                  <div>
                    <p className="text-sm font-medium text-gray-900">
                      {user.fullName}
                    </p>
                    <p className="text-xs text-gray-500">{user.email}</p>
                  </div>
                  <span
                    className={`px-2 py-1 text-xs rounded-full ${
                      user.isActive
                        ? "bg-green-100 text-green-800"
                        : "bg-gray-100 text-gray-800"
                    }`}
                  >
                    {user.isActive ? "Активен" : "Неактивен"}
                  </span>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-500 text-sm">Нет данных</p>
          )}
        </div>

        {/* Recent Organizations */}
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">
            Организации
          </h3>
          {orgsData && orgsData.length > 0 ? (
            <div className="space-y-3">
              {orgsData.slice(0, 5).map((org) => (
                <div key={org.id} className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-gray-900">
                      {org.name}
                    </p>
                    <p className="text-xs text-gray-500">{org.dbName}</p>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-500 text-sm">Нет данных</p>
          )}
        </div>
      </div>
    </div>
  );
}
