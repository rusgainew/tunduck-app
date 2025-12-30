"use client";

import { useParams, useRouter } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { documentsApi } from "@/lib/api";
import { ArrowLeft, Copy, Check } from "lucide-react";
import { useState } from "react";
import { format } from "date-fns";
import { ru } from "date-fns/locale";
import { useToast } from "@/hooks/useToast";

export default function DocumentDetailsPage() {
  const params = useParams();
  const router = useRouter();
  const { success } = useToast();
  const id = params.id as string;
  const [copiedId, setCopiedId] = useState(false);

  const { data: document, isLoading } = useQuery({
    queryKey: ["document", id],
    queryFn: () => documentsApi.getById(id),
    enabled: !!id,
  });

  const handleCopyId = async () => {
    if (document?.id) {
      await navigator.clipboard.writeText(document.id);
      setCopiedId(true);
      success("ID документа скопирован в буфер обмена");
      setTimeout(() => setCopiedId(false), 2000);
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (!document) {
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
          <p className="text-red-900">Документ не найден</p>
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
        <div className="flex-1">
          <h1 className="text-3xl font-bold text-gray-900">
            {document.foreignName || "Документ ЭСФ"}
          </h1>
          <p className="text-gray-600 mt-1">
            {document.createdAt
              ? format(new Date(document.createdAt), "dd MMMM yyyy HH:mm", {
                  locale: ru,
                })
              : "—"}
          </p>
        </div>
      </div>

      {/* Main Info */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Main Content */}
        <div className="lg:col-span-2 space-y-6">
          {/* Основная информация */}
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">
              Основная информация
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  ИНН контрагента
                </label>
                <p className="text-gray-900 font-mono">
                  {document.contractorTin || "—"}
                </p>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Валюта
                </label>
                <p className="text-gray-900">{document.currencyCode || "—"}</p>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Сумма
                </label>
                <p className="text-gray-900 font-semibold">
                  {document.totalCurrencyValue
                    ? document.totalCurrencyValue.toLocaleString("ru-RU")
                    : "—"}{" "}
                  {document.currencyCode || ""}
                </p>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Тип резидентности
                </label>
                <span
                  className={`px-3 py-1 text-sm rounded-full ${
                    document.isResident
                      ? "bg-green-100 text-green-800"
                      : "bg-blue-100 text-blue-800"
                  }`}
                >
                  {document.isResident ? "Резидент КР" : "Нерезидент"}
                </span>
              </div>
            </div>
          </div>

          {/* Даты и операции */}
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">
              Даты и операции
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Дата доставки
                </label>
                <p className="text-gray-900">
                  {document.deliveryDate
                    ? format(new Date(document.deliveryDate), "dd MMMM yyyy", {
                        locale: ru,
                      })
                    : "—"}
                </p>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Дата контракта
                </label>
                <p className="text-gray-900">
                  {document.contractStartDate
                    ? format(
                        new Date(document.contractStartDate),
                        "dd MMMM yyyy",
                        {
                          locale: ru,
                        }
                      )
                    : "—"}
                </p>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Код типа доставки
                </label>
                <p className="text-gray-900 font-mono">
                  {document.deliveryTypeCode || "—"}
                </p>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Код типа платежа
                </label>
                <p className="text-gray-900 font-mono">
                  {document.paymentCode || "—"}
                </p>
              </div>
            </div>
          </div>

          {/* Финансовая информация */}
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">
              Финансовая информация
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Сумма без налогов
                </label>
                <p className="text-gray-900">
                  {document.totalCurrencyValueWithoutTaxes
                    ? document.totalCurrencyValueWithoutTaxes.toLocaleString(
                        "ru-RU"
                      )
                    : "—"}
                </p>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Курс валюты
                </label>
                <p className="text-gray-900">{document.currencyRate || "—"}</p>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Оплачено
                </label>
                <p className="text-gray-900">{document.paidAmount || "—"}</p>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  К оплате
                </label>
                <p className="text-gray-900 font-semibold">
                  {document.amountToBePaid || "—"}
                </p>
              </div>
            </div>
          </div>

          {/* Комментарий */}
          {document.comment && (
            <div className="bg-white rounded-lg shadow p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-4">
                Комментарий
              </h2>
              <p className="text-gray-700 whitespace-pre-wrap">
                {document.comment}
              </p>
            </div>
          )}

          {/* Записи товаров/услуг */}
          {document.catalogEntries && document.catalogEntries.length > 0 && (
            <div className="bg-white rounded-lg shadow p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-4">
                Товары и услуги ({document.catalogEntries.length})
              </h2>
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-gray-200">
                      <th className="px-4 py-2 text-left text-gray-700 font-semibold">
                        Код каталога
                      </th>
                      <th className="px-4 py-2 text-left text-gray-700 font-semibold">
                        Название
                      </th>
                      <th className="px-4 py-2 text-right text-gray-700 font-semibold">
                        Кол-во
                      </th>
                      <th className="px-4 py-2 text-right text-gray-700 font-semibold">
                        Цена
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200">
                    {document.catalogEntries.map((entry, idx) => (
                      <tr key={idx} className="hover:bg-gray-50">
                        <td className="px-4 py-2 font-mono text-gray-900">
                          {entry.catalogCode}
                        </td>
                        <td className="px-4 py-2 text-gray-900">
                          {entry.catalogName}
                        </td>
                        <td className="px-4 py-2 text-right text-gray-900">
                          {entry.quantity}
                        </td>
                        <td className="px-4 py-2 text-right text-gray-900">
                          {entry.price?.toLocaleString("ru-RU")}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          )}
        </div>

        {/* Sidebar */}
        <div className="lg:col-span-1">
          <div className="bg-white rounded-lg shadow p-6 space-y-4 sticky top-20">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                ID документа
              </label>
              <div className="flex gap-2">
                <div className="flex-1 bg-gray-50 rounded p-3 font-mono text-xs text-gray-900 break-all">
                  {document.id}
                </div>
                <button
                  onClick={handleCopyId}
                  className="px-3 py-3 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors flex items-center justify-center"
                >
                  {copiedId ? (
                    <Check className="w-4 h-4" />
                  ) : (
                    <Copy className="w-4 h-4" />
                  )}
                </button>
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Создано
              </label>
              <p className="text-sm text-gray-600">
                {document.createdAt
                  ? format(new Date(document.createdAt), "dd MMMM yyyy HH:mm", {
                      locale: ru,
                    })
                  : "—"}
              </p>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Обновлено
              </label>
              <p className="text-sm text-gray-600">
                {document.updatedAt
                  ? format(new Date(document.updatedAt), "dd MMMM yyyy HH:mm", {
                      locale: ru,
                    })
                  : "—"}
              </p>
            </div>

            <div className="pt-4 border-t border-gray-200">
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-gray-600">Статус отправки:</span>
                  <span className="font-medium text-gray-900">
                    {document.isBranchDataSent ? "✓ Отправлен" : "Не отправлен"}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">Без налогов:</span>
                  <span className="font-medium text-gray-900">
                    {document.isPriceWithoutTaxes ? "Да" : "Нет"}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
