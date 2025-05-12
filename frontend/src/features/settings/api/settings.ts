export interface SettingsRequest {
  address: string;
  header: string;
  secret: string;
  demo: boolean;
}

export interface SettingsResponse {
  address: string;
  header: string;
  secret: string;
  demo: {
    team_id: string;
    enabled: boolean;
    started: string;
  };
}

export const saveSettings = async (settings: SettingsRequest) => {
  try {
    const { board: miroBoard } = window.miro;
    const boardPromise = miroBoard.getInfo();
    const tokenPromise = miroBoard.getIdToken();

    const [board, token] = await Promise.all([boardPromise, tokenPromise]);
    const path = `api/settings`;
    const response = await fetch(`${import.meta.env.VITE_MIRO_ONLYOFFICE_BACKEND}/${path}`, {
      method: 'POST',
      body: JSON.stringify({
        board_id: board.id,
        ...settings,
      }),
      headers: {
        'Content-Type': 'application/json',
        'x-miro-signature': token,
      },
    });

    if (response.ok) {
      await miroBoard.setAppData("settings", Date.now());
      return true;
    }

    return false;
  } catch {
    return false;
  }
};

export const fetchSettings: () => Promise<SettingsResponse> = async () => {
  const { board: miroBoard } = window.miro;
  const boardPromise = miroBoard.getInfo();
  const tokenPromise = miroBoard.getIdToken();

  const [board, token] = await Promise.all([boardPromise, tokenPromise]);
  const path = `api/settings?bid=${board.id}`;
  const response = await fetch(`${import.meta.env.VITE_MIRO_ONLYOFFICE_BACKEND}/${path}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      'x-miro-signature': token,
    }
  });

  if (response.ok) {
    const data = await response.json();
    return {
      ...data,
    }
  }

  if (response.status == 401)
    throw new Error("not authorized");

  if (response.status == 403)
    throw new Error("access denied");

  if (response.status == 404)
    return {
      address: "",
      header: "",
      secret: "",
    };

  throw new Error("failed to fetch settings");
}

export const checkSettings: () => Promise<boolean> = async () => {
  const { board: miroBoard } = window.miro;
  const settings = await miroBoard.getAppData("settings");
  return !!settings;
}; 