import { create } from 'zustand';
import { useSettingsStore } from '@features/settings/stores/useSettingsStore';
import { fetchAuthorization } from '@api/authorize';

interface ApplicationState {
  loading: boolean;
  authorized: boolean;
  admin: boolean;
  hasCookie: boolean;
  cookieExpiresAt: number | null;

  reload: () => Promise<void>;
  refresh: () => Promise<void>;
  authorize: () => Promise<void>;
  shouldRefreshCookie: () => boolean;
}

export const useApplicationStore = create<ApplicationState>((set, get) => ({
  loading: false,
  authorized: false,
  admin: false,
  hasCookie: false,
  cookieExpiresAt: null,

  reload: async () => {
    set({ loading: true, authorized: false, admin: false });

    try {
      const settingsStore = useSettingsStore.getState();
      await settingsStore.initializeSettings();
      set({
        loading: false,
        authorized: true,
        admin: true,
      });
    } catch (err) {
      const unauthorized = err instanceof Error && err.message === 'not authorized';
      const forbidden = err instanceof Error && err.message === 'access denied';
      set({ loading: false, authorized: !unauthorized, admin: (!unauthorized && !forbidden) });
    } finally {
      const settingsStore = useSettingsStore.getState()
      if (!settingsStore.hasSettings) {
        window.location.hash = '#/settings';
        return;
      }

      window.location.hash = '#/';
    }
  },
  refresh: async () => {
    try {
      const settingsStore = useSettingsStore.getState();
      await settingsStore.initializeSettings();
      set({
        authorized: true,
        admin: true,
      });
    } catch (err) {
      const unauthorized = err instanceof Error && err.message === 'not authorized';
      const forbidden = err instanceof Error && err.message === 'access denied';
      set({ authorized: !unauthorized, admin: (!unauthorized && !forbidden) });
    }
  },
  authorize: async () => {
    try {
      set({ hasCookie: false });
      const { expiresAt } = await fetchAuthorization();
      set({
        hasCookie: true,
        cookieExpiresAt: expiresAt
      });
    } catch (err) {
      const unauthorized = err instanceof Error && err.message === 'not authorized';
      const forbidden = err instanceof Error && err.message === 'access denied';
      set({
        hasCookie: false,
        authorized: !unauthorized,
        admin: (!unauthorized && !forbidden),
        cookieExpiresAt: null
      });
    }
  },
  shouldRefreshCookie: () => {
    const { hasCookie, cookieExpiresAt } = get();
    if (!hasCookie) return true;
    if (cookieExpiresAt === null) return true;
    return (cookieExpiresAt * 1000) - Date.now() <= 30000;
  }
}));

