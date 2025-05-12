import { create } from 'zustand';
import { checkSettings, fetchSettings, saveSettings } from '@features/settings/api/settings';

interface SettingsState {
  address: string;
  header: string;
  secret: string;
  persistedCredentials: boolean;
  demo: boolean;
  demoStarted: string;
  loading: boolean;
  hasSettings: boolean;
  error: string | null;

  setAddress: (value: string) => void;
  setHeader: (value: string) => void;
  setSecret: (value: string) => void;
  setDemo: (value: boolean) => void;
  saveSettings: () => Promise<void>;
  initializeSettings: () => Promise<void>;
}

export const useSettingsStore = create<SettingsState>((set, get) => ({
  address: '',
  header: '',
  secret: '',
  persistedCredentials: false,
  demo: false,
  demoStarted: '',
  loading: false,
  hasSettings: false,
  error: null,

  setAddress: (value) => set({ address: value }),
  setHeader: (value) => set({ header: value }),
  setSecret: (value) => set({ secret: value }),
  setDemo: (value ) => set({ demo: value }),
  initializeSettings: async () => {
    set({ loading: true, error: null });
    try {
      const settings = await fetchSettings();
      set({
        address: settings.address || '',
        header: settings.header || '',
        secret: settings.secret || '',
        demo: settings.demo.enabled,
        demoStarted: settings.demo.started,
        persistedCredentials: !!(settings.address && settings.header && settings.secret),
        loading: false,
        error: null,
      });
    } catch (error) {
      const isAccessDenied = error instanceof Error && (error.message === 'access denied' || error.message === 'not authorized');
      if (isAccessDenied)
        throw error;

      set({
        loading: false,
        error: isAccessDenied ? 'access denied' : null,
      });
    } finally {
      const hasSettings = await checkSettings();
      set({ hasSettings });
    }
  },

  saveSettings: async () => {
    const { address, header, secret, demo, demoStarted } = get();
    if (demoStarted && (address === '' || header === '' || secret === ''))
      return;

    set({ loading: true, error: null });
    try {
      await saveSettings({ address, header, secret, demo });
      set({ 
        address, header, secret,
        persistedCredentials: (address !== '' && header !== '' && secret !== ''),
        demoStarted: ((demoStarted == '' || !demoStarted) && demo) ? new Date().toISOString() : demoStarted,
        demo,
        loading: false,
       });
    } catch (error) {
      const isAccessDenied = error instanceof Error && error.message === 'access denied';
      set({
        loading: false,
        error: isAccessDenied ? 'access denied' : null,
      });
    }
  },
})); 
