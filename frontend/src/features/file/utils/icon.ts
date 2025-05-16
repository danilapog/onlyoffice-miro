const getIcon = (title?: string) => {
  if (!title) return '/other.svg';

  if (title.endsWith('.docx')) return '/word.svg';
  if (title.endsWith('.xlsx')) return '/cell.svg';
  if (title.endsWith('.pptx')) return '/slide.svg';
  if (title.endsWith('.pdf')) return '/pdf.svg';

  return '/other.svg';
};

export default getIcon;
