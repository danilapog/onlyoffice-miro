import { create } from 'zustand';
import { Document } from '@features/file/lib/type';
import { convertDocument, deleteDocument, fetchDocuments, navigateDocument } from '@features/file/api/file';

interface FilesState {
  searchQuery: string;
  filteredDocuments: Document[];
  documents: Document[];
  loading: boolean;
  cursor: string | null;
  activeDropdown: string | null;
  initialized: boolean;
  authError: boolean;
  serverConfigError: boolean;

  refreshDocuments: () => Promise<void>;
  toggleDropdown: (id: string | null) => void;
  navigateDocument: (document: Document) => void;
  downloadPdf: (document: Document) => void;
  deleteDocument: (document: Document) => Promise<void>;
  setSearchQuery: (query: string) => void;
  searchDocuments: () => void;
  loadMoreDocuments: () => Promise<void>;
  setObserverRef: (node: HTMLElement | null) => void;
  updateOnCreate: (documents: Document[]) => void;
  updateOnUpdate: (documents: Document[]) => void;
  updateOnDelete: (documentIds: string[]) => void;
}

const filterDocuments = (documents: Document[], query: string): Document[] => {
  if (!query.trim()) {
    return documents;
  }
  
  const lowerCaseQuery = query.toLowerCase();
  return documents.filter(doc => 
    doc.data?.title?.toLowerCase().includes(lowerCaseQuery)
  );
};

export const useFilesStore = create<FilesState>((set, get) => ({
  searchQuery: '',
  documents: [],
  filteredDocuments: [],
  loading: false,
  cursor: null,
  activeDropdown: null,
  initialized: false,
  authError: false,
  serverConfigError: false,

  refreshDocuments: async () => {
    const { loading, initialized } = get();
    if (loading) return;

    set({ loading: true, authError: false, serverConfigError: false });
    try {
      const pageable = await fetchDocuments();
      
      // If this is the first load or a refresh after document creation/deletion
      if (!initialized || !get().cursor) {
        set({ 
          documents: pageable.data,
          loading: false, 
          cursor: pageable.cursor,
          initialized: true
        });
        
        // If we have a cursor on initial load, automatically load more
        if (pageable.cursor) {
          await get().loadMoreDocuments();
        }
      } else {
        // For subsequent refreshes (after document creation/deletion)
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
        } else if (error.message === "document_server_not_configured") {
          set({ loading: false, serverConfigError: true });
        } else {
          set({ loading: false });
        }
      } else {
        set({ loading: false });
      }
    }
  },

  loadMoreDocuments: async () => {
    const { loading, cursor, initialized } = get();
    if (loading || !cursor || !initialized) return;

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

      // If we still have a cursor, continue loading
      if (pageable.cursor) {
        await get().loadMoreDocuments();
      }
    } catch (error) {
      if (error instanceof Error) {
        if (error.message === "not authorized" || error.message === "access denied") {
          set({ loading: false, authError: true });
        } else if (error.message === "document_server_not_configured") {
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

  setObserverRef: (node: HTMLElement | null) => {
    if (!node) return;

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && !get().loading && get().cursor) {
          get().loadMoreDocuments();
        }
      },
      { threshold: 0.1, rootMargin: '100px' }
    );

    observer.observe(node);
    return () => observer.disconnect();
  },

  toggleDropdown: (id: string | null) => {
    const { activeDropdown } = get();
    if (activeDropdown === id) {
      set({ activeDropdown: null });
    } else {
      set({ activeDropdown: id });
    }
  },

  navigateDocument: async (document: Document) => {
    await navigateDocument(document.id);
  },

  downloadPdf: async (document: Document) => {
    const response = await convertDocument(document.id);
    const { url, token } = response;
    const cresponse = await fetch(`${url}/converter`, {
      method: 'POST',
      body: JSON.stringify({
        token,
      }),
    });
    const { fileUrl } = await cresponse.json();
    window.open(fileUrl, '_blank');
    set({ activeDropdown: null });
  },

  deleteDocument: async (document: Document) => {
    await deleteDocument(document.id);
    set({ activeDropdown: null });
    await get().refreshDocuments();
  },

  setSearchQuery: (searchQuery: string) => {
    set({ searchQuery });
    get().searchDocuments();
  },

  searchDocuments: () => {
    const { documents, searchQuery } = get();
    const filteredDocuments = filterDocuments(documents, searchQuery);
    set({ filteredDocuments });
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

  updateOnDelete: (documentIds: string[]) => {
    set(state => {
      const docs = state.documents.filter(doc => !documentIds.includes(doc.id));
      return {
        documents: docs,
        filteredDocuments: filterDocuments(docs, state.searchQuery),
      };
    });
  }
}));
