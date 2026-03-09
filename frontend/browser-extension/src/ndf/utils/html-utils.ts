export const stripHtml = (html: string | undefined): string | undefined => {
  if (html === undefined) return undefined;
  if (!html) return '';
  const doc = new DOMParser().parseFromString(html, 'text/html');
  const text = doc.body.textContent || '';

  // If the text still contains HTML-like tags, strip again to handle escaped HTML
  if (text.includes('<') && text.includes('>')) {
    const doc2 = new DOMParser().parseFromString(text, 'text/html');
    const text2 = doc2.body.textContent || '';
    if (text2 !== text) return text2;
  }
  return text;
};
