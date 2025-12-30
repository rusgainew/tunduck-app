"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { organizationsApi } from "@/lib/api";
import { EsfOrganization } from "@/lib/types";
import { Search, Plus, Pencil, Trash2 } from "lucide-react";
import { format } from "date-fns";

export default function OrganizationsPage() {
  const queryClient = useQueryClient();
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
      } catch (error) {
        console.error("Error deleting organization:", error);
        alert("Ошибка при удалении организации");
      }
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
                    {format(new Date(org.createdAt), "dd.MM.yyyy")}
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

      {/* Create/Edit Dialog - To be implemented */}
      {isDialogOpen && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-md w-full">
            <h3 className="text-lg font-semibold mb-4">
              {editingOrg ? "Редактировать" : "Создать"} организацию
            </h3>
            <p className="text-gray-600 mb-4">
              Форма создания/редактирования будет реализована позже
            </p>
            <button
              onClick={() => setIsDialogOpen(false)}
              className="w-full px-4 py-2 bg-gray-200 text-gray-800 rounded-lg hover:bg-gray-300"
            >
              Закрыть
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
