"use client";

import { useEffect } from "react";
import { AlertCircle, CheckCircle, Info, X } from "lucide-react";

export type ToastType = "success" | "error" | "info" | "warning";

interface Toast {
  id: string;
  type: ToastType;
  message: string;
  duration?: number;
}

interface ToastProps {
  toasts: Toast[];
  onRemove: (id: string) => void;
}

const TOAST_ICONS: Record<ToastType, React.ReactNode> = {
  success: <CheckCircle className="w-5 h-5 text-green-600" />,
  error: <AlertCircle className="w-5 h-5 text-red-600" />,
  warning: <AlertCircle className="w-5 h-5 text-yellow-600" />,
  info: <Info className="w-5 h-5 text-blue-600" />,
};

const TOAST_COLORS: Record<ToastType, string> = {
  success: "bg-green-50 border-green-200",
  error: "bg-red-50 border-red-200",
  warning: "bg-yellow-50 border-yellow-200",
  info: "bg-blue-50 border-blue-200",
};

const TOAST_TEXT_COLORS: Record<ToastType, string> = {
  success: "text-green-900",
  error: "text-red-900",
  warning: "text-yellow-900",
  info: "text-blue-900",
};

export function ToastContainer({ toasts, onRemove }: ToastProps) {
  return (
    <div className="fixed bottom-4 right-4 z-[1000] space-y-2">
      {toasts.map((toast) => (
        <ToastItem key={toast.id} toast={toast} onRemove={onRemove} />
      ))}
    </div>
  );
}

interface ToastItemProps {
  toast: Toast;
  onRemove: (id: string) => void;
}

function ToastItem({ toast, onRemove }: ToastItemProps) {
  useEffect(() => {
    const duration = toast.duration || 4000;
    const timer = setTimeout(() => {
      onRemove(toast.id);
    }, duration);

    return () => clearTimeout(timer);
  }, [toast, onRemove]);

  return (
    <div
      className={`flex items-center gap-3 px-4 py-3 rounded-lg border shadow-md max-w-sm animate-in slide-in-from-right-4 fade-in ${
        TOAST_COLORS[toast.type]
      }`}
    >
      {TOAST_ICONS[toast.type]}
      <span
        className={`flex-1 text-sm font-medium ${
          TOAST_TEXT_COLORS[toast.type]
        }`}
      >
        {toast.message}
      </span>
      <button
        onClick={() => onRemove(toast.id)}
        className={`text-gray-400 hover:text-gray-600 transition-colors`}
      >
        <X className="w-4 h-4" />
      </button>
    </div>
  );
}
