import { remove as removeDiacritics } from 'diacritics';

export const containsWithoutAccents = (haystack, needle) =>
  removeDiacritics(haystack).toLowerCase().includes(removeDiacritics(needle).toLowerCase());
