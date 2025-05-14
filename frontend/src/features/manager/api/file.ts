export type FileCreatedResponse = {
  id: string;
  createdAt: string;
  modifiedAt: string;
  links: {
    self: string;
  };
}

export type FileCreatedEvent = {
  name: string;
  type: string;
} & FileCreatedResponse;

export type FileDeletedEvent = {
  id: string;
};

export type FilesAddedEvent = {
  items: FileCreatedResponse[];
};

export const createFile = async (name: string, type: string): Promise<FileCreatedResponse | null> => {
  try {
    const { board: miroBoard } = window.miro;
    const userPromise = miroBoard.getUserInfo();
    const boardPromise = miroBoard.getInfo();
    const tokenPromise = miroBoard.getIdToken();

    const [user, board, token] = await Promise.all([userPromise, boardPromise, tokenPromise]);
    const path = `api/files/create?uid=${user.id}&bid=${board.id}`;
    const response = await fetch(`${import.meta.env.VITE_MIRO_ONLYOFFICE_BACKEND}/${path}`, {
      method: 'POST',
      body: JSON.stringify({
        board_id: board.id,
        file_name: name,
        file_type: type,
        file_lang: board.locale,
      }),
      headers: {
        'Content-Type': 'application/json',
        'x-miro-signature': token,
      }
    });

    if (!response.ok)
      throw new Error('Failed to create a new document');

    return (await response.json()).data;
  } catch {
    return null;
  }
};

export const fetchSupportedFileTypes = () => {
  return [
    'docx',
    'xlsx',
    'pptx',
  ];
};
