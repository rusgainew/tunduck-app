import { useCallback } from 'react';
import { apiClient } from '@/lib/api-client';

export function useOrganizationToken() {
  const setOrganizationToken = useCallback((token: string) => {
    apiClient.setOrganizationToken(token);
  }, []);

  const getOrganizationToken = useCallback((): string | null => {
    return apiClient.getOrganizationToken();
  }, []);

  const clearOrganizationToken = useCallback(() => {
    apiClient.clearOrganizationToken();
  }, []);

  return {
    setOrganizationToken,
    getOrganizationToken,
    clearOrganizationToken,
  };
}
