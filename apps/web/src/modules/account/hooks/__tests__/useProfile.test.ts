import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import React from 'react';
import { useProfile } from '../useProfile';

vi.mock('../../services/profile.api', () => ({
  profileApi: {
    getProfile: vi.fn(),
    updateProfile: vi.fn(),
  },
}));

import { profileApi } from '../../services/profile.api';

const mockProfileApi = profileApi as unknown as {
  getProfile: ReturnType<typeof vi.fn>;
  updateProfile: ReturnType<typeof vi.fn>;
};

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, gcTime: 0 } },
  });
  return ({ children }: { children: React.ReactNode }) =>
    React.createElement(QueryClientProvider, { client: queryClient }, children);
}

describe('useProfile', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('fetches profile data on mount', async () => {
    const mockUser = { id: 'u1', email: 'test@example.com', first_name: 'John' };
    mockProfileApi.getProfile.mockResolvedValue(mockUser);

    const { result } = renderHook(() => useProfile(), { wrapper: createWrapper() });

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    expect(result.current.data).toEqual(mockUser);
    expect(mockProfileApi.getProfile).toHaveBeenCalledOnce();
  });

  it('exposes updateProfile mutation', async () => {
    mockProfileApi.getProfile.mockResolvedValue({ id: 'u1' });
    mockProfileApi.updateProfile.mockResolvedValue({ id: 'u1', first_name: 'Jane' });

    const { result } = renderHook(() => useProfile(), { wrapper: createWrapper() });

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    expect(result.current.updateProfile).toBeDefined();
    expect(typeof result.current.updateProfile.mutate).toBe('function');
  });
});
