"use client";

import { useParams, useRouter, useSearchParams } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { organizationsApi, documentsApi } from "@/lib/api";
import { ArrowLeft, Copy, Check } from "lucide-react";
import { useState, useEffect } from "react";
import { format } from "date-fns";
import { ru } from "date-fns/locale";
import { useToast } from "@/hooks/useToast";
import { useOrganizationToken } from "@/hooks/useOrganizationToken";

export default function OrganizationDetailsPage() {
  const params = useParams();
  const router = useRouter();
  const { success } = useToast();
  const { setOrganizationToken } = useOrganizationToken();
  const id = params.id as string;
  const [copiedToken, setCopiedToken] = useState(false);

  const { data: organization, isLoading: orgLoading } = useQuery({
    queryKey: ["organization", id],
    queryFn: () => organizationsApi.getById(id),
    enabled: !!id,
  });

  // Сохраняем токен организации в localStorage при загрузке
  useEffect(() => {
    if (organization?.token) {
      setOrganizationToken(organization.token);
    }
  }, [organization?.token, setOrganizationToken]);

  const { data: docs, isLoading: docsLoading } = useQuery({
    queryKey: ["documents", id],
    queryFn: () => documentsApi.getAll(id),
    enabled: !!id,
  });

  const handleCopyToken = async () => {
    if (organization?.token) {
      await navigator.clipboard.writeText(organization.token);
      setCopiedToken(true);
      success("Токен скопирован в буфер обмена");
      setTimeout(() => setCopiedToken(false), 2000);
    }
  };

  if (orgLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (!organization) {
    return (
      <div className="space-y-4">
        <button
          onClick={() => router.back()}
          className="flex items-center text-blue-600 hover:text-blue-700"
        >
          <ArrowLeft className="w-5 h-5 mr-2" />
          Назад
        </button>
        <div className="bg-red-50 border border-red-200 rounded-lg p-6 text-center">
          <p className="text-red-900">Организация не найдена</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <button
          onClick={() => router.back()}
          className="text-gray-600 hover:text-gray-900 transition-colors"
        >
          <ArrowLeft className="w-6 h-6" />
        </button>
        <div>
          <h1 className="text-3xl font-bold text-gray-900">
            {organization.name}
          </h1>
          <p className="text-gray-600 mt-1">
            Создано:{" "}
            {organization.createdAt
              ? format(new Date(organization.createdAt), "dd MMMM yyyy HH:mm", {
                  locale: ru,
                })
              : "—"}
          </p>
        </div>
      </div>

      {/* Organization Details */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Main Info */}
        <div className="lg:col-span-2 space-y-6">
          {/* Description Card */}
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">
              Описание
            </h2>
            <p className="text-gray-600">
              {organization.description || "Описание не указано"}
            </p>
          </div>

          {/* Technical Info Card */}
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">
              Техническая информация
            </h2>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Имя базы данных
                </label>
                <div className="bg-gray-50 rounded p-3 font-mono text-sm text-gray-900">
                  {organization.dbName}
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Токен доступа
                </label>
                <div className="flex gap-2">
                  <div className="flex-1 bg-gray-50 rounded p-3 font-mono text-sm text-gray-900 break-all">
                    {organization.token}
                  </div>
                  <button
                    onClick={handleCopyToken}
                    className="px-4 py-3 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors flex items-center gap-2 whitespace-nowrap"
                  >
                    {copiedToken ? (
                      <>
                        <Check className="w-4 h-4" />
                        Скопирован
                      </>
                    ) : (
                      <>
                        <Copy className="w-4 h-4" />
                        Копировать
                      </>
                    )}
                  </button>
                </div>
              </div>
            </div>
          </div>

          {/* Documents Section */}
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">
              Документы ЭСФ
            </h2>
            {docsLoading ? (
              <div className="flex justify-center py-8">
                <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-500"></div>
              </div>
            ) : docs?.data && docs.data.length > 0 ? (
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b border-gray-200">
                      <th className="px-4 py-2 text-left text-sm font-semibold text-gray-900">
                        Дата
                      </th>
                      <th className="px-4 py-2 text-left text-sm font-semibold text-gray-900">
                        Контрагент
                      </th>
                      <th className="px-4 py-2 text-left text-sm font-semibold text-gray-900">
                        Сумма
                      </th>
                      <th className="px-4 py-2 text-left text-sm font-semibold text-gray-900">
                        Валюта
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200">
                    {docs.data.map((doc) => (
                      <tr
                        key={doc.id}
                        className="hover:bg-gray-50 transition-colors"
                      >
                        <td className="px-4 py-3 text-sm text-gray-600">
                          {format(new Date(doc.createdAt), "dd.MM.yyyy", {
                            locale: ru,
                          })}
                        </td>
                        <td className="px-4 py-3 text-sm text-gray-900">
                          {doc.contractorTin || "—"}
                        </td>
                        <td className="px-4 py-3 text-sm text-gray-900">
                          {doc.totalCurrencyValue?.toLocaleString("ru-RU") ||
                            "0"}
                        </td>
                        <td className="px-4 py-3 text-sm text-gray-600">
                          {doc.currencyCode || "—"}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            ) : (
              <div className="text-center py-8 text-gray-500">
                <p>Документы не найдены</p>
              </div>
            )}
          </div>
        </div>

        {/* Sidebar */}
        <div className="lg:col-span-1">
          <div className="bg-white rounded-lg shadow p-6 space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                ID организации
              </label>
              <div className="bg-gray-50 rounded p-3 font-mono text-xs text-gray-900 break-all">
                {organization.id}
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Обновлено
              </label>
              <p className="text-sm text-gray-600">
                {organization.updatedAt
                  ? format(
                      new Date(organization.updatedAt),
                      "dd MMMM yyyy HH:mm",
                      {
                        locale: ru,
                      }
                    )
                  : "—"}
              </p>
            </div>

            {organization.deletedAt && (
              <div className="bg-red-50 border border-red-200 rounded p-3">
                <p className="text-xs font-medium text-red-900">Удалено</p>
                <p className="text-xs text-red-700 mt-1">
                  {format(
                    new Date(organization.deletedAt),
                    "dd.MM.yyyy HH:mm",
                    {
                      locale: ru,
                    }
                  ) || "—"}
                </p>
              </div>
            )}

            <div className="pt-4 border-t border-gray-200">
              <p className="text-xs text-gray-500">
                Всего документов:{" "}
                <span className="font-semibold text-gray-900">
                  {docs?.count || 0}
                </span>
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
