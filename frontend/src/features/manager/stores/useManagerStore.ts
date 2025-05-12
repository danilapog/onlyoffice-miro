import { create } from 'zustand';
import { createFile as createMiroFile, fetchSupportedFileTypes } from '@features/manager/api/file';

interface ManagerState {
  selectedName: string;
  selectedType: string;
  loading: boolean;
  error: boolean;

  setSelectedName: (value: string) => void;
  setSelectedType: (value: string) => void;

  resetSelected: () => void;
  createFile: () => Promise<void>;
  getSupportedTypes: () => string[];
}

export const useManagerStore = create<ManagerState>((set, get) => ({
  selectedName: '',
  selectedType: '',
  loading: false,
  error: false,

  setSelectedName: (value) => set({ selectedName: value }),
  setSelectedType: (value) => set({ selectedType: value }),

  resetSelected: () => {
    set({ selectedName: '', selectedType: '' });
  },
  createFile: async () => {
    set({ loading: true, error: false });
    const { selectedName, selectedType } = get();
    if (!selectedName || !selectedType) return;
    if (!await createMiroFile(selectedName, selectedType)) {
      set({ loading: false, error: true });
      setTimeout(() => set({ error: false }), 2500);
      return;
    }

    set({ loading: false });
  },
  getSupportedTypes(): string[] {
    return fetchSupportedFileTypes();
  },
}));

