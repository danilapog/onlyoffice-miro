export const fetchAuthorization: () => Promise<{ expiresAt: number }> = async () => {
  const { board: miroBoard } = window.miro;
  const token = await miroBoard.getIdToken();

  const path = `api/authorize`;
  const response = await fetch(`${import.meta.env.VITE_MIRO_ONLYOFFICE_BACKEND}/${path}`, {
    method: 'GET',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
      'x-miro-signature': token,
    }
  });

  if (response.ok) {
    const data = await response.json();
    if (!data.expires_at)
      throw new Error("cookie not set");

    return {
      expiresAt: data.expires_at,
    }
  }

  if (response.status == 401)
    throw new Error("not authorized");

  if (response.status == 403)
    throw new Error("access denied");

  throw new Error("failed to authorize the request");
}
