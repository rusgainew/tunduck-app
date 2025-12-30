"use client";

import { useState, useEffect } from "react";
import { EsfOrganization } from "@/lib/types";
import { X, RefreshCw } from "lucide-react";
import { v4 as uuidv4 } from "uuid";

interface OrganizationFormProps {
  organization?: EsfOrganization | null;
  isOpen: boolean;
  isLoading?: boolean;
  onClose: () => void;
  onSubmit: (
    data: Omit<EsfOrganization, "id" | "createdAt" | "updatedAt" | "deletedAt">
  ) => Promise<void>;
}

export default function OrganizationForm({
  organization,
  isOpen,
  isLoading = false,
  onClose,
  onSubmit,
}: OrganizationFormProps) {
  const [formData, setFormData] = useState({
    name: "",
    description: "",
    dbName: "",
    token: "",
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  const generateToken = () => {
    return uuidv4() + uuidv4().replace(/-/g, "");
  };

  useEffect(() => {
    if (organization) {
      setFormData({
        name: organization.name || "",
        description: organization.description || "",
        dbName: organization.dbName || "",
        token: organization.token || "",
      });
    } else {
      setFormData({
        name: "",
        description: "",
        dbName: "",
        token: generateToken(),
      });
    }
    setErrors({});
  }, [organization, isOpen]);

  const validateForm = (): boolean => {
    const newErrors: Record<string, string> = {};

    if (!formData.name.trim()) {
      newErrors.name = "Название организации обязательно";
    }

    if (!formData.dbName.trim()) {
      newErrors.dbName = "Имя базы данных обязательно";
    }

    if (!formData.token.trim()) {
      newErrors.token = "Токен обязателен";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    try {
      await onSubmit(formData);
      setFormData({
        name: "",
        description: "",
        dbName: "",
        token: "",
      });
      setErrors({});
      onClose();
    } catch (error) {
      console.error("Form submission error:", error);
    }
  };

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
    // Clear error for this field when user starts typing
    if (errors[name]) {
      setErrors((prev) => ({
        ...prev,
        [name]: "",
      }));
    }
  };

  if (!isOpen) {
    return null;
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg shadow-xl max-w-md w-full max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="flex justify-between items-center p-6 border-b border-gray-200 sticky top-0 bg-white">
          <h3 className="text-lg font-semibold text-gray-900">
            {organization
              ? "Редактировать организацию"
              : "Создать новую организацию"}
          </h3>
          <button
            onClick={onClose}
            disabled={isLoading}
            className="text-gray-400 hover:text-gray-600 transition-colors disabled:opacity-50"
          >
            <X className="w-6 h-6" />
          </button>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="p-6 space-y-4">
          {/* Name Field */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Название организации *
            </label>
            <input
              type="text"
              name="name"
              value={formData.name}
              onChange={handleChange}
              disabled={isLoading}
              placeholder="Введите название организации"
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors disabled:bg-gray-100 disabled:cursor-not-allowed ${
                errors.name ? "border-red-500" : "border-gray-300"
              }`}
            />
            {errors.name && (
              <p className="text-red-500 text-sm mt-1">{errors.name}</p>
            )}
          </div>

          {/* Description Field */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Описание
            </label>
            <textarea
              name="description"
              value={formData.description}
              onChange={handleChange}
              disabled={isLoading}
              placeholder="Введите описание организации (опционально)"
              rows={3}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors resize-none disabled:bg-gray-100 disabled:cursor-not-allowed"
            />
          </div>

          {/* Database Name Field */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Имя базы данных *
            </label>
            <input
              type="text"
              name="dbName"
              value={formData.dbName}
              onChange={handleChange}
              disabled={isLoading}
              placeholder="Введите имя базы данных"
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors disabled:bg-gray-100 disabled:cursor-not-allowed ${
                errors.dbName ? "border-red-500" : "border-gray-300"
              }`}
            />
            {errors.dbName && (
              <p className="text-red-500 text-sm mt-1">{errors.dbName}</p>
            )}
          </div>

          {/* Token Field */}
          <div>
            <div className="flex justify-between items-center mb-2">
              <label className="block text-sm font-medium text-gray-700">
                Токен *
              </label>
              {!organization && (
                <button
                  type="button"
                  onClick={() =>
                    setFormData((prev) => ({ ...prev, token: generateToken() }))
                  }
                  disabled={isLoading}
                  className="text-xs text-blue-600 hover:text-blue-700 flex items-center gap-1 disabled:opacity-50"
                  title="Сгенерировать новый токен"
                >
                  <RefreshCw className="w-3 h-3" />
                  Генерировать
                </button>
              )}
            </div>
            <textarea
              name="token"
              value={formData.token}
              onChange={handleChange}
              disabled={isLoading}
              placeholder="Введите токен доступа"
              rows={2}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors resize-none font-mono text-sm disabled:bg-gray-100 disabled:cursor-not-allowed ${
                errors.token ? "border-red-500" : "border-gray-300"
              }`}
            />
            {errors.token && (
              <p className="text-red-500 text-sm mt-1">{errors.token}</p>
            )}
          </div>

          {/* Buttons */}
          <div className="flex gap-3 pt-4">
            <button
              type="button"
              onClick={onClose}
              disabled={isLoading}
              className="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Отмена
            </button>
            <button
              type="submit"
              disabled={isLoading}
              className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
            >
              {isLoading ? (
                <>
                  <div className="animate-spin rounded-full h-4 w-4 border-t-2 border-b-2 border-white"></div>
                  <span>Сохранение...</span>
                </>
              ) : organization ? (
                "Обновить"
              ) : (
                "Создать"
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
