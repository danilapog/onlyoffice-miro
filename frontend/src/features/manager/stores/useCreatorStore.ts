import { create } from 'zustand';

import { FileCreatedResponse } from '@features/manager/lib/types';

import { createFile as createMiroFile, fetchSupportedFileTypes } from '@features/manager/api/file';

interface CreatorState {
  selectedName: string;
  selectedType: string;

  loading: boolean;
  error: boolean;

  getSupportedTypes: () => string[];
  setSelectedName: (value: string) => void;
  setSelectedType: (value: string) => void;

  createFile: () => Promise<FileCreatedResponse | null>;
  resetSelected: () => void;
}

export const useCreatorStore = create<CreatorState>((set, get) => ({
  selectedName: '',
  selectedType: '',

  loading: false,
  error: false,

  getSupportedTypes(): string[] {
    return fetchSupportedFileTypes();
  },
  setSelectedName: (value) => set({ selectedName: value }),
  setSelectedType: (value) => set({ selectedType: value }),

  createFile: async (): Promise<FileCreatedResponse | null> => {
    set({ loading: true, error: false });

    const { selectedName, selectedType } = get();
    if (!selectedName || !selectedType) return null;

    const createdFile = await createMiroFile(selectedName, selectedType);
    if (!createdFile) {
      set({ loading: false, error: true });
      setTimeout(() => set({ error: false }), 2500);
      return null;
    }

    set({ loading: false });
    return createdFile;
  },
  resetSelected: () => {
    set({ selectedName: '', selectedType: '' });
  },
}));

