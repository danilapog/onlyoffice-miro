import { create } from 'zustand';

import { FileInfo, FileCreatedEvent, FilesAddedEvent, FilesDeletedEvent } from '@lib/types';
import { useFilesStore } from '@features/file/stores/useFileStore';

export const EmitterEvents = Object.freeze({
    MIRO_ITEMS_CREATED: "items:create",
    MIRO_ITEMS_DELETED: "items:delete",
    MIRO_ITEMS_UPDATED: "experimental:items:update",

    DOCUMENT_CREATED: "document_created",
    DOCUMENT_DELETED: "document_deleted",
    
    DOCUMENTS_ADDED: "documents_added",
    DOCUMENTS_DELETED: "documents_deleted",

    REFRESH_DOCUMENTS: "refresh_documents",
} as const);

interface EmitterState {
    emitDocumentCreated: (file: FileInfo) => Promise<void>;
    emitDocumentDeleted: (id: string) => Promise<void>;

    emitDocumentsAdded: (files: FileInfo[]) => Promise<void>;
    emitDocumentsDeleted: (ids: string[]) => Promise<void>;

    emitRefreshDocuments: () => Promise<void>;
    emitNotification: (message: string, type?: "info" | "error") => Promise<void>;
}

export const useEmitterStore = create<EmitterState>(() => ({
    emitDocumentCreated: async (file: FileInfo) => {
        const fileStore = useFilesStore.getState();
        fileStore.updateOnCreate([{
            id: file.id,
            data: {
                title: file.name,
                documentUrl: file.links.self,
            },
            createdAt: file.createdAt,
            modifiedAt: file.modifiedAt,
        }]);

        const event: FileCreatedEvent = file;
        await miro?.board.events.broadcast(EmitterEvents.DOCUMENT_CREATED, event);
    },
    async emitDocumentDeleted(id: string) {
        const event: FilesDeletedEvent = { ids: [id] };
        await miro?.board.events.broadcast(EmitterEvents.DOCUMENT_DELETED, event);
    },

    async emitDocumentsAdded(files: FileInfo[]) {
        const event: FilesAddedEvent = { files };
        await miro?.board.events.broadcast(EmitterEvents.DOCUMENTS_ADDED, event);
    },
    async emitDocumentsDeleted(ids: string[]) {
        const event: FilesDeletedEvent = { ids };
        await miro?.board.events.broadcast(EmitterEvents.DOCUMENTS_DELETED, event);
    },

    async emitRefreshDocuments() {
        const fileStore = useFilesStore.getState();
        await fileStore.refreshDocuments();
        await miro?.board.events.broadcast(EmitterEvents.REFRESH_DOCUMENTS);
    },
    async emitNotification(message: string, type: "info" | "error" = "info") {
        if (type === "info")
            await miro?.board.notifications.showInfo(message);
        else
            await miro?.board.notifications.showError(message);
    },
}));
