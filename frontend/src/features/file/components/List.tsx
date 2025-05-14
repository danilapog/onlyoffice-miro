import React, { forwardRef, useEffect } from 'react';
import { FileItem } from '@features/file/components/Item';
import { useFilesStore } from '@features/file/stores/useFileStore';
import { Spinner } from '@components/Spinner';

import { ItemsCreateEvent } from '@mirohq/websdk-types/stable/api/ui';
import { ItemsDeleteEvent } from '@mirohq/websdk-types/stable/api/index';
import { ItemsUpdateEvent } from '@mirohq/websdk-types';

import '@features/file/components/list.css';
import { FileCreatedEvent, FileDeletedEvent } from '@features/manager/api/file';

interface FilesListProps extends React.HTMLAttributes<HTMLDivElement> {
}

export const FilesList = forwardRef<HTMLDivElement, FilesListProps>(({
  className,
  ...props
}, ref) => {
  const { board: miroBoard } = window.miro;
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

  const listenDocumentAdded = async (e: ItemsCreateEvent) => {
    const docs = e.items.filter(doc => doc.type === 'document');
    if (docs.length > 0)
      updateOnCreate(docs as any);
  };

  const listenDocumentRemoved = async (event: FileDeletedEvent) => {
    updateOnDelete([event.id]);
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

    miroBoard.ui.on("items:create", listenDocumentAdded);
    miroBoard.events.on("document_deleted", listenDocumentRemoved);
    miroBoard.ui.on("experimental:items:update", listenDocumentUpdated);
    miroBoard.events.on("document_created", listenDocumentCreated);

    return () => {
      miroBoard.ui.off("items:create", listenDocumentAdded);
      miroBoard.events.off("document_deleted", listenDocumentRemoved);
      miroBoard.ui.off("experimental:items:update", listenDocumentUpdated);
      miroBoard.events.off("document_created", listenDocumentCreated);
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
