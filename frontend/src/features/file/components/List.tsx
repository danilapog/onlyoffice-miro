import React, { forwardRef, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { ItemsCreateEvent, ItemsDeleteEvent } from '@mirohq/websdk-types/stable/api/ui';
import { ItemsUpdateEvent } from '@mirohq/websdk-types';

import { FileItem } from '@features/file/components/Item';
import { useFilesStore } from '@features/file/stores/useFileStore';
import { Spinner } from '@components/Spinner';

import '@features/file/components/list.css';
import { FileCreatedEvent, FileDeletedEvent, FilesAddedEvent } from '@features/manager/api/file';

interface FilesListProps extends React.HTMLAttributes<HTMLDivElement> {
}

export const FilesList = forwardRef<HTMLDivElement, FilesListProps>(({
  className,
  ...props
}, ref) => {
  const { board: miroBoard } = window.miro;
  const { t } = useTranslation();
  const {
    searchQuery,
    filteredDocuments,
    documents,
    refreshDocuments,
    loading,
    cursor,
    setObserverRef, 
    initialized,
    updateOnCreate,
    updateOnDelete,
    updateOnUpdate
  } = useFilesStore();

  const listenDocumentAddedUI = async (e: ItemsCreateEvent) => {
    const docs = e.items.filter(doc => doc.type === 'document').map(doc => {
      const documentItem = doc as any;
      return {
        id: documentItem?.id,
        data: {
          title: documentItem?.name || '',
          documentUrl: documentItem?.links?.self || '',
        },
        createdAt: documentItem?.createdAt,
        modifiedAt: documentItem?.modifiedAt,
        type: "document",
      };
    });

    if (docs.length > 0) {
      await miroBoard.events.broadcast("documents_added", { items: docs });
      miroBoard.notifications.showInfo(t("notifications.documents_added"));
    }
  }

  const listenDocumentDeletedUI = async (e: ItemsDeleteEvent) => {
    const ids = e.items.map(item => item.id);
    if (ids.length > 0) {
      await miroBoard.events.broadcast("documents_deleted", { ids });
      updateOnDelete(ids);
    }
  }

  const listenDocumentsAdded = async (_: FilesAddedEvent) => {
    miroBoard.notifications.showInfo(t("notifications.documents_added"));
  };

  const listenDocumentsDeleted = async (event: FileDeletedEvent) => {
    updateOnDelete(event.ids);
  }

  const listenDocumentUpdated = async (e: ItemsUpdateEvent) => {
    const docs = e.items.filter(doc => doc.type === 'document');
    if (docs.length > 0)
      updateOnUpdate(docs as any);
  }

  const listenDocumentCreated = async (event: FileCreatedEvent) => {
    const newDocument = { 
      id: event.id, 
      data: { 
        title: event.name, 
        documentUrl: event.links.self 
      }, 
      createdAt: event.createdAt, 
      modifiedAt: event.modifiedAt,
      type: "document",
    };
    updateOnCreate([newDocument]);
  }

  useEffect(() => {
    if (!initialized)
      refreshDocuments();

    miroBoard.ui.on("items:create", listenDocumentAddedUI);
    miroBoard.ui.on("items:delete", listenDocumentDeletedUI);
    miroBoard.ui.on("experimental:items:update", listenDocumentUpdated);

    miroBoard.events.on("document_created", listenDocumentCreated);
    miroBoard.events.on("documents_added", listenDocumentsAdded);
    miroBoard.events.on("documents_deleted", listenDocumentsDeleted);

    return () => {
      miroBoard.ui.off("items:create", listenDocumentAddedUI);
      miroBoard.ui.off("items:delete", listenDocumentDeletedUI);
      miroBoard.ui.off("experimental:items:update", listenDocumentUpdated);

      miroBoard.events.off("document_created", listenDocumentCreated);
      miroBoard.events.off("documents_added", listenDocumentsAdded);
      miroBoard.events.off("documents_deleted", listenDocumentsDeleted);
    };
  }, []);

  const docsToRender = searchQuery ? filteredDocuments : documents;
  
  return (
    <div
      ref={ref}
      className={`files-list-container ${className || ''}`}
      {...props}
    >
      {loading && !initialized && (
        <div className="files-list-container_overlay">
          <Spinner size="medium" />
        </div>
      )}

      {docsToRender.map((doc) => (
        <FileItem key={doc.id} document={doc} />
      ))}
      {!searchQuery && documents.length > 0 && cursor && (
        <div 
          ref={setObserverRef} 
          style={{ 
            height: '10px',
            margin: '20px 0'
          }}
        >
          {loading && initialized && (
            <div style={{ textAlign: 'center', padding: '10px' }}>
              <Spinner size="small" />
            </div>
          )}
        </div>
      )}
    </div>
  );
});
