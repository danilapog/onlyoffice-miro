const fetchAuthorization = async () => {
  const { board } = window.miro;
  const token = await board.getIdToken();

  const path = `api/authorize`;

  const controller = new AbortController();
  const tid = setTimeout(() => controller.abort(), 3500);

  try {
    const response = await fetch(
      `${import.meta.env.VITE_MIRO_ONLYOFFICE_BACKEND}/${path}`,
      {
        method: 'GET',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
          'x-miro-signature': token,
        },
        signal: controller.signal,
      }
    );

    if (response.ok) {
      const data = await response.json();
      if (!data.expires_at) throw new Error('failed to authorize the request');

      return {
        expiresAt: data.expires_at,
      };
    }

    if (response.status === 401) throw new Error('not authorized');

    if (response.status === 403) throw new Error('access denied');

    throw new Error('failed to authorize the request');
  } catch (error: unknown) {
    if (error instanceof Error && error.name === 'AbortError')
      throw new Error('request timeout');
    throw error;
  } finally {
    clearTimeout(tid);
  }
};

export default fetchAuthorization;
