export const validateAddress = (value: string): boolean => {
  if (!value) return false;

  try {
    const url = new URL(value);
    if (url.protocol !== 'http:' && url.protocol !== 'https:') {
      return false;
    }
    if (value.endsWith('/')) {
      return false;
    }
  } catch (e) {
    return false;
  }

  return true;
};

export const validateShortText = (value: string): boolean => {
  if (!value) return false;
  return value.length <= 255;
};
