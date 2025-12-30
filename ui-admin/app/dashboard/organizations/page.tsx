"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { organizationsApi } from "@/lib/api";
import { EsfOrganization } from "@/lib/types";
import { Search, Plus, Pencil, Trash2, Eye } from "lucide-react";
import { format } from "date-fns";
import OrganizationForm from "@/components/OrganizationForm";
import { useToast } from "@/hooks/useToast";
import { useOrganizationToken } from "@/hooks/useOrganizationToken";
import Link from "next/link";

export default function OrganizationsPage() {
  const queryClient = useQueryClient();
  const { success, error } = useToast();
  const { setOrganizationToken } = useOrganizationToken();
  const [searchTerm, setSearchTerm] = useState("");
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [editingOrg, setEditingOrg] = useState<EsfOrganization | null>(null);

  const { data: organizations, isLoading } = useQuery({
    queryKey: ["organizations"],
    queryFn: () => organizationsApi.getAll(),
  });

  const deleteMutation = useMutation({
    mutationFn: organizationsApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["organizations"] });
      success("Организация успешно удалена");
    },
    onError: () => {
      error("Ошибка при удалении организации");
    },
  });

  const createMutation = useMutation({
    mutationFn: (
      data: Omit<
        EsfOrganization,
        "id" | "createdAt" | "updatedAt" | "deletedAt"
      >
    ) => organizationsApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["organizations"] });
      success("Организация успешно создана");
      setIsDialogOpen(false);
    },
    onError: () => {
      error("Ошибка при создании организации");
    },
  });

  const updateMutation = useMutation({
    mutationFn: (data: { id: string; updates: Partial<EsfOrganization> }) =>
      organizationsApi.update(data.id, data.updates),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["organizations"] });
      success("Организация успешно обновлена");
      setIsDialogOpen(false);
    },
    onError: () => {
      error("Ошибка при обновлении организации");
    },
  });

  const filteredOrgs = organizations?.filter(
    (org) =>
      org.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      org.description.toLowerCase().includes(searchTerm.toLowerCase()) ||
      org.dbName.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleDelete = async (id: string) => {
    if (confirm("Вы уверены, что хотите удалить эту организацию?")) {
      try {
        await deleteMutation.mutateAsync(id);
      } catch (err) {
        console.error("Error deleting organization:", err);
      }
    }
  };

  const handleFormSubmit = async (
    data: Omit<EsfOrganization, "id" | "createdAt" | "updatedAt" | "deletedAt">
  ) => {
    if (editingOrg) {
      await updateMutation.mutateAsync({
        id: editingOrg.id,
        updates: data,
      });
    } else {
      await createMutation.mutateAsync(data);
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">Организации ЭСФ</h2>
          <p className="text-gray-600 mt-1">
            Всего организаций: {organizations?.length || 0}
          </p>
        </div>
        <button
          onClick={() => {
            setEditingOrg(null);
            setIsDialogOpen(true);
          }}
          className="flex items-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
        >
          <Plus className="w-5 h-5 mr-2" />
          Добавить организацию
        </button>
      </div>

      {/* Search */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
        <input
          type="text"
          placeholder="Поиск по названию, описанию или БД..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
      </div>

      {/* Organizations Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {filteredOrgs && filteredOrgs.length > 0 ? (
          filteredOrgs.map((org) => (
            <div
              key={org.id}
              className="bg-white rounded-lg shadow hover:shadow-lg transition-shadow p-6"
            >
              <div className="flex justify-between items-start mb-4">
                <h3 className="text-lg font-semibold text-gray-900">
                  {org.name}
                </h3>
                <div className="flex space-x-2">
                  <Link
                    href={`/dashboard/organizations/${org.id}`}
                    onClick={() => setOrganizationToken(org.token)}
                    className="text-green-600 hover:text-green-700 transition-colors"
                    title="Просмотреть детали"
                  >
                    <Eye className="w-4 h-4" />
                  </Link>
                  <button
                    onClick={() => {
                      setEditingOrg(org);
                      setIsDialogOpen(true);
                    }}
                    className="text-blue-600 hover:text-blue-700"
                  >
                    <Pencil className="w-4 h-4" />
                  </button>
                  <button
                    onClick={() => handleDelete(org.id)}
                    className="text-red-600 hover:text-red-700"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              </div>

              <p className="text-sm text-gray-600 mb-4 line-clamp-2">
                {org.description || "Нет описания"}
              </p>

              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-gray-500">База данных:</span>
                  <span className="font-medium text-gray-900">
                    {org.dbName}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-500">Токен:</span>
                  <span className="font-mono text-xs text-gray-900">
                    {org.token.substring(0, 16)}...
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-500">Создано:</span>
                  <span className="text-gray-900">
                    {org.createdAt
                      ? format(new Date(org.createdAt), "dd.MM.yyyy")
                      : "—"}
                  </span>
                </div>
              </div>
            </div>
          ))
        ) : (
          <div className="col-span-3 text-center py-12">
            <p className="text-gray-500">Организации не найдены</p>
          </div>
        )}
      </div>

      {/* Create/Edit Dialog */}
      <OrganizationForm
        organization={editingOrg}
        isOpen={isDialogOpen}
        isLoading={createMutation.isPending || updateMutation.isPending}
        onClose={() => {
          setIsDialogOpen(false);
          setEditingOrg(null);
        }}
        onSubmit={handleFormSubmit}
      />
    </div>
  );
}
