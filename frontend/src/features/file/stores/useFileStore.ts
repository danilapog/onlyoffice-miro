import { create } from 'zustand';

import { Document } from '@features/file/lib/types';

import { generateRandomString } from '@utils/random';

import {
  convertDocument,
  deleteDocument,
  fetchDocuments,
  navigateDocument,
} from '@features/file/api/file';

import { useEmitterStore } from '@stores/useEmitterStore';

interface FilesState {
  documents: Document[];
  filteredDocuments: Document[];

  activeDropdown: string | null;
  cursor: string | null;
  searchQuery: string;

  initialized: boolean;
  loading: boolean;
  converting: boolean;

  authError: boolean;
  serverConfigError: boolean;

  setSearchQuery: (query: string) => void;
  setObserverRef: (node: HTMLElement | null) => void;
  toggleDropdown: (id: string | null) => void;

  loadMoreDocuments: () => Promise<void>;
  refreshDocuments: () => Promise<void>;

  navigateDocument: (document: Document) => void;
  downloadPdf: (document: Document) => void;
  deleteDocument: (document: Document) => Promise<void>;

  updateOnCreate: (documents: Document[]) => void;
  updateOnUpdate: (documents: Document[]) => void;
  updateOnDelete: (documentIds: string[]) => void;
}

const filterDocuments = (documents: Document[], query: string): Document[] => {
  if (!query.trim())
    return documents;

  const lowerCaseQuery = query.toLowerCase();
  return documents.filter(doc =>
    doc.data?.title?.toLowerCase().includes(lowerCaseQuery)
  );
};

export const useFilesStore = create<FilesState>((set, get) => ({
  documents: [],
  filteredDocuments: [],

  activeDropdown: null,
  cursor: null,
  searchQuery: '',

  initialized: false,
  loading: false,
  converting: false,

  authError: false,
  serverConfigError: false,

  setSearchQuery: (searchQuery: string) => {
    set({ searchQuery });
    const { documents } = get();
    const filteredDocuments = filterDocuments(documents, searchQuery);
    set({ filteredDocuments });
  },
  setObserverRef: (node: HTMLElement | null) => {
    if (!node) return;

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && !get().loading && get().cursor)
          get().loadMoreDocuments();
      },
      { threshold: 0.1, rootMargin: '100px' }
    );

    observer.observe(node);
    return () => observer.disconnect();
  },
  toggleDropdown: (id: string | null) => {
    const { activeDropdown } = get();
    if (activeDropdown === id)
      set({ activeDropdown: null });
    else
      set({ activeDropdown: id });
  },

  loadMoreDocuments: async () => {
    const { loading, cursor, initialized } = get();
    if (loading || !cursor || !initialized)
      return;

    set({ loading: true });
    try {
      const pageable = await fetchDocuments(cursor);
      if (!pageable.data.length) {
        set({ loading: false, cursor: null });
        return;
      }

      set(state => ({
        documents: [...state.documents, ...pageable.data],
        loading: false,
        cursor: pageable.cursor,
      }));

      if (pageable.cursor)
        await get().loadMoreDocuments();
    } catch (error) {
      if (error instanceof Error) {
        if (error.message === "not authorized" || error.message === "access denied") {
          set({ loading: false, authError: true });
        } else if (error.message === "document server configuration error") {
          set({ loading: false, serverConfigError: true });
        } else {
          set({ loading: false });
          console.error('Error loading more documents:', error);
        }
      } else {
        set({ loading: false });
        console.error('Error loading more documents:', error);
      }
    }
  },
  refreshDocuments: async () => {
    const { loading, initialized } = get();
    if (loading) return;

    set({ loading: true, authError: false, serverConfigError: false });
    try {
      const pageable = await fetchDocuments();
      if (!initialized || !get().cursor) {
        set({
          documents: pageable.data,
          loading: false,
          cursor: pageable.cursor,
          initialized: true
        });

        if (pageable.cursor)
          await get().loadMoreDocuments();
      } else {
        set(state => ({
          documents: [...pageable.data, ...state.documents],
          loading: false,
          cursor: pageable.cursor || state.cursor
        }));
      }
    } catch (error) {
      if (error instanceof Error) {
        if (error.message === "not authorized" || error.message === "access denied") {
          set({ loading: false, authError: true });
        } else if (error.message === "document server configuration error") {
          set({ loading: false, serverConfigError: true });
        } else {
          set({ loading: false });
        }
      } else {
        set({ loading: false });
      }
    }
  },

  navigateDocument: async (document: Document) => {
    await navigateDocument(document.id);
  },
  downloadPdf: async (document: Document) => {
    try {
      set({ converting: true });
      const response = await convertDocument(document.id);
      const { url, token } = response;
      const cresponse = await fetch(`${url}/converter?shardKey=${generateRandomString(8)}`, {
        method: 'POST',
        body: JSON.stringify({
          token,
        }),
      });
      const { fileUrl } = await cresponse.json();
      window.open(fileUrl, '_blank');
      set({ activeDropdown: null, converting: false });
    } catch (error) {
      console.error('Error converting the document to pdf:', error);
      set({ converting: false });
    }
  },
  deleteDocument: async (document: Document) => {
    const emitterStore = useEmitterStore.getState();
    await deleteDocument(document.id);
    await emitterStore.emitDocumentDeleted(document.id);
    set({ activeDropdown: null });
    get().updateOnDelete([document.id]);
  },

  updateOnCreate: (documents: Document[]) => {
    set(state => {
      const existing = new Set(state.documents.map(doc => doc.id));
      const docs = documents.filter(doc => !existing.has(doc.id));

      if (docs.length === 0)
        return state;

      const merged = [...state.documents, ...docs];
      return {
        documents: merged,
        filteredDocuments: filterDocuments(merged, state.searchQuery),
      };
    });
  },
  updateOnUpdate: (documents: Document[]) => {
    set(state => {
      const docsMap = new Map(documents.map(doc => [doc.id, doc]));
      const docs = [...state.documents];
      docs.forEach(doc => {
        const updatedDoc = docsMap.get(doc.id);
        if (updatedDoc) {
          doc.createdAt = updatedDoc.createdAt || doc.createdAt;
          doc.modifiedAt = updatedDoc.modifiedAt || doc.modifiedAt;
        }
      });
      return {
        documents: docs,
        filteredDocuments: filterDocuments(docs, state.searchQuery),
      };
    });
  },
  updateOnDelete: (ids: string[]) => {
    set(state => {
      const docs = state.documents.filter(doc => !ids.includes(doc.id));
      return {
        documents: docs,
        filteredDocuments: filterDocuments(docs, state.searchQuery),
      };
    });
  }
}));
