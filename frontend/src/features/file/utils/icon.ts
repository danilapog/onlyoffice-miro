export const getIcon = (title?: string) => {
  if (!title) return '/word.svg';
  if (title.endsWith('.xlsx')) return '/cell.svg';
  if (title.endsWith('.pptx')) return '/slide.svg';

  return '/word.svg';
};
