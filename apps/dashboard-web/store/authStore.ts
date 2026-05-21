"use client";

import { create } from "zustand";

interface AuthState {
  accessToken: string | null;
  refreshToken: string | null;
  userId: string | null;
  tenantId: string | null;
  setTokens: (args: {
    accessToken: string;
    refreshToken: string;
    userId: string;
    tenantId: string;
  }) => void;
  clear: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  accessToken: null,
  refreshToken: null,
  userId: null,
  tenantId: null,
  setTokens: ({ accessToken, refreshToken, userId, tenantId }) =>
    set({ accessToken, refreshToken, userId, tenantId }),
  clear: () =>
    set({
      accessToken: null,
      refreshToken: null,
      userId: null,
      tenantId: null,
    }),
}));
