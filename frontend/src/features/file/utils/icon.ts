export const getIcon = (title?: string) => {
  if (!title) return '/other.svg';

  if (title.endsWith('.docx')) return '/word.svg';
  if (title.endsWith('.xlsx')) return '/cell.svg';
  if (title.endsWith('.pptx')) return '/slide.svg';

  return '/other.svg';
};
