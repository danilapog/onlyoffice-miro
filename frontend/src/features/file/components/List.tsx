import React, { forwardRef, useEffect } from 'react';
import { FileItem } from '@features/file/components/Item';
import { useFilesStore } from '@features/file/stores/useFileStore';
import { Spinner } from '@components/Spinner';

import { ItemsCreateEvent } from '@mirohq/websdk-types/stable/api/ui';
import { ItemsDeleteEvent } from '@mirohq/websdk-types/stable/api/index';

import '@features/file/components/list.css';

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
    initialized
  } = useFilesStore();

  const listenDocumentAdded = async (e: ItemsCreateEvent) => {
    const doc = e.items.find(doc => doc.type === 'document');
    if (doc) await refreshDocuments();
  };

  const listenDocumentRemoved = async (e: ItemsDeleteEvent) => {
    const doc = e.items.find(doc => doc.type === 'document');
    if (doc) await refreshDocuments();
  }

  useEffect(() => {
    miroBoard.ui.on("items:create", listenDocumentAdded);
    miroBoard.ui.on("items:delete", listenDocumentRemoved);
    refreshDocuments();

    return () => {
      miroBoard.ui.off("items:create", listenDocumentAdded);
      miroBoard.ui.off("items:delete", listenDocumentRemoved);
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
