"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { documentsApi } from "@/lib/api";
import { EsfDocument } from "@/lib/types";
import { Search, FileText, Calendar, Eye, Trash2 } from "lucide-react";
import { format } from "date-fns";
import { ru } from "date-fns/locale";
import { useToast } from "@/hooks/useToast";
import { useOrganizationToken } from "@/hooks/useOrganizationToken";
import Link from "next/link";

export default function DocumentsPage() {
  const { success, error: showError } = useToast();
  const { getOrganizationToken } = useOrganizationToken();
  const [searchTerm, setSearchTerm] = useState("");

  const { data, isLoading, error } = useQuery({
    queryKey: ["documents"],
    queryFn: () => documentsApi.getAll(),
  });

  const filteredDocs = data?.data?.filter(
    (doc: EsfDocument) =>
      doc.contractorTin?.toLowerCase().includes(searchTerm.toLowerCase()) ||
      doc.foreignName?.toLowerCase().includes(searchTerm.toLowerCase()) ||
      doc.comment?.toLowerCase().includes(searchTerm.toLowerCase())
  );

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-lg p-4">
        <p className="text-red-800">Ошибка загрузки документов</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">Документы ЭСФ</h2>
          <p className="text-gray-600 mt-1">
            Всего документов: {data?.count || 0}
          </p>
        </div>
      </div>

      {/* Search */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
        <input
          type="text"
          placeholder="Поиск по ИНН, названию или комментарию..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
      </div>

      {/* Documents Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {filteredDocs && filteredDocs.length > 0 ? (
          filteredDocs.map((doc: EsfDocument) => (
            <div
              key={doc.id}
              className="bg-white rounded-lg shadow hover:shadow-lg transition-shadow p-6"
            >
              <div className="flex items-start justify-between mb-4">
                <div className="flex items-center flex-1">
                  <FileText className="w-5 h-5 text-blue-600 mr-2 flex-shrink-0" />
                  <h3 className="text-sm font-semibold text-gray-900 line-clamp-2">
                    {doc.foreignName || "Без названия"}
                  </h3>
                </div>
                <div className="flex space-x-2 ml-2">
                  <Link
                    href={`/dashboard/documents/${doc.id}`}
                    onClick={() => {
                      const orgToken = getOrganizationToken();
                      if (orgToken) {
                        localStorage.setItem("doc_org_token", orgToken);
                      }
                    }}
                    className="text-green-600 hover:text-green-700 transition-colors flex-shrink-0"
                    title="Просмотреть детали"
                  >
                    <Eye className="w-4 h-4" />
                  </Link>
                  <button
                    onClick={() => {
                      if (
                        confirm("Вы уверены, что хотите удалить этот документ?")
                      ) {
                        // Delete will be implemented
                        showError("Удаление еще не реализовано");
                      }
                    }}
                    className="text-red-600 hover:text-red-700 transition-colors flex-shrink-0"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              </div>

              <div className="space-y-3 text-sm">
                <div className="flex justify-between">
                  <span className="text-gray-500">ИНН контрагента:</span>
                  <span className="font-medium text-gray-900">
                    {doc.contractorTin || "—"}
                  </span>
                </div>

                <div className="flex justify-between">
                  <span className="text-gray-500">Валюта:</span>
                  <span className="font-medium text-gray-900">
                    {doc.currencyCode || "—"}
                  </span>
                </div>

                {doc.totalCurrencyValue && (
                  <div className="flex justify-between">
                    <span className="text-gray-500">Сумма:</span>
                    <span className="font-semibold text-gray-900">
                      {doc.totalCurrencyValue.toLocaleString("ru-RU")}{" "}
                      {doc.currencyCode}
                    </span>
                  </div>
                )}

                <div className="flex items-center text-gray-500 pt-2 border-t">
                  <Calendar className="w-4 h-4 mr-1" />
                  <span className="text-xs">
                    {doc.deliveryDate
                      ? format(new Date(doc.deliveryDate), "dd.MM.yyyy", {
                          locale: ru,
                        })
                      : "—"}
                  </span>
                </div>

                {doc.isResident !== undefined && (
                  <div className="pt-2">
                    <span
                      className={`px-2 py-1 text-xs rounded-full ${
                        doc.isResident
                          ? "bg-green-100 text-green-800"
                          : "bg-blue-100 text-blue-800"
                      }`}
                    >
                      {doc.isResident ? "Резидент КР" : "Нерезидент"}
                    </span>
                  </div>
                )}

                {doc.comment && (
                  <div className="pt-2">
                    <p className="text-xs text-gray-600 line-clamp-2">
                      {doc.comment}
                    </p>
                  </div>
                )}
              </div>
            </div>
          ))
        ) : (
          <div className="col-span-3 text-center py-12">
            <p className="text-gray-500">Документы не найдены</p>
          </div>
        )}
      </div>
    </div>
  );
}
