export type FileInfo = {
  id: string;
  name: string;
  createdAt: string;
  modifiedAt: string;
  links: {
    self: string;
  };
  type: string;
};

export type FileCreatedEvent = FileInfo;

export type FileUpdatedEvent = FileInfo;

export type FilesDeletedEvent = {
  ids: string[];
};

export type FilesAddedEvent = {
  files: FileInfo[];
};
